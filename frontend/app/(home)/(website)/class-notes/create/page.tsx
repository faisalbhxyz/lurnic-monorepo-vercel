import CreateClassNoteForm from "@/components/shared/home/website/class-notes/CreateClassNoteForm";
import { auth } from "@/lib/auth";
import Link from "next/link";
import React from "react";

export default async function CreateClassNotePage() {
  const session = await auth();
  if (!session) return null;

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href="/class-notes" className="text-gray-500">
          Class-wise Notes
        </Link>
        /<span>Create</span>
      </div>
      <h3 className="font-medium mt-5 mb-5">Create Class</h3>
      <CreateClassNoteForm session={session} />
    </>
  );
}
