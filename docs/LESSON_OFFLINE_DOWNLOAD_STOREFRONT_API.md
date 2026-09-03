# Lesson Offline Download — Storefront Guide

**API base:** `https://<api-host>/v1`  
**Auth:** `app-key` + student `Authorization: Bearer <jwt>`  
**Last updated:** September 2026

লেসন **দেখা** আর **অফলাইন সেভ** দুটো আলাদা জিনিস।

| কাজ | কোথায় আসে | Storefront কী করবে |
|-----|-----------|-------------------|
| অ্যাপে ভিডিও প্লে | `source_type` + `source.data.data` (YouTube embed / Vimeo / upload) | আগের মতোই প্লেয়ার |
| অফলাইন ডাউনলোড | `offline_downloadable` + download API | শুধু তখনই Download বাটন |

Admin বাঁদিকে YouTube রাখে, ডান পাশে আলাদা **Google Drive** লিংক দেয়। Drive লিংক প্লেব্যাক সোর্স নয়।

---

## 1. কখন Download বাটন দেখাবেন

```ts
const canDownloadOffline =
  isEnrolled &&
  lesson.lesson_type === "video" &&
  lesson.offline_downloadable === true;
```

| সিগন্যাল | মানে |
|----------|------|
| `offline_downloadable: false` | বাটন লুকাবেন। YouTube/Vimeo প্লেব্যাক থাকলেও অফলাইন নয়। |
| `offline_downloadable: true` | বাটন দেখাবেন। |
| `download_url` | Drive **share** লিংক (preview)। ফাইল সেভের জন্য এটা ব্যবহার করবেন না। |
| `source.data.drive_url` | Admin যে Drive লিংক সেভ করেছে। ডিবাগের জন্য; ডাউনলোডে proxy API ব্যবহার করুন। |
| `source.data.data` | প্লেব্যাক (YouTube iframe ইত্যাদি)। এখান থেকে Drive খুঁজবেন না। |
| `resources[]` | PDF/Word/ZIP — লেসন ম্যাটেরিয়াল, অফলাইন ভিডিও নয়। |

লগইন নেই / এনরোল নেই → Download দেখাবেন না।

---

## 2. Lesson object (enrolled course)

`GET /v1/enrolled/courses` এবং `GET /v1/course/{slug}` — `course_chapters[].course_lessons[]`

### YouTube প্লে + Drive অফলাইন (সাধারণ কেস)

```json
{
  "id": 42,
  "title": "Chapter 1 — Introduction",
  "lesson_type": "video",
  "source_type": "youtube",
  "source": {
    "data": {
      "data": "<iframe src=\"https://www.youtube.com/embed/vHXkebRshWk\"></iframe>",
      "is_file": false,
      "playback_times": "00:15:30",
      "drive_url": "https://drive.google.com/file/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view?usp=sharing",
      "drive_file_id": "1AbCdEfGhIjKlMnOpQrStUvWxYz"
    }
  },
  "resources": [],
  "offline_downloadable": true,
  "download_url": "https://drive.google.com/file/d/1AbCdEfGhIjKlMnOpQrStUvWxYz/view?usp=sharing"
}
```

প্লেয়ার: `source_type` + `source.data.data`  
ডাউনলোড: নিচের endpoint

### অফলাইন নেই

```json
{
  "source_type": "youtube",
  "source": { "data": { "data": "<iframe ...></iframe>", "is_file": false } },
  "offline_downloadable": false
}
```

`download_url` নাও থাকতে পারে।

### লিগ্যাসি: পুরো সোর্সই Drive

পুরনো ডাটায় `source_type: "google_drive"` এবং `source.data.data`-তেই Drive URL থাকতে পারে। নতুন অ্যাডমিন ফ্লো এভাবে সেভ করে না। Storefront শুধু `offline_downloadable` দেখলেই যথেষ্ট।

---

## 3. ফাইল কীভাবে নামবেন (এটাই ব্যবহার করুন)

সরাসরি Google Drive share লিংক (`/view`) প্রায়ই HTML পেজ দেয়, `.mp4` না। তাই **Lurnic download API** কল করুন।

```
GET /v1/course/{slug}/lessons/{lessonId}/download?format=json
```

**Headers**

```
app-key: <tenant_app_key>
Authorization: Bearer <student_jwt>
```

**Path**

| Param | উদাহরণ |
|-------|--------|
| `slug` | কোর্স স্লাগ, যেমন `hsc-physics` |
| `lessonId` | লেসন numeric id, যেমন `42` |

**Query**

| Param | Default | কী হয় |
|-------|---------|--------|
| `format=json` | না থাকলে 302 | JSON-এ সরাসরি ফাইল URL |
| (কোনো `format` না) | 302 | `Location` এ Drive download URL |

মোবাইল/SPA-তে **`?format=json` ব্যবহার করুন** — 302 ফলো করা কষ্টকর, আর Bearer হেডার রেডাইরেক্টে যায় না।

### Success `200`

```json
{
  "data": {
    "download_url": "https://drive.usercontent.google.com/download?id=1AbCdEfGhIjKlMnOpQrStUvWxYz&export=download&confirm=t",
    "file_name": "chapter-1-introduction.mp4",
    "content_type": "video/mp4"
  }
}
```

এই `data.download_url` থেকে ফাইল ফেচ করে ডিভাইসে সেভ করুন। `file_name` লোকাল ফাইল নাম হিসেবে ব্যবহার করুন।

### Errors

