import BannerList from "@/components/shared/home/products/banner/BannerList";
import CategoryList from "@/components/shared/home/products/categories/CategoryList";
import Button from "@/components/ui/Button";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import Link from "next/link";
import React from "react";
import { LuPlus } from "react-icons/lu";

export const getAllBanners = async (
  session: Session
): Promise<IBanner[] | null> => {
  try {
    const res = await axiosInstance.get("/private/banner", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    console.log("[ERROR]", error);

    return null;
  }
};

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const banners = await getAllBanners(session);

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href={""} className="text-gray-500">
          Products
        </Link>
        /<Link href={""}>Course Banners</Link>
      </div>
      <div className="flex-between mt-5">
        <h3 className="font-medium">Course Banners</h3>
        <Button link src="/banner/create">
          <LuPlus /> Add New
        </Button>
      </div>

      <div className="mt-5">
        {banners && banners.length > 0 && <BannerList data={banners} />}
      </div>
    </>
  );
}
