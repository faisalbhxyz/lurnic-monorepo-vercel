"use client";

import React from "react";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import Modal from "@/components/ui/Modal";
import { RxCross2 } from "react-icons/rx";
import NewChapter from "./NewChapter";

export default function EditChapter() {
  const { isEditChapter, closeEditChapter } = useCoursesStore();

  return (
    <Modal isOpen={isEditChapter} onClose={closeEditChapter} className="p-0">
      <div className="p-4 flex items-center justify-between border-b border-gray-300">
        <p className="font-semibold text-lg">Edit Chapter</p>
        <button type="button" onClick={closeEditChapter}>
          <RxCross2 />
        </button>
      </div>
      <NewChapter isEdit/>
    </Modal>
  );
}
