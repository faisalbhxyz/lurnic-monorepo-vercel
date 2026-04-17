"use client";

import React from "react";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import Modal from "@/components/ui/Modal";
import { RxCross2 } from "react-icons/rx";
import NewChapter from "./NewChapter";

export default function CreateNewChapterModal() {
  const { isNewChapter, closeNewChapter } = useCoursesStore();

  return (
    <Modal isOpen={isNewChapter} onClose={closeNewChapter} className="p-0">
      <div className="p-4 flex items-center justify-between border-b border-gray-300">
        <p className="font-semibold text-lg">New Chapter</p>
        <button type="button" onClick={closeNewChapter}>
          <RxCross2 />
        </button>
      </div>
      <NewChapter isEdit={false}/>
    </Modal>
  );
}
