import React from "react";
import CreateStudent from "@/components/shared/home/users/students/CreateStudent";
import StudentList from "@/components/shared/home/users/students/StudentList";
import Link from "next/link";
import { auth } from "@/lib/auth";
import { Session } from "next-auth";
import axiosInstance from "@/lib/axiosInstance";
import UpdateStudent from "@/components/shared/home/users/students/UpdateStudent";
import { getGeneralSettings } from "@/app/actions/actions";

export const getAllStudents = async (session: Session): Promise<IStudent[]> => {
  try {
    const res = await axiosInstance.get("/private/student", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    return [];
  }
};

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const students = await getAllStudents(session);
  const generalSettings = await getGeneralSettings(session);

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href={""} className="text-gray-500">
          Users
        </Link>
        /<Link href={""}>Students</Link>
      </div>
      <div className="flex-between my-5">
        <h3 className="font-medium text-2xl">Students</h3>
        <CreateStudent />
        <UpdateStudent />
      </div>
      <StudentList studentPrefix={generalSettings?.student_prefix} data={students} />
    </>
  );
}
