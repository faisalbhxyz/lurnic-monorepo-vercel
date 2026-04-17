"use client";

import React from "react";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import { RxCross2 } from "react-icons/rx";
import QuizEdit from "./QuizEdit";

export default function UpdateQuizModal() {
  const { closeEditQuiz, isEditQuiz } = useCoursesStore();

  return (
    <div
      className={`fixed inset-0 flex flex-col justify-end transition-opacity duration-300 ${
        isEditQuiz ? "opacity-100" : "opacity-0 pointer-events-none"
      } bg-black/50`}
    >
      <div className="p-3 flex items-center justify-end">
        <button onClick={closeEditQuiz} className="text-white">
          <RxCross2 size={22} />
        </button>
      </div>
      <div
        className={`flex-1 bg-white transition-transform duration-300 rounded-t-2xl overflow-y-auto ${
          isEditQuiz ? "translate-y-0" : "translate-y-full"
        }`}
      >
        <QuizEdit isEdit />
      </div>
    </div>
  );
}
