"use client";

import { useCoursesStore } from "@/hooks/useCoursesStore";
import { items } from "@/lib/constants";
import React, { useEffect, useState } from "react";
import { RxCross2 } from "react-icons/rx";
import AssignmentEdit from "./AssignmentEdit";
import QuizEdit from "./QuizEdit";
import LessonItem from "./LessonItem";

export default function EditItemModal() {
  const { editItemId, setEditItem } = useCoursesStore();
  const [show, setShow] = useState(false);

  useEffect(() => {
    if (editItemId) {
      setShow(true);
    }
  }, [editItemId]);

  const handleClose = () => {
    setShow(false);
    setTimeout(() => setEditItem(null), 300);
  };

  if (!editItemId) return null;

  const item = items.find((i) => i.id === editItemId);

  return (
    <div
      className={`fixed inset-0 flex flex-col justify-end transition-opacity duration-300 ${
        show ? "opacity-100" : "opacity-0 pointer-events-none"
      } bg-black/50`}
    >
      <div className="p-3 flex items-center justify-end">
        <button onClick={handleClose} className="text-white">
          <RxCross2 size={22} />
        </button>
      </div>
      <div
        className={`flex-1 bg-white transition-transform duration-300 rounded-t-2xl overflow-y-auto ${
          show ? "translate-y-0" : "translate-y-full"
        }`}
      >
        {item?.type === "lesson" && <LessonItem isEdit />}
        {item?.type === "quiz" && <QuizEdit type="edit" />}
        {item?.type === "assignment" && <AssignmentEdit isEdit />}
      </div>
    </div>
  );
}
