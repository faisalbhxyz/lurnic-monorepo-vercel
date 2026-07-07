# Quiz — Storefront API, Admin & Submission Guide

**API base:** `https://<api-host>/v1`  
**Admin dashboard:** `https://<dashboard-host>/courses/{id}/update?tab=Quizzes`

এই doc এ quiz এর পুরো flow আছে — admin এ quiz তৈরি, storefront এ student দেখা/submit/skip, result, admin এ submission review + instructor feedback, আর **backend ↔ storefront status**।

---

## Quick reference — সব endpoint

| Who | Method | Path | Headers |
|-----|--------|------|---------|
| Public | `GET` | `/course/{slug}` | `app-key` |
| Student | `GET` | `/course/{slug}/quizzes/{quizId}` | `app-key` + `Bearer` |
| Student | `GET` | `/course/{slug}/quizzes/{quizId}/questions/{questionIndex}` | `app-key` + `Bearer` |
| Student | `POST` | `/course/{slug}/quizzes/{quizId}/submit` | `app-key` + `Bearer` |
| Student | `POST` | `/course/{slug}/quizzes/{quizId}/skip` | `app-key` + `Bearer` |
| Student | `GET` | `/student/quiz-submissions?course_id=` | `app-key` + `Bearer` |
| Student | `GET` | `/student/quiz-submissions/{submissionId}` | `app-key` + `Bearer` |
| Student | `GET` | `/course/{slug}/progress` | `app-key` + `Bearer` |
| Student | `POST` | `/student/login` | `app-key` |
| Admin | `POST` | `/private/course/create` | `Bearer` (multipart) |
| Admin | `PUT` | `/private/course/update/{id}` | `Bearer` (multipart) |
| Admin | `GET` | `/private/course/{courseId}/quiz-submissions` | `Bearer` |
| Admin | `GET` | `/private/course/{courseId}/quiz-submissions/{submissionId}` | `Bearer` |
| Admin | `POST` | `/private/course/{courseId}/quiz-submissions/{submissionId}/feedback` | `Bearer` |

- **Student routes:** `Authorization: Bearer <student_jwt>` (from `/student/login`)
- **Admin routes:** `Authorization: Bearer <admin_jwt>` (from `/user/login`)
- **Public/student tenant routes:** `app-key: <tenant_app_key>`

---

## End-to-end flow

```mermaid
sequenceDiagram
  participant Admin
  participant API
  participant Storefront
  participant Student

  Admin->>API: Create/update course + quizzes (Curriculum tab)
  Student->>API: POST /student/login
  Student->>API: Enrolled in course
  Storefront->>API: GET /course/{slug}
  Storefront->>Student: Quiz list
  Student->>API: GET /course/{slug}/quizzes/{id}
  Note over API: attempt session starts (timer + question order fixed)
  alt single_quiz_view
    Storefront->>API: GET .../quizzes/{id}/questions/{index}
  end
  Student->>API: POST .../submit
  API->>Student: score, percentage, passed, answers
  Admin->>API: GET quiz-submissions (Quizzes tab)
  Admin->>API: GET quiz-submissions/{id} (detail page)
  Admin->>API: POST quiz-submissions/{id}/feedback (optional)
  Student->>API: GET /student/quiz-submissions/{id}
  Storefront->>Student: Show instructor_feedback on result / history detail
```

---

## Part 1 — Admin: Quiz তৈরি ও manage

### 1.1 Curriculum তে quiz add

1. **Courses → Create/Edit → Curriculum**
2. Chapter এ **+ → Quiz**
3. Title, optional instructions, settings, **Add Questions**
4. প্রতিটি question এ:
   - **Single choice / Multiple choice:** options + correct answer (radio/checkbox)
   - **True/False:** True বা False select
5. Quiz **Save** → পুরো course **Save**

Quiz **published** (`is_published: true`) হলে storefront এ দেখা যাবে।  
প্রতিটি question এ **correct answer** থাকলে auto-grade কাজ করবে।

### 1.2 Quiz settings (API + UI)

