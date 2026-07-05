"use client";

import React, { useEffect, useState } from "react";
import axiosInstance from "@/lib/axiosInstance";
import { useSession } from "next-auth/react";
import axios from "axios";
import AssignmentSubmissionDetailModal from "./AssignmentSubmissionDetailModal";

type AssignmentSubmissionRow = {
  id: number;
  assignment_title: string;
  chapter_title: string;
  student_name: string;
  student_email: string;
  submitted_at: string;
  score: number;
  max_score: number;
  percentage: number;
  passed: boolean;
  status: "submitted" | "graded" | "pending_review";
  file_count: number;
};

export default function AssignmentTable({ courseId }: { courseId: number }) {
  const { data: session } = useSession();
  const [filter, setFilter] = useState("All");
  const [rows, setRows] = useState<AssignmentSubmissionRow[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedSubmissionId, setSelectedSubmissionId] = useState<
    number | null
  >(null);

  const loadRows = () => {
    if (!session?.accessToken || !courseId) return;

    setLoading(true);
    setError(null);
    axiosInstance
      .get(`/private/course/${courseId}/assignment-submissions`, {
        headers: {
          Authorization: `Bearer ${session.accessToken}`,
        },
      })
      .then((res) => {
        setRows(res.data.data ?? []);
      })
      .catch((err) => {
        if (axios.isAxiosError(err)) {
          setError(
            (err.response?.data as { error?: string })?.error ??
              "Failed to load assignment submissions."
          );
        } else {
          setError("Failed to load assignment submissions.");
        }
      })
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    loadRows();
  }, [session?.accessToken, courseId]);

  const filteredData =
    filter === "All"
      ? rows
      : rows.filter((d) => {
          if (filter === "Evaluate") return d.status === "pending_review";
          if (filter === "Pending") return d.status === "submitted";
          return true;
        });

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
      {loading ? (
        <p className="text-sm text-gray-500 py-6 text-center">
          Loading submissions...
        </p>
      ) : error ? (
        <p className="text-sm text-red-500 py-6 text-center">{error}</p>
      ) : (
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
                  Files
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
              {filteredData.map((row) => (
                <tr
                  key={row.id}
                  onClick={() => setSelectedSubmissionId(row.id)}
                  className="cursor-pointer hover:bg-gray-50 transition-colors"
                >
                  <td className="px-4 py-2 text-sm">{row.assignment_title}</td>
                  <td className="px-4 py-2 text-sm">{row.chapter_title}</td>
                  <td className="px-4 py-2 text-sm">
                    <div>{row.student_name}</div>
                    <div className="text-xs text-gray-500">
                      {row.student_email}
                    </div>
                  </td>
                  <td className="px-4 py-2 text-sm">
                    {new Date(row.submitted_at).toLocaleString()}
                  </td>
                  <td className="px-4 py-2 text-sm">{row.file_count}</td>
                  <td className="px-4 py-2 text-sm">
                    {row.status === "graded"
                      ? `${row.score}/${row.max_score} (${row.percentage}%)`
                      : "—"}
                  </td>
                  <td className="px-4 py-2 text-sm">
                    <span
                      className={`inline-block px-2 py-1 rounded text-xs font-medium ${
                        row.status === "pending_review"
                          ? "bg-yellow-200 text-yellow-800"
                          : row.status === "graded"
                          ? "bg-green-200 text-green-800"
                          : "bg-gray-200 text-gray-800"
                      }`}
                    >
                      {row.status}
                    </span>
                  </td>
                </tr>
              ))}
              {filteredData.length === 0 && (
                <tr>
                  <td
                    colSpan={7}
                    className="text-center py-4 text-sm text-gray-500"
                  >
                    No assignment submissions found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
          <p className="text-xs text-gray-500 mt-3">
            Click a row to review and grade the submission.
          </p>
        </div>
      )}

      <AssignmentSubmissionDetailModal
        courseId={courseId}
        submissionId={selectedSubmissionId}
        accessToken={session?.accessToken}
        onClose={() => setSelectedSubmissionId(null)}
        onGraded={loadRows}
      />
    </>
  );
}
