"use client";

import Button from "@/components/ui/Button";
import { GenericSelectWithSearch } from "@/components/ui/GenericSelectWithSearch";
import Modal from "@/components/ui/Modal";
import axiosInstance from "@/lib/axiosInstance";
import { zodResolver } from "@hookform/resolvers/zod";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";
import { toast } from "sonner";
import { z } from "zod";

const EnrollmentSchema = z.object({
  student_id: z.coerce.number().min(1),
  course_id: z.coerce.number().min(1),
});

type IEnrollment = z.infer<typeof EnrollmentSchema>;

export default function AddNewEnrollment({
  students,
  courses,
}: {
  students: Pick<IStudent, "id" | "first_name" | "last_name">[];
  courses: Pick<CourseDetails, "id" | "title">[];
}) {
  const router = useRouter();
  const { data: session } = useSession();
  const [isOpen, setIsOpen] = useState(false);

  const mappedStds = students.map((student) => ({
    id: student.id,
    name: `${student.first_name} ${student.last_name ?? ""}`,
  }));

  const mappedCourses = courses.map((course) => ({
    id: course.id,
    name: course.title,
  }));

  const { handleSubmit, control, reset } = useForm<IEnrollment>({
    resolver: zodResolver(EnrollmentSchema),
  });

  const handleSave = (data: IEnrollment) => {
    axiosInstance
      .post("/private/enrollment/create", data, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${session?.accessToken}`,
        },
      })
      .then((res) => {
        toast.success(res.data.message);
        setIsOpen(false);
        router.refresh();
        reset();
      })
      .catch((error) => {
        toast.error(error.response.data.error || "Something went wrong.");
      });
  };

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>
        <LuPlus /> Add New Enrollment
      </Button>
      <Modal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        className="p-0 max-w-md"
      >
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-300">
          <p className="font-medium">Create Enrollment</p>
          <button onClick={() => setIsOpen(false)}>
            <RxCross2 />
          </button>
        </div>
        <form onSubmit={handleSubmit(handleSave)}>
          <div className="p-4">
            <div className="mb-5">
              <label className="text-sm block mb-1 font-medium">
                Select Student <span className="text-red-500">*</span>
              </label>
              <Controller
                control={control}
                name="student_id"
                render={({ field: { onChange, value } }) => (
                  <GenericSelectWithSearch
                    items={mappedStds}
                    selectedItem={mappedStds.find(
                      (student) => student.id === value
                    )}
                    onSelect={(student) => onChange(student.id)}
                    getLabel={(student) => student.name}
                    className="w-full"
                  />
                )}
              />
            </div>
            <div>
              <label className="text-sm block mb-1 font-medium">
                Select Course <span className="text-red-500">*</span>
              </label>
              <Controller
                control={control}
                name="course_id"
                render={({ field: { onChange, value } }) => (
                  <GenericSelectWithSearch
                    items={mappedCourses}
                    selectedItem={mappedCourses.find(
                      (course) => course.id === value
                    )}
                    onSelect={(course) => onChange(course.id)}
                    getLabel={(course) => course.name}
                    placeholder="Select course"
                    className="w-full"
                  />
                )}
              />
            </div>
            <div className="flex items-center justify-end mt-5 gap-3">
              <button
                type="button"
                className="px-4 py-2 text-sm border rounded-full"
                onClick={() => setIsOpen(false)}
              >
                Cancel
              </button>
              <Button type="submit">Create</Button>
            </div>
          </div>
        </form>
      </Modal>
    </>
  );
}