| HTTP | `error` / code | UI |
|------|----------------|-----|
| `400` | App key missing | Retry / config |
| `401` | `UNAUTHORIZED` | Login |
| `403` | `NOT_ENROLLED` | Enroll CTA |
| `404` | `LESSON_NOT_FOUND` | Lesson missing / unpublished |
| `422` | `NOT_DOWNLOADABLE` | বাটন হাইড; “Offline not available” |

---

## 4. Storefront implementation steps

1. Curriculum লোড: `GET /enrolled/courses` বা `GET /course/{slug}` (`app-key` + Bearer)।
2. প্লেয়ার আগের মতো `source_type` / `source.data.data` দিয়ে চালাবেন।
3. `offline_downloadable === true` হলে Download আইকন।
4. ট্যাপে:

```ts
const res = await fetch(
  `${API_BASE}/v1/course/${courseSlug}/lessons/${lesson.id}/download?format=json`,
  {
    headers: {
      "app-key": APP_KEY,
      Authorization: `Bearer ${studentJwt}`,
    },
  }
);

if (!res.ok) {
  // 401 → login, 403 → enroll, 422 → hide button
  throw new Error("download_failed");
}

const body = await res.json();
const fileUrl: string = body.data.download_url;
const fileName: string = body.data.file_name; // e.g. chapter-1-introduction.mp4
const contentType: string = body.data.content_type; // video/mp4
```

5. `fileUrl` থেকে বাইনারি ডাউনলোড (React Native: `FileSystem.downloadAsync` / `expo-file-system`; Web: blob — web অফলাইন স্কোপের বাইরে থাকতে পারে)।
6. লোকাল DB-তে রাখুন: `lessonId`, `courseSlug`, `fileName`, local path, `contentType`।
7. অফলাইন প্লে: লোকাল ফাইল। অনলাইন প্লে: এখনও YouTube/Vimeo সোর্স।

### React Native স্কেচ

```ts
async function downloadLessonOffline(opts: {
  apiBase: string;
  appKey: string;
  token: string;
  slug: string;
  lessonId: number;
  destPath: string;
}) {
  const metaRes = await fetch(
    `${opts.apiBase}/v1/course/${opts.slug}/lessons/${opts.lessonId}/download?format=json`,
    {
      headers: {
        "app-key": opts.appKey,
        Authorization: `Bearer ${opts.token}`,
      },
    }
  );
  if (!metaRes.ok) {
    const err = await metaRes.json().catch(() => ({}));
    throw Object.assign(new Error("not_downloadable"), {
      status: metaRes.status,
      err,
    });
  }
  const { data } = await metaRes.json();
  // FileSystem.downloadAsync(data.download_url, `${opts.destPath}/${data.file_name}`)
  return data;
}
```

---

## 5. যা করবেন না

- YouTube/Vimeo URL অফলাইনে ডাউনলোড করার চেষ্টা।
- শুধু `source.data.data` স্ক্যান করে Drive খোঁজা — এখন সেখানে embed কোড থাকে।
- লেসন `download_url` (share `/view`) সরাসরি ফাইল হিসেবে সেভ।
- `resources[]` কে অফলাইন ভিডিও ভাবা।
- ডাউনলোড API তে Bearer ছাড়া কল — 401।
- ওয়েবে `window.open(download_url)` — হেডার যায় না; মোবাইলে authenticated `fetch` + file save।

---

## 6. Google Drive শর্ত (অ্যাডমিন)

স্টোরফ্রন্টে ফিক্স করার কিছু নেই; ফাইল না নামলে সাধারণত শেয়ার সেটিং।

1. ফাইল ভিডিও (`.mp4` / `.m4v` / `.webm` / `.mov`)
2. Sharing: **Anyone with the link → Viewer**
3. লিংক: `https://drive.google.com/file/d/{FILE_ID}/view`

---

## 7. QA

```bash
API_URL="https://<api-host>/v1"
TOKEN="<student_jwt>"
APP_KEY="<tenant_app_key>"
SLUG="<course-slug>"
LESSON_ID="<lesson-id>"

# কোন লেসনে অফলাইন আছে
curl -s "$API_URL/course/$SLUG" \
  -H "Authorization: Bearer $TOKEN" \
  -H "app-key: $APP_KEY" \
  | jq '.data.course_chapters[].course_lessons[]
      | {id, title, source_type, offline_downloadable, download_url, drive_url: .source.data.drive_url}'

# আসল ফাইল URL
curl -s "$API_URL/course/$SLUG/lessons/$LESSON_ID/download?format=json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "app-key: $APP_KEY"
```

চেকলিস্ট

- [ ] YouTube লেসনে প্লেয়ার কাজ করে
- [ ] Drive ফিল্ড খালি → `offline_downloadable: false`, বাটন নেই
- [ ] Drive ফিল্ড আছে → বাটন আছে
- [ ] `format=json` → `drive.usercontent.google.com/download?...`
- [ ] সেভ করা ফাইল অফলাইনে প্লে হয়
- [ ] আনএনরোল্ড → `403 NOT_ENROLLED`
- [ ] YouTube-only লেসনে download API → `422 NOT_DOWNLOADABLE`

---

## Related

- Master API: [`MINAR_ACADEMY_STOREFRONT_API.md`](./MINAR_ACADEMY_STOREFRONT_API.md) — §6.3 `/enrolled/courses`
- Video progress: [`LESSON_VIDEO_PROGRESS_STOREFRONT_API.md`](./LESSON_VIDEO_PROGRESS_STOREFRONT_API.md)
