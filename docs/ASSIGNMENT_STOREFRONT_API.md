# Assignment — Storefront API, Admin & Submission Guide

**API base:** `https://<api-host>/v1`  
**Admin dashboard:** `https://<dashboard-host>/courses/{id}/update?tab=Assignments`

এই doc এ assignment এর পুরো flow আছে — admin এ assignment তৈরি, storefront এ student দেখা/submit, result, admin এ grade, আর **frontend ↔ backend gap** (কী implement হয়েছে, কী এখনো নেই)।

---

## Quick reference — সব endpoint

| Who | Method | Path | Headers |
|-----|--------|------|---------|
| Public | `GET` | `/course/{slug}` | `app-key` |
| Student | `GET` | `/course/{slug}/assignments/{assignmentId}` | `app-key` + `Bearer` |
| Student | `POST` | `/course/{slug}/assignments/{assignmentId}/submit` | `app-key` + `Bearer` (multipart) |
| Student | `GET` | `/student/assignment-submissions?course_id=` | `app-key` + `Bearer` |
| Student | `GET` | `/student/assignment-submissions/{submissionId}` | `app-key` + `Bearer` |
| Student | `POST` | `/student/login` | `app-key` |
| Admin | `POST` | `/private/course/create` | `Bearer` (multipart) |
| Admin | `PUT` | `/private/course/update/{id}` | `Bearer` (multipart) |
| Admin | `GET` | `/private/course/{courseId}/assignment-submissions` | `Bearer` |
| Admin | `GET` | `/private/course/{courseId}/assignment-submissions/{submissionId}` | `Bearer` |
| Admin | `POST` | `/private/course/{courseId}/assignment-submissions/{submissionId}/grade` | `Bearer` |

- **Student routes:** `Authorization: Bearer <student_jwt>` (from `/student/login`)
- **Admin routes:** `Authorization: Bearer <admin_jwt>` (from `/user/login`)
- **Public/student tenant routes:** `app-key: <tenant_app_key>`

**Resubmit:** Same `POST .../submit` upserts while `status === "pending_review"` and timer not expired (no separate `PUT` route).

---

## End-to-end flow

```mermaid
sequenceDiagram
  participant Admin
  participant API
  participant Storefront
  participant Student

  Admin->>API: Create/update course + assignments (Curriculum tab)
  Student->>API: POST /student/login
  Student->>API: Enrolled in course
  Storefront->>API: GET /course/{slug}
  Storefront->>Student: Assignment list (chapter items)
  Student->>API: GET /course/{slug}/assignments/{id}
  Note over API: timer starts (deadline_at); full submission if already submitted
  Student->>API: POST .../submit (multipart, create or resubmit)
  API->>Student: full submission (response_text + files)
  Note over Storefront: revisit page → full submission on detail GET
  Admin->>API: GET assignment-submissions (Assignments tab)
  Admin->>API: POST .../grade
  Student->>API: GET assignment detail again → graded score
```

---

## Implemented features (storefront)

| Feature | Status |
|---------|--------|
| Live countdown (`deadline_at`, `seconds_remaining`, `started_at`) | ✅ Timer starts on first `GET` |
| Full submission on detail `GET` (`response_text`, `files[]`) | ✅ |
| Student submission detail `GET /student/assignment-submissions/{id}` | ✅ |
| Resubmit before grade (same `POST`) + `can_edit` | ✅ |
| Server file size (2 MB) + MIME allowlist | ✅ |
| HTML sanitize + max length (`50000`) | ✅ |
| Upload policy fields on assignment `GET` | ✅ `max_file_size_bytes`, `allowed_mime_types` |

---

## Part 1 — Admin: Assignment তৈরি ও manage

### 1.1 Curriculum তে assignment add

1. **Courses → Create/Edit → Curriculum**
2. Chapter এ **+ → Assignment**
3. Title, instructions (HTML), attachments, time limit, file limit, marks
4. Assignment **Save** → পুরো course **Save**

Published assignment (`is_published: true`) storefront এ দেখা যায়।

