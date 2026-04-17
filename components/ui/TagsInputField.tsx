"use client";

import React, { useState, useEffect } from "react";

interface TagInputFieldProps {
  value?: string | null; // semicolon-separated string
  onChange: (value: string | null) => void;
  placeholder?: string;
}

export default function TagInputField({
  value = "",
  onChange,
  placeholder = "Add tags separated by ;",
}: TagInputFieldProps) {
  const [inputValue, setInputValue] = useState("");
  const [tags, setTags] = useState<string[]>([]);

  useEffect(() => {
    const initialTags = value
      ? value
          .split(";")
          .map((t) => t.trim())
          .filter(Boolean)
      : [];
    setTags(initialTags);
  }, [value]);

  const addTagsFromInput = (val: string) => {
    const parts = val
      .split(";")
      .map((t) => t.trim())
      .filter(Boolean);
    const newTags = [...tags, ...parts.filter((t) => !tags.includes(t))];
    setTags(newTags);
    setInputValue("");
    onChange(newTags.length > 0 ? newTags.join(";") : null);
  };

  const removeTag = (tagToRemove: string) => {
    const updatedTags = tags.filter((t) => t !== tagToRemove);
    setTags(updatedTags);
    onChange(updatedTags.length > 0 ? updatedTags.join(";") : null);
  };

  return (
    <div className="relative w-full">
      <input
        type="text"
        placeholder={placeholder}
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onKeyDown={(e) => {
          if (e.key === "Enter" || e.key === ";") {
            e.preventDefault();
            if (inputValue.trim()) addTagsFromInput(inputValue);
          }
        }}
        className="w-full bg-white border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none"
      />

      {tags.length > 0 && (
        <div className="flex flex-wrap gap-2 mt-2">
          {tags.map((tag, idx) => (
            <div
              key={idx}
              className="bg-blue-100 text-blue-800 px-3 py-[0.2rem] rounded-full flex items-center"
            >
              {tag}
              <button
                type="button"
                onClick={() => removeTag(tag)}
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
