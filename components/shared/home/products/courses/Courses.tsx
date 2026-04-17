"use client";

import React, { useState } from "react";
import Button from "@/components/ui/Button";
import { LuPlus } from "react-icons/lu";
import CoursesList from "./CoursesList";
import { FiSearch } from "react-icons/fi";
import Link from "next/link";

const statuses = ["All", "Public", "Private", "Protected"];

export default function Courses({ courses }: { courses: CourseDetails[] }) {
  const [activeTab, setActiveTab] = useState("All");
  const [search, setSearch] = useState("");

  const filteredData = courses.filter((course) => {
    const matchesStatus =
      activeTab === "All" ||
      course.visibility.toLowerCase() === activeTab.toLowerCase();

    const searchLower = search.toLowerCase();

    const matchesSearch =
      course.title.toLowerCase().includes(searchLower) ||
      (course.author.name &&
        course.author.name.toLowerCase().includes(searchLower));

    return matchesStatus && matchesSearch;
  });

  return (
    <>
      <div className="flex items-center text-sm gap-1 mb-4">
        <Link href={""} className="text-gray-500">
          Products
        </Link>
        /<Link href={""}>Courses</Link>
      </div>
      <div className="w-full flex items-center justify-between mb-4">
        <h3 className="font-medium text-2xl">Courses</h3>
        <Button link src="/courses/create">
          <LuPlus /> Add New
        </Button>
      </div>
      <div className="flex items-center justify-between mb-4">
        <div className="flex gap-3">
          {statuses.map((status) => {
            const count =
              status === "All"
                ? courses.length
                : filteredData.filter(
                    (c) => c.visibility === status.toLowerCase()
                  ).length;
            return (
              <button
                key={status}
                onClick={() => setActiveTab(status)}
                className={`px-3 py-3 border-b text-sm ${
                  activeTab === status
                    ? "border-primary text-primary"
                    : "border-transparent text-gray-700 hover:text-primary"
                }`}
              >
                {status}({count})
              </button>
            );
          })}
        </div>
        <div className="relative w-full max-w-60">
          <FiSearch className="absolute top-1/2 left-3 transform -translate-y-1/2 text-gray-500" />
          <input
            type="text"
            placeholder="Search..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full border text-sm px-3 py-1.5 rounded-md pl-8"
          />
        </div>
      </div>
      <CoursesList data={filteredData} />
    </>
  );
}
