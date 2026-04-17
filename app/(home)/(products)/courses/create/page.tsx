import {
  getAllCategories,
  getAllInstructorsLite,
  getAllSubCategories,
} from "@/app/actions/actions";
import BackButton from "@/components/shared/BackButton";
import CoursesTabs from "@/components/shared/home/products/courses/create/CoursesTabs";
import DripSettingsModal from "@/components/shared/home/products/courses/create/curriculum/DripSettingsModal";
import EditItemModal from "@/components/shared/home/products/courses/create/curriculum/EditItemModal";
import { auth } from "@/lib/auth";
import React from "react";

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const categories = await getAllCategories(session);
  const subcategories = await getAllSubCategories(session);
  const instructors = await getAllInstructorsLite(session);

  return (
    <div className="w-full max-w-6xl mx-auto py-5">
      <BackButton buttonText={"Courses"} />
      <CoursesTabs
        categories={categories}
        subcategories={subcategories}
        instructors={instructors}
      />
      <EditItemModal />
      <DripSettingsModal />
    </div>
  );
}
