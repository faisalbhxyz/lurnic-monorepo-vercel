# Student Certificate — Storefront Implementation Guide

এই doc অনুসরণ করলে storefront-এর **student dashboard** থেকে course complete হওয়ার পর student তার certificate **দেখতে** এবং **PDF হিসেবে download** করতে পারবে।

**API base:** `https://<api-host>/v1`  
**Last updated:** July 2026  
**API reference:** [CERTIFICATE_STOREFRONT_API.md](./CERTIFICATE_STOREFRONT_API.md)  
**Related:** [LESSON_VIDEO_PROGRESS_STOREFRONT_API.md](./LESSON_VIDEO_PROGRESS_STOREFRONT_API.md) · [QUIZ_STOREFRONT_API.md](./QUIZ_STOREFRONT_API.md) · [ASSIGNMENT_STOREFRONT_API.md](./ASSIGNMENT_STOREFRONT_API.md)

---

## Status

| Layer | Status |
|-------|--------|
| API auto-issue (lesson / quiz / assignment) | ✅ Ready |
| API `GET /student/certificates` (list) | ✅ Ready |
| API `GET /student/certificates/{id}` (detail JSON) | ✅ Ready |
| API `GET /student/certificates/{id}/html` (view + print) | ✅ Ready |
| API `GET /course/{slug}/certificate` (course page JSON) | ✅ Ready |
| Storefront student dashboard UI | আপনার storefront এ implement করতে হবে |

---

## 1. যে endpoint ব্যবহার করবেন (এবং যেটা করবেন না)

| ✅ Storefront (student) | ❌ Admin — storefront এ ব্যবহার করবেন না |
|------------------------|------------------------------------------|
| `GET /v1/student/certificates` | `PUT /v1/private/course/{id}` (certificate settings) |
| `GET /v1/student/certificates/{id}` | Admin `Bearer` token |
| `GET /v1/student/certificates/{id}/html` | |
| `GET /v1/course/{slug}/certificate` | |
| `GET /v1/course/{slug}/progress` | |

**404 এড়াতে:**

1. URL শেষে **`/v1`** — `NEXT_PUBLIC_API_URL=https://api.example.com/v1`
2. Path **`/student/certificates`** — `private` বা admin prefix যোগ করবেন না
3. Header **`app-key`** + **`Authorization: Bearer <student_jwt>`** — admin token দিলে 401/403
4. `GET /student/details` দিয়ে certificate পাবেন না — আলাদা `GET /student/certificates` লাগবে

---

## 2. Environment variables

```env
NEXT_PUBLIC_API_URL=https://api.yourdomain.com/v1
NEXT_PUBLIC_APP_KEY=your-tenant-app-key
```

Local dev:

```env
NEXT_PUBLIC_API_URL=http://localhost:5000/v1
NEXT_PUBLIC_APP_KEY=your-local-tenant-app-key
```

---

## 3. Certificate কখন issue হয় (auto)

Admin **Courses → Edit → Settings → Certificates**:

| Setting | Field |
|---------|-------|
| Enable | `is_enabled` |
| Minimum % | `completion_percent` (default 100) |
| Count items | `count_lessons`, `count_quizzes`, `count_assignments` |

**Triggers** (threshold cross হলে একবার issue):

| Activity | Endpoint |
|----------|----------|
| Lesson complete | `POST /course/{slug}/lessons/{lessonId}/complete` |
| Quiz submit | `POST /course/{slug}/quizzes/{quizId}/submit` |
| Quiz skip / forfeit | `POST /course/{slug}/quizzes/{quizId}/skip` |
| Assignment submit | `POST /course/{slug}/assignments/{assignmentId}/submit` |

Certificate number: **14-char hex** (e.g. `a1b2c3d4e5f607`) — `CERT-` prefix নেই।

---

## 4. Suggested storefront routes

| Route | কাজ | API |
|-------|-----|-----|
| `/dashboard/certificates` | সব certificate list | `GET /student/certificates` |
| `/dashboard/certificates/{id}` | View + download | `GET /student/certificates/{id}/html` |
| `/courses/{slug}` | Progress + certificate CTA | progress + certificate endpoints |

**Dashboard sidebar:** “My Certificates” / “সার্টিফিকেট” → list page।

---

## 5. TypeScript types

