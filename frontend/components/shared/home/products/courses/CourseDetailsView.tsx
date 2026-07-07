"use client";

import React, { Suspense, useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { BiEditAlt } from "react-icons/bi";
import QuizTable from "./create/quiz-evaluation/QuizTable";
import AssignmentTable from "./create/assignment-evaluation/AssignmentTable";

const tabs = [
  { id: 1, label: "Quizzes" },
  { id: 2, label: "Assignments" },
];

export default function CourseDetailsView({
  course,
}: {
  course: Pick<CourseDetails, "id" | "title" | "featured_image">;
}) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const defaultTab = tabs[0].label;
  const tabFromQuery = searchParams.get("tab");
  const isValidTab = tabs.some((tab) => tab.label === tabFromQuery);

  const [activeTab, setActiveTab] = useState(
    isValidTab ? tabFromQuery! : defaultTab
  );

  const handleTabChange = (tabLabel: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("tab", tabLabel);
    params.delete("submission");
    router.replace(`?${params.toString()}`);
    setActiveTab(tabLabel);
  };

  useEffect(() => {
    if (tabFromQuery && tabFromQuery !== activeTab && isValidTab) {
      setActiveTab(tabFromQuery);
    }
  }, [tabFromQuery, activeTab, isValidTab]);

  const renderContent = () => {
    switch (activeTab) {
      case "Quizzes":
        return <QuizTable courseId={course.id} />;
      case "Assignments":
        return (
          <AssignmentTable courseId={course.id} courseTitle={course.title} />
        );
      default:
        return null;
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between my-5">
        <div className="flex items-center gap-5">
          <Image
            src={course.featured_image || "/images/placeholder.svg"}
            alt={course.title}
            width={130}
            height={130}
            className="w-32 h-20 object-contain rounded-md"
          />
          <div>
            <p className="font-semibold text-lg">{course.title}</p>
            <p className="text-sm text-gray-500 mt-1">Quiz & Assignment submissions</p>
          </div>
        </div>
        <Link
          href={`/courses/${course.id}/update`}
          className="flex items-center gap-2 border px-4 py-2 rounded-full text-sm text-gray-600 font-medium hover:bg-gray-50"
        >
          <BiEditAlt size={18} />
          Edit Course
        </Link>
      </div>

      <div className="border-b border-gray-300 flex space-x-4 mb-4">
        {tabs.map((tab) => (
          <button
            type="button"
            key={tab.id}
            onClick={() => handleTabChange(tab.label)}
            className={`px-3 py-2 border-b-2 font-medium text-sm transition-all focus:outline-none ${
              activeTab === tab.label
                ? "border-primary text-primary"
                : "border-transparent text-gray-500 hover:text-gray-700"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      <div className="border p-5 rounded-lg">
        <Suspense fallback={<p className="text-sm text-gray-500">Loading...</p>}>
          {renderContent()}
        </Suspense>
      </div>
    </div>
  );
}