| Field | Type | Required | Special values | Behaviour |
|-------|------|----------|----------------|-----------|
| `title` | string | ✅ | — | Quiz title |
| `instructions` | string | ❌ | empty allowed | Rich text; optional |
| `is_published` | boolean | — | — | `false` হলে student attempt করতে পারবে না |
| `randomize_questions` | boolean | — | — | প্রতি attempt-এ question order shuffle হয় (session-এ fixed) |
| `single_quiz_view` | boolean | — | — | `true` হলে এক সময়ে একটা question; paginated endpoint ব্যবহার করো |
| `time_limit` | number | ✅ | `0` = no limit | Timer শুরু হয় attempt session create হলে |
| `time_limit_option` | enum | — | `minutes`, `hours`, `days`, `weeks`, `months` | `time_limit` এর unit |
| `total_visible_questions` | number | — | `0` = all | Attempt-এ কতটা question দেখাবে |
| `reveal_answers` | boolean | — | — | Submit-এর পর correct answer + explanation |
| `enable_retry` | boolean | — | — | `false` = শুধু ১ attempt |
| `retry_attempts` | number | — | `0` = unlimited | শুধু `enable_retry: true` হলে কার্যকর; UI-তো তখনই দেখায় |
| `minimum_pass_percentage` | number | ✅ | `0`–`100` | `passed = percentage >= minimum_pass_percentage` |

**Admin UI notes:**

- **Instructions** — optional (no asterisk)
- **Retry Attempts** — শুধু **Enable Retry** on থাকলে দেখায়
- **Time Limit `0`** — no time limit
- **Total Visible Questions `0`** — সব question দেখাবে

### 1.3 Question fields (API + UI)

| Field | Type | UI | Notes |
|-------|------|-----|-------|
| `title` | string | ✅ | Required |
| `type` | `single_choice` \| `multiple_choice` \| `true_false` | ✅ | |
| `marks` | number | ✅ | |
| `options` | `[{ id, text }]` | ✅ | MCQ তে required (min 2) |
| `correct_answer` | JSON | ✅ | Auto-grade এর জন্য |
| `answer_explanation` | HTML string | ✅ | Optional; reveal হলে student দেখবে |
| `answer_required` | boolean | ✅ | Submit-এ empty থাকলে reject |

**`correct_answer` format:**

```json
// single_choice
{ "value": "a" }

// true_false
{ "value": true }

// multiple_choice
{ "values": ["a", "c"] }
```

**`options` example:**

```json
[
  { "id": "a", "text": "Markup language" },
  { "id": "b", "text": "Programming language" }
]
```

Course save এ `course_chapters` JSON এর ভিতরে `quizzes[].questions[]` হিসেবে যায় (`PUT /private/course/update/{id}` বা `POST /private/course/create`).

**Example quiz payload (course chapter JSON):**

```json
{
  "title": "HTML Basics Quiz",
  "instructions": "",
  "is_published": true,
  "randomize_questions": true,
  "single_quiz_view": false,
  "time_limit": 30,
  "time_limit_option": "minutes",
  "total_visible_questions": 0,
  "reveal_answers": true,
  "enable_retry": true,
  "retry_attempts": 2,
  "minimum_pass_percentage": 60,
  "questions": []
}
```

### 1.4 Admin: Submission review (Dashboard)

**Course Edit → Quizzes tab** (`?tab=Quizzes`)

| Feature | আছে |
|---------|-----|
| সব submission list | ✅ |
| Quiz / chapter / student name / email | ✅ |
| Marks `score/max_score (%)` | ✅ |
| Status filter (All / Evaluate / Pending) | ✅ |
| Row click → **detail page** (`?tab=Quizzes&submission={id}`) | ✅ |
| Summary row: attempt by, date, question count, quiz/attempt time, marks, pass marks, correct/incorrect, result | ✅ |
| **Quiz Overview** table: given answer vs correct answer, per-question result | ✅ |
| **Instructor Feedback** (rich text + Update) | ✅ |
| Manual mark edit | ❌ (auto-grade; feedback only) |

**Evaluate** filter = `status: pending_review` (যেখানে `correct_answer` ছিল না বা grade হয়নি)।

Detail page API call করে:

```http
GET /v1/private/course/{courseId}/quiz-submissions/{submissionId}
Authorization: Bearer <admin_token>
```

Admin detail এ সবসময় **correct answer + explanation** দেখায় (review এর জন্য)।

**Detail response (admin) — extra fields:**

| Field | Meaning |
|-------|---------|
| `total_questions` | Attempt-এ মোট question |
| `correct_count` / `incorrect_count` / `unanswered_count` | Per-question summary |
| `pass_marks` | `max_score × minimum_pass_percentage / 100` |
| `minimum_pass_percentage` | Quiz pass threshold |
| `quiz_time_seconds` | Allowed time (`0` = unlimited) |
| `attempt_time_seconds` | Student-এর নেওয়া সময় (session `started_at` → submit) |
| `answers[].options` | Option `id` + `text` (readable Given/Correct columns) |
| `answers[].question_marks` | Max marks per question |
| `instructor_feedback` | Admin-এর HTML feedback (null যদি না দেওয়া হয়) |

### 1.5 Admin: Instructor feedback

