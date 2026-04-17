"use client";

import Button from "@/components/ui/Button";
import { GenericSelectWithSearch } from "@/components/ui/GenericSelectWithSearch";
import Modal from "@/components/ui/Modal";
import React, { useState } from "react";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";

type Student = {
  id: number;
  name: string;
};

type Course = {
  id: number;
  name: string;
};

const students: Student[] = [
  { id: 1, name: "Arifin" },
  { id: 2, name: "Borhan" },
];

const courses: Course[] = [
  { id: 1, name: "Mathematics" },
  { id: 2, name: "Physics" },
];

export default function AddNewDigitalAccess() {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedStudent, setSelectedStudent] = useState<Student | null>(null);
  const [selectedCourse, setSelectedCourse] = useState<Course | null>(null);

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>
        <LuPlus /> Add New Digital Access
      </Button>
      <Modal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        className="p-0 max-w-md"
      >
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-300">
          <p className="font-medium">Create Digital Access</p>
          <button onClick={() => setIsOpen(false)}>
            <RxCross2 />
          </button>
        </div>
        <div className="p-4">
          <div className="mb-5">
            <label className="text-sm block mb-1 font-medium">
              Select Student <span className="text-red-500">*</span>
            </label>
            <GenericSelectWithSearch
              items={students}
              selectedItem={selectedStudent}
              onSelect={setSelectedStudent}
              getLabel={(student) => student.name}
              className="w-full"
            />
          </div>
          <div>
            <label className="text-sm block mb-1 font-medium">
              Select Digital Download <span className="text-red-500">*</span>
            </label>
            <GenericSelectWithSearch
              items={courses}
              selectedItem={selectedCourse}
              onSelect={setSelectedCourse}
              getLabel={(course) => course.name}
              className="w-full"
            />
          </div>
          <div className="flex items-center justify-end mt-5 gap-3">
            <button
              className="px-4 py-2 text-sm border rounded-full"
              onClick={() => setIsOpen(false)}
            >
              Cancel
            </button>
            <Button>Create</Button>
          </div>
        </div>
      </Modal>
    </>
  );
}
