"use client";

import React, { useState } from "react";
import { Session } from "next-auth";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import Label from "@/components/ui/Label";
import { LuPlus, LuTrash2 } from "react-icons/lu";
import Image from "next/image";

type Props = {
  session: Session;
  classData: IAcademicNoteClassDetail;
};

function authHeaders(session: Session) {
  return { Authorization: `Bearer ${session.accessToken}` };
}

export default function ClassNoteManager({ session, classData }: Props) {
  const router = useRouter();
  const [expandedSubject, setExpandedSubject] = useState<number | null>(null);
  const [expandedPaper, setExpandedPaper] = useState<number | null>(null);

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

  const handleAddNote = async (
    e: React.FormEvent<HTMLFormElement>,
    paperId: number
  ) => {
    e.preventDefault();
    if (!session?.accessToken) return;
    const form = e.currentTarget;
    const fd = new FormData(form);
    fd.append("paper_id", String(paperId));
    try {
      const res = await axiosInstance.post(
        "/private/academic-notes/notes/create",
        fd,
        { headers: authHeaders(session) }
      );
      toast.success(res.data.message);
      form.reset();
      refresh();
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } } };
      toast.error(err.response?.data?.error || "Something went wrong.");
    }
  };

  return (
    <div className="mt-8 space-y-6">
      <div className="border rounded-xl p-4">
        <h4 className="font-medium mb-3">Add Subject</h4>
        <form onSubmit={handleAddSubject} className="grid grid-cols-1 md:grid-cols-4 gap-3 items-end">
          <div>
            <Label>Subject Title</Label>
            <InputField name="title" placeholder="Bangla" required />
          </div>
          <div>
            <Label>Slug (optional)</Label>
            <InputField name="slug" placeholder="bangla" />
          </div>
          <div>
            <Label>Position</Label>
            <InputField name="position" type="number" defaultValue={0} />
          </div>
          <Button type="submit">
            <LuPlus size={16} /> Add Subject
          </Button>
        </form>
      </div>

      {classData.subjects?.map((subject) => (
        <div key={subject.id} className="border rounded-xl overflow-hidden">
          <div
            className="w-full flex items-center justify-between p-4 bg-gray-50 hover:bg-gray-100 cursor-pointer"
            onClick={() =>
              setExpandedSubject(expandedSubject === subject.id ? null : subject.id)
            }
          >
            <div>
              <p className="font-medium">{subject.title}</p>
              <p className="text-xs text-gray-500">
                {subject.note_count ?? 0} notes · {subject.slug}
              </p>
            </div>
            <button
              type="button"
              className="text-red-500 p-2 hover:bg-red-50 rounded"
              onClick={(e) => {
                e.stopPropagation();
                deleteItem(
                  `/private/academic-notes/subjects/delete/${subject.id}`,
                  "subject"
                );
              }}
            >
              <LuTrash2 size={16} />
            </button>
          </div>

          {expandedSubject === subject.id && (
            <div className="p-4 space-y-4 border-t">
              <form
                onSubmit={(e) => handleAddPaper(e, subject.id)}
                className="grid grid-cols-1 md:grid-cols-5 gap-3 items-end bg-slate-50 p-3 rounded-lg"
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

              {subject.papers?.map((paper) => (
                <div key={paper.id} className="border rounded-lg ml-2">
                  <div
                    className="w-full flex items-center justify-between p-3 hover:bg-gray-50 cursor-pointer"
                    onClick={() =>
                      setExpandedPaper(expandedPaper === paper.id ? null : paper.id)
                    }
                  >
                    <div className="flex items-center gap-2">
                      {paper.icon_color && (
                        <span
                          className="w-7 h-7 rounded-full text-white text-xs flex items-center justify-center"
                          style={{ backgroundColor: paper.icon_color }}
                        >
                          {paper.icon_label}
                        </span>
                      )}
                      <div>
                        <p className="font-medium text-sm">{paper.title}</p>
                        <p className="text-xs text-gray-500">
                          {paper.note_count ?? 0} notes
                        </p>
                      </div>
                    </div>
                    <button
                      type="button"
                      className="text-red-500 p-1"
                      onClick={(e) => {
                        e.stopPropagation();
                        deleteItem(
                          `/private/academic-notes/papers/delete/${paper.id}`,
                          "paper"
                        );
                      }}
                    >
                      <LuTrash2 size={14} />
                    </button>
                  </div>

                  {expandedPaper === paper.id && (
                    <div className="p-3 border-t space-y-3">
                      <form
                        onSubmit={(e) => handleAddNote(e, paper.id)}
                        className="grid grid-cols-1 md:grid-cols-2 gap-3 bg-slate-50 p-3 rounded-lg"
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
                          />
                        </div>
                        <div>
                          <Label>Position</Label>
                          <InputField name="position" type="number" defaultValue={0} />
                        </div>
                        <div className="flex items-end">
                          <Button type="submit" className="text-sm">
                            <LuPlus size={14} /> Add Note
                          </Button>
                        </div>
                      </form>

                      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                        {paper.notes?.map((note) => (
                          <div
                            key={note.id}
                            className="border rounded-lg overflow-hidden"
                          >
                            <div className="h-28 bg-blue-50 flex items-center justify-center">
                              {note.thumbnail ? (
                                <Image
                                  src={note.thumbnail}
                                  alt={note.title}
                                  width={80}
                                  height={100}
                                  className="h-24 w-auto object-contain"
                                />
                              ) : (
                                <span className="text-xs text-gray-400">No thumbnail</span>
                              )}
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
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
