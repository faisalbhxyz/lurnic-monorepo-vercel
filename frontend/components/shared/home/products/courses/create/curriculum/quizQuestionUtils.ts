interface QuizOption {
  id: string;
  text: string;
}

interface QuizQuestionLike {
  type: "multiple_choice" | "single_choice" | "true_false";
  options?: QuizOption[] | null;
  correct_answer?:
    | { value?: string | boolean; values?: string[] }
    | null
    | undefined;
}

export function getCorrectAnswerLabel(question: QuizQuestionLike): string {
  if (question.type === "true_false") {
    if (question.correct_answer?.value === true) return "True";
    if (question.correct_answer?.value === false) return "False";
    return "—";
  }

  const options = question.options ?? [];

  if (question.type === "single_choice") {
    const value = question.correct_answer?.value;
    if (typeof value !== "string") return "—";
    return options.find((o) => o.id === value)?.text?.trim() || value;
  }

  const values = question.correct_answer?.values ?? [];
  if (values.length === 0) return "—";

  return values
    .map((id) => options.find((o) => o.id === id)?.text?.trim() || id)
    .join(", ");
}

export function getQuestionTypeBadge(type: QuizQuestionLike["type"]): {
  label: string;
  className: string;
} {
  switch (type) {
    case "single_choice":
      return {
        label: "MCQ",
        className: "bg-emerald-50 text-emerald-700 border-emerald-100",
      };
    case "multiple_choice":
      return {
        label: "Multi MCQ",
        className: "bg-sky-50 text-sky-700 border-sky-100",
      };
    case "true_false":
      return {
        label: "True/False",
        className: "bg-violet-50 text-violet-700 border-violet-100",
      };
    default:
      return {
        label: type,
        className: "bg-gray-50 text-gray-700 border-gray-100",
      };
  }
}

export function getOptionsCount(question: QuizQuestionLike): number {
  if (question.type === "true_false") return 2;
  return question.options?.length ?? 0;
}
