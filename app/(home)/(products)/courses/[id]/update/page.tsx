import {
  getAllCategories,
  getAllInstructorsLite,
  getAllSubCategories,
} from "@/app/actions/actions";
import { getCourseByID } from "@/app/actions/course_actions";
import BackButton from "@/components/shared/BackButton";
import CoursesTabs from "@/components/shared/home/products/courses/create/CoursesTabs";
import AddNewQuestion from "@/components/shared/home/products/courses/create/curriculum/AddNewQuestion";
import DripSettingsModal from "@/components/shared/home/products/courses/create/curriculum/DripSettingsModal";
import EditItemModal from "@/components/shared/home/products/courses/create/curriculum/EditItemModal";
import { auth } from "@/lib/auth";
import React from "react";

export default async function page({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const session = await auth();
  if (!session) return null;

  const { id } = await params;

  const categories = await getAllCategories(session);
  const subcategories = await getAllSubCategories(session);
  const instructors = await getAllInstructorsLite(session);

  const courseDetails = await getCourseByID(session, id);

  // console.log("courseDetails", courseDetails);

  return (
    <div className="w-full max-w-6xl mx-auto py-5">
      <BackButton buttonText={"Courses"} />
      <CoursesTabs
        isEdit={true}
        categories={categories}
        subcategories={subcategories}
        instructors={instructors}
        courseDetails={courseDetails}
      />
      {/* <CreateQuizModal /> */}
      {/* <CreateAssignmentModal/> */}
      <EditItemModal />
      <AddNewQuestion />
      <DripSettingsModal />
    </div>
  );
}
