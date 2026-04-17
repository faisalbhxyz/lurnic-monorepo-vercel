import React from "react";
import AddNewEnrollment from "@/components/shared/home/sales/enrollments/AddNewEnrollment";
import EnrollmentList from "@/components/shared/home/sales/enrollments/EnrollmentList";
import Link from "next/link";
import { auth } from "@/lib/auth";
import { getStudentsLite } from "@/app/actions/actions";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import { getCoursesLite } from "@/app/actions/course_actions";

const getEnrollments = async (session: Session) => {
  try {
    const res = await axiosInstance.get("/private/enrollment", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    console.log("[ERROR]", error);
    return [];
  }
};

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const students = await getStudentsLite(session);
  const courses = await getCoursesLite(session);

  const enrollments = await getEnrollments(session);

  // console.log("enrollments", enrollments);

  return (
    <>
      <div className="flex items-center text-sm gap-1 mb-5">
        <Link href={""} className="text-gray-500">
          Sales
        </Link>
        /<Link href={""}>Enrollments</Link>
      </div>
      <div className="flex-between mb-5">
        <h3 className="font-medium text-2xl">Enrollments</h3>
        <AddNewEnrollment students={students} courses={courses} />
      </div>

      <EnrollmentList data={enrollments} session={session} />
    </>
  );
}
