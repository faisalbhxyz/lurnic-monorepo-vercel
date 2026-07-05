"use client";

import React, { useEffect, useState } from "react";
import { Session } from "next-auth";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import Label from "@/components/ui/Label";
import { LuBookOpen, LuFileText, LuLayers, LuPlus, LuTrash2 } from "react-icons/lu";
import { FileText } from "lucide-react";

type Props = {
  session: Session;
  classData: IAcademicNoteClassDetail;
};

function authHeaders(session: Session) {
  return { Authorization: `Bearer ${session.accessToken}` };
}

function slugifyTitle(title: string) {
  return title
    .toLowerCase()
    .replace(/\s+/g, "-")
    .replace(/[^a-z0-9-]/g, "");
}

function StepBadge({ step, label }: { step: number; label: string }) {
  return (
    <div className="flex items-center gap-2 text-sm">
      <span className="w-6 h-6 rounded-full bg-primary text-white text-xs font-semibold flex items-center justify-center shrink-0">
        {step}
      </span>
      <span className="font-medium text-gray-700">{label}</span>
    </div>
  );
}

function NotePreview({ note }: { note: IAcademicNoteItem }) {
  if (note.thumbnail) {
    return (
      // eslint-disable-next-line @next/next/no-img-element
      <img
        src={note.thumbnail}
        alt={note.title}
        className="h-full w-full object-contain"
      />
    );
  }

  if (note.pdf_url) {
    return (
      <a
        href={note.pdf_url}
        target="_blank"
        rel="noopener noreferrer"
        className="group relative block h-full w-full"
        title="Open PDF"
      >
        <iframe
          src={`${note.pdf_url}#page=1&view=FitH&toolbar=0&navpanes=0`}
          title={`${note.title} PDF preview`}
          className="h-full w-full border-0 pointer-events-none bg-white"
        />
        <div className="absolute inset-0 flex flex-col items-center justify-center gap-1 bg-blue-50/90 opacity-0 transition-opacity group-hover:opacity-100">
          <FileText className="h-8 w-8 text-blue-600" />
          <span className="text-xs font-medium text-blue-700">Open PDF</span>
        </div>
      </a>
    );
  }

  return <span className="text-xs text-gray-400">No preview</span>;
}

