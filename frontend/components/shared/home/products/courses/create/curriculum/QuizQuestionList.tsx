"use client";

import React, { useState } from "react";
import {
  closestCenter,
  DndContext,
  DragEndEvent,
  PointerSensor,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import {
  CheckSquare,
  Copy,
  Eye,
  GripVertical,
  Pencil,
  Trash2,
} from "lucide-react";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import QuizQuestionPreviewModal from "./QuizQuestionPreviewModal";
import {
  getCorrectAnswerLabel,
  getOptionsCount,
  getQuestionTypeBadge,
} from "./quizQuestionUtils";

interface QuizQuestion {
  title: string;
  _id: number;
  type: "multiple_choice" | "single_choice" | "true_false";
  marks: number;
  answer_required: boolean;
  details?: string | null | undefined;
  id?: number | null | undefined;
  media?: unknown[] | null | undefined;
  options?: { id: string; text: string }[] | null | undefined;
  correct_answer?:
    | { value?: string | boolean; values?: string[] }
    | null
    | undefined;
  answer_explanation?: string | null | undefined;
}

function SortableQuestionCard({
  question,
  index,
  onPreview,
}: {
  question: QuizQuestion;
  index: number;
  onPreview: (question: QuizQuestion) => void;
}) {
  const { openEditQuestion, removeQuestion, duplicateQuestion } =
    useCoursesStore();

  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: question._id });

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  const typeBadge = getQuestionTypeBadge(question.type);
  const optionsCount = getOptionsCount(question);
  const correctAnswer = getCorrectAnswerLabel(question);

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm"
    >
      <div className="flex items-start gap-3">
        <button
          type="button"
          className="mt-1 shrink-0 cursor-grab text-gray-400 hover:text-gray-600 active:cursor-grabbing"
          aria-label="Drag to reorder"
          {...attributes}
          {...listeners}
        >
          <GripVertical size={18} />
        </button>

        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-blue-50 text-sm font-bold text-blue-600">
          {String(index + 1).padStart(2, "0")}
        </div>

        <div className="min-w-0 flex-1">
          <div className="flex items-start justify-between gap-3">
            <h4 className="font-semibold text-gray-900 leading-snug">
              {question.title}
            </h4>
            <span
              className={`shrink-0 rounded-full border px-2.5 py-1 text-xs font-semibold ${typeBadge.className}`}
            >
              {typeBadge.label}
            </span>
          </div>

          <p className="mt-1 text-sm text-gray-500">
            Options: {optionsCount} • Correct Answer: {correctAnswer}
          </p>

          <div className="mt-2 flex flex-wrap items-center gap-4 text-sm text-gray-600">
            <span className="inline-flex items-center gap-1.5">
              <CheckSquare size={15} className="text-gray-500" />
              {question.marks} Marks
            </span>
            {question.answer_required && (
              <span className="inline-flex items-center gap-1.5">
                <span className="h-2 w-2 rounded-full bg-amber-500" />
                Required
              </span>
            )}
          </div>

          <div className="mt-3 flex flex-wrap items-center justify-end gap-1 sm:gap-2">
            <button
              type="button"
              onClick={() => onPreview(question)}
              className="inline-flex items-center gap-1.5 rounded-md px-2 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-100"
            >
              <Eye size={15} />
              Preview
            </button>
            <button
              type="button"
              onClick={() => openEditQuestion(question._id)}
              className="inline-flex items-center gap-1.5 rounded-md px-2 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-100"
            >
              <Pencil size={15} />
              Edit
            </button>
            <button
              type="button"
              onClick={() => duplicateQuestion(question)}
              className="inline-flex items-center gap-1.5 rounded-md px-2 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-100"
            >
              <Copy size={15} />
              Duplicate
            </button>
            <button
              type="button"
              onClick={() => removeQuestion(question)}
              className="inline-flex items-center gap-1.5 rounded-md px-2 py-1.5 text-sm font-medium text-red-600 hover:bg-red-50"
            >
              <Trash2 size={15} />
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function QuizQuestionList({
  questions,
}: {
  questions: QuizQuestion[];
}) {
  const { reorderQuestions } = useCoursesStore();
  const [previewQuestion, setPreviewQuestion] = useState<QuizQuestion | null>(
    null
  );

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: { distance: 8 },
    })
  );

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    const oldIndex = questions.findIndex((q) => q._id === active.id);
    const newIndex = questions.findIndex((q) => q._id === over.id);
    if (oldIndex === -1 || newIndex === -1) return;

    reorderQuestions(arrayMove(questions, oldIndex, newIndex));
  };

  if (questions.length === 0) {
    return (
      <p className="py-6 text-center text-sm text-gray-500">
        No questions added yet.
      </p>
    );
  }

  return (
    <>
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragEnd={handleDragEnd}
      >
        <SortableContext
          items={questions.map((q) => q._id)}
          strategy={verticalListSortingStrategy}
        >
          <div className="space-y-3">
            {questions.map((question, index) => (
              <SortableQuestionCard
                key={question._id}
                question={question}
                index={index}
                onPreview={setPreviewQuestion}
              />
            ))}
          </div>
        </SortableContext>
      </DndContext>

      <QuizQuestionPreviewModal
        question={previewQuestion}
        onClose={() => setPreviewQuestion(null)}
      />
    </>
  );
}
