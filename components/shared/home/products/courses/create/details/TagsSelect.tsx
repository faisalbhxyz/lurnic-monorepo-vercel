import React, { useState } from "react";
import { useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";

export default function TagsSelect() {
  const { setValue, watch } = useFormContext<TCourseSchema>();
  const tags = watch("tags") || [];
  const [searchValue, setSearchValue] = useState<string>("");

  const handleAddTag = () => {
    const trimmed = searchValue.trim();
    if (!trimmed || tags.includes(trimmed)) return;

    setValue("tags", [...tags, trimmed]);
    setSearchValue("");
  };

  const handleRemoveTag = (tagToRemove: string) => {
    const updatedTags = tags.filter((tag) => tag !== tagToRemove);
    setValue("tags", updatedTags);
  };

  return (
    <div className="relative w-full">
      <input
        type="text"
        placeholder="Add tags"
        value={searchValue}
        onChange={(e) => setSearchValue(e.target.value)}
        onKeyDown={(e) => {
          if ((e.key === "Enter" || e.key === ",") && searchValue.trim()) {
            e.preventDefault();
            handleAddTag();
          }
        }}
        className="w-full bg-white border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none"
      />

      {tags.length > 0 && (
        <div className="flex items-center gap-2 flex-wrap mt-2">
          {tags.map((tag, index) => (
            <div
              key={index}
              className="bg-blue-100 text-blue-800 px-3 py-[0.2rem] rounded-full flex items-center"
            >
              {tag}
              <button
                type="button"
                onClick={() => handleRemoveTag(tag)}
                className="ml-2 text-blue-800"
              >
                ×
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
