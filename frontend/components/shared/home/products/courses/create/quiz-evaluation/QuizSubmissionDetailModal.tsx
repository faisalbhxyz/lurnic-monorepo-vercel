"use client";

import React, { useEffect, useState } from "react";
import Modal from "@/components/ui/Modal";
import { RxCross2 } from "react-icons/rx";
import axiosInstance from "@/lib/axiosInstance";
import axios from "axios";

type QuizSubmissionAnswer = {
  question_id: number;
  question_title: string;
  question_type: string;
  submitted_answer: unknown;
  is_correct: boolean | null;
  marks_awarded: number;
  answer_explanation?: string | null;
  correct_answer?: unknown;
};

type QuizSubmissionDetail = {
  id: number;
  quiz_id: number;
  quiz_title: string;
  chapter_id: number;
  chapter_title: string;
  student_id: number;
  student_name: string;
  student_email: string;
  attempt_number: number;
  score: number;
  max_score: number;
  percentage: number;
  passed: boolean;
  status: "submitted" | "graded" | "pending_review";
  submitted_at: string;
  reveal_answers: boolean;
  answers: QuizSubmissionAnswer[];
};

function formatAnswer(value: unknown): string {
  if (value === null || value === undefined) return "—";
  if (typeof value === "boolean") return value ? "True" : "False";
  if (Array.isArray(value)) return value.map(String).join(", ");
  if (typeof value === "object") {
    try {
      return JSON.stringify(value);
    } catch {
      return String(value);
    }
  }
  return String(value);
}

function formatCorrectAnswer(value: unknown): string {
  if (value === null || value === undefined) return "—";
  if (typeof value === "object" && value !== null) {
    const obj = value as { value?: unknown; values?: unknown[] };
    if (obj.values && Array.isArray(obj.values)) {
      return obj.values.map(String).join(", ");
    }
    if (obj.value !== undefined) return formatAnswer(obj.value);
  }
  return formatAnswer(value);
}

export default function QuizSubmissionDetailModal({
  courseId,
  submissionId,
  accessToken,
  onClose,
}: {
  courseId: number;
  submissionId: number | null;
  accessToken?: string;
  onClose: () => void;
}) {
  const [detail, setDetail] = useState<QuizSubmissionDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!submissionId || !accessToken) return;

    setLoading(true);
    setError(null);
    setDetail(null);

    axiosInstance
      .get(`/private/course/${courseId}/quiz-submissions/${submissionId}`, {
        headers: { Authorization: `Bearer ${accessToken}` },
      })
      .then((res) => setDetail(res.data.data ?? null))
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

  return (
    <Modal
      isOpen={submissionId != null}
      onClose={onClose}
      className="p-0 max-w-3xl"
    >
      <div className="p-4 flex items-center justify-between border-b border-gray-200">
        <p className="font-semibold text-lg">Quiz Submission</p>
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
                <p className="text-gray-500">Quiz</p>
                <p className="font-medium">{detail.quiz_title}</p>
                <p className="text-gray-600">{detail.chapter_title}</p>
              </div>
              <div>
                <p className="text-gray-500">Submitted</p>
                <p className="font-medium">
                  {new Date(detail.submitted_at).toLocaleString()}
                </p>
                <p className="text-gray-600">Attempt #{detail.attempt_number}</p>
              </div>
              <div>
                <p className="text-gray-500">Score</p>
                <p className="font-medium">
                  {detail.score}/{detail.max_score} ({detail.percentage}%)
                </p>
                <p className="text-gray-600">
                  {detail.passed ? "Passed" : "Not passed"}
                </p>
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

            <div className="border rounded-lg overflow-hidden">
              <div className="bg-gray-100 px-4 py-2 text-sm font-semibold">
                Answers
              </div>
              <ul className="divide-y">
                {detail.answers.map((answer, index) => (
                  <li key={answer.question_id} className="p-4 text-sm space-y-2">
                    <p className="font-medium">
                      {index + 1}. {answer.question_title}
                    </p>
                    <p className="text-gray-500 text-xs uppercase">
                      {answer.question_type.replace(/_/g, " ")}
                    </p>
                    <div className="grid sm:grid-cols-2 gap-2">
                      <div>
                        <p className="text-gray-500 text-xs">Submitted</p>
                        <p>{formatAnswer(answer.submitted_answer)}</p>
                      </div>
                      {detail.reveal_answers && (
                        <div>
                          <p className="text-gray-500 text-xs">Correct</p>
                          <p>{formatCorrectAnswer(answer.correct_answer)}</p>
                        </div>
                      )}
                    </div>
                    <div className="flex flex-wrap items-center gap-3 text-xs">
                      <span>
                        Marks: {answer.marks_awarded}
                      </span>
                      {answer.is_correct === true && (
                        <span className="text-green-700 font-medium">Correct</span>
                      )}
                      {answer.is_correct === false && (
                        <span className="text-red-600 font-medium">Incorrect</span>
                      )}
                      {answer.is_correct == null && (
                        <span className="text-yellow-700 font-medium">
                          Pending review
                        </span>
                      )}
                    </div>
                    {detail.reveal_answers && answer.answer_explanation && (
                      <div className="bg-gray-50 border rounded p-2 text-xs text-gray-700">
                        <p className="font-medium mb-1">Explanation</p>
                        <div
                          dangerouslySetInnerHTML={{
                            __html: answer.answer_explanation,
                          }}
                        />
                      </div>
                    )}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
}