Dashboard detail page-এর নিচে **Instructor Feedback** rich text editor আছে। Save করতে:

```http
POST /v1/private/course/{courseId}/quiz-submissions/{submissionId}/feedback
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "feedback": "<p>Well done! Review Q3 again.</p>"
}
```

| Field | Type | Notes |
|-------|------|-------|
| `feedback` | string (HTML) \| `null` | Empty string পাঠালে clear করতে `null` ব্যবহার করো |

**Success `200`:**

```json
{
  "message": "Instructor feedback updated successfully",
  "data": { "...full QuizSubmissionDetail..." }
}
```

Student storefront এ feedback দেখাতে হলে `GET /student/quiz-submissions/{submissionId}` (বা submit `201` `data`) থেকে `instructor_feedback` bind করো — নিচে §3.7।

---

## Part 2 — Storefront: Quiz দেখানো

### 2.1 Course page — quiz list

```http
GET /v1/course/{course-slug}
app-key: <tenant_app_key>
```

Published quiz গুলো `data.course_chapters[].quizzes[]` তে আসে।

**পাঠানো হয়:** `id`, `title`, `instructions`, settings, `questions[]` (with `options`)  
**পাঠানো হয় না:** `correct_answer`, `answer_explanation` (security)

> Public course response সব question দেখায় (preview/listing)। Actual attempt-এ student authenticated endpoint ব্যবহার করবে — সেখানে session, timer, randomize/limit apply হয়।

**Storefront UI:**
1. Course load করো
2. Chapter অনুযায়ী quiz cards দেখাও
3. Intro screen (optional) — `attempts_used` / `can_retry` এর জন্য `GET /student/quiz-submissions?course_id=` ব্যবহার করো (নিচে §2.2 note)
4. “Start quiz” → enrolled student হলে `GET .../quizzes/{id}` call করো (**attempt session + timer শুরু**)

### 2.2 Quiz attempt — start session (enrolled student only)

```http
GET /v1/course/{course-slug}/quizzes/{quizId}
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

এই call **attempt session** তৈরি/পুনরায় ব্যবহার করে (`quiz_attempt_sessions` table)। এক attempt-এর মধ্যে:

- Question order **fixed** থাকে (randomize/limit একবার apply হয়)
- Timer শুরু হয় (`time_limit > 0` হলে)
- Page refresh করলে same session reuse হয় (timer reset হয় না)

> **⚠️ Intro screen gap:** কোনো আলাদা “quiz preview” endpoint নেই। `GET .../quizzes/{id}` call করলেই session create/reuse হয় এবং timer চালু হয়। Intro-তে attempt info দেখাতে **session শুরু করো না** — `GET /student/quiz-submissions?course_id=` থেকে `quiz_id` অনুযায়ী attempt count derive করো। “Start Quiz” ক্লিকে তবেই `GET .../quizzes/{id}` call করো।

**Response fields (quiz object + attempt metadata):**

| Field | Meaning |
|-------|---------|
| `attempts_used` | আগে কতবার submit/forfeit হয়েছে |
| `can_retry` | আবার নতুন attempt করা যাবে কিনা |
| `attempt_number` | বর্তমান attempt (1-based) |
| `display_mode` | `"all"` বা `"single"` |
| `total_questions` | এই attempt-এ মোট কত question |
| `current_question_index` | Single view-এ প্রথম question = `0` |
| `started_at` | ISO 8601 — session শুরুর সময় |
| `expires_at` | ISO 8601 — time limit থাকলে |
| `seconds_remaining` | Timer countdown (null = no limit) |
| `questions[]` | Attempt-এ visible questions (`single` mode-এ প্রথমটা only) |

**Example response (abbreviated):**

```json
{
  "data": {
    "id": 9,
    "title": "HTML Basics Quiz",
    "instructions": "",
    "single_quiz_view": true,
    "time_limit": 30,
    "time_limit_option": "minutes",
    "randomize_questions": true,
    "total_visible_questions": 0,
    "reveal_answers": true,
    "enable_retry": true,
    "retry_attempts": 2,
    "minimum_pass_percentage": 60,
    "attempts_used": 0,
    "can_retry": true,
    "attempt_number": 1,
    "display_mode": "single",
    "total_questions": 5,
    "current_question_index": 0,
    "started_at": "2026-07-07T00:30:00Z",
    "expires_at": "2026-07-07T01:00:00Z",
    "seconds_remaining": 1800,
    "questions": [
      {
        "id": 41,
        "title": "What is HTML?",
        "type": "single_choice",
        "options": [{ "id": "a", "text": "Markup language" }]
      }
    ]
  }
}
```

**Errors:**

| HTTP | Reason |
|------|--------|
| `403` | Enrolled নয় |
| `404` | Course/quiz নেই বা unpublished |
| `400` | `quiz retry is disabled` |
| `400` | `maximum quiz attempts reached` |

### 2.3 Single quiz view — paginated questions

`single_quiz_view: true` হলে পরের question-গুলো এই endpoint দিয়ে load করো:

```http
GET /v1/course/{course-slug}/quizzes/{quizId}/questions/{questionIndex}
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

