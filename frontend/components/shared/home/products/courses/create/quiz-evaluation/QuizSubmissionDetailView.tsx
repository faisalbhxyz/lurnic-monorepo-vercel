"use client";

import React, { useEffect, useState } from "react";
import { IoArrowBackOutline } from "react-icons/io5";
import { HiOutlineCheckCircle, HiOutlineXCircle } from "react-icons/hi";
import { LuCircle, LuListChecks } from "react-icons/lu";
import { TbToggleLeft } from "react-icons/tb";
import Button from "@/components/ui/Button";
import TextEditor from "@/components/shared/text-editor/TextEditor";
import axiosInstance from "@/lib/axiosInstance";
import axios from "axios";
import { toast } from "sonner";

type QuizOption = { id: string; text: string };

type QuizSubmissionAnswer = {
  question_id: number;
  question_title: string;
  question_type: "multiple_choice" | "single_choice" | "true_false";
  options?: QuizOption[] | null;
  submitted_answer?: unknown;
  is_correct: boolean | null;
  marks_awarded: number;
  question_marks: number;
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
  total_questions: number;
  correct_count: number;
  incorrect_count: number;
  unanswered_count: number;
  pass_marks: number;
  minimum_pass_percentage: number;
  quiz_time_seconds: number;
  attempt_time_seconds?: number | null;
  attempt_started_at?: string;
  instructor_feedback?: string | null;
  answers: QuizSubmissionAnswer[];
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

function formatDuration(totalSeconds: number) {
  if (totalSeconds <= 0) return "—";
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}m ${seconds}s`;
}

function optionLabel(
  options: QuizOption[] | null | undefined,
  value: unknown
): string {
  if (value === null || value === undefined) return "—";
  if (typeof value === "boolean") return value ? "True" : "False";
  const id = String(value);
  const match = options?.find((opt) => opt.id === id);
  return match?.text?.trim() ? match.text : id;
}

function extractAnswerValues(answer: unknown): string[] {
  if (answer === null || answer === undefined) return [];
  if (typeof answer === "boolean") return [answer ? "True" : "False"];
  if (typeof answer === "object" && answer !== null) {
    const obj = answer as { value?: unknown; values?: unknown[] };
    if (Array.isArray(obj.values)) {
      return obj.values.map((item) => String(item));
    }
    if (obj.value !== undefined) return [String(obj.value)];
  }
  if (Array.isArray(answer)) return answer.map((item) => String(item));
  return [String(answer)];
}

function formatAnswerDisplay(
  answer: unknown,
  options?: QuizOption[] | null,
  questionType?: QuizSubmissionAnswer["question_type"]
): string[] {
  const values = extractAnswerValues(answer);
  if (values.length === 0) return ["—"];

  if (questionType === "true_false") {
    const raw = values[0];
    if (raw === "true" || raw === "True") return ["True"];
    if (raw === "false" || raw === "False") return ["False"];
  }

  return values.map((value) => optionLabel(options, value));
}

function QuestionTypeIcon({
  type,
}: {
  type: QuizSubmissionAnswer["question_type"];
}) {
  if (type === "multiple_choice") {
    return <LuListChecks className="text-gray-500" size={18} />;
  }
  if (type === "true_false") {
    return <TbToggleLeft className="text-gray-500" size={18} />;
  }
  return <LuCircle className="text-gray-500" size={18} />;
}

function ResultBadge({ isCorrect }: { isCorrect: boolean | null }) {
  if (isCorrect === true) {
    return (
      <span className="inline-flex items-center rounded-full bg-green-100 px-3 py-1 text-xs font-medium text-green-800">
        Correct
      </span>
    );
  }
  if (isCorrect === false) {
    return (
      <span className="inline-flex items-center rounded-full bg-red-100 px-3 py-1 text-xs font-medium text-red-700">
        Incorrect
      </span>
    );
  }
  return (
    <span className="inline-flex items-center rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-800">
      Pending
    </span>
  );
}

function ResultSummaryBadge({
  status,
  passed,
}: {
  status: QuizSubmissionDetail["status"];
  passed: boolean;
}) {
  if (status === "pending_review") {
    return (
      <span className="inline-flex items-center rounded-full bg-amber-100 px-3 py-1 text-xs font-medium text-amber-800">
        Pending
      </span>
    );
  }
  if (passed) {
    return (
      <span className="inline-flex items-center rounded-full bg-green-100 px-3 py-1 text-xs font-medium text-green-800">
        Passed
      </span>
    );
  }
  return (
    <span className="inline-flex items-center rounded-full bg-red-100 px-3 py-1 text-xs font-medium text-red-700">
      Failed
    </span>
  );
}

function AnswerLines({ lines }: { lines: string[] }) {
  return (
    <div className="space-y-1">
      {lines.map((line, index) => (
        <p key={`${line}-${index}`} className="text-sm text-gray-800">
          {line}
        </p>
      ))}
    </div>
  );
}

export default function QuizSubmissionDetailView({
  courseId,
  submissionId,
  accessToken,
  onBack,
}: {
  courseId: number;
  submissionId: number;
  accessToken?: string;
  onBack: () => void;
}) {
  const [detail, setDetail] = useState<QuizSubmissionDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [feedback, setFeedback] = useState("");
  const [savingFeedback, setSavingFeedback] = useState(false);

  useEffect(() => {
    if (!accessToken) return;

    setLoading(true);
    setError(null);
    setDetail(null);
    setFeedback("");

    axiosInstance
      .get(`/private/course/${courseId}/quiz-submissions/${submissionId}`, {
        headers: { Authorization: `Bearer ${accessToken}` },
      })
      .then((res) => {
        const data = res.data.data as QuizSubmissionDetail;
        setDetail(data);
        setFeedback(data.instructor_feedback ?? "");
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

  const handleUpdateFeedback = async () => {
    if (!accessToken || !detail) return;

    setSavingFeedback(true);
    try {
      const res = await axiosInstance.post(
        `/private/course/${courseId}/quiz-submissions/${submissionId}/feedback`,
        {
          feedback: feedback.trim() || null,
        },
        {
          headers: { Authorization: `Bearer ${accessToken}` },
        }
      );
      const data = res.data.data as QuizSubmissionDetail;
      setDetail(data);
      setFeedback(data.instructor_feedback ?? "");
      toast.success(res.data.message ?? "Instructor feedback updated.");
    } catch (err) {
      if (axios.isAxiosError(err)) {
        toast.error(
          (err.response?.data as { error?: string })?.error ??
            "Failed to update instructor feedback."
        );
      } else {
        toast.error("Failed to update instructor feedback.");
      }
    } finally {
      setSavingFeedback(false);
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

  const quizTimeLabel =
    detail.quiz_time_seconds > 0
      ? formatDuration(detail.quiz_time_seconds)
      : "Unlimited";

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

      <div className="space-y-2">
        <h1 className="text-2xl font-semibold text-gray-900">{detail.quiz_title}</h1>
        <p className="text-sm text-gray-600">{detail.chapter_title}</p>
      </div>

      <div className="overflow-x-auto border border-gray-200 rounded-lg">
        <table className="min-w-full">
          <thead className="bg-gray-100">
            <tr>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Attempt By
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Date
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Question
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Quiz Time
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Attempt Time
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Total Marks
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Pass Marks
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Correct Answer
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Incorrect Answer
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Earned Marks
              </th>
              <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                Result
              </th>
            </tr>
          </thead>
          <tbody>
            <tr className="divide-y divide-gray-200">
              <td className="px-4 py-3 text-sm">
                <div className="font-medium text-gray-900">{detail.student_name}</div>
                <div className="text-xs text-gray-500">{detail.student_email}</div>
              </td>
              <td className="px-4 py-3 text-sm text-gray-800">
                {formatSubmittedAt(detail.submitted_at)}
              </td>
              <td className="px-4 py-3 text-sm text-gray-800">
                {detail.total_questions}
              </td>
              <td className="px-4 py-3 text-sm text-gray-800">{quizTimeLabel}</td>
              <td className="px-4 py-3 text-sm text-gray-800">
                {detail.attempt_time_seconds != null
                  ? formatDuration(detail.attempt_time_seconds)
                  : "—"}
              </td>
              <td className="px-4 py-3 text-sm text-gray-800">
                {detail.max_score.toFixed(2)}
              </td>
              <td className="px-4 py-3 text-sm text-gray-800">
                {detail.pass_marks.toFixed(2)} ({detail.minimum_pass_percentage}%)
              </td>
              <td className="px-4 py-3 text-sm text-gray-800">
                {detail.correct_count}
              </td>
              <td className="px-4 py-3 text-sm text-gray-800">
                {detail.incorrect_count}
              </td>
              <td className="px-4 py-3 text-sm text-gray-800">
                {detail.score.toFixed(2)} ({detail.percentage}%)
              </td>
              <td className="px-4 py-3 text-sm">
                <ResultSummaryBadge status={detail.status} passed={detail.passed} />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <section className="space-y-4">
        <h2 className="text-lg font-semibold text-gray-900">Quiz Overview</h2>

        <div className="overflow-x-auto border border-gray-200 rounded-lg">
          <table className="min-w-full">
            <thead className="bg-gray-100">
              <tr>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 w-12">
                  No
                </th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 w-14">
                  Type
                </th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 min-w-[220px]">
                  Questions
                </th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 min-w-[160px]">
                  Given Answer
                </th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 min-w-[160px]">
                  Correct Answer
                </th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                  Result
                </th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">
                  Review
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {detail.answers.map((answer, index) => {
                const givenLines = formatAnswerDisplay(
                  answer.submitted_answer,
                  answer.options,
                  answer.question_type
                );
                const correctLines = detail.reveal_answers
                  ? formatAnswerDisplay(
                      answer.correct_answer,
                      answer.options,
                      answer.question_type
                    )
                  : ["Hidden"];

                return (
                  <tr key={answer.question_id} className="align-top">
                    <td className="px-4 py-4 text-sm text-gray-800">{index + 1}</td>
                    <td className="px-4 py-4">
                      <QuestionTypeIcon type={answer.question_type} />
                    </td>
                    <td className="px-4 py-4 text-sm text-gray-800">
                      <p className="font-medium">{answer.question_title}</p>
                      <p className="text-xs text-gray-500 mt-1">
                        {answer.question_type.replace(/_/g, " ")} ·{" "}
                        {answer.marks_awarded}/{answer.question_marks} marks
                      </p>
                    </td>
                    <td className="px-4 py-4">
                      <AnswerLines lines={givenLines} />
                    </td>
                    <td className="px-4 py-4">
                      <AnswerLines lines={correctLines} />
                    </td>
                    <td className="px-4 py-4">
                      <ResultBadge isCorrect={answer.is_correct} />
                    </td>
                    <td className="px-4 py-4">
                      <div className="flex items-center gap-2">
                        {answer.is_correct === true ? (
                          <HiOutlineCheckCircle
                            className="text-green-600"
                            size={22}
                            aria-label="Correct"
                          />
                        ) : answer.is_correct === false ? (
                          <HiOutlineXCircle
                            className="text-red-500"
                            size={22}
                            aria-label="Incorrect"
                          />
                        ) : (
                          <span className="text-xs text-amber-700">Pending</span>
                        )}
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>

        {detail.answers.some((answer) => answer.answer_explanation) && (
          <div className="space-y-3">
            <h3 className="text-sm font-semibold text-gray-900">Explanations</h3>
            {detail.answers.map((answer, index) =>
              answer.answer_explanation ? (
                <div
                  key={`explanation-${answer.question_id}`}
                  className="rounded-lg border border-gray-200 bg-gray-50 p-4"
                >
                  <p className="text-sm font-medium text-gray-900 mb-2">
                    {index + 1}. {answer.question_title}
                  </p>
                  <div
                    className="text-sm text-gray-700 prose max-w-none"
                    dangerouslySetInnerHTML={{
                      __html: answer.answer_explanation,
                    }}
                  />
                </div>
              ) : null
            )}
          </div>
        )}
      </section>

      <section className="space-y-3">
        <h2 className="text-sm font-medium text-gray-800">Instructor Feedback</h2>
        <TextEditor value={feedback} onChange={setFeedback} />
        <Button
          type="button"
          onClick={handleUpdateFeedback}
          disabled={savingFeedback}
          className="justify-center"
        >
          {savingFeedback ? "Updating..." : "Update"}
        </Button>
      </section>
    </div>
  );
}