### 1.2 Assignment settings (API + UI)

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `title` | string | ✅ | |
| `instructions` | string (HTML) | ✅ | Instructor brief; rich text |
| `attachments` | `Attachment[]` | ❌ | Instructor files (PDF, etc.) |
| `is_published` | boolean | — | `false` → student `404` on dedicated GET |
| `time_limit` | number | ✅ | `0` = no deadline enforcement |
| `time_limit_option` | enum | — | `minutes`, `hours`, `days`, `weeks`, `months` |
| `file_upload_limit` | number | — | Max files per submission (default `1` if unset) |
| `total_marks` | number | ✅ | `max_score` on submission |
| `minimum_pass_marks` | number | ✅ | `passed = score >= minimum_pass_marks` after grade |
| `position` | number | — | Chapter item order |

**Attachment object (instructor + student submission files):**

```json
{
  "url": "https://cdn.example.com/brief.pdf",
  "file_name": "brief.pdf",
  "mime_type": "application/pdf",
  "size": 20480
}
```

### 1.3 Instructor attachment upload (course save)

Multipart field pattern (course create/update):

```
assignment_attachments[{chapterIndex}][{assignmentIndex}][] = <file>
```

Existing attachment metadata stays in JSON; new files upload and merge on save.

---

## Part 2 — Storefront: Assignment list (public course)

```http
GET /v1/course/{course-slug}
app-key: <tenant_app_key>
```

**Relevant response:**

```json
{
  "data": {
    "slug": "react-masterclass",
    "course_chapters": [
      {
        "id": 3,
        "title": "Chapter 1",
        "assignments": [
          {
            "id": 15,
            "course_id": 12,
            "chapter_id": 3,
            "title": "Build a Landing Page",
            "instructions": "<p>Submit your HTML/CSS project</p>",
            "attachments": [],
            "is_published": true,
            "time_limit": 2,
            "time_limit_option": "weeks",
            "file_upload_limit": 3,
            "total_marks": 10,
            "minimum_pass_marks": 6,
            "created_at": "2026-06-01T08:00:00Z",
            "updated_at": "2026-06-01T08:00:00Z"
          }
        ]
      }
    ]
  }
}
```

**UI:** Render chapter assignments; navigate to submit screen with `course slug` + `assignment id`.

**Course viewer:** Same chapter payload — enrolled student can call dedicated assignment endpoints without a separate “inline” API.

---

## Part 3 — Student: Assignment detail (enrolled)

