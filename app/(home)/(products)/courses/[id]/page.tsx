import CoursesDetails from "@/components/shared/home/products/courses/CoursesDetails";
import React from "react";
import { LuHeart } from "react-icons/lu";
import { IoMdShareAlt } from "react-icons/io";

export default function page() {
  return (
    <div className="p-5 max-w-6xl mx-auto">
      <p className="text-lg font-semibold mb-3">New Course</p>
      <div className="flex items-center justify-between text-gray-500">
        <p className="text-sm">Uncategorized</p>
        <div className="flex items-center gap-3 text-sm">
          <button className="flex items-center gap-1">
            <LuHeart /> Wishlist
          </button>
          <button className="flex items-center gap-1">
            <IoMdShareAlt /> Share
          </button>
        </div>
      </div>
      <CoursesDetails />
    </div>
  );
}
