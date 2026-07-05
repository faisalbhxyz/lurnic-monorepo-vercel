import AddClassNote from "@/components/shared/home/website/class-notes/AddClassNote";
import ClassNoteList from "@/components/shared/home/website/class-notes/ClassNoteList";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import Link from "next/link";
import React from "react";
export const getAllClassNotes = async (
  session: Session
): Promise<IAcademicNoteClass[] | null> => {
  try {
    const res = await axiosInstance.get("/private/academic-notes/classes", {
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

export default async function ClassNotesPage() {
  const session = await auth();
  if (!session) return null;

  const classes = await getAllClassNotes(session);

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href="/" className="text-gray-500">
          Website
        </Link>
        /<span>Class-wise Notes</span>
      </div>
      <div className="flex-between mt-5">
        <div>
          <h3 className="font-medium">Class-wise Notes</h3>
          <p className="text-sm text-gray-500 mt-1">
            Create a class, then add subjects, papers, and PDF notes in one place.
          </p>
        </div>
        <AddClassNote session={session} />
      </div>

      <div className="mt-5">
        {classes && classes.length > 0 ? (
          <ClassNoteList data={classes} />
        ) : (
          <p className="text-sm text-gray-500 text-center py-10">
            No classes yet. Add your first class (e.g. HSC, ৮ম শ্রেণি).
          </p>
        )}
      </div>
    </>
  );
}
