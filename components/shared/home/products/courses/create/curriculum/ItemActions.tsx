import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { BsThreeDotsVertical } from "react-icons/bs";
import { LuTrash } from "react-icons/lu";
import { AiOutlineEdit } from "react-icons/ai";
import { LiaFillDripSolid } from "react-icons/lia";
import { useCoursesStore } from "@/hooks/useCoursesStore";
import { useFieldArray, useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";
import { toast } from "sonner";

export default function ItemActions({
  chapterID,
  itemID,
  itemType,
}: {
  chapterID: number;
  itemID: number;
  itemType: "lesson" | "quiz" | "assignment";
}) {
  const { openEditLesson, setChapterId, openEditAssignment, openEditQuiz } =
    useCoursesStore();

  const { control, watch } = useFormContext<TCourseSchema>();

  const chapterIndex = watch("course_chapters", []).findIndex(
    (chapter) => chapter._id === chapterID
  );
  const safeChapterIndex = chapterIndex === -1 ? 0 : chapterIndex;

  // one fieldArray per type
  const lessonsFieldArray = useFieldArray({
    control,
    name: `course_chapters.${safeChapterIndex}.course_lessons`,
    keyName: "uid",
  });
  const quizzesFieldArray = useFieldArray({
    control,
    name: `course_chapters.${safeChapterIndex}.quizzes`,
    keyName: "uid",
  });
  const assignmentsFieldArray = useFieldArray({
    control,
    name: `course_chapters.${safeChapterIndex}.assignments`,
    keyName: "uid",
  });

  const handleDeleteItem = () => {
    let items: any[] = [];
    let removeFn: ((index: number) => void) | null = null;

    switch (itemType) {
      case "lesson":
        items =
          watch(`course_chapters.${safeChapterIndex}.course_lessons`) || [];
        removeFn = lessonsFieldArray.remove;
        break;
      case "quiz":
        items = watch(`course_chapters.${safeChapterIndex}.quizzes`) || [];
        removeFn = quizzesFieldArray.remove;
        break;
      case "assignment":
        items = watch(`course_chapters.${safeChapterIndex}.assignments`) || [];
        removeFn = assignmentsFieldArray.remove;
        break;
      default:
        console.error("Unknown itemType", itemType);
        return;
    }

    const index = items.findIndex((item) => item._id === itemID);

    if (index !== -1 && removeFn) {
      removeFn(index);
    } else {
      toast.error("Failed to delete item.");
    }
  };

  const handleEditItem = () => {
    switch (itemType) {
      case "lesson":
        setChapterId(chapterID);
        openEditLesson(itemID);
        break;
      case "quiz":
        setChapterId(chapterID);
        openEditQuiz(itemID);
        break;
      case "assignment":
        setChapterId(chapterID);
        openEditAssignment(itemID);
        break;
      default:
        toast.error("Something went wrong.");
        return;
    }
  };

  return (
    <Menu>
      <MenuButton className="p-1.5 rounded-md outline-none">
        <BsThreeDotsVertical />
      </MenuButton>

      <MenuItems
        transition
        anchor="bottom end"
        className="w-48 origin-top-right rounded-md border bg-white text-gray-600 p-1 text-sm font-medium transition duration-100 ease-out [--anchor-gap:--spacing(1)] focus:outline-none data-closed:scale-95 data-closed:opacity-0"
      >
        <MenuItem>
          <button
            type="button"
            onClick={handleEditItem}
            className="group flex w-full items-center gap-2 rounded-lg px-3 py-1.5 data-focus:bg-gray-100"
          >
            <AiOutlineEdit size={18} />
            Edit
          </button>
        </MenuItem>
        {/* <MenuItem>
          <button
            onClick={() => openDripSettings(item.type)}
            className="group flex w-full items-center gap-2 rounded-lg px-3 py-1.5 data-focus:bg-gray-100"
          >
            <LiaFillDripSolid size={18} />
            Drip Settings
          </button>
        </MenuItem> */}
        <MenuItem>
          <button
            type="button"
            className="group flex w-full items-center gap-2 rounded-lg px-3 py-1.5 data-focus:bg-gray-100 text-red-500"
            onClick={handleDeleteItem}
          >
            <LuTrash />
            Delete
          </button>
        </MenuItem>
      </MenuItems>
    </Menu>
  );
}
