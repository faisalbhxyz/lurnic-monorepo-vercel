"use client";

import InputField from "@/components/ui/InputField";
import Image from "next/image";
import React, { useRef } from "react";
import { Controller, useFormContext } from "react-hook-form";
import { FiUpload } from "react-icons/fi";
import { IoIosArrowDown, IoIosArrowUp } from "react-icons/io";
import { TCourseSchema } from "@/schema/course.schema";

const certificates = [
  { id: 1, name: "/images/Certificat-14.jpg" },
  { id: 2, name: "/images/Certificat-15.jpg" },
  { id: 3, name: "/images/Certificat-16.jpg" },
  { id: 4, name: "/images/Certificat-17.jpg" },
];

export default function CertificatesTab() {
  const [isOpen, setIsOpen] = React.useState(true);
  const {
    register,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useFormContext<TCourseSchema>();

  const ownerFileInputRef = useRef<HTMLInputElement>(null);
  const instructorFileInputRef = useRef<HTMLInputElement>(null);

  const isEnabled = watch("certificate_settings.is_enabled");
  const selectedTemplate = watch("certificate_settings.template_path");
  const ownerSignature = watch("certificate_settings.owner_signature");
  const instructorSignature = watch("certificate_settings.instructor_signature");

  const ownerPreview =
    ownerSignature && typeof ownerSignature === "object" && "isDBImg" in ownerSignature
      ? (ownerSignature.name as string)
      : ownerSignature instanceof File
        ? URL.createObjectURL(ownerSignature)
        : null;

  const instructorPreview =
    instructorSignature &&
    typeof instructorSignature === "object" &&
    "isDBImg" in instructorSignature
      ? (instructorSignature.name as string)
      : instructorSignature instanceof File
        ? URL.createObjectURL(instructorSignature)
        : null;

  return (
    <div className="flex items-start gap-10">
      <div className="w-full">
        <div className="mb-6 flex items-center justify-between rounded-md border p-4">
          <div>
            <p className="font-medium">Enable certificate</p>
            <p className="text-sm text-gray-500 mt-1">
              Auto-issue when a student reaches the completion threshold.
            </p>
          </div>
          <label className="inline-flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
              {...register("certificate_settings.is_enabled")}
            />
            <span className="text-sm font-medium">Enabled</span>
          </label>
        </div>

        <div className="mb-5">
          <label className="text-sm font-medium mb-1 block">
            Minimum completion (%)
          </label>
          <InputField
            type="number"
            min={1}
            max={100}
            className="w-full max-w-xs"
            disabled={!isEnabled}
            {...register("certificate_settings.completion_percent")}
          />
          {errors.certificate_settings?.completion_percent && (
            <p className="text-sm text-red-500 mt-1">
              {errors.certificate_settings.completion_percent.message}
            </p>
          )}
          <p className="text-xs text-gray-500 mt-1">
            Example: 80 means the student must complete 80% of selected items
            below.
          </p>
        </div>

        <div className="mb-6 rounded-md border p-4">
          <p className="font-medium mb-1">Count toward progress</p>
          <p className="text-sm text-gray-500 mb-3">
            Choose what counts for this course. Different courses can use
            different rules.
          </p>
          <div className="space-y-2">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                {...register("certificate_settings.count_lessons")}
              />
              Lessons (published)
            </label>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                {...register("certificate_settings.count_quizzes")}
              />
              Quizzes (submitted)
            </label>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                {...register("certificate_settings.count_assignments")}
              />
              Assignments (submitted)
            </label>
          </div>
          {errors.certificate_settings?.count_lessons && (
            <p className="text-sm text-red-500 mt-2">
              {errors.certificate_settings.count_lessons.message}
            </p>
          )}
        </div>

        <p className="font-medium mb-2">Choose a design</p>
        <div className="border rounded-md overflow-hidden mb-4">
          <Image
            src={selectedTemplate || certificates[0].name}
            alt="Selected Certificate"
            width={500}
            height={500}
            className="w-full h-auto object-cover"
          />
        </div>

        <div className="flex items-center justify-center gap-3">
          {certificates.map((image) => (
            <button
              type="button"
              key={image.id}
              disabled={!isEnabled}
              onClick={() =>
                setValue("certificate_settings.template_path", image.name, {
                  shouldDirty: true,
                })
              }
              className={`w-12 h-10 border-2 rounded-md overflow-hidden p-0 disabled:opacity-50 ${
                selectedTemplate === image.name
                  ? "border-primary"
                  : "border-transparent"
              }`}
            >
              <Image
                src={image.name}
                alt={`Certificate ${image.id}`}
                width={100}
                height={100}
                className="w-full h-full"
              />
            </button>
          ))}
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium mb-1">Certificate Title</label>
          <InputField
            className="w-full"
            disabled={!isEnabled}
            {...register("certificate_settings.title")}
          />
        </div>
        <div className="mt-5">
          <label className="text-sm font-medium mb-1">
            Certificate Subtitle One
          </label>
          <InputField
            className="w-full"
            disabled={!isEnabled}
            {...register("certificate_settings.subtitle_one")}
          />
        </div>
        <div className="mt-5">
          <label className="text-sm font-medium mb-1">
            Certificate Subtitle Two
          </label>
          <InputField
            className="w-full"
            disabled={!isEnabled}
            {...register("certificate_settings.subtitle_two")}
          />
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium block mb-2">
            School Owner Signature (150x250 px)
          </label>
          <div className="flex items-center gap-4">
            {ownerPreview && (
              <Image
                src={ownerPreview}
                alt="Owner signature"
                width={100}
                height={100}
                className="w-20 h-20 object-cover rounded"
              />
            )}
            <button
              type="button"
              disabled={!isEnabled}
              onClick={() => ownerFileInputRef.current?.click()}
              className="border px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline disabled:opacity-50"
            >
              <FiUpload /> Upload
            </button>
            <Controller
              control={control}
              name="certificate_settings.owner_signature"
              render={({ field: { onChange } }) => (
                <input
                  ref={ownerFileInputRef}
                  type="file"
                  accept="image/*"
                  hidden
                  disabled={!isEnabled}
                  onChange={(event) => {
                    const file = event.target.files?.[0];
                    if (file) onChange(file);
                  }}
                />
              )}
            />
          </div>
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium block mb-2">
            Instructor Signature (150x250 px)
          </label>
          <div className="flex items-center gap-4">
            {instructorPreview && (
              <Image
                src={instructorPreview}
                alt="Instructor signature"
                width={100}
                height={100}
                className="w-20 h-20 object-cover rounded"
              />
            )}
            <button
              type="button"
              disabled={!isEnabled}
              onClick={() => instructorFileInputRef.current?.click()}
              className="border px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline disabled:opacity-50"
            >
              <FiUpload /> Upload
            </button>
            <Controller
              control={control}
              name="certificate_settings.instructor_signature"
              render={({ field: { onChange } }) => (
                <input
                  ref={instructorFileInputRef}
                  type="file"
                  accept="image/*"
                  hidden
                  disabled={!isEnabled}
                  onChange={(event) => {
                    const file = event.target.files?.[0];
                    if (file) onChange(file);
                  }}
                />
              )}
            />
          </div>
        </div>
      </div>

      <div className="w-80 min-w-80">
        <div className="border p-5 rounded-md">
          <button
            type="button"
            onClick={() => setIsOpen(!isOpen)}
            className={`flex items-center justify-between w-full ${
              isOpen ? "text-primary" : "text-gray-700"
            }`}
          >
            <p className="font-semibold text-start">How certificates work</p>
            {isOpen ? <IoIosArrowUp /> : <IoIosArrowDown />}
          </button>

          {isOpen && (
            <p className="text-sm mt-5 text-gray-600">
              Progress is calculated from the items you select above. When a
              student reaches your completion percentage, a certificate is
              issued automatically. Lessons are marked complete from the
              storefront when the student finishes watching/reading.
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
