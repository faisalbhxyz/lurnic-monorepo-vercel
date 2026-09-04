# Student Class Profile — Storefront API

**API base:** `https://<api-host>/v1`  
**Related:** [FREE_LESSONS_STOREFRONT_API.md](./FREE_LESSONS_STOREFRONT_API.md) · [CLASS_NOTE_STOREFRONT_API.md](./CLASS_NOTE_STOREFRONT_API.md) · [STUDENT_PROFILE_STOREFRONT_API.md](./STUDENT_PROFILE_STOREFRONT_API.md)

প্রতিটি student-এর **class / grade preference** সার্ভারে থাকে — mobile + website একই profile দেখে।  
আগে শুধু device-local (`AsyncStorage`) ছিল; অন্য device / web এ “কোন ক্লাস” restore হত না।

সব request এ **`app-key`** লাগবে। Protected route এ **`Authorization: Bearer <student_jwt>`**।

---

## Status

| Layer | Status |
|-------|--------|
| Migration `student_class_profiles` | ✅ Ready |
| `GET /student/class-profile` | ✅ Ready |
| `PUT /student/class-profile` | ✅ Ready |
| Login + `GET /student/details` nested `class_profile` | ✅ Ready |

---

## Quick reference

| Action | Method | Path | Auth |
|--------|--------|------|------|
| Get class profile | `GET` | `/student/class-profile` | `app-key` + Bearer |
| Upsert / merge class profile | `PUT` | `/student/class-profile` | `app-key` + Bearer |
| Login (nested profile) | `POST` | `/student/login` | `app-key` |
| Details (nested profile) | `GET` | `/student/details` | `app-key` + Bearer |

---

## Canonical fields

| Field | Type | Notes |
|-------|------|-------|
| `class_level` | string | Onboarding enum (required on full save) |
| `hsc_batch` | string \| null | Required when `class_level = hsc_beyond` |
| `department` | string \| null | Required when HSC batch ≠ `after_hsc` |
| `preferred_class_slug` | string \| null | CMS `academic_note_classes.slug` (home / free lessons filter) |
| `onboarding_completed` | boolean | Client sets `true` after onboarding |
| `updated_at` | RFC3339 | Last write wins for sync |

### `class_level`

| Value | Meaning |
|-------|---------|
| `class_five` | Class 5 |
| `class_six` | Class 6 |
| `class_seven` | Class 7 |
| `class_eight` | Class 8 |
| `class_9_10_ssc` | Class 9–10 & SSC |
| `hsc_beyond` | HSC & beyond |

### `hsc_batch` (only `hsc_beyond`)

`hsc_2026` · `hsc_2027` · `hsc_2028` · `after_hsc`

### `department` (HSC batch ≠ `after_hsc`)

`science` · `business_studies` · `humanities`

### Default `preferred_class_slug` mapping

Backend auto-fills when client omits slug **and** the CMS class exists:

| `class_level` | Default slug |
|---------------|--------------|
| `class_five` | `class-5` |
| `class_six` | `class-6` |
| `class_seven` | `class-7` |
| `class_eight` | `class-8` |
| `class_9_10_ssc` | `ssc` |
| `hsc_beyond` | `hsc` |

If mapped slug missing in CMS → `class_level` still saved, `preferred_class_slug` stays `null` until home picker.

---

## 1. `GET /student/class-profile`

```http
GET /v1/student/class-profile
app-key: <tenant_app_key>
Authorization: Bearer <student_jwt>
```

### Never set → `200` with `data: null`

```json
{ "data": null }
```

### Saved → `200`

```json
{
  "data": {
    "class_level": "class_eight",
    "hsc_batch": null,
    "department": null,
    "preferred_class_slug": "class-8",
    "onboarding_completed": true,
    "updated_at": "2026-09-04T08:00:00Z"
  }
}
```

---

## 2. `PUT /student/class-profile`