```http
GET /v1/course/{course-slug}/assignments/{assignmentId}
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

### 3.1 Not yet submitted

```json
{
  "data": {
    "id": 15,
    "course_id": 12,
    "chapter_id": 3,
    "title": "Build a Landing Page",
    "instructions": "<p>Submit your HTML/CSS project</p>",
    "attachments": [],
    "is_published": true,
    "time_limit": 2,
    "time_limit_option": "weeks",
    "file_upload_limit": 3,
    "total_marks": 10,
    "minimum_pass_marks": 6,
    "created_at": "2026-06-01T08:00:00Z",
    "updated_at": "2026-06-01T08:00:00Z",
    "has_submitted": false,
    "can_submit": true,
    "can_edit": false,
    "started_at": "2026-07-07T08:00:00Z",
    "deadline_at": "2026-07-14T08:00:00Z",
    "seconds_remaining": 604740,
    "max_file_size_bytes": 2097152,
    "allowed_mime_types": ["application/pdf", "image/jpeg", "image/png", "image/gif", "image/webp", "application/zip", "application/x-zip-compressed", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "text/plain"],
    "max_response_text_length": 50000,
    "submission": null
  }
}
```

### 3.2 Already submitted (pending review — editable)

```json
{
  "data": {
    "has_submitted": true,
    "can_submit": true,
    "can_edit": true,
    "deadline_at": "2026-07-14T08:00:00Z",
    "seconds_remaining": 604000,
    "submission": {
      "id": 88,
      "score": 0,
      "max_score": 10,
      "percentage": 0,
      "passed": false,
      "status": "pending_review",
      "submitted_at": "2026-06-25T10:30:00Z",
      "response_text": "<p>Here is my solution...</p>",
      "files": [
        {
          "id": 201,
          "url": "https://cdn.example.com/project.zip",
          "file_name": "project.zip",
          "mime_type": "application/zip",
          "size": 1048576
        }
      ]
    }
  }
}
```

### 3.3 Graded (locked)

```json
{
  "data": {
    "has_submitted": true,
    "can_submit": false,
    "can_edit": false,
    "submission": {
      "id": 88,
      "score": 8,
      "max_score": 10,
      "percentage": 80,
      "passed": true,
      "status": "graded",
      "submitted_at": "2026-06-25T10:30:00Z"
    }
  }
}
```

### 3.4 Student-only fields (logic)

| Field | Type | When set | Meaning |
|-------|------|----------|---------|
| `has_submitted` | boolean | always | Student has a row in `assignment_submissions` |
| `can_submit` | boolean | always | `true` if not submitted and timer valid, **or** pending review + timer valid (resubmit) |
| `can_edit` | boolean | always | `true` when `status === "pending_review"` and timer not expired |
| `submission` | object \| null | `has_submitted` | Includes `response_text` and `files[]` |

**`can_submit` / `can_edit` rules:**

```
can_edit    = has_submitted && status == pending_review && !timer_expired
can_submit  = (!has_submitted && !timer_expired) || can_edit
```

After `graded` → both `false`. After timer expires → both `false` (existing submission remains visible).

### 3.5 Deadline / countdown — **implemented**

Quiz-style attempt session (`assignment_attempt_sessions`). First `GET .../assignments/{id}` creates session; timer derived from `time_limit` + `time_limit_option` (`0` = unlimited).

| Field | Meaning |
|-------|---------|
| `started_at` | When student first opened assignment |
| `deadline_at` | Expiry time (omitted if no time limit) |
| `seconds_remaining` | Live countdown seconds (`0` if expired) |

Submit/resubmit rejected with `assignment time limit exceeded` when expired.

### 3.6 Errors

| Status | `error` | Meaning |
|--------|---------|---------|
| `403` | `enrollment required` | Not enrolled |
| `404` | `course not found` | Bad slug |
| `404` | `assignment not found` | Unpublished or wrong id |

---

## Part 4 — Student: Submit assignment

```http
POST /v1/course/{course-slug}/assignments/{assignmentId}/submit
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
Content-Type: multipart/form-data
```

### 4.1 Form fields

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `response_text` | string | optional* | HTML or plain text (see §7) |
| `files` | file[] | optional* | Same field name `files` for each file |

\* At least one of `response_text` (non-empty after trim) **or** one file.

### 4.2 Server validation (implemented)

| Rule | Enforced |
|------|----------|
| Enrollment | ✅ |
| Assignment published | ✅ |
| One submission per student per assignment | ✅ |
| `len(files) <= file_upload_limit` | ✅ (default limit `1` if unset) |
| Non-empty content | ✅ |
| Per-file max size 2 MB | ✅ |
| MIME allowlist | ✅ |
| Time limit on submit | ✅ |
| HTML sanitize (bluemonday UGC) | ✅ |
| Max response text 50000 chars | ✅ |

### 4.3 Success `201`

```json
{
  "message": "Assignment submitted successfully",
  "data": {
    "id": 88,
    "assignment_id": 15,
    "assignment_title": "Build a Landing Page",
    "chapter_id": 3,
    "chapter_title": "Chapter 1",
    "student_id": 42,
    "student_name": "Rahim Uddin",
    "student_email": "rahim@example.com",
    "score": 0,
    "max_score": 10,
    "percentage": 0,
    "passed": false,
    "status": "pending_review",
    "submitted_at": "2026-06-25T10:30:00Z",
    "file_count": 2,
    "response_text": "<p>Here is my solution...</p>",
    "files": [
      {
        "id": 201,
        "url": "https://cdn.example.com/submission.zip",
        "file_name": "project.zip",
        "mime_type": "application/zip",
        "size": 1048576
      }
    ]
  }
}
```

**Note:** This is the **only** student response that includes full `response_text` + `files[]` today. Cache it client-side if you need to show content after navigation without an extra API.

### 4.4 Errors

| `error` | Meaning |
|---------|---------|
| `assignment time limit exceeded` | Past `deadline_at` |
| `assignment already graded` | Resubmit after grade |
| `file "x" exceeds maximum size of 2 MB` | File too large |
| `file type "text/html" is not allowed` | MIME not allowed |
| `response text exceeds maximum length of 50000 characters` | Text too long |
| `response text or at least one file is required` | Empty body |
| `maximum N file(s) allowed` | Exceeds `file_upload_limit` |
| `enrollment required` | `403` |
| `course not found` / `assignment not found` | `404` |

### 4.5 Resubmit — **supported**

Same `POST .../submit` while `can_edit: true` (`pending_review` + timer valid):

- Send new `response_text` and/or new `files`
- Omitting `files` keeps existing files
- Omitting `response_text` keeps existing text
- After grade → `assignment already graded`

---

## Part 5 — Student: Submission history (list)

```http
GET /v1/student/assignment-submissions?course_id=12
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