- `questionIndex` = **0-based** (`0` = first, `total_questions - 1` = last)
- Same active attempt session ব্যবহার করে (timer + order consistent)
- `single_quiz_view: false` হলে `400 single quiz view is disabled for this quiz`

**Response fields:**

| Field | Meaning |
|-------|---------|
| `question_index` | যে index request করা হয়েছে |
| `total_questions` | attempt-এ মোট question |
| `display_mode` | `"single"` |
| `started_at`, `expires_at`, `seconds_remaining` | Timer info |
| (question fields) | `id`, `title`, `type`, `options`, … |

**Storefront UI flow (`single_quiz_view`):**

1. `GET .../quizzes/{id}` → question `0` + metadata
2. Next button → `GET .../questions/1`, then `2`, …
3. Last question-এ Submit → `POST .../submit`

### 2.4 Quiz settings → storefront behaviour

| Setting | API enforcement | Storefront UI |
|---------|-------------------|---------------|
| `is_published` | Unpublished = 404 on attempt | Hide or disable unpublished quizzes |
| `randomize_questions` | ✅ Per-attempt shuffle (session-fixed) | No extra work |
| `total_visible_questions` | ✅ `0` = all; `N` = subset | Show `total_questions` from response |
| `time_limit` + `time_limit_option` | ✅ Server enforces on submit; auto-forfeit on expiry | Render countdown from `seconds_remaining` |
| `single_quiz_view` | ✅ First Q on main GET; rest via `/questions/{index}` | Paginate UI |
| `reveal_answers` | ✅ On submit response | Show/hide correct answers |
| `enable_retry` + `retry_attempts` | ✅ Server enforces | Show retry if `can_retry: true` |
| `minimum_pass_percentage` | ✅ Sets `passed` on submit | Pass/fail badge |

### 2.5 Time limit behaviour

| `time_limit` | Behaviour |
|--------------|-----------|
| `0` | No timer; `expires_at` / `seconds_remaining` = null |
| `> 0` | Session `started_at` থেকে countdown; unit = `time_limit_option` |

**On expiry (server-side):**

- Active session expire হলে পরের `GET` attempt **auto-forfeit** করে (0 score, `passed: false`)
- `POST .../submit` after expiry → `400 quiz time limit exceeded`
- Forfeit-ও একটা attempt হিসেবে গণনা হয় (`attempts_used` বাড়ে)

**Storefront recommendation:** `seconds_remaining` দিয়ে countdown দেখাও; `0` হলে submit disable করো বা auto-submit করো।

---

## Part 3 — Student: Submit ও Result

### 3.1 Submit

```http
POST /v1/course/{course-slug}/quizzes/{quizId}/submit
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
Content-Type: application/json
```

```json
{
  "answers": [
    { "question_id": 41, "value": "a" },
    { "question_id": 42, "value": true },
    { "question_id": 43, "value": ["a", "c"] }
  ]
}
```

| Question type | `value` |
|---------------|---------|
| `single_choice` | `"a"` (option `id`) |
| `true_false` | `true` or `false` |
| `multiple_choice` | `["a", "c"]` |

**Submit rules:**

- Answers শুধু **current attempt session-এর question order/subset** থেকে হতে হবে
- `answer_required: true` question খালি থাকলে `400`
- Time limit expired হলে `400 quiz time limit exceeded`
- Successful submit session `submitted_at` mark করে

### 3.2 Submit response — student কী দেখবে

**Success `201` — storefront এ result screen এ এগুলো render করো:**

```json
{
  "message": "Quiz submitted successfully",
  "data": {
    "id": 101,
    "quiz_title": "HTML Basics Quiz",
    "attempt_number": 1,
    "score": 4,
    "max_score": 5,
    "percentage": 80,
    "passed": true,
    "status": "graded",
    "submitted_at": "2026-07-07T01:00:00Z",
    "reveal_answers": true,
    "total_questions": 5,
    "correct_count": 4,
    "incorrect_count": 1,
    "unanswered_count": 0,
    "pass_marks": 3,
    "minimum_pass_percentage": 60,
    "quiz_time_seconds": 1800,
    "attempt_time_seconds": 253,
    "instructor_feedback": "<p>Well done! Review question 3.</p>",
    "answers": [
      {
        "question_id": 41,
        "question_title": "What is HTML?",
        "question_type": "single_choice",
        "submitted_answer": "a",
        "is_correct": true,
        "marks_awarded": 1,
        "correct_answer": { "value": "a" },
        "answer_explanation": "HTML is a markup language."
      }
    ]
  }
}
```

