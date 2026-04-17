import React from "react";
import AddNewInstructor from "@/components/shared/home/users/instructors/AddNewInstructor";
import InstructorList from "@/components/shared/home/users/instructors/InstructorList";
import Link from "next/link";
import { Session } from "next-auth";
import axiosInstance from "@/lib/axiosInstance";
import { auth } from "@/lib/auth";
import UpdateInstructor from "@/components/shared/home/users/instructors/UpdateInstructor";
import { getGeneralSettings } from "@/app/actions/actions";

export const getAllInstructors = async (session: Session) => {
  try {
    const res = await axiosInstance.get("/private/instructor", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    return null;
  }
};

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const sampleStudents = await getAllInstructors(session);
  const generalSettings = await getGeneralSettings(session);

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href={""} className="text-gray-500">
          Users
        </Link>
        /<Link href={""}>Instructors</Link>
      </div>
      <div className="flex-between my-5">
        <h3 className="font-medium text-2xl">Instructors</h3>
        <AddNewInstructor />
        <UpdateInstructor />
      </div>
      <InstructorList teacherPrefix={generalSettings?.teacher_prefix} data={sampleStudents} />
    </>
  );
}
