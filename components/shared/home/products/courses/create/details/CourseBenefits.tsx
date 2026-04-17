import { TCourseSchema } from "@/schema/course.schema";
import React from "react";
import { useFormContext } from "react-hook-form";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";

export default function CourseBenefits() {
  const { setValue, watch, register } = useFormContext<TCourseSchema>();
  const overview = watch("overview") || [];

  const handleCreate = () => {
    setValue("overview", [...overview, ""]);
  };

  const handleRemove = (indexToRemove: number) => {
    const updated = overview.filter((_, index) => index !== indexToRemove);
    setValue("overview", updated);
  };

  return (
    <>
      <div className="space-y-3">
        {overview.map((_, index) => (
          <div
            key={index}
            className="bg-white flex items-center border rounded-md px-2 focus-within:border-primary"
          >
            <input
              type="text"
              {...register(`overview.${index}`)}
              className="w-full py-2 outline-none text-sm"
              placeholder="Exciting Benefits..."
            />
            <button
              type="button"
              onClick={() => handleRemove(index)}
              className="p-1 ml-2 rounded-full bg-gray-200 hover:bg-gray-300"
            >
              <RxCross2 size={12} />
            </button>
          </div>
        ))}
      </div>

      <button
        type="button"
        onClick={handleCreate}
        className="flex items-center gap-2 border text-primary font-medium mt-5 px-3 py-1.5 text-sm rounded-md"
      >
        <LuPlus /> Add more
      </button>
    </>
  );
}