```ts
export type StudentCertificate = {
  id: number;
  course_id: number;
  course_title: string;
  certificate_number: string;
  student_name: string;
  progress_percent: number;
  template_path: string;
  title?: string | null;
  subtitle_one?: string | null;
  subtitle_two?: string | null;
  brand_logo?: string | null;
  watermark_image?: string | null;
  watermark_opacity?: number;
  organization_name?: string | null;
  signer_name?: string | null;
  signer_role?: string | null;
  signer_org?: string | null;
  dual_signers_enabled?: boolean;
  signer2_name?: string | null;
  signer2_role?: string | null;
  signer2_org?: string | null;
  pricing_model?: "free" | "paid";
  owner_signature?: string | null;
  instructor_signature?: string | null;
  issued_at: string;
  download_url?: string;
};
```

---

## 6. Copy-paste: API client (`lib/studentCertificateApi.ts`)

```ts
const API_URL = process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, "") ?? "";
const APP_KEY = process.env.NEXT_PUBLIC_APP_KEY ?? "";

export type StudentCertificate = {
  id: number;
  course_id: number;
  course_title: string;
  certificate_number: string;
  student_name: string;
  progress_percent: number;
  template_path: string;
  title?: string | null;
  subtitle_one?: string | null;
  subtitle_two?: string | null;
  brand_logo?: string | null;
  watermark_image?: string | null;
  watermark_opacity?: number;
  organization_name?: string | null;
  signer_name?: string | null;
  signer_role?: string | null;
  signer_org?: string | null;
  dual_signers_enabled?: boolean;
  signer2_name?: string | null;
  signer2_role?: string | null;
  signer2_org?: string | null;
  pricing_model?: "free" | "paid";
  owner_signature?: string | null;
  instructor_signature?: string | null;
  issued_at: string;
  download_url?: string;
};

export class CertificateApiError extends Error {
  status: number;
  constructor(message: string, status: number) {
    super(message);
    this.status = status;
  }
}

function requireConfig() {
  if (!API_URL) throw new Error("NEXT_PUBLIC_API_URL is not set");
  if (!APP_KEY) throw new Error("NEXT_PUBLIC_APP_KEY is not set");
}

function studentHeaders(token: string): HeadersInit {
  if (!token) throw new CertificateApiError("Not logged in", 401);
  return {
    "app-key": APP_KEY,
    Authorization: `Bearer ${token}`,
  };
}

async function parseError(res: Response): Promise<CertificateApiError> {
  let body: { error?: string; message?: string } = {};
  try {
    body = await res.json();
  } catch {
    // HTML error page from /html route
  }
  return new CertificateApiError(
    body.message ?? body.error ?? `Request failed (${res.status})`,
    res.status
  );
}

export async function fetchStudentCertificates(
  token: string
): Promise<StudentCertificate[]> {
  requireConfig();
  const res = await fetch(`${API_URL}/student/certificates`, {
    headers: studentHeaders(token),
    cache: "no-store",
  });
  if (!res.ok) throw await parseError(res);
  const json = (await res.json()) as { data: StudentCertificate[] };
  return json.data ?? [];
}

export async function fetchStudentCertificate(
  token: string,
  certificateId: number
): Promise<StudentCertificate> {
  requireConfig();
  const res = await fetch(`${API_URL}/student/certificates/${certificateId}`, {
    headers: studentHeaders(token),
    cache: "no-store",
  });
  if (!res.ok) throw await parseError(res);
  const json = (await res.json()) as { data: StudentCertificate };
  return json.data;
}

export async function fetchCourseCertificate(
  token: string,
  courseSlug: string
): Promise<StudentCertificate> {
  requireConfig();
  const res = await fetch(`${API_URL}/course/${courseSlug}/certificate`, {
    headers: studentHeaders(token),
    cache: "no-store",
  });
  if (!res.ok) throw await parseError(res);
  const json = (await res.json()) as { data: StudentCertificate };
  return json.data;
}

export async function fetchCertificateHTML(
  token: string,
  certificateId: number
): Promise<string> {
  requireConfig();
  const res = await fetch(`${API_URL}/student/certificates/${certificateId}/html`, {
    headers: studentHeaders(token),
    cache: "no-store",
  });
  if (!res.ok) throw await parseError(res);
  return res.text();
}

/** Recommended: auth fetch → blob → new tab (built-in Download PDF button). */
export async function openCertificateInNewTab(
  token: string,
  certificateId: number
): Promise<void> {
  const html = await fetchCertificateHTML(token, certificateId);
  const blob = new Blob([html], { type: "text/html;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const win = window.open(url, "_blank", "noopener,noreferrer");
  if (!win) {
    URL.revokeObjectURL(url);
    throw new Error("Pop-up blocked. Allow pop-ups for this site.");
  }
  setTimeout(() => URL.revokeObjectURL(url), 60_000);
}
```

