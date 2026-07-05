# Class-wise Notes Storefront API

Lurnic API base: `https://<api-host>/v1`

সব public request এ **`app-key`** header লাগবে (tenant key)।  
Authentication লাগে না — এগুলো public academic resource।

---

## Data hierarchy

Storefront এ notes ৩-লেভেল hierarchy তে দেখানো হয়:

```
Class (HSC, ৮ম শ্রেণি)
  └── Subject (Bangla, English)
        └── Paper (Bangla 1st Paper, Bangla 2nd Paper)
              └── Notes (individual PDF lecture sheets)
```

Dashboard থেকে manage করার path: **Website → Class-wise Notes**

---

## 1. Homepage — class cards (একাডেমিক পড়াশোনার সবকিছু)

Published classes list করো, প্রতিটিতে total note count সহ।

```http
GET /v1/academic-notes
app-key: <tenant_app_key>
```

**Success `200`:**

```json
{
  "data": [
    {
      "id": 1,
      "title": "HSC",
      "slug": "hsc",
      "icon_label": "H",
      "icon_color": "#E91E63",
      "position": 0,
      "note_count": 88
    },
    {
      "id": 2,
      "title": "৯ম-১০ম শ্রেণি",
      "slug": "class-9-10",
      "icon_label": "S",
      "icon_color": "#FFC107",
      "position": 1,
      "note_count": 44
    }
  ]
}
```

**UI flow:**
1. Section heading: `একাডেমিক পড়াশোনার সবকিছু`
2. Grid card render করো — icon (color + label), title, `note_count`
3. Subtitle: `• {note_count} টি লেকচার শীট [PDF Download]`
4. Card click → `/resource/academic/{slug}` (storefront route তোমার app অনুযায়ী)

**Note:** শুধু published class/subject/paper/note count হয়; empty branch storefront এ আসে না।

---

## 2. Class page — subjects + papers (HSC page)

```http
GET /v1/academic-notes/{classSlug}
app-key: <tenant_app_key>
```

**Example:** `GET /v1/academic-notes/hsc`

**Success `200`:**

```json
{
  "data": {
    "id": 1,
    "title": "HSC",
    "slug": "hsc",
    "icon_label": "H",
    "icon_color": "#E91E63",
    "position": 0,
    "subjects": [
      {
        "id": 1,
        "class_id": 1,
        "title": "Bangla",
        "slug": "bangla",
        "position": 0,
        "note_count": 24,
        "papers": [
          {
            "id": 1,
            "subject_id": 1,
            "title": "Bangla 1st Paper",
            "slug": "bangla-1st-paper",
            "icon_label": "১ম",
            "icon_color": "#42A5F5",
            "position": 0,
            "note_count": 12
          },
          {
            "id": 2,
            "subject_id": 1,
            "title": "Bangla 2nd Paper",
            "slug": "bangla-2nd-paper",
            "icon_label": "২য়",
            "icon_color": "#FF7043",
            "position": 1,
            "note_count": 12
          }
        ]
      }
    ]
  }
}
```

**UI flow:**
1. Breadcrumb: `Resource > Academic Resource > {class.title}`
2. Subject heading (bold green): `{subject.title} →`
3. Paper cards side-by-side:
   - Circular icon: `icon_color` background + `icon_label` text
   - Title: `{paper.title}`
   - Subtitle: `• {paper.note_count} টি লেকচার শীট [PDF Download]`
4. Paper card click → notes list page

---

## 3. Notes list page — individual PDF cards

```http
GET /v1/academic-notes/{classSlug}/{subjectSlug}/{paperSlug}
app-key: <tenant_app_key>
```

**Example:** `GET /v1/academic-notes/hsc/bangla/bangla-1st-paper`

**Success `200`:**

```json
{
  "data": {
    "class": {
      "id": 1,
      "title": "HSC",
      "slug": "hsc",
      "icon_label": "H",
      "icon_color": "#E91E63"
    },
    "subject": {
      "id": 1,
      "class_id": 1,
      "title": "Bangla",
      "slug": "bangla",
      "note_count": 24
    },
    "paper": {
      "id": 1,
      "subject_id": 1,
      "title": "Bangla 1st Paper",
      "slug": "bangla-1st-paper",
      "icon_label": "১ম",
      "icon_color": "#42A5F5",
      "note_count": 12
    },
    "notes": [
      {
        "id": 1,
        "paper_id": 1,
        "title": "অপরিচিতা",
        "subtitle": "অপরিচিতা",
        "thumbnail": "https://cdn.example.com/thumb.jpg",
        "pdf_url": "https://cdn.example.com/oporichita.pdf",
        "pdf_file_name": "oporichita.pdf",
        "position": 0
      },
      {
        "id": 2,
        "paper_id": 1,
        "title": "মানব কল্যাণ",
        "subtitle": "মানব কল্যাণ",
        "thumbnail": null,
        "pdf_url": "https://cdn.example.com/manob-kollyan.pdf",
        "pdf_file_name": "manob-kollyan.pdf",
        "position": 1
      }
    ]
  }
}
```