function AddNoteForm({
  paperTitle,
  paperId,
  session,
  onSuccess,
}: {
  paperTitle: string;
  paperId: number;
  session: Session;
  onSuccess: () => void;
}) {
  const [thumbnailPreview, setThumbnailPreview] = useState<string | null>(null);
  const [pdfPreviewUrl, setPdfPreviewUrl] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formKey, setFormKey] = useState(0);

  useEffect(() => {
    return () => {
      if (pdfPreviewUrl) URL.revokeObjectURL(pdfPreviewUrl);
    };
  }, [pdfPreviewUrl]);

  const clearPreviews = () => {
    setThumbnailPreview(null);
    if (pdfPreviewUrl) URL.revokeObjectURL(pdfPreviewUrl);
    setPdfPreviewUrl(null);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!session?.accessToken) return;

    const form = e.currentTarget;
    const pdfInput = form.elements.namedItem("pdf") as HTMLInputElement | null;
    if (!pdfInput?.files?.[0]) {
      toast.error("PDF file is required.");
      return;
    }

    const fd = new FormData(form);
    fd.append("paper_id", String(paperId));

    setIsSubmitting(true);
    try {
      const res = await axiosInstance.post(
        "/private/academic-notes/notes/create",
        fd,
        { headers: authHeaders(session) }
      );
      toast.success(res.data.message);
      clearPreviews();
      setFormKey((k) => k + 1);
      onSuccess();
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err.response?.data?.error || "Something went wrong.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="rounded-lg border border-dashed border-blue-200 bg-blue-50/40 p-4">
      <div className="flex items-center gap-2 mb-3">
        <LuFileText className="text-primary" size={18} />
        <div>
          <p className="text-sm font-medium">Add PDF Note</p>
          <p className="text-xs text-gray-500">Upload a lecture sheet to {paperTitle}</p>
        </div>
      </div>
      <form
        key={formKey}
        onSubmit={handleSubmit}
        className="grid grid-cols-1 md:grid-cols-2 gap-3"
        encType="multipart/form-data"
      >
        <div>
          <Label>Note Title</Label>
          <InputField name="title" placeholder="অপরিচিতা" required />
        </div>
        <div>
          <Label>Subtitle (optional)</Label>
          <InputField name="subtitle" placeholder="অপরিচিতা" />
        </div>
        <div>
          <Label>Thumbnail (optional)</Label>
          <input
            type="file"
            name="thumbnail"
            accept="image/*"
            className="text-sm w-full"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (!file) {
                setThumbnailPreview(null);
                return;
              }
              const reader = new FileReader();
              reader.onloadend = () => setThumbnailPreview(reader.result as string);
              reader.readAsDataURL(file);
            }}
          />
        </div>
        <div>
          <Label>PDF File</Label>
          <input
            type="file"
            name="pdf"
            accept="application/pdf"
            required
            className="text-sm w-full"
            onChange={(e) => {
              const file = e.target.files?.[0];
              if (pdfPreviewUrl) URL.revokeObjectURL(pdfPreviewUrl);
              setPdfPreviewUrl(file ? URL.createObjectURL(file) : null);
            }}
          />
        </div>
        {(thumbnailPreview || pdfPreviewUrl) && (
          <div className="md:col-span-2 grid grid-cols-1 sm:grid-cols-2 gap-3">
            {thumbnailPreview && (
              <div>
                <p className="text-xs text-gray-500 mb-1">Thumbnail preview</p>
                <div className="h-32 border rounded-lg overflow-hidden bg-white flex items-center justify-center">
                  {/* eslint-disable-next-line @next/next/no-img-element */}
                  <img
                    src={thumbnailPreview}
                    alt="Thumbnail preview"
                    className="max-h-full max-w-full object-contain"
                  />
                </div>
              </div>
            )}
            {pdfPreviewUrl && (
              <div className={thumbnailPreview ? "" : "sm:col-span-2"}>
                <p className="text-xs text-gray-500 mb-1">PDF preview</p>
                <div className="h-32 border rounded-lg overflow-hidden bg-white">
                  <iframe
                    src={pdfPreviewUrl}
                    title="PDF preview"
                    className="h-full w-full border-0"
                  />
                </div>
              </div>
            )}
          </div>
        )}
        <div className="md:col-span-2 flex items-end justify-between gap-3">
          <div className="w-28">
            <Label>Position</Label>
            <InputField name="position" type="number" defaultValue={0} />
          </div>
          <Button type="submit" disabled={isSubmitting}>
            <LuPlus size={14} /> Add Note
          </Button>
        </div>
      </form>
    </div>
  );
}

