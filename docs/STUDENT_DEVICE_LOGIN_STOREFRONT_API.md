# Student Single-Device Login — Storefront API

**API base:** `https://<api-host>/v1`

প্রতিটি student **এক সময়ে শুধু একটা device** থেকে class/quiz/assignment access করতে পারবে।  
নতুন device এ login করলে **পুরনো device এর session invalidate** হয়ে যাবে — পুরনো device থেকে পরের API call এ `401` পাবে।

Admin dashboard → **Students → Details** এ active device info দেখা যায়।

---

## Quick reference

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| `POST` | `/student/login` | `app-key` | `device_id` **required** |
| `POST` | `/student/logout` | `app-key` + `Bearer` | Current device session end |
| Protected routes | e.g. `/course/...`, `/student/...` | `app-key` + `Bearer` | Session must match DB |

---

## Login (device required)

```http
POST /v1/student/login
app-key: <tenant_app_key>
Content-Type: application/json

{
  "email": "student@example.com",
  "password": "secret",
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "device_name": "Chrome on Windows"
}
```

| Field | Required | Notes |
|-------|----------|-------|
| `email` | ✅ | |
| `password` | ✅ | min 6 chars |
| `device_id` | ✅ | Stable per-browser ID (8–128 chars). Store in `localStorage` |
| `device_name` | ❌ | Display label for admin; API falls back to User-Agent |

### Success `200`

```json
{
  "token": "<jwt>",
  "user": {
    "user_id": "...",
    "name": "Student Name",
    "email": "student@example.com",
    "phone": null
  }
}
```

JWT contains `session_id` — server tracks one active session per student.

### New device login behaviour

1. Student already logged in on **Phone A**
2. Same account logs in on **Laptop B** with different `device_id`
3. **Phone A** — next API call → `401` with `code: "SESSION_REPLACED"`
4. **Laptop B** — works normally
5. Admin sees **Laptop B** in student details → Active Device

---

## Logout

```http
POST /v1/student/logout
app-key: <tenant_app_key>
Authorization: Bearer <student_token>
```

Success: `{ "message": "Logged out successfully" }`

---

## Session replaced / expired (`401`)

```json
{
  "error": "Session expired or logged in on another device",
  "code": "SESSION_REPLACED",
  "message": "Your account was logged in on another device. Please sign in again."
}
```

**Storefront action:** clear stored token → redirect to login → show message.

---

## Storefront implementation checklist

### 1. Device ID (once per browser)

```javascript
const STORAGE_KEY = "lurnic_device_id";

function getDeviceId() {
  let id = localStorage.getItem(STORAGE_KEY);
  if (!id) {
    id = crypto.randomUUID();
    localStorage.setItem(STORAGE_KEY, id);
  }
  return id;
}
```

### 2. Login request

```javascript
const res = await fetch(`${API}/student/login`, {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    "app-key": APP_KEY,
  },
  body: JSON.stringify({
    email,
    password,
    device_id: getDeviceId(),
    device_name: optionalLabel, // e.g. from navigator.userAgent parsing
  }),
});
```

### 3. Global 401 handler

On any student API `401` with `code === "SESSION_REPLACED"` (or missing session):

- Remove token from storage
- Redirect to `/login`
- Toast: logged in on another device

### 4. Password reset

`POST /student/reset-password` success এর পর সব session invalidate হয় — user কে আবার login করতে হবে (`device_id` সহ)।

---

## Admin: student details

`GET /v1/private/student/details/{id}` (admin Bearer) response এ:

```json
{
  "data": {
    "active_device": {
      "device_id": "550e8400-...",
      "device_name": "Chrome on Windows",
      "ip_address": "103.x.x.x",
      "user_agent": "Mozilla/5.0 ...",
      "logged_in_at": "2026-07-05T10:00:00Z",
      "last_seen_at": "2026-07-05T12:30:00Z"
    }
  }
}
```

`active_device` is `null` when student is not logged in.

---

## Migration

```bash
cd api
goose -dir migrations mysql "<GOOSE_DBSTRING>" up
```

Requires migration `00052_create_student_sessions_table.sql`.

---

## Related docs

- [STUDENT_PASSWORD_RESET_STOREFRONT_API.md](./STUDENT_PASSWORD_RESET_STOREFRONT_API.md)
- [QUIZ_STOREFRONT_API.md](./QUIZ_STOREFRONT_API.md)
- [CERTIFICATE_STOREFRONT_API.md](./CERTIFICATE_STOREFRONT_API.md)
