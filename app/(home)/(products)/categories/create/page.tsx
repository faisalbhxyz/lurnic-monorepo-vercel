import CreateCategory from "@/components/shared/home/products/categories/CreateCategory";
import Button from "@/components/ui/Button";
import { auth } from "@/lib/auth";
import React from "react";

export default async function page() {
  const session = await auth();
  if (!session) return null;

  return (
    <>
      <CreateCategory session={session} />
    </>
  );
}
