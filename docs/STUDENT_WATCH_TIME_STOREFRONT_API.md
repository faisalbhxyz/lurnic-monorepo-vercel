# Student Watch Time & Learning Report — Storefront API

**API base:** `https://<api-host>/v1`  
**Related:** [LESSON_VIDEO_PROGRESS_STOREFRONT_API.md](./LESSON_VIDEO_PROGRESS_STOREFRONT_API.md) · [FREE_LESSONS_STOREFRONT_API.md](./FREE_LESSONS_STOREFRONT_API.md) · [MINAR_ACADEMY_STOREFRONT_API.md](./MINAR_ACADEMY_STOREFRONT_API.md)

**Actual video play seconds** (report-card learning time) সার্ভারে থাকে — mobile, website, admin student details একই নম্বর দেখে।  
Resume position (`max_position_seconds`) আলাদা API — সেটাকে “minutes watched” হিসেবে ব্যবহার করা যাবে না।

সব request এ **`app-key`** লাগবে। Student route এ **`Authorization: Bearer <student_jwt>`**। Admin route এ admin JWT।

---

## Status

| Layer | Status |
|-------|--------|
| Migration `student_daily_watch` + `student_watch_events` | ✅ Ready (`00068`) |
| `POST /student/watch-time` | ✅ Ready |
| `POST /student/watch-time/batch` | ✅ Ready |
| `GET /student/learning-report` (authoritative daily seconds) | ✅ Ready |
| `GET /admin/students/{id}/learning-report` | ✅ Ready |
| Admin student details nested `learning_time` | ✅ Ready |

---

## Quick reference

| Action | Method | Path | Auth |
|--------|--------|------|------|
| Ingest play delta | `POST` | `/student/watch-time` | `app-key` + Bearer (student) |
| Offline batch flush | `POST` | `/student/watch-time/batch` | `app-key` + Bearer (student) |
| Learning report chart | `GET` | `/student/learning-report?period=7d\|30d\|90d` | `app-key` + Bearer (student) |
| Admin full report | `GET` | `/admin/students/{studentId}/learning-report?period=7d\|30d\|90d\|all` | admin JWT |
| Admin alias | `GET` | `/private/student/{studentId}/learning-report` | admin JWT |
| Admin details embed | `GET` | `/private/student/details/{id}` → `learning_time` | admin JWT |

---

## Two metrics (do not mix)

| Metric | API | Meaning |
|--------|-----|---------|
| Resume position | `PATCH /course/{slug}/lessons/{id}/progress` | Playhead farthest point |
| Actual watch time | `POST /student/watch-time` | Real play seconds (rewatch counts again) |
| Lesson complete | `POST /course/{slug}/lessons/{id}/complete` | ≥ ~80% position |

---

## 1. Ingest — `POST /v1/student/watch-time`

Clients flush ~every 15s while playing, and on pause / background / unmount.

```http
POST /v1/student/watch-time
Authorization: Bearer <token>
app-key: <key>
Content-Type: application/json

{
  "client_event_id": "550e8400-e29b-41d4-a716-446655440000",
  "watched_seconds": 14,
  "watch_date": "2026-09-04",
  "timezone": "Asia/Dhaka",
  "watched_at": "2026-09-04T08:15:22Z",
  "course_id": 12,
  "lesson_id": 88,
  "source": "enrolled",
  "device_platform": "android"
}
```

### Required fields

| Field | Notes |
|-------|--------|
| `client_event_id` | UUID; unique per student — idempotent |
| `watched_seconds` | `> 0`; values `> 300` clamped to `300` |
| `watch_date` | Student local calendar day `YYYY-MM-DD` |
| `timezone` | IANA tz, e.g. `Asia/Dhaka` |

### Optional

| Field | Notes |
|-------|--------|
| `watched_at` | RFC3339 UTC; default = server now |
| `course_id` / `lesson_id` | If `lesson_id` set → must be enrolled **or** public free lesson |
| `source` | `enrolled` (default) · `free_lesson` · `offline` |
| `device_platform` | `ios` · `android` · `web` |

### Success `200`