`course_id` query param optional — omit for all courses.

### 5.1 Response `200` (full detail per item)

```json
{
  "data": [
    {
      "id": 88,
      "assignment_id": 15,
      "assignment_title": "Build a Landing Page",
      "status": "graded",
      "submitted_at": "2026-06-25T10:30:00Z",
      "score": 8,
      "max_score": 10,
      "file_count": 2,
      "response_text": "<p>Here is my solution...</p>",
      "files": [
        {
          "id": 201,
          "url": "https://cdn.example.com/project.zip",
          "file_name": "project.zip",
          "mime_type": "application/zip",
          "size": 1048576
        }
      ]
    }
  ]
}
```

### 5.2 Single submission detail

```http
GET /v1/student/assignment-submissions/{submissionId}
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

Returns same shape as one list item (full `response_text` + `files[]`). Only the owning student can access.

### 5.3 List vs detail

| Field | List endpoint | Assignment `GET` `submission` | Admin detail |
|-------|---------------|------------------------------|--------------|
| `response_text` | ✅ | ✅ | ✅ |
| `files[]` | ✅ | ✅ | ✅ |
| `file_count` | ✅ | via `files.length` | ✅ |

---

## Part 6 — Submission status & grading

### 6.1 Status values

| Status | Set when | Meaning |
|--------|----------|---------|
| `pending_review` | Default on submit | Awaiting instructor grade |
| `graded` | After `POST .../grade` | Score finalized |
| `submitted` | Legacy enum value | **Not used** on new submits; treat like received if seen in old data |

### 6.2 Score fields

| Field | Before grade | After grade |
|-------|--------------|-------------|
| `score` | `0` | Instructor `score` |
| `max_score` | `total_marks` | `total_marks` |
| `percentage` | `0` | `round(score/max*100, 2)` |
| `passed` | `false` | `score >= minimum_pass_marks` |

Graded submissions always have numeric `score` (not null) in current implementation.

### 6.3 Storefront results table mapping

| UI column | API field |
|-----------|-----------|
| Date | `submission.submitted_at` |
| Total Marks | `total_marks` (assignment) or `submission.max_score` |
| Pass Marks | `minimum_pass_marks` |
| Earned Marks | `submission.score` (when `status === "graded"`) |
| Result | `submission.passed` or derive from score vs minimum |

Show “Pending review” when `status === "pending_review"` and `score === 0`.

---

## Part 7 — `response_text` (HTML)

| Topic | Current behavior |
|-------|------------------|
| Format | Sanitized HTML (bluemonday UGC policy) |
| Max length | `50000` characters (`max_response_text_length` on assignment GET) |
| XSS | `<script>` and unsafe tags stripped server-side |

---

## Part 8 — File upload

### 8.1 Implemented

- Multiple files via `files` form key
- Count capped by `file_upload_limit`
- Upload to CDN (Bunny); metadata saved in `assignment_submission_files`
- MIME from `Content-Type` header (fallback `application/octet-stream`)

### 8.2 Enforced server-side

Returned on assignment `GET`:

```json
{
  "max_file_size_bytes": 2097152,
  "allowed_mime_types": [
    "application/pdf",
    "image/jpeg",
    "image/png",
    "image/gif",
    "image/webp",
    "application/zip",
    "application/x-zip-compressed",
    "application/msword",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "text/plain"
  ]
}
```

Errors: `file "foo.pdf" exceeds maximum size of 2 MB`, `file type "text/html" is not allowed`.

---

## Part 9 — Admin: Submissions & grading

**Dashboard:** Course Edit → **Assignments** tab (after course saved).

### 9.1 List submissions

```http
GET /v1/private/course/{courseId}/assignment-submissions
Authorization: Bearer <admin_token>
```

Same list shape as student list (all students in course).

### 9.2 Submission detail

```http
GET /v1/private/course/{courseId}/assignment-submissions/{submissionId}
Authorization: Bearer <admin_token>
```

Includes `response_text`, `files[]`, `instructor_feedback` when graded.

### 9.3 Grade submission

```http
POST /v1/private/course/{courseId}/assignment-submissions/{submissionId}/grade
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "score": 8,
  "feedback": "Good structure, improve responsive layout."
}
```

| Field | Required | Validation |
|-------|----------|------------|
| `score` | ✅ | `>= 0`, `<= assignment.total_marks` |
| `feedback` | ❌ | Optional string |

**Success `200`:** Updated submission detail (`status: "graded"`, `graded_at` set).

Grading may trigger certificate issue if course certificate rules require assignments.

---

## Part 10 — Auth & enrollment

```http
POST /v1/student/login
app-key: <tenant_app_key>
Content-Type: application/json

