"use client";

import Modal from "@/components/ui/Modal";
import { RxCross2 } from "react-icons/rx";
import {
  getCorrectAnswerLabel,
  getQuestionTypeBadge,
} from "./quizQuestionUtils";

interface QuizQuestion {
  title: string;
  _id: number;
  type: "multiple_choice" | "single_choice" | "true_false";
  marks: number;
  answer_required: boolean;
  details?: string | null | undefined;
  options?: { id: string; text: string }[] | null | undefined;
  correct_answer?:
    | { value?: string | boolean; values?: string[] }
    | null
    | undefined;
  answer_explanation?: string | null | undefined;
}

interface Props {
  question: QuizQuestion | null;
  onClose: () => void;
}

function isCorrectOption(
  question: QuizQuestion,
  optionId: string
): boolean {
  if (question.type === "single_choice") {
    return question.correct_answer?.value === optionId;
  }
  if (question.type === "multiple_choice") {
    return (question.correct_answer?.values ?? []).includes(optionId);
  }
  return false;
}

export default function QuizQuestionPreviewModal({ question, onClose }: Props) {
  if (!question) return null;

  const typeBadge = getQuestionTypeBadge(question.type);
  const correctLabel = getCorrectAnswerLabel(question);

  return (
    <Modal isOpen={!!question} onClose={onClose} className="max-w-2xl p-0">
      <div className="border-b border-gray-200 px-6 py-4 flex items-center justify-between">
        <p className="font-semibold text-lg">Question Preview</p>
        <button
          type="button"
          onClick={onClose}
          className="text-gray-500 hover:text-gray-700"
          aria-label="Close preview"
        >
          <RxCross2 size={20} />
        </button>
      </div>

      <div className="p-6 space-y-4">
        <div className="flex items-start justify-between gap-3">
          <h3 className="font-semibold text-gray-900 text-base leading-snug">
            {question.title}
          </h3>
          <span
            className={`shrink-0 text-xs font-semibold px-2.5 py-1 rounded-full border ${typeBadge.className}`}
          >
            {typeBadge.label}
          </span>
        </div>

        <div className="flex flex-wrap items-center gap-3 text-sm text-gray-600">
          <span>{question.marks} Marks</span>
          {question.answer_required && (
            <>
              <span className="text-gray-300">•</span>
              <span>Required answer</span>
            </>
          )}
        </div>

        {question.details && (
          <div
            className="text-sm text-gray-700 prose prose-sm max-w-none"
            dangerouslySetInnerHTML={{ __html: question.details }}
          />
        )}

        {question.type === "true_false" ? (
          <div className="space-y-2">
            {(["True", "False"] as const).map((label) => {
              const isCorrect =
                (label === "True" && question.correct_answer?.value === true) ||
                (label === "False" && question.correct_answer?.value === false);

              return (
                <div
                  key={label}
                  className={`rounded-lg border px-4 py-3 text-sm ${
                    isCorrect
                      ? "border-emerald-300 bg-emerald-50 text-emerald-800 font-medium"
                      : "border-gray-200 bg-gray-50 text-gray-700"
                  }`}
                >
                  {label}
                </div>
              );
            })}
          </div>
        ) : (
          <div className="space-y-2">
            {(question.options ?? []).map((option) => {
              const isCorrect = isCorrectOption(question, option.id);
              return (
                <div
                  key={option.id}
                  className={`rounded-lg border px-4 py-3 text-sm ${
                    isCorrect
                      ? "border-emerald-300 bg-emerald-50 text-emerald-800 font-medium"
                      : "border-gray-200 bg-gray-50 text-gray-700"
                  }`}
                >
                  {option.text || `Option ${option.id.toUpperCase()}`}
                </div>
              );
            })}
          </div>
        )}

        <div className="rounded-lg bg-gray-50 border border-gray-200 px-4 py-3 text-sm">
          <span className="font-medium text-gray-700">Correct answer: </span>
          <span className="text-gray-900">{correctLabel}</span>
        </div>

        {question.answer_explanation && (
          <div className="rounded-lg bg-blue-50 border border-blue-100 px-4 py-3 text-sm text-blue-900">
            <p className="font-medium mb-1">Explanation</p>
            <p>{question.answer_explanation}</p>
          </div>
        )}
      </div>
    </Modal>
  );
}
