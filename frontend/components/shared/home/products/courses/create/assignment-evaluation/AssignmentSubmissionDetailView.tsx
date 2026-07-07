"use client";

import React, { useEffect, useState } from "react";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import { IoArrowBackOutline } from "react-icons/io5";
import { LuCloudDownload } from "react-icons/lu";
import axiosInstance from "@/lib/axiosInstance";
import axios from "axios";
import { toast } from "sonner";

type AssignmentSubmissionFile = {
  id: number;
  url: string;
  file_name: string;
  mime_type: string;
  size: number;
};

type AssignmentSubmissionDetail = {
  id: number;
  assignment_id: number;
  assignment_title: string;
  chapter_id: number;
  chapter_title: string;
  student_id: number;
  student_name: string;
  student_email: string;
  score: number;
  max_score: number;
  percentage: number;
  passed: boolean;
  status: "submitted" | "graded" | "pending_review";
  submitted_at: string;
  response_text?: string | null;
  instructor_feedback?: string | null;
  files: AssignmentSubmissionFile[];
};

function formatSubmittedAt(value: string) {
  return new Date(value).toLocaleString("en-US", {
    month: "long",
    day: "numeric",
    year: "numeric",
    hour: "numeric",
    minute: "2-digit",
    hour12: true,
  });
}

export default function AssignmentSubmissionDetailView({
  courseId,
  courseTitle,
  submissionId,
  accessToken,
  onBack,
  onGraded,
}: {
  courseId: number;
  courseTitle: string;
  submissionId: number;
  accessToken?: string;
  onBack: () => void;
  onGraded?: () => void;
}) {
  const [detail, setDetail] = useState<AssignmentSubmissionDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [score, setScore] = useState("");
  const [feedback, setFeedback] = useState("");
  const [grading, setGrading] = useState(false);

  useEffect(() => {
    if (!accessToken) return;

    setLoading(true);
    setError(null);
    setDetail(null);
    setScore("");
    setFeedback("");

    axiosInstance
      .get(`/private/course/${courseId}/assignment-submissions/${submissionId}`, {
        headers: { Authorization: `Bearer ${accessToken}` },
      })
      .then((res) => {
        const data = res.data.data as AssignmentSubmissionDetail;
        setDetail(data);
        if (data.status === "graded") {
          setScore(String(data.score));
          setFeedback(data.instructor_feedback ?? "");
        }
      })
      .catch((err) => {
        if (axios.isAxiosError(err)) {
          setError(
            (err.response?.data as { error?: string })?.error ??
              "Failed to load submission details."
          );
        } else {
          setError("Failed to load submission details.");
        }
      })
      .finally(() => setLoading(false));
  }, [courseId, submissionId, accessToken]);

  const handleGrade = async () => {
    if (!accessToken || !detail) return;

    const parsedScore = Number(score);
    if (Number.isNaN(parsedScore) || parsedScore < 0) {
      toast.error("Enter a valid score.");
      return;
    }
    if (parsedScore > detail.max_score) {
      toast.error(`Score cannot exceed ${detail.max_score}.`);
      return;
    }

    setGrading(true);
    try {
      const res = await axiosInstance.post(
        `/private/course/${courseId}/assignment-submissions/${submissionId}/grade`,
        {
          score: parsedScore,
          feedback: feedback.trim() || null,
        },
        {
          headers: { Authorization: `Bearer ${accessToken}` },
        }
      );
      setDetail(res.data.data ?? null);
      toast.success(res.data.message ?? "Assignment graded.");
      onGraded?.();
    } catch (err) {
      if (axios.isAxiosError(err)) {
        toast.error(
          (err.response?.data as { error?: string })?.error ??
            "Failed to grade assignment."
        );
      } else {
        toast.error("Failed to grade assignment.");
      }
    } finally {
      setGrading(false);
    }
  };

  if (loading) {
    return (
      <p className="text-sm text-gray-500 py-12 text-center">
        Loading submission...
      </p>
    );
  }

  if (error) {
    return (
      <div className="space-y-4">
        <button
          type="button"
          onClick={onBack}
          className="flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900"
        >
          <IoArrowBackOutline />
          Back to submissions
        </button>
        <p className="text-sm text-red-500 py-8 text-center">{error}</p>
      </div>
    );
  }

  if (!detail) return null;

  return (
    <div className="space-y-6">
      <button
        type="button"
        onClick={onBack}
        className="flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900"
      >
        <IoArrowBackOutline />
        Back to submissions
      </button>

      <div className="space-y-3">
        <h1 className="text-2xl font-semibold text-gray-900">
          {detail.assignment_title}
        </h1>
        <div className="flex flex-wrap gap-x-8 gap-y-2 text-sm text-gray-600">
          <p>
            <span className="font-medium text-gray-800">Course:</span>{" "}
            {courseTitle}
          </p>
          <p>
            <span className="font-medium text-gray-800">Student:</span>{" "}
            {detail.student_name}
          </p>
          <p>
            <span className="font-medium text-gray-800">Submitted Date:</span>{" "}
            {formatSubmittedAt(detail.submitted_at)}
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
        <div className="lg:col-span-2 space-y-6">
          <section className="space-y-3">
            <h2 className="text-lg font-semibold text-gray-900">Assignment</h2>
            {detail.response_text ? (
              <div
                className="text-sm text-gray-700 prose max-w-none bg-gray-50 border border-gray-200 rounded-lg p-4"
                dangerouslySetInnerHTML={{ __html: detail.response_text }}
              />
            ) : (
              <p className="text-sm text-gray-500 bg-gray-50 border border-gray-200 rounded-lg p-4">
                No written response provided.
              </p>
            )}
          </section>

          {detail.files.length > 0 && (
            <section className="space-y-3">
              <h2 className="text-lg font-semibold text-gray-900">
                Assignment File(s)
              </h2>
              <div className="space-y-3">
                {detail.files.map((file) => (
                  <a
                    key={file.id}
                    href={file.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center justify-between gap-4 border border-gray-200 rounded-lg bg-white px-4 py-3 hover:border-primary/40 hover:bg-gray-50 transition-colors"
                  >
                    <span className="text-sm font-medium text-gray-800 truncate">
                      {file.file_name}
                    </span>
                    <LuCloudDownload className="shrink-0 text-gray-500" size={20} />
                  </a>
                ))}
              </div>
            </section>
          )}
        </div>

        <aside className="lg:col-span-1 border border-gray-200 rounded-lg bg-white p-5 space-y-5">
          <h2 className="text-lg font-semibold text-gray-900">Evaluation</h2>

          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-800">
              Your Points
            </label>
            <div className="flex flex-col sm:flex-row sm:items-center gap-2">
              <InputField
                type="number"
                min={0}
                max={detail.max_score}
                step="0.01"
                value={score}
                onChange={(e) => setScore(e.target.value)}
                className="w-full sm:max-w-[120px]"
              />
              <span className="text-sm text-gray-500">
                Evaluate this assignment out of {detail.max_score}
              </span>
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-800">
              Write a feedback
            </label>
            <textarea
              value={feedback}
              onChange={(e) => setFeedback(e.target.value)}
              rows={6}
              className="w-full border border-gray-200 rounded-lg p-3 text-sm focus:outline-none focus:ring-2 focus:ring-primary/20"
              placeholder="Write feedback for the student..."
            />
          </div>

          <Button
            type="button"
            onClick={handleGrade}
            disabled={grading}
            className="w-full justify-center"
          >
            {grading
              ? "Saving..."
              : detail.status === "graded"
              ? "Update evaluation"
              : "Evaluate this submission"}
          </Button>
        </aside>
      </div>
    </div>
  );
}
