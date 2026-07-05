"use client";

import React, { useEffect, useState } from "react";
import Modal from "@/components/ui/Modal";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import { RxCross2 } from "react-icons/rx";
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

export default function AssignmentSubmissionDetailModal({
  courseId,
  submissionId,
  accessToken,
  onClose,
  onGraded,
}: {
  courseId: number;
  submissionId: number | null;
  accessToken?: string;
  onClose: () => void;
  onGraded?: () => void;
}) {
  const [detail, setDetail] = useState<AssignmentSubmissionDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [score, setScore] = useState("");
  const [feedback, setFeedback] = useState("");
  const [grading, setGrading] = useState(false);

  useEffect(() => {
    if (!submissionId || !accessToken) return;

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
    if (!submissionId || !accessToken || !detail) return;

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

  return (
    <Modal
      isOpen={submissionId != null}
      onClose={onClose}
      className="p-0 max-w-3xl"
    >
      <div className="p-4 flex items-center justify-between border-b border-gray-200">
        <p className="font-semibold text-lg">Assignment Submission</p>
        <button type="button" onClick={onClose} aria-label="Close">
          <RxCross2 size={20} />
        </button>
      </div>

      <div className="p-4 max-h-[70vh] overflow-y-auto">
        {loading && (
          <p className="text-sm text-gray-500 py-8 text-center">Loading...</p>
        )}
        {error && (
          <p className="text-sm text-red-500 py-8 text-center">{error}</p>
        )}
        {detail && !loading && (
          <div className="space-y-5">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
              <div>
                <p className="text-gray-500">Student</p>
                <p className="font-medium">{detail.student_name}</p>
                <p className="text-gray-600">{detail.student_email}</p>
              </div>
              <div>
                <p className="text-gray-500">Assignment</p>
                <p className="font-medium">{detail.assignment_title}</p>
                <p className="text-gray-600">{detail.chapter_title}</p>
              </div>
              <div>
                <p className="text-gray-500">Submitted</p>
                <p className="font-medium">
                  {new Date(detail.submitted_at).toLocaleString()}
                </p>
              </div>
              <div>
                <p className="text-gray-500">Max Marks</p>
                <p className="font-medium">{detail.max_score}</p>
              </div>
            </div>

            <div>
              <span
                className={`inline-block px-2 py-1 rounded text-xs font-medium ${
                  detail.status === "pending_review"
                    ? "bg-yellow-200 text-yellow-800"
                    : detail.status === "graded"
                    ? "bg-green-200 text-green-800"
                    : "bg-gray-200 text-gray-800"
                }`}
              >
                {detail.status}
              </span>
            </div>

            {detail.response_text && (
              <div className="border rounded-lg overflow-hidden">
                <div className="bg-gray-100 px-4 py-2 text-sm font-semibold">
                  Student Response
                </div>
                <div
                  className="p-4 text-sm prose max-w-none"
                  dangerouslySetInnerHTML={{ __html: detail.response_text }}
                />
              </div>
            )}

            {detail.files.length > 0 && (
              <div className="border rounded-lg overflow-hidden">
                <div className="bg-gray-100 px-4 py-2 text-sm font-semibold">
                  Uploaded Files
                </div>
                <ul className="divide-y">
                  {detail.files.map((file) => (
                    <li
                      key={file.id}
                      className="p-4 text-sm flex items-center justify-between gap-3"
                    >
                      <div>
                        <p className="font-medium">{file.file_name}</p>
                        <p className="text-xs text-gray-500">{file.mime_type}</p>
                      </div>
                      <a
                        href={file.url}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-primary text-sm font-medium hover:underline"
                      >
                        Download
                      </a>
                    </li>
                  ))}
                </ul>
              </div>
            )}

            <div className="border rounded-lg p-4 space-y-3">
              <p className="text-sm font-semibold">Grade Submission</p>
              <div>
                <label className="text-sm font-medium block mb-1">
                  Score (out of {detail.max_score})
                </label>
                <InputField
                  type="number"
                  min={0}
                  max={detail.max_score}
                  step="0.01"
                  value={score}
                  onChange={(e) => setScore(e.target.value)}
                  className="w-full"
                />
              </div>
              <div>
                <label className="text-sm font-medium block mb-1">
                  Feedback (optional)
                </label>
                <textarea
                  value={feedback}
                  onChange={(e) => setFeedback(e.target.value)}
                  rows={4}
                  className="w-full border rounded-md p-2 text-sm"
                  placeholder="Write feedback for the student..."
                />
              </div>
              <Button
                type="button"
                onClick={handleGrade}
                disabled={grading}
                className="w-full sm:w-auto"
              >
                {grading
                  ? "Saving..."
                  : detail.status === "graded"
                  ? "Update Grade"
                  : "Submit Grade"}
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
}
