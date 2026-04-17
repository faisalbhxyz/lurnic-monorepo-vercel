import CreateCategory from "@/components/shared/home/products/categories/CreateCategory";
import Button from "@/components/ui/Button";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import React from "react";

const getCategoryByID = async (session: Session, id: string) => {
  try {
    const res = await axiosInstance.get(`/private/category/${id}`, {
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

export default async function page({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const session = await auth();
  if (!session) return null;
  const { id } = await params;

  const category = await getCategoryByID(session, id);

  return (
    <>
      <CreateCategory isEdit session={session} category={category} />
    </>
  );
}
