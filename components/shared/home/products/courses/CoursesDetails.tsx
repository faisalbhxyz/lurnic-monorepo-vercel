"use client";

import Button from "@/components/ui/Button";
import Image from "next/image";
import React, { useState } from "react";
import CourseInfo from "./CourseInfo";
import Reviews from "./Reviews";

export default function CoursesDetails() {
  const [activeTab, setActiveTab] = useState<"course" | "reviews">("course");

  return (
    <div className="flex items-start mt-5 gap-10">
      <div className="w-full">
        <Image
          src={"/images/placeholder.svg"}
          alt={"image"}
          width={500}
          height={400}
          className="w-full h-auto rounded-xl"
        />
        <div className="bg-white border-b border-gray-300 mt-5 flex">
          <button
            className={`px-4 py-2 border-b-2 ${
              activeTab === "course"
                ? "border-primary text-primary"
                : "border-transparent"
            }`}
            onClick={() => setActiveTab("course")}
          >
            Course Info
          </button>
          <button
            className={`px-4 py-2 border-b-2 ${
              activeTab === "reviews"
                ? "border-primary text-primary"
                : "border-transparent"
            }`}
            onClick={() => setActiveTab("reviews")}
          >
            Reviews
          </button>
        </div>
        <div className="mt-4">
          {activeTab === "course" && <CourseInfo />}
          {activeTab === "reviews" && <Reviews />}
        </div>
      </div>
      <div className="w-96 min-w-96">
        <div className="border rounded-lg">
          <div className="p-10 rounded-md bg-gray-100 rounded-t-lg">
            <p className="font-bold text-2xl mb-5">Free</p>
            <Button className="w-full flex justify-center">Enroll Now</Button>
            <p className="mt-3 text-xs font-medium text-center text-gray-500">
              Free access this course
            </p>
          </div>
          <div className="border-t bg-white rounded-b-lg border-gray-300 text-gray-700 px-10 py-8 space-y-2">
            <div>
              <p>Intermediate</p>
            </div>
            <div>
              <p>0 Total Enroll</p>
            </div>
            <div>
              <p>0 April 21, 2025 Last Updated</p>
            </div>
          </div>
        </div>
        <div className="border mt-5 rounded-lg bg-white text-gray-700 px-10 py-8 space-y-2">
          <p>A course by</p>
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 flex items-center justify-center rounded-full text-white bg-blue-700">
              A
            </div>
            <p className="font-medium">Amerrajjo</p>
          </div>
        </div>
      </div>
    </div>
  );
}