| Field | UI তে দেখাও |
|-------|-------------|
| `score` / `max_score` | **“You scored 4/5”** |
| `percentage` | **“80%”** |
| `passed` | Pass / Fail badge |
| `status` | `graded` = final; `pending_review` = “Under review” |
| `answers[].marks_awarded` | Per-question marks |
| `answers[].is_correct` | ✓ / ✗ (null = pending) |
| `correct_answer` | শুধু যখন quiz `reveal_answers: true` |
| `answer_explanation` |同上 |
| `instructor_feedback` | Admin feedback (HTML); null = এখনো দেওয়া হয়নি — §3.7 |

> **Note:** এই repo তে student storefront UI নেই — তোমার learner app এ submit response bind করতে হবে। `instructor_feedback` থাকলে result screen-এ **Instructor Feedback** block দেখাও (`dangerouslySetInnerHTML` বা sanitizer দিয়ে)।

### 3.3 Grading rules

| Condition | `status` | Score |
|-----------|----------|-------|
| সব question এ `correct_answer` আছে | `graded` | Auto-calculated |
| কোনো question এ `correct_answer` নেই | `pending_review` | Partial/0 until review |
| `passed` | — | `percentage >= minimum_pass_percentage` |
| Time limit expired (auto-forfeit) | `graded` | `0` / `passed: false` |

### 3.4 Retry

| Quiz setting | Behaviour |
|--------------|-----------|
| `enable_retry: false` | ১ বার attempt (submit বা forfeit) |
| `enable_retry: true`, `retry_attempts: N` (N > 0) | সর্বোচ্চ N attempt |
| `enable_retry: true`, `retry_attempts: 0` | **Unlimited** attempts |
| Retry শেষ | `GET` quiz → `400 maximum quiz attempts reached` |

### 3.5 Skip / forfeit (manual)

Intro থেকে quiz না দিয়ে এগিয়ে যেতে:

```http
POST /v1/course/{course-slug}/quizzes/{quizId}/skip
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

**Behaviour (server-side):**

- Active unsubmitted session থাকলে → **forfeit** (0 score, `passed: false`, `status: graded`)
- Session না থাকলে → forfeit submission create হয় (intro থেকে skip; timer শুরু হয় না)
- Auto-forfeit (timer expire) এর মতোই **course progress**-এ quiz complete হিসেবে count হয়
- Certificate threshold cross করলে `TryIssueCertificate` চলে (submit + skip + auto-forfeit সবখানে)

**Success `201`:** Submit response-এর মতো same shape (`data` = `QuizSubmissionDetail`).

**Errors:** same as `GET .../quizzes/{id}` (`403` enrollment, `404` not found, `400` retry disabled / max attempts).

**Storefront UI:** Intro-তে “Skip quiz” বাটন এই endpoint call করবে — শুধু redirect নয়।

### 3.6 Submission history

```http
GET /v1/student/quiz-submissions?course_id=12
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

`course_id` optional। List এ summary: `score`, `max_score`, `percentage`, `passed`, `status`, `submitted_at`.

### 3.7 Submission detail (past attempts)

```http
GET /v1/student/quiz-submissions/{submissionId}
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

Returns same shape as submit `201` `data` (per-question `answers[]` + summary stats). শুধু owning student access করতে পারবে।

| Field | List | Detail `GET` | Submit `201` |
|-------|------|--------------|--------------|
| `score`, `percentage`, `passed` | ✅ | ✅ | ✅ |
| `total_questions`, `correct_count`, `incorrect_count` | ❌ | ✅ | ✅ |
| `answers[]` | ❌ | ✅ | ✅ |
| `correct_answer` | ❌ | ✅ if `reveal_answers` | ✅ if `reveal_answers` |
| `instructor_feedback` | ❌ | ✅ | ✅ |

**`instructor_feedback` (storefront):**

- Admin `POST .../feedback` দিয়ে save করলে student detail এ দেখা যাবে
- Submit-এর ঠিক পরে সাধারণত `null` (admin পরে লিখবে) — student **পুনরায়** `GET /student/quiz-submissions/{id}` call করলে updated feedback পাবে
- HTML string; empty হলে section hide করো

**Storefront UI — কোথায় দেখাবে:**

1. **Result screen** (submit/skip `201` `data`) — score + answers-এর নিচে “Instructor Feedback” (যদি `instructor_feedback` non-null)
2. **Past attempt detail** (`GET /student/quiz-submissions/{id}`) — same block
3. Admin পরে feedback দিলে student list refresh বা detail re-fetch করলে দেখা যাবে

**Example (student detail, feedback সহ):**

```json
{
  "data": {
    "id": 101,
    "quiz_title": "HTML Basics Quiz",
    "score": 4,
    "max_score": 5,
    "percentage": 80,
    "passed": true,
    "instructor_feedback": "<p>Well done! Focus on semantic HTML next time.</p>",
    "answers": [ "..."]
  }
}
```

**Storefront:** Result table-এর **Details** বাটন `GET /student/quiz-submissions/{id}` call করবে — client cache optional fallback হিসেবে রাখা যায়।

### 3.8 Course progress (quiz completion)

```http
GET /v1/course/{course-slug}/progress
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

