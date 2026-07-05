# Quiz — Storefront API, Admin & Submission Guide

**API base:** `https://<api-host>/v1`  
**Admin dashboard:** `https://<dashboard-host>/courses/{id}/update?tab=Quizzes`

এই doc এ quiz এর পুরো flow আছে — admin এ quiz তৈরি, storefront এ student দেখা/submit, result, আর admin এ submission review।

---

## Quick reference — সব endpoint

| Who | Method | Path | Headers |
|-----|--------|------|---------|
| Public | `GET` | `/course/{slug}` | `app-key` |
| Student | `GET` | `/course/{slug}/quizzes/{quizId}` | `app-key` + `Bearer` |
| Student | `POST` | `/course/{slug}/quizzes/{quizId}/submit` | `app-key` + `Bearer` |
| Student | `GET` | `/student/quiz-submissions?course_id=` | `app-key` + `Bearer` |
| Student | `POST` | `/student/login` | `app-key` |
| Admin | `POST` | `/private/course/create` | `Bearer` (multipart) |
| Admin | `PUT` | `/private/course/update/{id}` | `Bearer` (multipart) |
| Admin | `GET` | `/private/course/{courseId}/quiz-submissions` | `Bearer` |
| Admin | `GET` | `/private/course/{courseId}/quiz-submissions/{submissionId}` | `Bearer` |

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
  Student->>API: POST .../submit
  API->>Student: score, percentage, passed, answers
  Admin->>API: GET quiz-submissions (Quizzes tab)
  Admin->>API: GET quiz-submissions/{id} (detail modal)
```

---

## Part 1 — Admin: Quiz তৈরি ও manage

### 1.1 Curriculum তে quiz add

1. **Courses → Create/Edit → Curriculum**
2. Chapter এ **+ → Quiz**
3. Title, instructions, settings, **Add Questions**
4. প্রতিটি question এ:
   - **Single choice / Multiple choice:** options + correct answer (radio/checkbox)
   - **True/False:** True বা False select
5. Quiz **Save** → পুরো course **Save**

Quiz **published** (`is_published: true`) হলে storefront এ দেখা যাবে।  
প্রতিটি question এ কমপক্ষে ১টা question + **correct answer** থাকলে auto-grade কাজ করবে।

### 1.2 Question fields (API + UI)

| Field | Type | UI | Notes |
|-------|------|-----|-------|
| `title` | string | ✅ | Required |
| `type` | `single_choice` \| `multiple_choice` \| `true_false` | ✅ | |
| `marks` | number | ✅ | |
| `options` | `[{ id, text }]` | ✅ | MCQ তে required (min 2) |
| `correct_answer` | JSON | ✅ | Auto-grade এর জন্য |
| `answer_explanation` | HTML string | ✅ | Optional; reveal হলে student দেখবে |
| `answer_required` | boolean | ✅ | |

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

### 1.3 Admin: Submission review (Dashboard)

**Course Edit → Quizzes tab** (`?tab=Quizzes`)

| Feature | আছে |
|---------|-----|
| সব submission list | ✅ |
| Quiz / chapter / student name / email | ✅ |
| Marks `score/max_score (%)` | ✅ |
| Status filter (All / Evaluate / Pending) | ✅ |
| Row click → **detail modal** | ✅ |
| Per-question answers, correct/incorrect, explanation | ✅ |
| Manual mark edit | ❌ (দরকার নেই যদি সব question auto-grade হয়) |

**Evaluate** filter = `status: pending_review` (যেখানে `correct_answer` ছিল না বা grade হয়নি)।

Detail modal API call করে:

```http
GET /v1/private/course/{courseId}/quiz-submissions/{submissionId}
Authorization: Bearer <admin_token>
```

Admin detail এ সবসময় **correct answer + explanation** দেখায় (review এর জন্য)।

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

**Storefront UI:**
1. Course load করো
2. Chapter অনুযায়ী quiz cards দেখাও
3. “Start quiz” → enrolled student হলে attempt screen এ যাও

### 2.2 Quiz attempt screen (enrolled student only)

```http
GET /v1/course/{course-slug}/quizzes/{quizId}
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

**Extra fields:**

| Field | Meaning |
|-------|---------|
| `attempts_used` | আগে কতবার submit করেছে |
| `can_retry` | আবার attempt করা যাবে কিনা |

**Errors:**

| HTTP | Reason |
|------|--------|
| `403` | Enrolled নয় |
| `404` | Course/quiz নেই বা unpublished |
| `400` | Retry disabled বা max attempts reached |

**UI hints from quiz settings:**

