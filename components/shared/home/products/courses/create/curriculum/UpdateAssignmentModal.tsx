"use client";

import React from "react";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import { RxCross2 } from "react-icons/rx";
import AssignmentEdit from "./AssignmentEdit";

export default function UpdateAssignmentModal() {
  const { closeEditAssignment, isEditAssignment } = useCoursesStore();

  return (
    <div
      className={`fixed inset-0 flex flex-col justify-end transition-opacity duration-300 ${
        isEditAssignment ? "opacity-100" : "opacity-0 pointer-events-none"
      } bg-black/50`}
    >
      <div className="p-3 flex items-center justify-end">
        <button onClick={closeEditAssignment} className="text-white">
          <RxCross2 size={22} />
        </button>
      </div>
      <div
        className={`flex-1 bg-white transition-transform duration-300 rounded-t-2xl overflow-y-auto ${
          isEditAssignment ? "translate-y-0" : "translate-y-full"
        }`}
      >
        <AssignmentEdit isEdit />
      </div>
    </div>
  );
}
