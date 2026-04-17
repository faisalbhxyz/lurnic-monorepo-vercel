import TextEditor from "@/components/shared/text-editor/TextEditor";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import SelectList from "@/components/ui/SelectList";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import ValidationErrorMsg from "@/components/ValidationErrorMsg";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import {
  CourseQuizSchema,
  TCourseQuizSchema,
  TCourseSchema,
} from "@/schema/course.schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { LucideEdit } from "lucide-react";
import React, { useEffect, useState } from "react";
import {
  Controller,
  useFieldArray,
  useForm,
  useFormContext,
} from "react-hook-form";
import { LuTrash2 } from "react-icons/lu";
import { toast } from "sonner";

const options = [
  { id: 1, name: "Days", value: "days" },
  { id: 2, name: "minutes", value: "minutes" },
  { id: 3, name: "hours", value: "hours" },
  { id: 4, name: "weeks", value: "weeks" },
  { id: 5, name: "months", value: "months" },
];

export default function QuizEdit({ isEdit = false }: { isEdit?: boolean }) {
  const {
    openNewQuestion,
    chapterId,
    clearChapterId,
    closeNewQuiz,
    closeEditQuiz,
    questions,
    setQuestions,
    removeQuestion,
    editQuizID,
    openEditQuestion,
  } = useCoursesStore();

  const { watch, control } = useFormContext<TCourseSchema>();

  const chapterIndex = watch("course_chapters", []).findIndex(
    (chapter) => chapter._id === chapterId
  );
  const safeChapterIndex = chapterIndex === -1 ? 0 : chapterIndex;

  const { append, update } = useFieldArray({
    control,
    name: `course_chapters.${safeChapterIndex}.quizzes`,
    keyName: "uid",
  });

  const formMethods = useForm<TCourseQuizSchema>({
    resolver: zodResolver(CourseQuizSchema),
    defaultValues: {
      _id: Date.now(),
      type: "quiz",
      title: "",
      instructions: "",
      enable_retry: false,
      retry_attempts: 0,
      is_published: false,
      minimum_pass_percentage: 0,
      single_quiz_view: false,
      reveal_answers: false,
      randomize_questions: false,
      time_limit: 1,
      time_limit_option: "weeks",
      total_visible_questions: 1,
      questions: [],
    },
  });

  useEffect(() => {
    const watchedLesson =
      watch(`course_chapters.${safeChapterIndex}.quizzes`) || [];
    if (isEdit && editQuizID && watchedLesson?.length > 0) {
      const index = watchedLesson?.findIndex((c) => c._id === editQuizID);
      if (index !== -1) {
        formMethods.reset(watchedLesson[index]);
        // console.log("[QS] ", watchedLesson[index].questions);

        setQuestions(watchedLesson[index].questions);
      }
    }
  }, [isEdit, editQuizID, watch]);

  const handleSave = (data: TCourseQuizSchema) => {
    if (questions.length === 0) {
      toast.error("Please add at least one question.");
      return;
    }

    if (isEdit && editQuizID) {
      const watchedLesson =
        watch(`course_chapters.${safeChapterIndex}.quizzes`) || [];

      const index = watchedLesson?.findIndex((c) => c._id === editQuizID);
      if (index !== -1) {
        update(index, {
          ...watchedLesson[index],
          title: data.title,
          instructions: data.instructions,
          is_published: data.is_published,
          enable_retry: data.enable_retry,
          retry_attempts: data.retry_attempts,
          minimum_pass_percentage: data.minimum_pass_percentage,
          single_quiz_view: data.single_quiz_view,
          reveal_answers: data.reveal_answers,
          randomize_questions: data.randomize_questions,
          time_limit: data.time_limit,
          time_limit_option: data.time_limit_option,
          total_visible_questions: data.total_visible_questions,
          questions: questions,
        });
      }
      clearChapterId();
      closeEditQuiz();
    } else {
      closeNewQuiz();
      clearChapterId();
      append({
        ...data,
        questions: questions,
      });
      formMethods.reset();
    }
  };

  return (
    <>
      <div className="border-b border-gray-300 z-10 sticky top-0 bg-white mb-5">
        <div className="max-w-7xl mx-auto flex items-center justify-between py-4">
          <p className="font-semibold text-xl">
            {!isEdit ? "New" : "Edit"} Quiz
          </p>
          <Button type="button" onClick={formMethods.handleSubmit(handleSave)}>
            Save
          </Button>
        </div>
      </div>
      {/* {JSON.stringify(formMethods.formState.errors)} */}
      <div className="max-w-7xl mx-auto flex items-start gap-5 pb-5">
        <div className="w-full border rounded-xl">
          <div className="p-4 border-b border-gray-300">
            <p className="font-semibold">Quiz Details</p>
          </div>
          <div className="p-4">
            <div className="mb-5">
              <label
                htmlFor="title"
                className="block mb-1 text-sm font-semibold"
              >
                Title *
              </label>
              <InputField
                className="w-full"
                {...formMethods.register("title")}
                error={formMethods.formState.errors.title?.message}
              />
            </div>
            <div className="mb-5">
              <label
                htmlFor="title"
                className="block mb-1 text-sm font-semibold"
              >
                Instructions *
              </label>
              <Controller
                control={formMethods.control}
                name="instructions"
                render={({ field }) => (
                  <TextEditor
                    value={field.value || ""}
                    onChange={field.onChange}
                  />
                )}
              />
              {formMethods.formState.errors.instructions && (
                <ValidationErrorMsg
                  error={formMethods.formState.errors.instructions.message}
                />
              )}
            </div>
            <div className="border rounded-xl">
              <div className="p-4 border-b border-gray-300">
                <p className="font-semibold">Questions</p>
              </div>
              <div className="p-2">
                <div className="mb-3 p-4">
                  <ul>
                    {questions.map((question, index) => (
                      <li key={index} className="flex items-center gap-3 mb-2">
                        {index + 1}. {question.title}{" "}
                        {isEdit && (
                          <button
                            type="button"
                            className="text-blue-500"
                            onClick={() => openEditQuestion(question._id)}
                          >
                            <LucideEdit size={16} />
                          </button>
                        )}
                        <button
                          type="button"
                          className="text-red-500"
                          onClick={() => removeQuestion(question)}
                        >
                          <LuTrash2 size={16} />
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
                <div className="flex items-center justify-center">
                  <Button
                    onClick={() => {
                      openNewQuestion(formMethods.watch("_id"));
                    }}
                  >
                    Add Questions
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div className="w-96 min-w-96">
          <div className="border rounded-xl">
            <div className="p-4 border-b border-gray-300">
              <p className="font-semibold">Quiz Settings</p>
            </div>
            <div className="p-4">
              <div className="flex items-center justify-between">
                <label className="text-sm font-semibold">Publish</label>
                <Controller
                  name="is_published"
                  control={formMethods.control}
                  render={({ field }) => (
                    <ToggleSwitch
                      checked={field.value}
                      onChange={() => field.onChange(!field.value)}
                    />
                  )}
                />
              </div>
              <p className="text-sm text-gray-600 font-medium">
                Make lesson as publish
              </p>
              <div className="mt-4">
                <div className="flex items-center justify-between">
                  <label className="text-sm font-semibold">
                    Randomize Questions
                  </label>
                  <Controller
                    name="randomize_questions"
                    control={formMethods.control}
                    render={({ field }) => (
                      <ToggleSwitch
                        checked={field.value}
                        onChange={() => field.onChange(!field.value)}
                      />
                    )}
                  />
                </div>
                <p className="text-sm text-gray-600 font-medium">
                  This will randomize questions for the quiz
                </p>
              </div>
              <div className="mt-4">
                <div className="flex items-center justify-between">
                  <label className="text-sm font-semibold">
                    Single Quiz View
                  </label>
                  <Controller
                    name="single_quiz_view"
                    control={formMethods.control}
                    render={({ field }) => (
                      <ToggleSwitch
                        checked={field.value}
                        onChange={() => field.onChange(!field.value)}
                      />
                    )}
                  />
                </div>
                <p className="text-sm text-gray-600 font-medium">
                  Show only one question at a time
                </p>
              </div>
              <div className="my-4">
                <label className="text-sm font-semibold block">
                  Time Limit <span className="text-red-500">*</span>
                </label>
                <div className="flex items-center gap-5">
                  <InputField
                    type="number"
                    className="w-full"
                    {...formMethods.register("time_limit")}
                    error={formMethods.formState.errors.time_limit?.message}
                  />
                  <Controller
                    control={formMethods.control}
                    name="time_limit_option"
                    defaultValue="weeks"
                    render={({ field }) => (
                      <SelectList
                        options={options}
                        className="w-full"
                        value={options.find(
                          (item) =>
                            item.value.toLowerCase() ===
                            field.value.toLowerCase()
                        )}
                        onChange={(d) => field.onChange(d.value)}
                      />
                    )}
                  />
                </div>
              </div>
              <label className="text-sm font-semibold block">
                Total Visible Questions
              </label>
              <p className="text-sm text-gray-600 font-medium mb-1">
                Number of questions to be displayed in the quiz
              </p>
              <InputField
                type="number"
                className="w-full"
                {...formMethods.register("total_visible_questions")}
                error={
                  formMethods.formState.errors.total_visible_questions?.message
                }
              />
            </div>
          </div>
          <div className="border rounded-xl mt-5">
            <div className="p-4 border-b border-gray-300">
              <p className="font-semibold">Answer Settings</p>
            </div>
            <div className="p-4">
              <div className="flex items-center justify-between">
                <label className="text-sm font-semibold">Reveal Answers</label>
                <Controller
                  name="reveal_answers"
                  control={formMethods.control}
                  render={({ field }) => (
                    <ToggleSwitch
                      checked={field.value}
                      onChange={() => field.onChange(!field.value)}
                    />
                  )}
                />
              </div>
              <p className="text-sm text-gray-600 font-medium">
                Show answers after the quiz submission.
              </p>
              <div className="mt-4">
                <div className="flex items-center justify-between">
                  <label className="text-sm font-semibold">Enable Retry</label>
                  <Controller
                    name="enable_retry"
                    control={formMethods.control}
                    render={({ field }) => (
                      <ToggleSwitch
                        checked={field.value}
                        onChange={() => field.onChange(!field.value)}
                      />
                    )}
                  />
                </div>
                <p className="text-sm text-gray-600 font-medium">
                  Enable retry on quiz submission.
                </p>
              </div>
              <div className="mt-4">
                <label className="text-sm font-semibold block">
                  Retry Attempts
                  <span className="text-red-500">*</span>
                </label>
                <p className="text-sm text-gray-600 font-medium mb-1">
                  Number of times a student can attempt the quiz.
                </p>
                <InputField
                  type="number"
                  className="w-full"
                  {...formMethods.register("retry_attempts")}
                  error={formMethods.formState.errors.retry_attempts?.message}
                />
              </div>
            </div>
          </div>
          <div className="border rounded-xl mt-5">
            <div className="p-4 border-b border-gray-300">
              <p className="font-semibold">Quiz Grading</p>
            </div>
            <div className="p-4">
              <div>
                <label className="text-sm font-semibold block">
                  Minimum Pass Percentage{" "}
                  <span className="text-red-500">*</span>
                </label>
                <p className="text-sm text-gray-600 font-medium mb-1">
                  Minimum pass percentage for the quiz
                </p>
                <InputField
                  type="number"
                  className="w-full"
                  {...formMethods.register("minimum_pass_percentage")}
                  error={
                    formMethods.formState.errors.minimum_pass_percentage
                      ?.message
                  }
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
