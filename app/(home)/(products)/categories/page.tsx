import { getAllCategories } from "@/app/actions/actions";
import CategoryList from "@/components/shared/home/products/categories/CategoryList";
import Button from "@/components/ui/Button";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import Link from "next/link";
import React from "react";
import { LuPlus } from "react-icons/lu";

export default async function page() {
  const session = await auth();
  if (!session) return null;
  
  const category = await getAllCategories(session);

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href={""} className="text-gray-500">
          Products
        </Link>
        /<Link href={""}>Course Categories</Link>
      </div>
      <div className="flex-between mt-5">
        <h3 className="font-medium">Course Categories</h3>
        <Button link src="/categories/create">
          <LuPlus /> Add New
        </Button>
      </div>

      <div className="mt-5">
        {category && category.length > 0 && <CategoryList data={category} />}
      </div>
    </>
  );
}
