import CreateClassNoteForm from "@/components/shared/home/website/class-notes/CreateClassNoteForm";
import ClassNoteManager from "@/components/shared/home/website/class-notes/ClassNoteManager";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import Link from "next/link";
import { notFound } from "next/navigation";
import React from "react";

export const getClassNoteById = async (
  session: Session,
  id: string
): Promise<IAcademicNoteClassDetail | null> => {
  try {
    const res = await axiosInstance.get(`/private/academic-notes/classes/${id}`, {
      headers: {
        Authorization: `Bearer ${session.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    console.log("[ERROR]", error);
    return null;
  }
};

export default async function EditClassNotePage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const session = await auth();
  if (!session) return null;

  const { id } = await params;
  const classData = await getClassNoteById(session, id);
  if (!classData) notFound();

  const subjectCount = classData.subjects?.length ?? 0;
  const noteCount = classData.note_count ?? 0;

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href="/class-notes" className="text-gray-500 hover:text-primary">
          Class-wise Notes
        </Link>
        /<span>{classData.title}</span>
      </div>

      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mt-5">
        <div>
          <h3 className="font-medium text-lg">{classData.title}</h3>
          <p className="text-sm text-gray-500 mt-1">
            Add subjects, papers, and PDF notes — everything is on this page.
          </p>
        </div>
        <div className="flex gap-2 text-sm">
          <span className="px-3 py-1 rounded-full bg-slate-100 text-gray-700">
            {subjectCount} subject{subjectCount === 1 ? "" : "s"}
          </span>
          <span className="px-3 py-1 rounded-full bg-blue-50 text-blue-700">
            {noteCount} note{noteCount === 1 ? "" : "s"}
          </span>
        </div>
      </div>

      <ClassNoteManager session={session} classData={classData} />

      <details className="mt-8 border rounded-xl overflow-hidden group">
        <summary className="cursor-pointer px-4 py-3 bg-gray-50 font-medium text-sm hover:bg-gray-100">
          Class settings (title, icon, publish)
        </summary>
        <div className="p-5 border-t">
          <CreateClassNoteForm session={session} isEdit classData={classData} />
        </div>
      </details>
    </>
  );
}