**Success `200` (relevant fields):**

```json
{
  "data": {
    "lessons_completed": 4,
    "lessons_total": 5,
    "quizzes_completed": 2,
    "quizzes_total": 2,
    "assignments_completed": 1,
    "assignments_total": 1,
    "progress_percent": 87.5,
    "count_lessons": true,
    "count_quizzes": true,
    "count_assignments": true,
    "completed_lesson_ids": [12, 13, 14, 15],
    "completed_quiz_ids": [9, 11]
  }
}
```

| Field | আছে? | Notes |
|-------|-------|-------|
| `quizzes_completed` / `quizzes_total` | ✅ | Count |
| `progress_percent` | ✅ | Certificate threshold ([CERTIFICATE_STOREFRONT_API.md](./CERTIFICATE_STOREFRONT_API.md)) |
| `count_quizzes` | ✅ | Certificate settings — quiz progress count হয় কিনা |
| `completed_lesson_ids` | ✅ | Lesson sidebar tick |
| `completed_quiz_ids` | ✅ | Quiz sidebar tick — যেকোনো submission (submit/skip/forfeit) |

**Quiz → progress update:**

- Progress on-the-fly calculate হয়
- যেকোনো submission (`submit`, `skip`, auto-forfeit) `DISTINCT quiz_id` হিসেবে count হয় — pass/fail matter করে না
- Submit / skip / auto-forfeit-এর পর `GET /course/{slug}/progress` refresh করো
- Certificate auto-issue: submit, skip, এবং auto-forfeit সবখানে trigger হয়

---

## Part 4 — Auth & enrollment

### Student login

```http
POST /v1/student/login
app-key: <tenant_app_key>
Content-Type: application/json

{
  "email": "student@example.com",
  "password": "secret"
}
```

Response: `{ "token": "...", "user": { ... } }`

### Enrollment

Quiz attempt/submit এর আগে student কে course এ **enroll** করতে হবে (`enrollments` table / admin enrollment API)।  
না থাকলে: `403 enrollment required`.

---

## Part 5 — Backend ↔ Storefront status

### Implemented (latest)

| # | Feature | Endpoint / field | Status |
|---|---------|------------------|--------|
| 1 | Skip / forfeit quiz | `POST /course/{slug}/quizzes/{id}/skip` | ✅ |
| 2 | Past attempt Details | `GET /student/quiz-submissions/{submissionId}` | ✅ |
| 3 | Quiz sidebar IDs | `completed_quiz_ids` in `GET /course/{slug}/progress` | ✅ |
| 4 | Quiz → progress % | On-the-fly from `quiz_submissions` | ✅ |
| 5 | Certificate after forfeit/skip | `TryIssueCertificate` on submit, skip, auto-forfeit | ✅ |
| 6 | Admin submission detail page (summary + quiz overview) | Dashboard `QuizSubmissionDetailView` | ✅ |
| 7 | Admin instructor feedback | `POST .../quiz-submissions/{id}/feedback` | ✅ |
| 8 | Student reads `instructor_feedback` | `GET /student/quiz-submissions/{id}`, submit `201` | ✅ API |

### Remaining gaps / storefront notes

| # | Area | Status | Storefront action |
|---|------|--------|-------------------|
| 1 | **Intro attempt state** | ⚠️ No preview API | `GET .../quizzes/{id}` = session + timer start। Intro-তে `GET /student/quiz-submissions` দিয়ে derive করো; Start Quiz-এ server truth |
| 2 | **Text question types** | ❌ Not supported | শুধু `single_choice`, `multiple_choice`, `true_false` — admin/UI-তেও নেই |
| 3 | **Manual grading** | ❌ Admin grade API নেই | `correct_answer` ছাড়া question → `pending_review`; admin view only |
| 4 | **Student feedback UI** | ⚠️ Learner app | API ready — bind `instructor_feedback` on result + detail; re-fetch after admin updates |