| Setting | Storefront behaviour |
|---------|---------------------|
| `single_quiz_view` | এক সময়ে একটা question |
| `time_limit` + `time_limit_option` | Timer UI |
| `randomize_questions` | API question order shuffle করে |
| `total_visible_questions` | কতটা question দেখাবে (subset) |

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
    "submitted_at": "2026-06-25T10:30:00Z",
    "reveal_answers": true,
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

> **Note:** এই repo তে student storefront UI নেই — তোমার learner app এ submit response bind করতে হবে।

### 3.3 Grading rules

| Condition | `status` | Score |
|-----------|----------|-------|
| সব question এ `correct_answer` আছে | `graded` | Auto-calculated |
| কোনো question এ `correct_answer` নেই | `pending_review` | Partial/0 until review |
| `passed` | — | `percentage >= minimum_pass_percentage` |

### 3.4 Retry

| Quiz setting | Behaviour |
|--------------|-----------|
| `enable_retry: false` | ১ বার submit |
| `enable_retry: true`, `retry_attempts: N` | সর্বোচ্চ N attempt |
| Retry শেষ | `GET` quiz → `400 maximum quiz attempts reached` |

### 3.5 Submission history (পরে আবার দেখতে)

```http
GET /v1/student/quiz-submissions?course_id=12
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

`course_id` optional। List এ `score`, `max_score`, `percentage`, `passed`, `submitted_at` পাবে।

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

## Part 5 — Deploy

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

তারপর **API + web** redeploy।

### Env (storefront)

| Var | Example |
|-----|---------|
| `NEXT_PUBLIC_API_URL` | `https://api.example.com/v1` |
| Storefront `app-key` | Tenant এর `app_key` |

---

## Part 6 — cURL examples

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

# 3) Load quiz for attempt
curl -s \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  "https://api.example.com/v1/course/react-masterclass/quizzes/9"

# 4) Submit
curl -s -X POST \
  -H "app-key: TENANT_KEY" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"answers":[{"question_id":41,"value":"a"}]}' \
  "https://api.example.com/v1/course/react-masterclass/quizzes/9/submit"

# 5) Admin — list submissions
curl -s \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  "https://api.example.com/v1/private/course/12/quiz-submissions"

# 6) Admin — submission detail
curl -s \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  "https://api.example.com/v1/private/course/12/quiz-submissions/101"
```

---

## Feature matrix (final)

| Feature | Status | Where |
|---------|--------|-------|
| Admin quiz create/edit (Curriculum) | ✅ | Dashboard |
| Question options + correct answer UI | ✅ | Add/Edit Question modal |
| Course save with quizzes | ✅ | Create/Update course API |
| Storefront quiz list | ✅ API | `GET /course/{slug}` |
| Student quiz attempt | ✅ API | `GET .../quizzes/{id}` |
| Student submit + instant result | ✅ API | `POST .../submit` |
| Student submission history | ✅ API | `GET /student/quiz-submissions` |
| Auto-grading | ✅ | `correct_answer` on questions |
| Admin submission list | ✅ | Course Edit → **Quizzes** tab |
| Admin submission detail (modal) | ✅ | Click row in Quizzes tab |
| Manual grading / mark override | ❌ | Not built (not required for auto-grade MCQ) |
| Student storefront result UI | ❌ | Build in learner app using submit API |

---

## Storefront implementation checklist

- [ ] `GET /course/{slug}` → render `course_chapters[].quizzes[]`
- [ ] Student login; store JWT
- [ ] Verify enrollment (or handle `403`)
- [ ] `GET /course/{slug}/quizzes/{id}` before attempt
- [ ] Build answer form from `questions[]` + `options`
- [ ] `POST .../submit` on finish
- [ ] Result screen: `score`, `max_score`, `percentage`, `passed`
- [ ] If `reveal_answers`: show `correct_answer`, `answer_explanation` per question
- [ ] Retry button only if `can_retry` on next `GET` quiz
- [ ] Optional: history via `GET /student/quiz-submissions`

---

## Related files (codebase)

| Area | Path |
|------|------|
| Quiz API module | `api/internal/modules/quiz/` |
| Course quiz CRUD | `api/internal/modules/course/service.go` |
| Admin Quizzes tab | `frontend/.../CoursesTabs.tsx` |
| Submission table | `frontend/.../quiz-evaluation/QuizTable.tsx` |
| Submission detail modal | `frontend/.../quiz-evaluation/QuizSubmissionDetailModal.tsx` |
| Question form | `frontend/.../curriculum/AddNewQuestion.tsx`, `QuizQuestionAnswerFields.tsx` |
| Migrations | `api/migrations/00040` – `00043` |
