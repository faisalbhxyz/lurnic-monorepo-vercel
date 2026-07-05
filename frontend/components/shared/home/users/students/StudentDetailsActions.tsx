"use client";

import Button from "@/components/ui/Button";
import { useEditStore } from "@/hooks/useEditStore";
import { BiEditAlt } from "react-icons/bi";
import { HiOutlineArrowLeft } from "react-icons/hi";

export default function StudentDetailsActions({
  studentId,
}: {
  studentId: number;
}) {
  const { setEditID, toggleStudentEdit } = useEditStore();

  return (
    <div className="flex items-center gap-2">
      <Button link src="/students" variant="outline">
        <HiOutlineArrowLeft /> Back
      </Button>
      <button
        type="button"
        className="inline-flex items-center gap-2 text-sm border border-primary text-primary px-3 py-1.5 rounded-md hover:bg-primary hover:text-white"
        onClick={() => {
          setEditID(studentId);
          toggleStudentEdit(true);
        }}
      >
        <BiEditAlt size={16} />
        Edit Student
      </button>
    </div>
  );
}
