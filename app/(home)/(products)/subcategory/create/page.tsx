import { getAllCategories } from "@/app/actions/actions";
import CreateSubCategory from "@/components/shared/home/products/categories/CreateSubCategory";
import { auth } from "@/lib/auth";
import React from "react";

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const categories = await getAllCategories(session);

  return (
    <>
      <CreateSubCategory session={session} categories={categories}/>
    </>
  );
}