### Axios variant

```ts
import axios from "axios";

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: { "app-key": process.env.NEXT_PUBLIC_APP_KEY },
});

export async function fetchStudentCertificatesAxios(token: string) {
  if (!token) throw new Error("Not logged in");
  const res = await api.get("/student/certificates", {
    headers: { Authorization: `Bearer ${token}` },
  });
  return res.data.data as StudentCertificate[];
}

export async function fetchCertificateHTMLAxios(
  token: string,
  certificateId: number
) {
  if (!token) throw new Error("Not logged in");
  const res = await api.get(`/student/certificates/${certificateId}/html`, {
    headers: { Authorization: `Bearer ${token}` },
    responseType: "text",
  });
  return res.data as string;
}
```

---

## 7. View + Download — তিনটা approach

### Approach A (Recommended): `/html` endpoint

- Admin preview-এর মতো **Minar Academy** design
- API-built **Download PDF** button (`window.print()`)
- Storefront-এ শুধু `openCertificateInNewTab(token, cert.id)` call

```tsx
<button type="button" onClick={() => openCertificateInNewTab(token, cert.id)}>
  View / Download
</button>
```

Student নতুন tab-এ certificate দেখে → **Download PDF** → print dialog → “Save as PDF”।

### Approach B: In-page iframe (`srcdoc`)

Dashboard modal/page-এ embed:

```tsx
"use client";

import { useEffect, useState } from "react";
import { fetchCertificateHTML } from "@/lib/studentCertificateApi";

export function CertificateEmbed({
  token,
  certificateId,
}: {
  token: string;
  certificateId: number;
}) {
  const [html, setHtml] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    fetchCertificateHTML(token, certificateId)
      .then((h) => { if (!cancelled) setHtml(h); })
      .catch((e) => { if (!cancelled) setError(e.message); });
    return () => { cancelled = true; };
  }, [token, certificateId]);

  if (error) return <p className="text-red-600">{error}</p>;
  if (!html) return <p>Loading certificate…</p>;

  return (
    <iframe
      title="Certificate"
      srcDoc={html}
      className="w-full border-0"
      style={{ minHeight: "70vh", aspectRatio: "297 / 210" }}
      sandbox="allow-scripts allow-same-origin"
    />
  );
}
```

### Approach C: Client-side React (separate repo)

Admin dashboard files copy করুন:

| File | Purpose |
|------|---------|
| `frontend/components/shared/certificates/StudentCertificateView.tsx` | Template router |
| `frontend/components/shared/home/products/courses/create/settings/MinarAcademyCertificate.tsx` | Minar layout |
| `frontend/styles/minar-certificate.css` | Styles |
| `frontend/lib/printCertificate.ts` | Print without pop-up blocker |
| `frontend/lib/certificate-format.ts` | Date formatting |
| `frontend/components/shared/certificates/certificate-templates.ts` | `isMinarCertificateTemplate()` |

`printCertificatePreview(ref)` দিয়ে in-page download — Approach A-র চেয়ে বেশি setup।

---

## 8. Copy-paste: Dashboard list page

