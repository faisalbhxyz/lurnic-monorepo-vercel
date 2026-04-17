"use client";

import React from "react";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import { RxCross2 } from "react-icons/rx";
import AssignmentEdit from "./AssignmentEdit";

export default function CreateAssignmentModal() {
  const { closeNewAssignment, isNewAssignment } = useCoursesStore();

  return (
    <div
      className={`fixed inset-0 flex flex-col justify-end transition-opacity duration-300 ${
        isNewAssignment ? "opacity-100" : "opacity-0 pointer-events-none"
      } bg-black/50`}
    >
      <div className="p-3 flex items-center justify-end">
        <button onClick={closeNewAssignment} className="text-white">
          <RxCross2 size={22} />
        </button>
      </div>
      <div
        className={`flex-1 bg-white transition-transform duration-300 rounded-t-2xl overflow-y-auto ${
          isNewAssignment ? "translate-y-0" : "translate-y-full"
        }`}
      >
        <AssignmentEdit />
      </div>
    </div>
  );
}
