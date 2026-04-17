"use client";

import React, { useState } from "react";

const sampleData = [
  {
    id: 1,
    quizName: "Algebra Basics",
    chapterName: "Chapter 1",
    submittedBy: "Alice",
    submittedAt: "2025-06-01",
    marks: 8,
    status: "Evaluate",
  },
  {
    id: 2,
    quizName: "Geometry Intro",
    chapterName: "Chapter 2",
    submittedBy: "Bob",
    submittedAt: "2025-06-02",
    marks: 10,
    status: "Pending",
  },
  {
    id: 3,
    quizName: "Advanced Algebra",
    chapterName: "Chapter 3",
    submittedBy: "Charlie",
    submittedAt: "2025-06-03",
    marks: 7,
    status: "Evaluate",
  },
];

export default function AssignmentTable() {
  const [filter, setFilter] = useState("All");

  const filteredData =
    filter === "All"
      ? sampleData
      : sampleData.filter((d) => d.status === filter);

  return (
    <>
      <div className="mb-4 flex space-x-1">
        {["All", "Evaluate", "Pending"].map((status) => (
          <button
            key={status}
            onClick={() => setFilter(status)}
            className={`px-4 py-2 border-b ${
              filter === status
                ? "border-primary text-primary"
                : "border-transparent text-gray-600 hover:text-primary"
            }`}
          >
            {status}
          </button>
        ))}
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full">
          <thead className="bg-gray-100">
            <tr>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Assignment Name
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Chapter Name
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Submitted By
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Submitted At
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Marks
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Status
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {filteredData.map((quiz) => (
              <tr key={quiz.id}>
                <td className="px-4 py-2 text-sm">{quiz.quizName}</td>
                <td className="px-4 py-2 text-sm">{quiz.chapterName}</td>
                <td className="px-4 py-2 text-sm">{quiz.submittedBy}</td>
                <td className="px-4 py-2 text-sm">{quiz.submittedAt}</td>
                <td className="px-4 py-2 text-sm">{quiz.marks}</td>
                <td className="px-4 py-2 text-sm">
                  <span
                    className={`inline-block px-2 py-1 rounded text-xs font-medium ${
                      quiz.status === "Evaluate"
                        ? "bg-yellow-200 text-yellow-800"
                        : "bg-green-200 text-green-800"
                    }`}
                  >
                    {quiz.status}
                  </span>
                </td>
              </tr>
            ))}
            {filteredData.length === 0 && (
              <tr>
                <td
                  colSpan={6}
                  className="text-center py-4 text-sm text-gray-500"
                >
                  No assignment found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </>
  );
}
