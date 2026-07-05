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

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href="/class-notes" className="text-gray-500">
          Class-wise Notes
        </Link>
        /<span>{classData.title}</span>
      </div>
      <h3 className="font-medium mt-5">Manage: {classData.title}</h3>

      <div className="mt-5 border rounded-xl p-5">
        <h4 className="font-medium mb-4">Class Settings</h4>
        <CreateClassNoteForm session={session} isEdit classData={classData} />
      </div>

      <ClassNoteManager session={session} classData={classData} />
    </>
  );
}
