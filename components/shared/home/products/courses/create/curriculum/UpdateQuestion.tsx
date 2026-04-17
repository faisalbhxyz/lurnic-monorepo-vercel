"use client";

import React, {
  ChangeEvent,
  useRef,
  useState,
  DragEvent,
  useEffect,
} from "react";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import Modal from "@/components/ui/Modal";
import { RxCross2 } from "react-icons/rx";
import InputField from "@/components/ui/InputField";
import TextEditor from "@/components/shared/text-editor/TextEditor";
import { IoIosArrowDown } from "react-icons/io";
import SelectList from "@/components/ui/SelectList";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import { RiInformationLine } from "react-icons/ri";
import Image from "next/image";
import Button from "@/components/ui/Button";
import { Controller, useForm } from "react-hook-form";
import {
  QuizQuestionSchema,
  TQuizQuestionSchema,
} from "@/schema/course.schema";
import { zodResolver } from "@hookform/resolvers/zod";

const questionType = [
  {
    id: 1,
    name: "Multiple choice",
    value: "multiple_choice",
  },
  {
    id: 2,
    name: "Single Choice",
    value: "single_choice",
  },
  {
    id: 3,
    name: "True/False",
    value: "true_false",
  },
];

export default function UpdateQuestion() {
  const {
    isEditQuestion,
    closeEditQuestion,
    questions,
    editQuestionID,
    updateQuestion,
  } = useCoursesStore();
  const [isMedia, setIsMedia] = useState(false);
  const [isAdditionalDetails, setIsAdditionalDetails] = useState(false);
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [previews, setPreviews] = useState<string[]>([]);

  const formMethods = useForm<TQuizQuestionSchema>({
    resolver: zodResolver(QuizQuestionSchema),
    defaultValues: {
      _id: Date.now(),
      title: "",
      details: "",
      type: "single_choice",
      answer_required: false,
      answer_explanation: "",
      marks: 1,
      media: [],
    },
  });

  useEffect(() => {
    if (editQuestionID) {
      const question = questions.find((q) => q._id === editQuestionID);
      if (question) {
        formMethods.reset({
          _id: question._id,
          title: question.title,
          details: question.details,
          type: question.type,
          answer_required: question.answer_required,
          answer_explanation: question.answer_explanation,
          marks: question.marks,
          media: question.Media,
        });
      }
    }
  }, [editQuestionID]);

  const handleSave = (data: TQuizQuestionSchema) => {
    updateQuestion({
      ...data,
      id: data._id,
    });
    closeEditQuestion();
    // formMethods.reset({
    //   _id: Date.now(),
    //   title: "",
    //   details: "",
    //   type: "single_choice",
    //   answer_required: false,
    //   answer_explanation: "",
    //   marks: 1,
    //   media: [],
    // });
    // setPreviews([]);
  };

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
        formMethods.setValue("media", [
          ...(formMethods.watch("media") ?? []),
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
    formMethods.setValue("media", [
      ...(formMethods.watch("media") ?? []).filter((_, i) => i !== index),
    ]);
    setPreviews((prev) => prev.filter((_, i) => i !== index));
  };

  return (
    <Modal isOpen={isEditQuestion} onClose={closeEditQuestion} className="p-0">
      {/* {JSON.stringify(formMethods.formState.errors)} */}
      <div className="p-4 flex items-center justify-between border-b border-gray-300">
        <p className="font-semibold text-lg">Update Question</p>
        <button type="button" onClick={closeEditQuestion}>
          <RxCross2 />
        </button>
      </div>
      <div className="p-4">
        <div className="mb-3">
          <label className="text-sm block font-medium mb-1">
            Question Title <span className="text-red-600">*</span>
          </label>
          <InputField
            className="w-full"
            {...formMethods.register("title")}
            error={formMethods.formState.errors.title?.message}
          />
        </div>
        <div className="mb-3">
          <label className="text-sm block font-medium mb-1">
            Question Details
          </label>
          <Controller
            control={formMethods.control}
            name="details"
            render={({ field }) => (
              <TextEditor value={field.value || ""} onChange={field.onChange} />
            )}
          />
        </div>
        <div className="mb-3">
          <button
            onClick={() => setIsMedia((prev) => !prev)}
            className="border border-primary text-sm px-3 py-1.5 text-primary font-medium rounded-md flex items-center gap-1"
          >
            Add Media <IoIosArrowDown />
          </button>
          {isMedia && (
            <div className="mt-3">
              <label className="text-sm block font-medium mb-1">
                Add Media
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
        <div className="mb-3">
          <label className="text-sm block font-medium mb-1">
            Question Type <span className="text-red-600">*</span>
          </label>
          <Controller
            control={formMethods.control}
            name="type"
            render={({ field: { value, onChange } }) => {
              const selected = questionType.find((o) => o.value === value);
              return (
                <SelectList
                  options={questionType}
                  value={selected}
                  onChange={(option) => {
                    onChange(option.value);
                  }}
                  className="w-full font-medium"
                />
              );
            }}
          />
        </div>
        <div className="mb-3 flex items-center gap-5">
          <div className="w-1/2">
            <label className="text-sm block font-medium mb-1">Marks</label>
            <InputField
              type="number"
              className="w-full"
              {...formMethods.register("marks")}
              error={formMethods.formState.errors.marks?.message}
            />
          </div>
          <div>
            <label className="text-sm block font-medium mb-1">
              Answer Required
            </label>
            <Controller
              name="answer_required"
              control={formMethods.control}
              render={({ field }) => (
                <ToggleSwitch
                  checked={field.value}
                  onChange={() => field.onChange(!field.value)}
                />
              )}
            />
          </div>
        </div>
        <div className="bg-gray-100 border p-3 rounded-md">
          <div
            onClick={() => setIsAdditionalDetails((prev) => !prev)}
            className="flex items-center justify-between cursor-pointer"
          >
            <div className="flex items-center gap-2">
              <RiInformationLine />
              <p className="font-medium">Additional Details</p>
            </div>
            <IoIosArrowDown />
          </div>
          {isAdditionalDetails && (
            <div className="mt-3">
              <label className="text-sm block font-medium mb-1">
                Answer Explanation
              </label>
              <Controller
                control={formMethods.control}
                name="answer_explanation"
                render={({ field }) => (
                  <TextEditor
                    value={field.value || ""}
                    onChange={field.onChange}
                  />
                )}
              />
            </div>
          )}
        </div>
        <div className="mt-4 flex justify-end">
          <Button onClick={formMethods.handleSubmit(handleSave)}>Save</Button>
        </div>
      </div>
    </Modal>
  );
}
