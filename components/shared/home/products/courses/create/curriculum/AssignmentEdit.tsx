import TextEditor from "@/components/shared/text-editor/TextEditor";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import SelectList from "@/components/ui/SelectList";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import ValidationErrorMsg from "@/components/ValidationErrorMsg";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import {
  CourseAssignmentSchema,
  TCourseAssignmentSchema,
  TCourseSchema,
} from "@/schema/course.schema";
import { zodResolver } from "@hookform/resolvers/zod";
import Image from "next/image";
import React, {
  ChangeEvent,
  DragEvent,
  useEffect,
  useRef,
  useState,
} from "react";
import {
  Controller,
  useFieldArray,
  useForm,
  useFormContext,
} from "react-hook-form";
import { IoIosArrowDown } from "react-icons/io";

const options = [
  { id: 2, name: "Minutes", value: "minutes" },
  { id: 3, name: "Hours", value: "hours" },
  { id: 1, name: "Days", value: "days" },
  { id: 4, name: "Weeks", value: "weeks" },
  { id: 5, name: "Months", value: "months" },
];

export default function AssignmentEdit({
  isEdit = false,
}: {
  isEdit?: boolean;
}) {
  const [isMedia, setIsMedia] = useState(false);
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [previews, setPreviews] = useState<string[]>([]);

  const {
    chapterId,
    closeNewAssignment,
    clearChapterId,
    assignmentID,
    closeEditAssignment,
  } = useCoursesStore();

  const { watch, control } = useFormContext<TCourseSchema>();

  const chapterIndex = watch("course_chapters", []).findIndex(
    (chapter) => chapter._id === chapterId
  );
  const safeChapterIndex = chapterIndex === -1 ? 0 : chapterIndex;

  const { append, update } = useFieldArray({
    control,
    name: `course_chapters.${safeChapterIndex}.assignments`,
    keyName: "uid",
  });

  const formMethods = useForm<TCourseAssignmentSchema>({
    resolver: zodResolver(CourseAssignmentSchema),
    defaultValues: {
      _id: Date.now(),
      type: "assignment",
      title: "",
      instructions: "",
      is_published: false,
      total_marks: 1,
      minimum_pass_marks: 0,
      time_limit: 1,
      time_limit_option: "weeks",
      attachments: null,
      file_upload_limit: 1,
    },
  });

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files) {
      const validTypes = [
        "image/jpeg",
        "image/jpg",
        "image/png",
        "image/gif",
        "image/svg+xml",
      ];
      const fileArray = Array.from(files).filter((file) =>
        validTypes.includes(file.type)
      );

      fileArray.forEach((file) => {
        const reader = new FileReader();
        reader.onloadend = () => {
          setPreviews((prev) => [...prev, reader.result as string]);
        };
        reader.readAsDataURL(file);
      });
    }
  };

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    const files = e.dataTransfer.files;
    if (files) {
      const validTypes = [
        "image/jpeg",
        "image/jpg",
        "image/png",
        "image/gif",
        "image/svg+xml",
      ];
      const fileArray = Array.from(files).filter((file) =>
        validTypes.includes(file.type)
      );

      fileArray.forEach((file) => {
        formMethods.setValue("attachments", [
          ...(formMethods.watch("attachments") ?? []),
          file,
        ]);

        const reader = new FileReader();
        reader.onloadend = () => {
          setPreviews((prev) => [...prev, reader.result as string]);
        };
        reader.readAsDataURL(file);
      });
    }
  };

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault();
  };

  const removePreview = (index: number) => {
    formMethods.setValue("attachments", [
      ...(formMethods.watch("attachments") ?? []).filter((_, i) => i !== index),
    ]);
    setPreviews((prev) => prev.filter((_, i) => i !== index));
  };

  useEffect(() => {
    const watchedLesson =
      watch(`course_chapters.${safeChapterIndex}.assignments`) || [];
    if (isEdit && assignmentID && watchedLesson?.length > 0) {
      const index = watchedLesson?.findIndex((c) => c._id === assignmentID);
      if (index !== -1) {
        formMethods.reset(watchedLesson[index]);
      }
    }
  }, [isEdit, assignmentID, watch]);

  const handleSave = (data: TCourseAssignmentSchema) => {
    if (isEdit && assignmentID) {
      const watchedAssignments =
        watch(`course_chapters.${safeChapterIndex}.assignments`) || [];

      const index = watchedAssignments?.findIndex(
        (c) => c._id === assignmentID
      );
      if (index !== -1) {
        update(index, {
          ...watchedAssignments[index],
          title: data.title,
          instructions: data.instructions,
          is_published: data.is_published,
          file_upload_limit: data.file_upload_limit,
          minimum_pass_marks: data.minimum_pass_marks,
          time_limit: data.time_limit,
          time_limit_option: data.time_limit_option,
          total_marks: data.total_marks,
          type: "assignment",
        });
      }
      closeEditAssignment();
    } else {
      closeNewAssignment();
      clearChapterId();
      append({
        ...data,
      });
      formMethods.reset();
    }
  };

  return (
    <>
      <div className="border-b border-gray-300 z-10 sticky top-0 bg-white mb-5">
        <div className="max-w-7xl mx-auto flex items-center justify-between py-4">
          <p className="font-semibold text-xl">
            {!isEdit ? "New" : "Edit"} Assignment
          </p>
          <Button type="button" onClick={formMethods.handleSubmit(handleSave)}>
            Save
          </Button>
        </div>
      </div>
      <div className="max-w-7xl mx-auto flex items-start gap-5 pb-5">
        <div className="w-full border rounded-xl">
          <div className="p-4 border-b border-gray-300">
            <p className="font-semibold">Assignment Details</p>
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
            <div className="mb-3">
              <button
                onClick={() => setIsMedia((prev) => !prev)}
                className="border border-primary text-sm px-3 py-1.5 text-primary font-medium rounded-md flex items-center gap-1"
              >
                Add Attachments <IoIosArrowDown />
              </button>
              {isMedia && (
                <div className="mt-3">
                  <label className="text-sm block font-medium mb-1">
                    Add Attachments
                  </label>
                  <div className="bg-gray-100 p-2 rounded-md">
                    {previews.length > 0 ? (
                      <div className="mt-2 space-y-4">
                        {previews.map((src, idx) => (
                          <div
                            key={idx}
                            className="flex items-center justify-between border rounded p-2"
                          >
                            <Image
                              src={src}
                              alt={`Preview ${idx}`}
                              width={80}
                              height={80}
                              className="object-contain rounded"
                            />
                            <button
                              type="button"
                              onClick={() => removePreview(idx)}
                              className="text-red-500 text-xl"
                            >
                              &times;
                            </button>
                          </div>
                        ))}
                        <Button className="w-full mt-5 flex items-center justify-center">
                          Upload {previews.length} files
                        </Button>
                      </div>
                    ) : (
                      <div
                        onDrop={handleDrop}
                        onDragOver={handleDragOver}
                        onClick={() => fileInputRef.current?.click()}
                        className="border border-dashed flex justify-center items-center p-5 rounded-md min-h-32 cursor-pointer text-center"
                      >
                        <p>
                          Drop images here or{" "}
                          <button
                            type="button"
                            className="text-sky-500 hover:underline"
                            onClick={(e) => {
                              e.stopPropagation();
                              fileInputRef.current?.click();
                            }}
                          >
                            Browse
                          </button>
                        </p>
                        <input
                          type="file"
                          accept="image/jpeg,image/jpg,image/png,image/gif,image/svg+xml"
                          multiple
                          onChange={handleFileChange}
                          ref={fileInputRef}
                          className="hidden"
                        />
                      </div>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
        <div className="w-96 min-w-96">
          <div className="border rounded-xl">
            <div className="p-4 border-b border-gray-300">
              <p className="font-semibold">Assignment Settings</p>
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
                File Upload Limit
              </label>
              <p className="text-sm text-gray-600 font-medium mb-1">
                Define the number of files that a student can upload in this
                assignment.
              </p>
              <InputField
                type="number"
                className="w-full"
                {...formMethods.register("file_upload_limit")}
                error={formMethods.formState.errors.file_upload_limit?.message}
              />
            </div>
          </div>
          <div className="border rounded-xl mt-5">
            <div className="p-4 border-b border-gray-300">
              <p className="font-semibold">Assignment Marks</p>
            </div>
            <div className="p-4">
              <div className="mb-3">
                <label className="text-sm font-semibold block">
                  Total Marks <span className="text-red-500">*</span>
                </label>
                <p className="text-sm text-gray-600 font-medium mb-1">
                  Maximum marks a student can score
                </p>
                <InputField
                  type="number"
                  className="w-full"
                  {...formMethods.register("total_marks")}
                  error={formMethods.formState.errors.total_marks?.message}
                />
              </div>
              <div>
                <label className="text-sm font-semibold block">
                  Minimum Pass Marks <span className="text-red-500">*</span>
                </label>
                <p className="text-sm text-gray-600 font-medium mb-1">
                  Minimum marks required for the student to pass this
                  assignment.
                </p>
                <InputField
                  type="number"
                  className="w-full"
                  {...formMethods.register("minimum_pass_marks")}
                  error={
                    formMethods.formState.errors.minimum_pass_marks?.message
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
