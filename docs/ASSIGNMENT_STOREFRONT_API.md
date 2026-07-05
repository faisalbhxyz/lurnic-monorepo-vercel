# Assignment Storefront & Submission API

Lurnic API base: `https://<api-host>/v1`

সব public/student request এ **`app-key`** header লাগবে (tenant key)।  
Student-protected route এ **`Authorization: Bearer <student_jwt>`** লাগবে।

---

## 1. Storefront এ assignment দেখানো

### Option A — Course page থেকে (catalog)

Published assignment গুলো course public API তে chapter এর ভিতরে আসে।

```http
GET /v1/course/{course-slug}
app-key: <tenant_app_key>
```

**Response (relevant part):**

```json
{
  "data": {
    "id": 12,
    "title": "React Masterclass",
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
            "attachments": [
              {
                "url": "https://cdn.example.com/brief.pdf",
                "file_name": "brief.pdf",
                "mime_type": "application/pdf",
                "size": 20480
              }
            ],
            "is_published": true,
            "time_limit": 2,
            "time_limit_option": "weeks",
            "file_upload_limit": 3,
            "total_marks": 10,
            "minimum_pass_marks": 6
          }
        ]
      }
    ]
  }
}
```

UI flow:
1. Course slug দিয়ে course load করো
2. `course_chapters[].assignments[]` list render করো
3. Assignment card এ `title`, `time_limit`, `total_marks`, `file_upload_limit` দেখাও
4. Start Assignment → নিচের dedicated assignment endpoint call করো

---

### Option B — Assignment submit screen (enrolled student)

```http
GET /v1/course/{course-slug}/assignments/{assignmentId}
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

**Success `200`:**

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
    "has_submitted": false,
    "can_submit": true,
    "submission": null
  }
}
```

**Already submitted example:**

```json
{
  "data": {
    "id": 15,
    "title": "Build a Landing Page",
    "has_submitted": true,
    "can_submit": false,
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

**Common errors:**

| Status | Meaning |
|--------|---------|
| `403` | Student enrolled নয় |
| `404` | Course/assignment নেই বা unpublished |

---

## 2. Student assignment submit

`multipart/form-data` ব্যবহার করো (text + files একসাথে)।

```http
POST /v1/course/{course-slug}/assignments/{assignmentId}/submit
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
Content-Type: multipart/form-data
```

**Form fields:**

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `response_text` | string | optional* | HTML/plain text উত্তর |
| `files` | file[] | optional* | একাধিক file same key `files` দিয়ে পাঠাও |

\* `response_text` অথবা অন্তত **১টি file** দিতে হবে।

**File rules:**
- সর্বোচ্চ `file_upload_limit` টি file (assignment settings থেকে)
- এক student প্রতি assignment এ **একবারই** submit করা যাবে
- Supported: images, PDF, DOC/DOCX, ZIP (অন্যান্য type server accept করতে পারে)

**Success `201`:**

```json
{
  "message": "Assignment submitted successfully",
  "data": {
    "id": 88,
    "assignment_id": 15,
    "assignment_title": "Build a Landing Page",
    "chapter_title": "Chapter 1",
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

**Grading:**
- Assignment সবসময় **manual review** → initial `status: "pending_review"`
- Instructor grade দিলে `status: "graded"`, `score`, `percentage`, `passed` update হয়
- `passed` = `score >= minimum_pass_marks`

**Common errors:**

| Error | Meaning |
|-------|---------|
| `assignment already submitted` | দ্বিতীয়বার submit করা যাবে না |
| `response text or at least one file is required` | খালি submission |
| `maximum N file(s) allowed` | `file_upload_limit` অতিক্রম |

---

## 3. Student নিজের submission history

```http
GET /v1/student/assignment-submissions?course_id=12
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

`course_id` optional।

---

## 4. Admin — student submissions দেখা ও grade করা

Dashboard: **Course Edit → Assignments tab**

### List submissions

```http
GET /v1/private/course/{courseId}/assignment-submissions
Authorization: Bearer <admin_token>
```

### Submission detail

```http
GET /v1/private/course/{courseId}/assignment-submissions/{submissionId}
Authorization: Bearer <admin_token>
```

### Grade submission

```http
POST /v1/private/course/{courseId}/assignment-submissions/{submissionId}/grade
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "score": 8,
  "feedback": "Good structure, improve responsive layout."
}
```

**Submission status:**

| Status | Meaning |
|--------|---------|
| `pending_review` | Instructor review অপেক্ষায় |
| `graded` | Score দেওয়া হয়েছে |
| `submitted` | Received (legacy/fallback) |

---

## 5. Admin — assignment setup (curriculum)

Course create/update এ `course_chapters[].assignments[]`:

```json
{
  "id": 15,
  "title": "Build a Landing Page",
  "instructions": "<p>Submit your work</p>",
  "attachments": [
    {
      "url": "https://cdn.example.com/brief.pdf",
      "file_name": "brief.pdf",
      "mime_type": "application/pdf",
      "size": 20480
    }
  ],
  "is_published": true,
  "time_limit": 2,
  "time_limit_option": "weeks",
  "file_upload_limit": 3,
  "total_marks": 10,
  "minimum_pass_marks": 6
}
```

নতুন instructor attachment upload (course save multipart):

```
assignment_attachments[{chapterIndex}][{assignmentIndex}][] = <file>
```

Existing URL metadata JSON তে থাকবে; নতুন file গুলো upload হয়ে merge হবে।

---

## 6. Auth quick reference

```http
POST /v1/student/login
app-key: <tenant_app_key>
Content-Type: application/json

{ "email": "student@example.com", "password": "secret" }
```

Enrollment ছাড়া assignment GET/submit → `403 enrollment required`.

---

## 7. Migrations (deploy)

- `00044` — `assignment_submissions`
- `00045` — `assignment_submission_files`

```bash
cd api
goose -dir migrations mysql "<GOOSE_DBSTRING>" up
```

---

## 8. Storefront checklist

- [ ] `GET /course/{slug}` → assignment list
- [ ] Student login + token
- [ ] `GET /course/{slug}/assignments/{id}` (enrolled)
- [ ] Submit form: `response_text` + `files[]` via `multipart/form-data`
- [ ] Show `pending_review` until graded
- [ ] Poll or revisit detail after instructor grades
- [ ] Disable submit when `can_submit: false`

---

## 9. Example: submit with fetch

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
```

---

## Summary

| Feature | Status |
|---------|--------|
| Storefront assignment list (public course API) | Ready |
| Enrolled student assignment fetch | New |
| Student assignment submit (multipart) | New |
| Manual grading (admin) | New |
| Admin Assignments tab | Wired |
| Instructor attachment upload on course save | New |
