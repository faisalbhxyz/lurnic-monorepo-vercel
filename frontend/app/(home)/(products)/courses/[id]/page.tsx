import { getCourseByID } from "@/app/actions/course_actions";
import BackButton from "@/components/shared/BackButton";
import CourseDetailsView from "@/components/shared/home/products/courses/CourseDetailsView";
import { auth } from "@/lib/auth";
import Link from "next/link";
import { notFound } from "next/navigation";
import React, { Suspense } from "react";

export default async function CourseDetailsPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const session = await auth();
  if (!session) return null;

  const { id } = await params;
  const courseDetails = await getCourseByID(session, id);

  if (!courseDetails) notFound();

  return (
    <div className="w-full max-w-6xl mx-auto py-5">
      <div className="flex items-center text-sm gap-1 mb-4">
        <Link href="/courses" className="text-gray-500 hover:text-gray-700">
          Products
        </Link>
        /
        <Link href="/courses" className="text-gray-500 hover:text-gray-700">
          Courses
        </Link>
        /<span>{courseDetails.title}</span>
      </div>

      <BackButton buttonText="Courses" />

      <Suspense fallback={<p className="text-sm text-gray-500 mt-5">Loading...</p>}>
        <CourseDetailsView course={courseDetails} />
      </Suspense>
    </div>
  );
}
