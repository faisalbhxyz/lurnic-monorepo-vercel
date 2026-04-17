import CreateBanner from "@/components/shared/home/products/banner/CreateBanner";
import { auth } from "@/lib/auth";
import React from "react";

export default async function page() {
  const session = await auth();
  if (!session) return null;

  return (
    <>
      <CreateBanner session={session} />
    </>
  );
}