### Skip vs progress (reference)

| Action | Backend record | Progress | Certificate |
|--------|----------------|----------|-------------|
| `POST .../submit` | Graded submission | ✅ Quiz counted | ✅ May issue |
| `POST .../skip` | Forfeit (0 score) | ✅ Quiz counted | ✅ May issue |
| Timer expire → `GET` | Auto-forfeit | ✅ Quiz counted | ✅ May issue |

### Assignment parity (submission APIs)

| Endpoint | Assignment | Quiz |
|----------|------------|------|
| `GET /student/...-submissions` (list) | ✅ summary + files | ✅ summary |
| `GET /student/...-submissions/{id}` | ✅ full detail | ✅ full detail (`answers[]`, `instructor_feedback`) |

---

## Part 6 — Deploy

### Migrations (required)

```bash
cd api
goose -dir migrations mysql "<GOOSE_DBSTRING>" up
```

| Migration | Table/column |
|-----------|----------------|
| `00040` | `quiz_questions.options` |
| `00041` | `quiz_questions.correct_answer` |
| `00042` | `quiz_submissions` |
| `00043` | `quiz_submission_answers` |
| `00055` | `quiz_attempt_sessions` (timer + fixed question order per attempt) |
| `00058` | `quiz_submissions.instructor_feedback` |

তারপর **API + web** redeploy।

### Env (storefront)

| Var | Example |
|-----|---------|
| `NEXT_PUBLIC_API_URL` | `https://api.example.com/v1` |
| Storefront `app-key` | Tenant এর `app_key` |

---

## Part 7 — cURL examples

```bash
# 1) Public course + quizzes
curl -s -H "app-key: TENANT_KEY" \
  "https://api.example.com/v1/course/react-masterclass"

# 2) Student login
TOKEN=$(curl -s -X POST \
  -H "app-key: TENANT_KEY" \
  -H "Content-Type: application/json" \
  -d '{"email":"s@example.com","password":"pass"}' \
  https://api.example.com/v1/student/login | jq -r .token)

# 3) Start quiz attempt (creates/reuses session + timer)
curl -s \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  "https://api.example.com/v1/course/react-masterclass/quizzes/9"

# 4) Single quiz view — load question 2 (0-based index)
curl -s \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  "https://api.example.com/v1/course/react-masterclass/quizzes/9/questions/1"

# 5) Submit
curl -s -X POST \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"answers":[{"question_id":41,"value":"a"}]}' \
  "https://api.example.com/v1/course/react-masterclass/quizzes/9/submit"

# 6) Skip quiz (intro — no timer if no active session)
curl -s -X POST \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  "https://api.example.com/v1/course/react-masterclass/quizzes/9/skip"

# 7) Admin — list submissions
curl -s \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  "https://api.example.com/v1/private/course/12/quiz-submissions"

# 8) Admin — submission detail
curl -s \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  "https://api.example.com/v1/private/course/12/quiz-submissions/101"

# 8b) Admin — save instructor feedback
curl -s -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"feedback":"<p>Well done!</p>"}' \
  "https://api.example.com/v1/private/course/12/quiz-submissions/101/feedback"

# 9) Course progress (includes completed_quiz_ids)
curl -s \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  "https://api.example.com/v1/course/react-masterclass/progress"

# 10) Student submission history (summary)
curl -s \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  "https://api.example.com/v1/student/quiz-submissions?course_id=12"

# 11) Student submission detail (past attempt answers)
curl -s \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  "https://api.example.com/v1/student/quiz-submissions/101"
```

---

## Part 8 — Feature matrix (latest)

