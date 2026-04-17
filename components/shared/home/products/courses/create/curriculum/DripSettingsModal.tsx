"use client";

import React, { useEffect, useState } from "react";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import Modal from "@/components/ui/Modal";
import { RxCross2 } from "react-icons/rx";
import LessonDrip from "./LessonDrip";
import QuizDrip from "./QuizDrip";
import AssignmentDrip from "./AssignmentDrip";

export default function DripSettingsModal() {
  const { dripSettings, closeDripSettings } = useCoursesStore();
  const [show, setShow] = useState(false);

  useEffect(() => {
    if (dripSettings) {
      setShow(true);
    }
  }, [dripSettings]);

  const handleClose = () => {
    setShow(false);
    setTimeout(() => closeDripSettings(null), 300);
  };

  if (!dripSettings) return null;

  return (
    <Modal isOpen={show} onClose={handleClose} className="p-0">
      <div className="p-4 flex items-center justify-between border-b border-gray-300">
        <h2 className="text-lg font-medium">Drip Settings</h2>
        <button
          type="button"
          onClick={handleClose}
          className="text-gray-600 hover:text-black transition"
          aria-label="Close"
        >
          <RxCross2 size={20} />
        </button>
      </div>

      <div className="p-4">
        {/* Replace with actual settings component or form */}
        {dripSettings === "Lesson" && <LessonDrip />}
        {dripSettings === "Quiz" && <QuizDrip />}
        {dripSettings === "Assignment" && <AssignmentDrip />}
      </div>
    </Modal>
  );
}
