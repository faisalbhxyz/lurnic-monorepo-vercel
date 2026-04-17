import CreateBanner from "@/components/shared/home/products/banner/CreateBanner";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import React from "react";

const getBannerByID = async (session: Session, id: string) => {
  try {
    const res = await axiosInstance.get(`/private/banner/${id}`, {
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

  const banner = await getBannerByID(session, id);

  return (
    <>
      <CreateBanner isEdit session={session} banner={banner} />
    </>
  );
}