**UI flow:**
1. Breadcrumb: `Resource > Academic Resource > HSC > Bangla > Bangla 1st Paper`
2. Grid of note cards:
   - Top: light blue area with `thumbnail` (fallback: generic PDF icon)
   - Bottom: `title` (bold), `subtitle` (smaller), green `PDF` label
3. Card click → open `pdf_url` in new tab or trigger download

---

## 4. Storefront routing suggestion

| Page | Suggested route | API |
|------|-----------------|-----|
| Class grid (home section) | `/` or `/resources` | `GET /academic-notes` |
| Class detail | `/resources/academic/{classSlug}` | `GET /academic-notes/{classSlug}` |
| Notes list | `/resources/academic/{classSlug}/{subjectSlug}/{paperSlug}` | `GET /academic-notes/{classSlug}/{subjectSlug}/{paperSlug}` |

---

## 5. Example fetch (Next.js / React)

```typescript
const API = process.env.NEXT_PUBLIC_API_URL; // must end with /v1
const APP_KEY = process.env.NEXT_PUBLIC_APP_KEY;

async function fetchAcademicNotes() {
  const res = await fetch(`${API}/academic-notes`, {
    headers: { "app-key": APP_KEY },
    next: { revalidate: 300 },
  });
  if (!res.ok) throw new Error("Failed to load classes");
  const json = await res.json();
  return json.data;
}

async function fetchClassDetail(classSlug: string) {
  const res = await fetch(`${API}/academic-notes/${classSlug}`, {
    headers: { "app-key": APP_KEY },
    next: { revalidate: 300 },
  });
  if (!res.ok) throw new Error("Class not found");
  const json = await res.json();
  return json.data;
}

async function fetchNotes(classSlug: string, subjectSlug: string, paperSlug: string) {
  const res = await fetch(
    `${API}/academic-notes/${classSlug}/${subjectSlug}/${paperSlug}`,
    { headers: { "app-key": APP_KEY }, next: { revalidate: 300 } }
  );
  if (!res.ok) throw new Error("Notes not found");
  const json = await res.json();
  return json.data;
}
```

---

## 6. Bengali number formatting (optional)

Storefront এ count বাংলায় দেখাতে:

```typescript
function toBnNumber(n: number): string {
  return n.toLocaleString("bn-BD");
}
// 88 → "৮৮"
```

Subtitle: `` `• ${toBnNumber(paper.note_count)} টি লেকচার শীট [PDF Download]` ``

---

## 7. Admin API (dashboard only)

Dashboard authenticated routes (`Authorization: Bearer <admin_jwt>`):

| Action | Method | Path |
|--------|--------|------|
| List classes | GET | `/private/academic-notes/classes` |
| Get class + tree | GET | `/private/academic-notes/classes/{id}` |
| Create class | POST | `/private/academic-notes/classes/create` |
| Update class | PUT | `/private/academic-notes/classes/update/{id}` |
| Delete class | DELETE | `/private/academic-notes/classes/delete/{id}` |
| Create subject | POST | `/private/academic-notes/subjects/create` |
| Update subject | PUT | `/private/academic-notes/subjects/update/{id}` |
| Delete subject | DELETE | `/private/academic-notes/subjects/delete/{id}` |
| Create paper | POST | `/private/academic-notes/papers/create` |
| Update paper | PUT | `/private/academic-notes/papers/update/{id}` |
| Delete paper | DELETE | `/private/academic-notes/papers/delete/{id}` |
| Create note | POST (multipart) | `/private/academic-notes/notes/create` |
| Update note | PUT (multipart) | `/private/academic-notes/notes/update/{id}` |
| Delete note | DELETE | `/private/academic-notes/notes/delete/{id}` |

**Create note** multipart fields:
- `paper_id` (required)
- `title` (required)
- `subtitle` (optional)
- `thumbnail` (optional image, max 2MB)
- `pdf` (required PDF, max 20MB)
- `position` (optional, default 0)
- `is_published` (optional, default true)

---

## 8. Migration

Production deploy এর আগে migration run করো:

```bash
goose -dir api/migrations up
```

Migration file: `00049_create_academic_notes_tables.sql`

Tables: `academic_note_classes`, `academic_note_subjects`, `academic_note_papers`, `academic_notes`

---

## 9. Error responses

| Status | Meaning |
|--------|---------|
| `404` | Class/subject/paper not found or unpublished |
| `400` | Invalid input (missing PDF, bad ID) |
| `401` | Missing/invalid `app-key` (public) or Bearer token (admin) |
| `500` | Server error |

Public endpoints empty list return `200` with `"data": []` — error নয়।
