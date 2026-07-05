import React from "react";
import Link from "next/link";
import Image from "next/image";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { notFound } from "next/navigation";
import StudentDetailsView from "@/components/shared/home/users/students/StudentDetailsView";
import UpdateStudent from "@/components/shared/home/users/students/UpdateStudent";
import { getGeneralSettings } from "@/app/actions/actions";

async function getStudentDetails(id: string): Promise<IStudentDetails | null> {
  const session = await auth();
  if (!session?.accessToken) return null;

  try {
    const res = await axiosInstance.get(`/private/student/details/${id}`, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session.accessToken}`,
      },
    });
    return res.data.data;
  } catch {
    return null;
  }
}

export default async function StudentDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  const session = await auth();
  if (!session) return null;

  const [student, generalSettings] = await Promise.all([
    getStudentDetails(id),
    getGeneralSettings(session),
  ]);

  if (!student) notFound();

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href="/students" className="text-gray-500">
          Users
        </Link>
        /<Link href="/students">Students</Link>/
        <span>
          {student.first_name} {student.last_name ?? ""}
        </span>
      </div>
      <StudentDetailsView
        student={student}
        studentPrefix={generalSettings?.student_prefix}
      />
      <UpdateStudent />
    </>
  );
}
