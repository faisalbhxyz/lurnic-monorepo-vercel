import React, { useEffect, useState } from "react";
import { useFieldArray, useFormContext } from "react-hook-form";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import { TCourseSchema } from "@/schema/course.schema";
import InputField from "@/components/ui/InputField";
import Button from "@/components/ui/Button";

const accessOptions = [
  { label: "Draft", value: "draft" },
  { label: "Published", value: "published" },
];

export default function NewChapter({ isEdit = false }: { isEdit: boolean }) {
  const { chapterId, closeNewChapter, closeEditChapter } = useCoursesStore();

  const { control } = useFormContext<TCourseSchema>();
  const { fields, append, update } = useFieldArray({
    control,
    name: "course_chapters",
    keyName: "uid",
  });

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [access, setAccess] = useState<"draft" | "published">("draft");

  // For editing: find the matching chapter and populate state
  useEffect(() => {
    if (isEdit && chapterId) {
      const chapter = fields.find((c) => c["_id"] === chapterId);
      if (chapter) {
        setTitle(chapter.title || "");
        setDescription(chapter.description || "");
        setAccess(chapter.access || "draft");
      }
    }
  }, [chapterId, isEdit, fields]);

  const onSubmit = () => {
    if (!title.trim()) {
      alert("Chapter title is required.");
      return;
    }

    if (isEdit && chapterId) {
      const index = fields.findIndex((c) => c._id === chapterId);
      if (index !== -1) {
        update(index, {
          ...fields[index],
          title,
          description,
          access,
        });
      }
      closeEditChapter();
    } else {
      append({
        _id: Date.now(),
        title,
        description,
        access,
        position: fields.length,
      });
      setTitle("");
      setDescription("");
      setAccess("draft");
      closeNewChapter();
    }
  };

  return (
    <div className="p-4">
      {/* Title */}
      <div className="mb-3">
        <label htmlFor="title" className="block text-sm mb-1 font-medium">
          Chapter Title *
        </label>
        <InputField
          id="title"
          placeholder="Ex. Introduction"
          className="w-full"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>

      {/* Description */}
      <div className="mb-3">
        <label htmlFor="description" className="block text-sm mb-1 font-medium">
          Description
        </label>
        <textarea
          id="description"
          placeholder="Write here"
          className="w-full min-h-28 bg-white border border-gray-300 resize-none rounded-md text-sm px-3 py-1.5 focus:outline-none focus:ring-1 focus:ring-primary"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
      </div>

      {/* Access */}
      <div className="mb-3">
        <label htmlFor="access" className="block text-sm mb-1 font-medium">
          Access *
        </label>
        <div className="flex items-center gap-5">
          {accessOptions.map((option) => (
            <label
              key={option.value}
              htmlFor={option.value}
              className={`w-full border rounded-md p-3 flex items-center justify-between cursor-pointer ${
                access === option.value ? "border-primary" : "border-gray-300"
              }`}
            >
              <p>{option.label}</p>
              <div
                className={`rounded-full p-1.5 ${
                  access === option.value ? "bg-primary" : "bg-gray-200"
                }`}
              >
                <div className="w-2 h-2 rounded-full bg-white" />
              </div>
              <input
                type="radio"
                id={option.value}
                value={option.value}
                hidden
                checked={access === option.value}
                onChange={() =>
                  setAccess(option.value as "draft" | "published")
                }
              />
            </label>
          ))}
        </div>
      </div>

      {/* Buttons */}
      <div className="flex items-center justify-between mt-10">
        <button
          type="button"
          onClick={() => {
            isEdit ? closeNewChapter() : closeEditChapter();
            !isEdit && setTitle("");
            !isEdit && setDescription("");
            !isEdit && setAccess("draft");
          }}
          className="border px-3 py-1.5 rounded-md text-sm"
        >
          Cancel
        </button>
        <Button type="button" onClick={onSubmit}>
          {isEdit ? "Update Chapter" : "Create Chapter"}
        </Button>
      </div>
    </div>
  );
}