```tsx
"use client";

import { useEffect, useState } from "react";
import {
  fetchStudentCertificates,
  openCertificateInNewTab,
  CertificateApiError,
  type StudentCertificate,
} from "@/lib/studentCertificateApi";

const TOKEN_KEY = "student_token";

function formatDate(iso: string) {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}

export default function MyCertificatesPage() {
  const [certs, setCerts] = useState<StudentCertificate[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = localStorage.getItem(TOKEN_KEY);
    if (!token) {
      window.location.href = "/login";
      return;
    }

    fetchStudentCertificates(token)
      .then(setCerts)
      .catch((e: unknown) => {
        if (e instanceof CertificateApiError && e.status === 401) {
          localStorage.removeItem(TOKEN_KEY);
          window.location.href = "/login";
          return;
        }
        setError(e instanceof Error ? e.message : "Failed to load certificates");
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p>Loading certificates…</p>;
  if (error) return <p className="text-red-600">{error}</p>;

  if (certs.length === 0) {
    return (
      <div>
        <h1>My Certificates</h1>
        <p>Complete a course to earn your first certificate.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">My Certificates</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {certs.map((cert) => (
          <article key={cert.id} className="rounded-lg border p-4 shadow-sm">
            <h2 className="font-medium">{cert.course_title}</h2>
            <p className="text-sm text-gray-600 mt-1">
              Issued {formatDate(cert.issued_at)}
            </p>
            <p className="text-xs text-gray-500 mt-1">
              ID: {cert.certificate_number}
            </p>
            <button
              type="button"
              className="mt-4 w-full rounded-md bg-teal-800 px-4 py-2 text-white text-sm font-medium"
              onClick={() => {
                const token = localStorage.getItem(TOKEN_KEY);
                if (!token) {
                  window.location.href = "/login";
                  return;
                }
                openCertificateInNewTab(token, cert.id).catch((e) =>
                  alert(e instanceof Error ? e.message : "Could not open certificate")
                );
              }}
            >
              View / Download
            </button>
          </article>
        ))}
      </div>
    </div>
  );
}
```

---

## 9. Copy-paste: Certificate detail page (iframe embed)

```tsx
"use client";

import { useParams, useRouter } from "next/navigation";
import { CertificateEmbed } from "@/components/CertificateEmbed";
// CertificateEmbed = section 7 Approach B

const TOKEN_KEY = "student_token";

export default function CertificateDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = Number(params.id);
  const token = typeof window !== "undefined" ? localStorage.getItem(TOKEN_KEY) : null;

  if (!token) {
    router.replace("/login");
    return null;
  }
  if (!Number.isFinite(id)) {
    return <p>Invalid certificate.</p>;
  }

  return (
    <div className="mx-auto max-w-5xl p-4">
      <h1 className="text-xl font-semibold mb-4">Your Certificate</h1>
      <CertificateEmbed token={token} certificateId={id} />
    </div>
  );
}
```

---

## 10. Course page: progress + certificate button

```ts
import {
  fetchCourseCertificate,
  openCertificateInNewTab,
  CertificateApiError,
  type StudentCertificate,
} from "@/lib/studentCertificateApi";

export async function getCourseCertificateOrLocked(
  token: string,
  courseSlug: string
): Promise<"locked" | StudentCertificate> {
  try {
    return await fetchCourseCertificate(token, courseSlug);
  } catch (e) {
    if (e instanceof CertificateApiError && e.status === 404) return "locked";
    throw e;
  }
}
```

```tsx
// Course page snippet
const cert = await getCourseCertificateOrLocked(token, slug);

{cert === "locked" ? (
  <p>Complete the course to unlock your certificate.</p>
) : (
  <button onClick={() => openCertificateInNewTab(token, cert.id)}>
    View Certificate
  </button>
)}
```

Progress bar: `GET /course/{slug}/progress` → `data.progress_percent`

---

## 11. Errors & auth guard

| HTTP | Action |
|------|--------|
| `401` / `403` | Clear token → `/login` |
| `404` on certificate | “Not available yet” |
| `400` enrollment | Enroll CTA |
| Pop-up blocked | Show in-page `CertificateEmbed` fallback |

```ts
function getStudentToken(): string | null {
  const token = localStorage.getItem("student_token");
  if (!token || token === "undefined") return null;
  return token;
}
```

---

## 12. Admin checklist

1. Certificate **enabled** + **completion %** set
2. Published lessons / quizzes / assignments আছে
3. Student **enrolled**
4. Correct **`app-key`**
5. Migrations `00059`–`00062` applied (brand, watermark, dual signers)

---

## 13. Quick test plan

1. Admin: enable certificate, 100% threshold
2. Student login → enroll → complete all items
3. `GET /student/certificates` → new entry
4. Dashboard **View / Download** → design OK → **Download PDF** works
5. Course page button only after issue
6. 401 → login redirect

---

## Related docs

- [CERTIFICATE_STOREFRONT_API.md](./CERTIFICATE_STOREFRONT_API.md)
- [LESSON_VIDEO_PROGRESS_STOREFRONT_API.md](./LESSON_VIDEO_PROGRESS_STOREFRONT_API.md)
- [STUDENT_DEVICE_LOGIN_STOREFRONT_API.md](./STUDENT_DEVICE_LOGIN_STOREFRONT_API.md)