```json
{
  "data": {
    "accepted": true,
    "watch_date": "2026-09-04",
    "day_video_seconds": 1842,
    "duplicate": false,
    "client_event_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

Same `client_event_id` again → `duplicate: true`, daily total unchanged.

### Errors

| Status | When |
|--------|------|
| `401` | Missing / invalid student JWT |
| `403` | Lesson not enrolled and not public free |
| `422` | `watched_seconds ≤ 0`, bad date/tz, future `watch_date` > today+1 |

---

## 2. Batch — `POST /v1/student/watch-time/batch`

Offline catch-up. Max **50** events per request.

```json
{
  "events": [
    {
      "client_event_id": "…",
      "watched_seconds": 14,
      "watch_date": "2026-09-04",
      "timezone": "Asia/Dhaka"
    }
  ]
}
```

**Success `200`**

```json
{
  "data": {
    "accepted_count": 1,
    "duplicate_count": 0,
    "results": [ /* per-event WatchTimeAcceptResponse */ ],
    "daily_totals": [
      { "date": "2026-09-04", "seconds": 1842 }
    ]
  }
}
```

---

## 3. Student report — `GET /v1/student/learning-report`

```http
GET /v1/student/learning-report?period=7d
```

`period`: `7d` · `30d` · `90d` (default `7d`).

```json
{
  "data": {
    "period": "7d",
    "daily_watch_seconds": [
      { "date": "2026-08-29", "seconds": 0 },
      { "date": "2026-09-04", "seconds": 3600 }
    ],
    "streak_days": 3,
    "quiz_accuracy_percent": 85.5,
    "courses_in_progress": 2,
    "courses_completed": 1
  }
}
```

- `daily_watch_seconds` = `student_daily_watch.video_seconds` (one row per day in period, zeros included).
- `streak_days` = consecutive days ending **today** (`Asia/Dhaka`) with `seconds > 0`.
- Quiz / course counts: existing quiz submissions + enrollment progress logic.

---

## 4. Admin report — `GET /v1/admin/students/{studentId}/learning-report`

```http
GET /v1/admin/students/42/learning-report?period=30d
```

`period`: `7d` · `30d` · `90d` · `all`.

```json
{
  "data": {
    "period": "30d",
    "daily_watch_seconds": [ /* … */ ],
    "streak_days": 5,
    "quiz_accuracy_percent": 72,
    "courses_in_progress": 1,
    "courses_completed": 2,
    "totals": {
      "video_seconds_period": 54000,
      "video_seconds_all_time": 180000,
      "last_watched_at": "2026-09-04T08:15:22Z"
    },
    "by_course": [
      {
        "course_id": 12,
        "course_title": "HSC Physics",
        "video_seconds": 12000
      }
    ]
  }
}
```

### Nested on admin student details

`GET /private/student/details/{id}`:

```json
{
  "learning_time": {
    "video_seconds_7d": 7200,
    "video_seconds_30d": 54000,
    "streak_days": 5,
    "last_watched_at": "2026-09-04T08:15:22Z"
  }
}
```

---

## Sync behaviour

```text
Playing (~15s) / pause / leave  →  POST /student/watch-time
Offline queue                   →  POST /student/watch-time/batch
Open Learning Report            →  GET /student/learning-report  (server wins)
Admin student details           →  GET admin learning-report / nested learning_time
```

Local AsyncStorage cache OK for offline UX only. After first successful GET, prefer server days.

---

## Validation & anti-abuse (v1)

| Rule | Behaviour |
|------|-----------|
| `watched_seconds ≤ 0` | `422` |
| `watched_seconds > 300` | clamp to 300 |
| Same `client_event_id` | `200` no-op (`duplicate: true`) |
| Future `watch_date` > today+1 (tz) | `422` |
| Lesson not playable | `403` |

Do **not** derive daily totals from `max_position_seconds`.

---

## Migration

```bash
cd api
goose -dir migrations mysql "$GOOSE_DBSTRING" up
```

| Migration | Tables |
|-----------|--------|
| `00068_create_student_watch_tables.sql` | `student_daily_watch`, `student_watch_events` |

**Backend:** `modules/student/watch_time.go`, `storefront_*.go`, `details.go`

---

## Client wiring (after deploy)

### Mobile
- `useLessonWatch` / `addDailyWatchSeconds` → also `POST /student/watch-time` with UUID `client_event_id`
- Learning report: prefer API after ingest reliable

### Website
- `LessonVideoPlayer`: track play ticks → `POST /student/watch-time` (keep existing progress PATCH)

### Admin
- Student details: show 7d/30d minutes + streak from `learning_time` or full admin GET chart

---

## Acceptance

- [x] `POST /student/watch-time` upserts daily total; duplicate `client_event_id` is idempotent
- [x] `GET /student/learning-report` reads `student_daily_watch` (not position)
- [x] Admin GET + nested `learning_time` on student details
- [x] Existing progress / complete APIs unchanged
- [ ] Client: playing 2 min → ~120s on report (mobile + web wiring)
- [ ] Cross-device: phone post → web GET shows same seconds

---

## Out of scope (v1)

- Live-class / quiz time categories (columns reserved; stay `0`)
- Realtime WebSocket report push
- Migrating historical AsyncStorage minutes
- Per-instructor analytics dashboards
