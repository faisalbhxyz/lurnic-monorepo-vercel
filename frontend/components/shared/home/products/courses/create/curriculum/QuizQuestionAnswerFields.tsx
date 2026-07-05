"use client";

import React, { useEffect } from "react";
import {
  Control,
  Controller,
  FieldErrors,
  useFieldArray,
  UseFormSetValue,
  UseFormWatch,
} from "react-hook-form";
import InputField from "@/components/ui/InputField";
import Button from "@/components/ui/Button";
import ValidationErrorMsg from "@/components/ValidationErrorMsg";
import { TQuizQuestionSchema } from "@/schema/course.schema";
import { LuPlus, LuTrash2 } from "react-icons/lu";

const nextOptionId = (existing: { id: string }[]) => {
  const used = new Set(existing.map((o) => o.id));
  for (let i = 0; i < 26; i++) {
    const id = String.fromCharCode(97 + i);
    if (!used.has(id)) return id;
  }
  return `opt_${Date.now()}`;
};

type Props = {
  control: Control<TQuizQuestionSchema>;
  watch: UseFormWatch<TQuizQuestionSchema>;
  setValue: UseFormSetValue<TQuizQuestionSchema>;
  errors: FieldErrors<TQuizQuestionSchema>;
};

export default function QuizQuestionAnswerFields({
  control,
  watch,
  setValue,
  errors,
}: Props) {
  const questionType = watch("type");
  const options = watch("options") ?? [];
  const correctAnswer = watch("correct_answer");

  const { fields, append, remove } = useFieldArray({
    control,
    name: "options",
  });

  useEffect(() => {
    if (questionType === "true_false") {
      setValue("options", null);
      setValue("correct_answer", { value: true });
      return;
    }

    const current = watch("options");
    if (!current || current.length < 2) {
      setValue("options", [
        { id: "a", text: "" },
        { id: "b", text: "" },
      ]);
    }

    if (questionType === "single_choice") {
      setValue("correct_answer", {
        value: watch("options")?.[0]?.id ?? "a",
      });
    } else if (questionType === "multiple_choice") {
      setValue("correct_answer", { values: [] });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [questionType, setValue]);

  if (questionType === "true_false") {
    return (
      <div className="mb-3 border rounded-md p-3 bg-gray-50">
        <p className="text-sm font-semibold mb-2">Correct Answer *</p>
        <div className="flex items-center gap-6">
          <label className="flex items-center gap-2 text-sm cursor-pointer">
            <input
              type="radio"
              name="true_false_answer"
              checked={correctAnswer?.value === true}
              onChange={() => setValue("correct_answer", { value: true })}
            />
            True
          </label>
          <label className="flex items-center gap-2 text-sm cursor-pointer">
            <input
              type="radio"
              name="true_false_answer"
              checked={correctAnswer?.value === false}
              onChange={() => setValue("correct_answer", { value: false })}
            />
            False
          </label>
        </div>
        {errors.correct_answer?.message && (
          <ValidationErrorMsg error={String(errors.correct_answer.message)} />
        )}
      </div>
    );
  }

  return (
    <div className="mb-3 border rounded-md p-3 bg-gray-50 space-y-3">
      <div className="flex items-center justify-between">
        <p className="text-sm font-semibold">Answer Options *</p>
        <Button
          type="button"
          className="p-1.5 text-sm flex items-center gap-1"
          onClick={() =>
            append({
              id: nextOptionId(options),
              text: "",
            })
          }
        >
          <LuPlus size={14} />
          Add option
        </Button>
      </div>

      {errors.options?.message && (
        <ValidationErrorMsg error={String(errors.options.message)} />
      )}

      <div className="space-y-2">
        {fields.map((field, index) => (
          <div key={field.id} className="flex items-start gap-2">
            <div className="pt-2">
              {questionType === "single_choice" ? (
                <input
                  type="radio"
                  name="single_correct"
                  checked={correctAnswer?.value === options[index]?.id}
                  onChange={() =>
                    setValue("correct_answer", { value: options[index]?.id })
                  }
                />
              ) : (
                <input
                  type="checkbox"
                  checked={(correctAnswer?.values ?? []).includes(
                    options[index]?.id
                  )}
                  onChange={(e) => {
                    const id = options[index]?.id;
                    if (!id) return;
                    const current = new Set(correctAnswer?.values ?? []);
                    if (e.target.checked) current.add(id);
                    else current.delete(id);
                    setValue("correct_answer", {
                      values: Array.from(current),
                    });
                  }}
                />
              )}
            </div>
            <InputField
              className="w-16"
              value={options[index]?.id ?? ""}
              onChange={(e) =>
                setValue(`options.${index}.id`, e.target.value, {
                  shouldValidate: true,
                })
              }
              placeholder="id"
            />
            <Controller
              control={control}
              name={`options.${index}.text`}
              render={({ field: textField }) => (
                <InputField
                  className="flex-1"
                  {...textField}
                  placeholder="Option text"
                  error={errors.options?.[index]?.text?.message}
                />
              )}
            />
            {fields.length > 2 && (
              <button
                type="button"
                className="text-red-500 pt-2"
                onClick={() => remove(index)}
              >
                <LuTrash2 size={16} />
              </button>
            )}
          </div>
        ))}
      </div>

      <p className="text-xs text-gray-600">
        {questionType === "single_choice"
          ? "Select one radio button as the correct answer."
          : "Check all correct answers for this question."}
      </p>
      {errors.correct_answer?.message && (
        <ValidationErrorMsg error={String(errors.correct_answer.message)} />
      )}
    </div>
  );
}