| Feature | Status | Where |
|---------|--------|-------|
| Admin quiz create/edit (Curriculum) | ✅ | Dashboard |
| Optional quiz instructions | ✅ | Quiz edit form |
| Quiz settings (publish, randomize, timer, retry, pass %) | ✅ | Quiz edit + API |
| Question options + correct answer UI | ✅ | Add/Edit Question modal |
| Course save with quizzes | ✅ | Create/Update course API |
| Storefront quiz list | ✅ API | `GET /course/{slug}` |
| Student quiz attempt + session | ✅ API | `GET .../quizzes/{id}` |
| Single-question pagination | ✅ API | `GET .../questions/{index}` |
| Time limit (server enforced) | ✅ API | Session + submit validation |
| Randomize + visible question limit | ✅ API | Per-attempt session order |
| Student submit + instant result | ✅ API | `POST .../submit` |
| Skip / manual forfeit | ✅ API | `POST .../skip` |
| Student submission history (summary) | ✅ API | `GET /student/quiz-submissions` |
| Student submission detail (past attempts) | ✅ API | `GET /student/quiz-submissions/{id}` |
| Course progress (`quizzes_completed`) | ✅ API | `GET /course/{slug}/progress` |
| `completed_quiz_ids` in progress | ✅ API | `GET /course/{slug}/progress` |
| Auto-grading (MCQ + T/F) | ✅ | `correct_answer` on questions |
| Text / open-ended question types | ❌ | Not in schema |
| Admin submission list | ✅ | Course Edit → **Quizzes** tab |
| Admin submission detail (summary + overview + feedback editor) | ✅ | `QuizSubmissionDetailView` |
| Admin instructor feedback API | ✅ | `POST .../quiz-submissions/{id}/feedback` |
| Student `instructor_feedback` on detail/submit | ✅ API | `GET /student/quiz-submissions/{id}` |
| Manual grading / mark override | ❌ | Not built |
| Certificate trigger (submit/skip/forfeit) | ✅ | `TryIssueCertificate` |
| Student storefront result UI | — | Build in learner app |

---

## Part 9 — Storefront implementation checklist

- [ ] `GET /course/{slug}` → render `course_chapters[].quizzes[]`
- [ ] Student login; store JWT (guard empty Bearer)
- [ ] Verify enrollment (or handle `403`)
- [ ] **Intro screen:** `GET /student/quiz-submissions?course_id=` → derive attempts; **do not** call `GET .../quizzes/{id}` until Start
- [ ] **Start attempt:** `GET /course/{slug}/quizzes/{id}` (session + timer begins)
- [ ] If `display_mode === "single"`: paginate with `GET .../questions/{index}`
- [ ] Else: build form from `questions[]` + `options` (MCQ/T-F only)
- [ ] Show countdown from `seconds_remaining` (if present)
- [ ] `POST .../submit` before timer hits 0
- [ ] Result screen: `score`, `max_score`, `percentage`, `passed`, `answers[]`
- [ ] If `instructor_feedback` present: render HTML block below answers (result + detail)
- [ ] Details (history): `GET /student/quiz-submissions/{submissionId}` — include feedback re-fetch
- [ ] If `reveal_answers`: show `correct_answer`, `answer_explanation` per question
- [ ] After submit/skip: `GET /course/{slug}/progress` refresh
- [ ] Quiz sidebar ✔️: `completed_quiz_ids` from progress response
- [ ] Intro Skip: `POST /course/{slug}/quizzes/{id}/skip` then progress refresh
- [ ] Retry button only if `can_retry: true` on next `GET` quiz
- [ ] Handle `400 quiz time limit exceeded` gracefully

---

## Summary

| Area | Status |
|------|--------|
| Quiz attempt + timer + submit + instant result | ✅ |
| Skip quiz (`POST .../skip`) | ✅ |
| Past attempt Details (`GET /student/quiz-submissions/{id}`) | ✅ |
| Instructor feedback (admin write → student read) | ✅ API |
| `completed_quiz_ids` in progress | ✅ |
| Course progress + certificate (submit/skip/forfeit) | ✅ |
| Single-view pagination | ✅ |
| Auto-grade (MCQ / T-F) | ✅ |
| Intro preview without starting timer | ⚠️ Use submission list on intro; no preview endpoint |
| Text question types | ❌ Not supported |
| Manual admin grading | ❌ View only |

**Deploy:** run migrations `00040`–`00043`, `00055`, `00058` before release.

**Source of truth (Go):** `api/internal/modules/quiz/` — `service.go`, `router.go`, `response.go`; progress: `api/internal/progress/course.go`, `api/internal/modules/courseprogress/`.

---

## Related files (codebase)

| Area | Path |
|------|------|
| Quiz API module | `api/internal/modules/quiz/` |
| Attempt session model | `api/internal/models/quiz_attempt_session.go` |
| Course progress API | `api/internal/modules/courseprogress/` |
| Progress calculation | `api/internal/progress/course.go` |
| Course quiz CRUD | `api/internal/modules/course/service.go` |
| Admin quiz form | `frontend/.../curriculum/QuizEdit.tsx` |
| Admin Quizzes tab | `frontend/.../CoursesTabs.tsx` |
| Submission table | `frontend/.../quiz-evaluation/QuizTable.tsx` |
| Submission detail view | `frontend/.../quiz-evaluation/QuizSubmissionDetailView.tsx` |
| Question form | `frontend/.../curriculum/QuizQuestionAnswerFields.tsx` |
| Migrations | `api/migrations/00040` – `00043`, `00055`, `00058` |
