"use client";

import React from "react";

const sampleData = [
  {
    id: 1,
    notice: "Your review has been approved.",
    status: "Approved",
    date: "2025-06-08",
  },
  {
    id: 2,
    notice: "Pending review for your submission.",
    status: "Pending",
    date: "2025-06-07",
  },
  {
    id: 3,
    notice: "Review has been rejected due to policy violation.",
    status: "Rejected",
    date: "2025-06-06",
  },
];

export default function ReviewsTable() {
  return (
    <div className="overflow-x-auto">
      <table className="min-w-full">
        <thead className="bg-gray-100">
          <tr>
            <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
              Notice
            </th>
            <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
              Status
            </th>
            <th className="px-4 py-2 text-left text-sm font-semibold text-gray-700">
              Date
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200">
          {sampleData.map((review) => (
            <tr key={review.id}>
              <td className="px-4 py-2 text-sm text-gray-800">
                {review.notice}
              </td>
              <td className="px-4 py-2 text-sm text-gray-800">
                {review.status}
              </td>
              <td className="px-4 py-2 text-sm text-gray-800">{review.date}</td>
            </tr>
          ))}
          {sampleData.length === 0 && (
            <tr>
              <td
                colSpan={3}
                className="text-center py-4 text-sm text-gray-500"
              >
                No reviews found.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
