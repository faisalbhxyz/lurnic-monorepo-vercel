import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { BsThreeDotsVertical } from "react-icons/bs";
import { LuTrash } from "react-icons/lu";
import { AiOutlineEdit } from "react-icons/ai";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import { useFieldArray, useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";

export default function ChapterActions({ id }: { id: number }) {
  const { openEditChapter, setChapterId } = useCoursesStore();

  const { control, watch } = useFormContext<TCourseSchema>();
  const { remove } = useFieldArray({
    control,
    name: "course_chapters",
    keyName: "uid",
  });

  const handleEditChapter = () => {
    openEditChapter();
    setChapterId(id);
  };

  const handleDeleteChapter = () => {
    const index = watch("course_chapters").findIndex((chapter) => chapter._id === id);
    if (index !== -1) remove(index);
  };

  return (
    <Menu>
      <MenuButton className="p-1.5 rounded-md outline-none">
        <BsThreeDotsVertical />
      </MenuButton>

      <MenuItems
        transition
        anchor="bottom end"
        className="w-44 origin-top-right rounded-md border bg-white text-gray-600 p-1 text-sm font-medium transition duration-100 ease-out [--anchor-gap:--spacing(1)] focus:outline-none data-closed:scale-95 data-closed:opacity-0"
      >
        <MenuItem>
          <button
            type="button"
            onClick={handleEditChapter}
            className="group flex w-full items-center gap-2 rounded-lg px-3 py-1.5 data-focus:bg-gray-100"
          >
            <AiOutlineEdit />
            Edit
          </button>
        </MenuItem>
        <MenuItem>
          <button
            type="button"
            className="group flex w-full items-center gap-2 rounded-lg px-3 py-1.5 data-focus:bg-gray-100 text-red-500"
            onClick={handleDeleteChapter}
          >
            <LuTrash />
            Delete Chapters
          </button>
        </MenuItem>
      </MenuItems>
    </Menu>
  );
}