Idempotent upsert. Omitted fields keep previous values (merge). Always bumps `updated_at`.

### Full onboarding (Class 8)

```http
PUT /v1/student/class-profile
Authorization: Bearer <token>
app-key: <key>
Content-Type: application/json

{
  "class_level": "class_eight",
  "onboarding_completed": true
}
```

### HSC

```json
{
  "class_level": "hsc_beyond",
  "hsc_batch": "hsc_2027",
  "department": "science",
  "preferred_class_slug": "hsc",
  "onboarding_completed": true
}
```

### Home picker only (partial)

```json
{ "preferred_class_slug": "class-9" }
```

- Profile exists → merge; enums unchanged.
- Profile missing + known default slug (`class-8`, `ssc`, …) → create with inferred `class_level`, `onboarding_completed: false`.
- Unknown slug + no profile → `422` (`class_level` required).

### Success `200`

Same shape as GET `data` (saved object).

### Errors

| Status | When |
|--------|------|
| `401` | Missing / invalid token |
| `422` | Invalid enum, bad HSC combo, unknown CMS slug, empty body |

---

## 3. Nested on auth responses

### `POST /student/login` → `user.class_profile`

```json
{
  "token": "...",
  "user": {
    "id": 1,
    "first_name": "...",
    "email": "...",
    "class_profile": {
      "class_level": "class_eight",
      "hsc_batch": null,
      "department": null,
      "preferred_class_slug": "class-8",
      "onboarding_completed": true,
      "updated_at": "2026-09-04T08:00:00Z"
    }
  }
}
```

Unset → `"class_profile": null`.

### `GET /student/details` → same nested `class_profile`

Clients should treat **details** as source of truth after login and on app foreground.

`PUT /student/update` response `data` also includes `class_profile`.

---

## Sync behaviour

Server is source of truth. No WebSocket in v1 — push on change, pull on login / resume.

| Event | Client action |
|-------|----------------|
| Finish onboarding | `PUT` full enums + `onboarding_completed: true` |
| Home header class change | `PUT` `{ preferred_class_slug }` only |
| Login | Read `user.class_profile`; if `null` → show onboarding |
| App resume / web load | `GET /student/details` (or dedicated GET); **server wins** over local cache |
| First launch after upgrade | If server `null` and local profile exists → one-time `PUT` migrate |

Filter free lessons with `GET /free-lessons?class_slug=<preferred_class_slug>`.

---

## Curl smoke test

```bash
API=https://<api-host>/v1
APP_KEY=<tenant_app_key>

TOKEN=$(curl -s -X POST "$API/student/login" \
  -H "app-key: $APP_KEY" -H "Content-Type: application/json" \
  -d '{"email":"student@example.com","password":"...","device_id":"device-12345678"}' \
  | jq -r .token)

curl -s -X PUT "$API/student/class-profile" \
  -H "app-key: $APP_KEY" -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"class_level":"class_eight","onboarding_completed":true}' | jq

curl -s "$API/student/class-profile" \
  -H "app-key: $APP_KEY" -H "Authorization: Bearer $TOKEN" | jq

curl -s "$API/student/details" \
  -H "app-key: $APP_KEY" -H "Authorization: Bearer $TOKEN" | jq .data.class_profile
```

---

## Acceptance checklist

- [ ] `PUT /student/class-profile` upserts and returns saved row
- [ ] Validation rejects invalid HSC combo
- [ ] `GET` returns `null` when never set
- [ ] Login + `GET /student/details` include `class_profile`
- [ ] Class set on phone → web login shows same filters
- [ ] Class changed on web → mobile resume shows new class after details refetch
- [ ] `preferred_class_slug` works with `GET /free-lessons?class_slug=…`
- [ ] CMS slugs exist for mapped defaults (`class-5` … `hsc`)

---

## Out of scope (v1)

- WebSocket push of class changes
- Admin CMS UI to edit student class
- Automatic AsyncStorage bulk migration (client one-time upload is enough)
