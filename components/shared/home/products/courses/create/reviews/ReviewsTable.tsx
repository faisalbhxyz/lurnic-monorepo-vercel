"use client";

import React, { useState } from "react";

const sampleData = [
  {
    id: 1,
    review: "Great course! Learned a lot.",
    rating: 5,
    student: "Alice Johnson",
    status: "Evaluate",
    featured: true,
  },
  {
    id: 2,
    review: "Too fast-paced for beginners.",
    rating: 3,
    student: "Bob Smith",
    status: "Pending",
    featured: false,
  },
  {
    id: 3,
    review: "Well-structured and informative.",
    rating: 4,
    student: "Carol Davis",
    status: "Evaluate",
    featured: true,
  },
];

export default function ReviewsTable() {
  const [filter, setFilter] = useState("All");

  const filteredData =
    filter === "All"
      ? sampleData
      : sampleData.filter((d) => d.status === filter);

  return (
    <>
      <div className="mb-4 flex space-x-1">
        {["All", "Approved", "Pending"].map((status) => (
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
                Review
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Rating
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Student
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Status
              </th>
              <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
                Featured
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {filteredData.map((review) => (
              <tr key={review.id}>
                <td className="px-4 py-2 text-sm text-gray-700">
                  {review.review}
                </td>
                <td className="px-4 py-2 text-sm text-gray-700">
                  {review.rating}
                </td>
                <td className="px-4 py-2 text-sm text-gray-700">
                  {review.student}
                </td>
                <td className="px-4 py-2 text-sm text-gray-700">
                  {review.status}
                </td>
                <td className="px-4 py-2 text-sm text-gray-700">
                  {review.featured ? "Yes" : "No"}
                </td>
              </tr>
            ))}
            {filteredData.length === 0 && (
              <tr>
                <td
                  colSpan={5}
                  className="text-center py-4 text-sm text-gray-500"
                >
                  No reviews found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </>
  );
}
