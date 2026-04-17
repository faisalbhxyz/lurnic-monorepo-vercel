"use client";

import Label from "@/components/ui/Label";
import React, { useEffect, useState } from "react";
import Visibility from "./Visibility";
import Schedule from "./Schedule";
import IntroVideo from "./IntroVideo";
import RadioButton from "@/components/ui/RadioButton";
import TagsSelect from "./TagsSelect";
import Author from "./Author";
import InputField from "@/components/ui/InputField";
import TextEditor from "@/components/shared/text-editor/TextEditor";
import CourseBenefits from "./CourseBenefits";
import { toast } from "sonner";
import { Controller, useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";
import CourseFeaturedImage from "./CourseFeaturedImage";

export default function Basics() {
  const {
    register,
    handleSubmit,
    formState: { errors },
    control,
    watch,
    setValue,
  } = useFormContext<TCourseSchema>();

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const value = e.target.value;

    if (value.length > 50) {
      toast.warning("Maximum of 50 characters allowed");
      return;
    }
  };

  const watchSummary = watch("summary");
  const watchIsSchedule = watch("is_scheduled");

  useEffect(() => {
    if (watchIsSchedule) {
      setValue("visibility", "protected");
    } else {
      setValue("visibility", "public");
    }
  }, [watchIsSchedule]);

  return (
    <div className="flex">
      <div className="w-full border-r border-gray-300 py-5 pr-5">
        <div className="mb-5">
          <Label htmlFor="title">Title</Label>
          <InputField
            placeholder="Learn Javascript"
            className="w-full"
            {...register("title")}
            error={errors.title?.message}
          />
        </div>
        <div className="mb-5">
          <div className="flex items-center justify-between">
            <Label htmlFor="title">Course Summary</Label>
            <span className="text-gray-500 text-sm">
              {watchSummary
                ? watchSummary.trim() === ""
                  ? 0
                  : watchSummary.trim().split(/\s+/).length
                : 0}
              /50 words
            </span>
          </div>
          <textarea
            id="course-summary"
            rows={4}
            className="bg-white border w-full min-h-20 rounded-md px-3 py-2 outline-none focus:border-primary"
            placeholder="Enter a short course summary..."
            {...register("summary")}
          />
          {errors.summary && (
            <span className="text-red-500 text-sm py-1">
              {errors.summary.message}
            </span>
          )}
          <p className="text-[13px] text-gray-500">
            Short summary of the course, this will be displayed on the course
            card
          </p>
        </div>
        <Label htmlFor="title">Description</Label>
        <Controller
          control={control}
          name="description"
          render={({ field }) => (
            <TextEditor
              value={field.value || ""}
              onChange={(html) => field.onChange(html)}
            />
          )}
        />
        <div className="mt-5">
          <Label htmlFor="title">What&apos;s in the Course?</Label>
          <p className="text-[13px] text-gray-500 mb-1">
            List the benefits that are included in this course.
          </p>
          <CourseBenefits />
        </div>
      </div>
      <div className="w-80 min-w-80 pl-5 py-5">
        {!watchIsSchedule && (
          <div className="mb-3">
            <Label htmlFor={""}>Visibility</Label>
            <Visibility />
          </div>
        )}
        <Schedule />
        <div className="my-5">
          <Label htmlFor={""}>Featured Image</Label>
          <CourseFeaturedImage
            onFileSelected={(file) => {
              setValue("featured_image", file, { shouldValidate: true });
            }}
          />
          {errors.featured_image && (
            <span className="text-red-500 text-sm py-1">
              {(errors.featured_image.message as string) || ""}
            </span>
          )}
        </div>
        <div>
          <Label htmlFor={""}>Intro Video</Label>
          <IntroVideo />
        </div>
        <div className="mt-5">
          <p className="text-sm mb-1 font-medium text-gray-600">
            Pricing Model
          </p>
          <div className="flex items-center gap-5">
            <Controller
              control={control}
              name="pricing_model"
              render={({ field: { onChange, value } }) => (
                <div className="flex items-center gap-2">
                  <RadioButton
                    id="free"
                    checked={value === "free"}
                    onChange={() => onChange("free")}
                  />
                  <label htmlFor="free" className="text-sm">
                    Free
                  </label>
                </div>
              )}
            />
            <Controller
              control={control}
              name="pricing_model"
              render={({ field: { onChange, value } }) => (
                <div className="flex items-center gap-2">
                  <RadioButton
                    id="paid"
                    checked={value === "paid"}
                    onChange={() => onChange("paid")}
                  />
                  <label htmlFor="paid" className="text-sm">
                    Paid
                  </label>
                </div>
              )}
            />
          </div>
        </div>
        {watch("pricing_model") === "paid" && (
          <div className="flex items-center gap-5 mt-3">
            <div>
              <p className="text-sm mb-1">Regular Price</p>
              <div className="border bg-white focus-within:border-primary flex items-center rounded-md">
                <span className="w-10 h-8 flex items-center justify-center border-r border-gray-300">
                  $
                </span>
                <input
                  type="text"
                  placeholder="0"
                  className="w-full outline-none text-sm px-3 py-1.5"
                  {...register("regular_price")}
                  onFocus={(e) => e.target.select()}
                />
              </div>
            </div>
            <div>
              <p className="text-sm mb-1">Sale Price</p>
              <div className="border bg-white focus-within:border-primary flex items-center rounded-md">
                <span className="w-10 h-8 flex items-center justify-center border-r border-gray-300">
                  $
                </span>
                <input
                  type="text"
                  placeholder="0"
                  className="w-full outline-none text-sm px-3 py-1.5"
                  {...register("sale_price")}
                  onFocus={(e) => e.target.select()}
                />
              </div>
            </div>
          </div>
        )}
        {/* <div className="mt-5">
            <Label htmlFor={""}>Categories</Label>
            <Categories />
          </div> */}
        <div className="mt-5">
          <Label>Tags</Label>
          <TagsSelect />
        </div>
        {/* <div className="mt-5">
          <Label htmlFor={""}>Author</Label>
          <Author />
        </div> */}
      </div>
    </div>
  );
}
