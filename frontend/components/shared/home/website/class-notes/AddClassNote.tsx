"use client";

import Button from "@/components/ui/Button";
import Modal from "@/components/ui/Modal";
import CreateClassNoteForm from "@/components/shared/home/website/class-notes/CreateClassNoteForm";
import { Session } from "next-auth";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";

export default function AddClassNote({ session }: { session: Session }) {
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);

  const handleClose = () => setIsOpen(false);

  const handleCreated = (id: number) => {
    setIsOpen(false);
    router.push(`/class-notes/${id}/edit`);
  };

  return (
    <>
      <Button type="button" onClick={() => setIsOpen(true)}>
        <LuPlus /> Add Class
      </Button>
      <Modal isOpen={isOpen} onClose={handleClose} className="p-0 max-w-2xl">
        <div className="flex items-center justify-between py-3 px-4 border-b border-gray-300">
          <div>
            <p className="font-medium text-lg">Create Class</p>
            <p className="text-xs text-gray-500">
              After saving, you&apos;ll go straight to add subjects and notes.
            </p>
          </div>
          <button type="button" onClick={handleClose}>
            <RxCross2 />
          </button>
        </div>
        <div className="p-5 max-h-[80vh] overflow-y-auto">
          <CreateClassNoteForm
            session={session}
            onCancel={handleClose}
            onCreated={handleCreated}
          />
        </div>
      </Modal>
    </>
  );
}