{ "email": "student@example.com", "password": "secret" }
```

Assignment `GET` / `POST .../submit` without enrollment → `403 enrollment required`.

---

## Part 11 — Migrations

| Migration | Table |
|-----------|-------|
| `00018` | `course_assignments` |
| `00044` | `assignment_submissions` |
| `00045` | `assignment_submission_files` |
| `00047` | `position` on assignments |
| `00056` | `assignment_attempt_sessions` (timer) |

```bash
cd api
goose -dir migrations mysql "<GOOSE_DBSTRING>" up
```

---

## Part 12 — Storefront integration checklist

- [ ] `GET /course/{slug}` → list assignments in chapters
- [ ] Student login + JWT
- [ ] `GET /course/{slug}/assignments/{id}` when opening assignment
- [ ] `GET /course/{slug}/assignments/{id}` → use `deadline_at` + `seconds_remaining` for countdown
- [ ] Submit / resubmit: same `POST` while `can_edit: true`
- [ ] After submit: reload detail GET for `submission.response_text` + `files[]`
- [ ] Client file validation: use `max_file_size_bytes` + `allowed_mime_types` from API
- [ ] Disable submit when `can_submit: false` (graded or timer expired)
- [ ] Show `pending_review` until `status === "graded"`
- [ ] Results: `score`, `max_score`, `passed`, `minimum_pass_marks`
- [ ] Optional: `GET /student/assignment-submissions?course_id=` for dashboard table (scores only)

---

## Part 13 — Example: submit with fetch

```javascript
const fd = new FormData();
fd.append("response_text", "<p>My answer</p>");
for (const file of selectedFiles) {
  fd.append("files", file);
}

const res = await fetch(
  `${API_URL}/course/${slug}/assignments/${assignmentId}/submit`,
  {
    method: "POST",
    headers: {
      "app-key": TENANT_APP_KEY,
      Authorization: `Bearer ${studentToken}`,
    },
    body: fd,
  }
);

if (!res.ok) {
  const err = await res.json();
  throw new Error(err.error ?? "Submit failed");
}

const { data } = await res.json();
// data includes response_text + files; detail GET also returns full submission
```

---

## Summary

| Feature | Status |
|---------|--------|
| Assignment list on public course API | ✅ |
| Enrolled student assignment GET + timer | ✅ |
| Student submit / resubmit (multipart) | ✅ |
| Full submission on detail GET + student list/detail | ✅ |
| `can_edit` / `can_submit` logic | ✅ |
| Manual grading (admin) | ✅ |
| Server file size / MIME validation | ✅ |
| HTML sanitize | ✅ |

**Deploy:** run migration `00056` before release.

**Source of truth (Go):** `api/internal/modules/assignment/` — `service.go`, `response.go`, `router.go`, `helpers.go`, `sanitize.go`.