export default function ClassNoteManager({ session, classData }: Props) {
  const router = useRouter();
  const [subjectTitle, setSubjectTitle] = useState("");
  const [subjectSlug, setSubjectSlug] = useState("");

  const subjects = classData.subjects ?? [];
  const totalNotes = subjects.reduce((sum, subject) => {
    const paperNotes = (subject.papers ?? []).reduce(
      (pSum, paper) => pSum + (paper.notes?.length ?? 0),
      0
    );
    return sum + paperNotes;
  }, 0);

  const refresh = () => router.refresh();

  const deleteItem = async (path: string, label: string) => {
    if (!session?.accessToken) return;
    if (!confirm(`Delete this ${label}?`)) return;
    try {
      const res = await axiosInstance.delete(path, {
        headers: authHeaders(session),
      });
      toast.success(res.data.message);
      refresh();
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err.response?.data?.error || "Something went wrong.");
    }
  };

  const submitForm = async (
    method: "post" | "put",
    path: string,
    fd: FormData
  ) => {
    if (!session?.accessToken) return;
    try {
      const res =
        method === "post"
          ? await axiosInstance.post(path, fd, { headers: authHeaders(session) })
          : await axiosInstance.put(path, fd, { headers: authHeaders(session) });
      toast.success(res.data.message);
      refresh();
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err.response?.data?.error || "Something went wrong.");
    }
  };

  const handleAddSubject = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const fd = new FormData(form);
    fd.append("class_id", String(classData.id));
    await submitForm("post", "/private/academic-notes/subjects/create", fd);
    form.reset();
    setSubjectTitle("");
    setSubjectSlug("");
  };

  const handleAddPaper = async (
    e: React.FormEvent<HTMLFormElement>,
    subjectId: number
  ) => {
    e.preventDefault();
    const form = e.currentTarget;
    const fd = new FormData(form);
    fd.append("subject_id", String(subjectId));
    await submitForm("post", "/private/academic-notes/papers/create", fd);
    form.reset();
  };

  return (
    <div className="mt-6 space-y-6">
      <div className="rounded-xl border bg-gradient-to-r from-slate-50 to-blue-50 p-4">
        <p className="text-sm font-medium mb-3">How to add notes</p>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <StepBadge step={1} label="Add a subject (e.g. Bangla)" />
          <StepBadge step={2} label="Add a paper (e.g. 1st Paper)" />
          <StepBadge step={3} label="Upload PDF lecture sheets" />
        </div>
        <p className="text-xs text-gray-500 mt-3">
          {subjects.length} subject{subjects.length === 1 ? "" : "s"} · {totalNotes} note
          {totalNotes === 1 ? "" : "s"} in this class
        </p>
      </div>

      <section className="rounded-xl border p-4">
        <div className="flex items-center gap-2 mb-4">
          <LuLayers className="text-primary" size={18} />
          <div>
            <h4 className="font-medium">Step 1 — Add Subject</h4>
            <p className="text-xs text-gray-500">Subjects group papers under this class</p>
          </div>
        </div>
        <form
          onSubmit={handleAddSubject}
          className="grid grid-cols-1 md:grid-cols-4 gap-3 items-end"
        >
          <div>
            <Label>Subject Title</Label>
            <InputField
              name="title"
              placeholder="Bangla"
              required
              value={subjectTitle}
              onChange={(e) => {
                const title = e.target.value;
                setSubjectTitle(title);
                setSubjectSlug(slugifyTitle(title));
              }}
            />
          </div>
          <div>
            <Label>Slug (optional)</Label>
            <InputField
              name="slug"
              placeholder="bangla"
              value={subjectSlug}
              onChange={(e) => setSubjectSlug(e.target.value)}
            />
          </div>
          <div>
            <Label>Position</Label>
            <InputField name="position" type="number" defaultValue={0} />
          </div>
          <Button type="submit">
            <LuPlus size={16} /> Add Subject
          </Button>
        </form>
      </section>

      {subjects.length === 0 ? (
        <div className="rounded-xl border border-dashed p-8 text-center text-sm text-gray-500">
          No subjects yet. Add your first subject above to continue.
        </div>
      ) : (
        subjects.map((subject) => (
          <section key={subject.id} className="rounded-xl border overflow-hidden">
            <div className="flex items-center justify-between gap-3 p-4 bg-gray-50 border-b">
              <div className="flex items-center gap-3">
                <div className="w-9 h-9 rounded-lg bg-primary/10 text-primary flex items-center justify-center">
                  <LuBookOpen size={18} />
                </div>
                <div>
                  <p className="font-medium">{subject.title}</p>
                  <p className="text-xs text-gray-500">
                    {subject.papers?.length ?? 0} paper
                    {(subject.papers?.length ?? 0) === 1 ? "" : "s"} · {subject.note_count ?? 0}{" "}
                    notes · {subject.slug}
                  </p>
                </div>
              </div>
              <button
                type="button"
                className="text-red-500 p-2 hover:bg-red-50 rounded-lg"
                onClick={() =>
                  deleteItem(
                    `/private/academic-notes/subjects/delete/${subject.id}`,
                    "subject"
                  )
                }
              >
                <LuTrash2 size={16} />
              </button>
            </div>

            <div className="p-4 space-y-4">
              <div className="rounded-lg border bg-slate-50 p-4">
                <div className="flex items-center gap-2 mb-3">
                  <StepBadge step={2} label={`Add Paper under ${subject.title}`} />
                </div>
                <form
                  onSubmit={(e) => handleAddPaper(e, subject.id)}
                  className="grid grid-cols-1 md:grid-cols-5 gap-3 items-end"
                >
                  <div>
                    <Label>Paper Title</Label>
                    <InputField name="title" placeholder="Bangla 1st Paper" required />
                  </div>
                  <div>
                    <Label>Icon Label</Label>
                    <InputField name="icon_label" placeholder="১ম" />
                  </div>
                  <div>
                    <Label>Icon Color</Label>
                    <InputField name="icon_color" placeholder="#42A5F5" />
                  </div>
                  <div>
                    <Label>Position</Label>
                    <InputField name="position" type="number" defaultValue={0} />
                  </div>
                  <Button type="submit" className="text-sm">
                    <LuPlus size={14} /> Add Paper
                  </Button>
                </form>
              </div>

              {(subject.papers?.length ?? 0) === 0 ? (
                <p className="text-sm text-gray-500 text-center py-4 border border-dashed rounded-lg">
                  No papers yet. Add a paper above, then upload PDF notes below it.
                </p>
              ) : (
                subject.papers?.map((paper) => (
                  <div key={paper.id} className="rounded-lg border ml-0 sm:ml-3">
                    <div className="flex items-center justify-between gap-3 p-3 border-b bg-white">
                      <div className="flex items-center gap-2">
                        {paper.icon_color ? (
                          <span
                            className="w-8 h-8 rounded-full text-white text-xs flex items-center justify-center shrink-0"
                            style={{ backgroundColor: paper.icon_color }}
                          >
                            {paper.icon_label}
                          </span>
                        ) : (
                          <span className="w-8 h-8 rounded-full bg-blue-100 text-blue-700 text-xs flex items-center justify-center shrink-0">
                            <LuFileText size={14} />
                          </span>
                        )}
                        <div>
                          <p className="font-medium text-sm">{paper.title}</p>
                          <p className="text-xs text-gray-500">
                            {paper.notes?.length ?? 0} note
                            {(paper.notes?.length ?? 0) === 1 ? "" : "s"}
                          </p>
                        </div>
                      </div>
                      <button
                        type="button"
                        className="text-red-500 p-1.5 hover:bg-red-50 rounded"
                        onClick={() =>
                          deleteItem(
                            `/private/academic-notes/papers/delete/${paper.id}`,
                            "paper"
                          )
                        }
                      >
                        <LuTrash2 size={14} />
                      </button>
                    </div>

                    <div className="p-4 space-y-4 bg-slate-50/50">
                      <AddNoteForm
                        paperTitle={paper.title}
                        paperId={paper.id}
                        session={session}
                        onSuccess={refresh}
                      />

                      {(paper.notes?.length ?? 0) > 0 ? (
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                          {paper.notes?.map((note) => (
                            <div
                              key={note.id}
                              className="border rounded-lg overflow-hidden bg-white"
                            >
                              <div className="h-36 bg-blue-50 overflow-hidden">
                                <NotePreview note={note} />
                              </div>
                              <div className="p-3">
                                <p className="font-medium text-sm">{note.title}</p>
                                {note.subtitle && (
                                  <p className="text-xs text-gray-500">{note.subtitle}</p>
                                )}
                                <div className="flex items-center justify-between mt-2">
                                  <a
                                    href={note.pdf_url}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-xs text-green-600 font-medium"
                                  >
                                    PDF
                                  </a>
                                  <button
                                    type="button"
                                    className="text-red-500 text-xs flex items-center gap-1"
                                    onClick={() =>
                                      deleteItem(
                                        `/private/academic-notes/notes/delete/${note.id}`,
                                        "note"
                                      )
                                    }
                                  >
                                    <LuTrash2 size={12} /> Delete
                                  </button>
                                </div>
                              </div>
                            </div>
                          ))}
                        </div>
                      ) : (
                        <p className="text-xs text-gray-500 text-center py-2">
                          No notes in this paper yet. Use the form above to upload your first PDF.
                        </p>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </section>
        ))
      )}
    </div>
  );
}
