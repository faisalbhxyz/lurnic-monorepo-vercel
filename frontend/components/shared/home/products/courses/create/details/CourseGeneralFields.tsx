"use client";

import React from "react";
import SelectList from "@/components/ui/SelectList";
import InputField from "@/components/ui/InputField";
import AddInstructor from "../settings/AddInstructor";
import MultiSelect from "@/components/ui/MultiSelect";
import { Controller, useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";
import Label from "@/components/ui/Label";

const DifficultyLevels = [
  { id: 1, value: "all", name: "All" },
  { id: 2, value: "beginner", name: "Beginner" },
  { id: 3, value: "intermediate", name: "Intermediate" },
  { id: 4, value: "expert", name: "Expert" },
];

export default function CourseGeneralFields({
  categories,
  subcategories,
  instructors,
}: {
  categories: ICategory[] | null;
  subcategories: ISubCategory[] | null;
  instructors: IInstructor[] | null;
}) {
  const {
    register,
    control,
    setValue,
    watch,
    formState: { errors },
  } = useFormContext<TCourseSchema>();

  const watchCategoryID = watch("general_settings.category_id");

  return (
    <div className="mt-5 space-y-5">
      <div>
        <div className="flex items-center justify-between mb-1">
          <Label>Instructors</Label>
          <AddInstructor />
        </div>
        <Controller
          control={control}
          name="course_instructors"
          render={({ field: { onChange, value = [] } }) => {
            const selected =
              instructors && instructors?.length > 0
                ? instructors
                    .filter((instructor) => value.includes(instructor.id))
                    .map((instructor) => ({
                      id: instructor.id,
                      name: `${instructor.first_name} ${
                        instructor.last_name ?? ""
                      }`,
                    }))
                : [];

            return (
              <MultiSelect
                options={
                  instructors && instructors?.length > 0
                    ? instructors.map((instructor) => ({
                        id: instructor.id,
                        name: `${instructor.first_name} ${instructor.last_name}`,
                      }))
                    : []
                }
                selected={selected}
                onChange={(data) => onChange(data.map((item) => item.id))}
                placeholder="Select options..."
              />
            );
          }}
        />
        {errors?.course_instructors?.message && (
          <p className="text-sm text-red-500 mt-1">
            {errors?.course_instructors?.message}
          </p>
        )}
      </div>

      <div>
        <Label>Difficulty Level</Label>
        <Controller
          control={control}
          name="general_settings.difficulty_level"
          render={({ field: { onChange, value } }) => (
            <SelectList
              options={DifficultyLevels}
              value={DifficultyLevels.find((item) => item.value === value)}
              onChange={(val) => onChange(val.value)}
              placeholder="Select option..."
              className="w-full"
            />
          )}
        />
        <p className="text-[13px] text-gray-500 mt-1">
          Pick a difficulty level for this course!
        </p>
      </div>

      <div>
        <Label>Maximum Students</Label>
        <InputField
          className="w-full"
          {...register("general_settings.maximum_student")}
          onFocus={(e) => e.target.select()}
          error={errors?.general_settings?.maximum_student?.message}
        />
        <p className="text-[13px] text-gray-500 mt-1">Leave 0 for unlimited.</p>
      </div>

      <div>
        <Label>Language</Label>
        <Controller
          control={control}
          name="general_settings.language"
          defaultValue={"english"}
          render={({ field: { onChange, value } }) => (
            <SelectList
              options={[{ id: 1, name: "English", value: "english" }]}
              value={{ id: 1, name: "English", value: "english" }}
              onChange={(val) => onChange(val.value)}
              placeholder="Select option..."
              className="w-full"
            />
          )}
        />
      </div>

      <div>
        <Label>
          Category <span className="text-red-500">*</span>
        </Label>
        <Controller
          control={control}
          name="general_settings.category_id"
          render={({ field }) => {
            const selected = categories
              ?.map((item) => ({
                id: item.id,
                name: item.name,
                value: String(item.id),
              }))
              .find((item) => item.id === field.value);

            return (
              <SelectList
                options={
                  categories?.map((item) => ({
                    id: item.id,
                    name: item.name,
                    value: String(item.id),
                  })) || []
                }
                value={selected}
                onChange={(val) => {
                  field.onChange(Number(val.id));
                  setValue("general_settings.sub_category_id", null);
                }}
                placeholder="Select option..."
                className="w-full"
              />
            );
          }}
        />
        {errors?.general_settings?.category_id?.message && (
          <p className="text-sm text-red-500 mt-1">
            {errors?.general_settings?.category_id?.message}
          </p>
        )}
      </div>

      <div>
        <Label>Sub Category</Label>
        <Controller
          control={control}
          name="general_settings.sub_category_id"
          render={({ field }) => {
            const selected = subcategories
              ?.filter((item) => item.category.id === watchCategoryID)
              ?.map((item) => ({
                id: item.id,
                name: item.name,
                value: String(item.id),
              }))
              .find((item) => item.id === field.value);

            return (
              <SelectList
                options={
                  subcategories
                    ?.filter((item) => item.category.id === watchCategoryID)
                    .map((item) => ({
                      id: item.id,
                      name: item.name,
                      value: String(item.id),
                    })) || []
                }
                value={selected}
                onChange={(val) => field.onChange(Number(val.id))}
                placeholder="Select option..."
                className="w-full"
              />
            );
          }}
        />
        {errors?.general_settings?.sub_category_id?.message && (
          <p className="text-sm text-red-500 mt-1">
            {errors?.general_settings?.sub_category_id?.message}
          </p>
        )}
      </div>

      <div>
        <Label>Course Duration</Label>
        <InputField
          className="w-full"
          {...register("general_settings.duration")}
          onFocus={(e) => e.target.select()}
          error={errors?.general_settings?.duration?.message}
        />
        <p className="text-[13px] text-gray-500 mt-1">
          Specify the total duration of the course.
        </p>
      </div>
    </div>
  );
}
