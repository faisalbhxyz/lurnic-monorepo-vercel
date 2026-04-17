import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { LuPlus } from "react-icons/lu";
import { FiBookOpen } from "react-icons/fi";
import { MdChecklist, MdOutlineFileCopy } from "react-icons/md";
import { useCoursesStore } from "@/hooks/useCoursesStore";

export default function AddItem({ id }: { id: number }) {
  const { addNewAssignment, addNewLesson, addNewQuiz, setChapterId } =
    useCoursesStore();

  return (
    <Menu>
      <MenuButton
        onClick={(e) => e.stopPropagation()}
        className="outline-none hover:bg-primary border border-primary text-primary hover:text-white text-sm p-1.5 rounded-md flex items-center gap-2 font-medium"
      >
        <LuPlus />
      </MenuButton>

      <MenuItems
        transition
        anchor="bottom end"
        className="w-40 origin-top-right rounded-md border bg-white text-gray-600 p-1 text-sm transition duration-100 ease-out [--anchor-gap:--spacing(1)] focus:outline-none data-closed:scale-95 data-closed:opacity-0"
      >
        <MenuItem>
          <button
            onClick={() => {
              addNewLesson();
              setChapterId(id);
            }}
            className="group flex w-full items-center gap-2 rounded-lg px-3 py-1.5 data-focus:bg-gray-100"
          >
            <FiBookOpen />
            Lesson
          </button>
        </MenuItem>
        <MenuItem>
          <button
            onClick={() => {
              addNewQuiz();
              setChapterId(id);
            }}
            className="group flex w-full items-center gap-2 rounded-lg px-3 py-1.5 data-focus:bg-gray-100"
          >
            <MdChecklist />
            Quiz
          </button>
        </MenuItem>
        <MenuItem>
          <button
            onClick={() => {
              addNewAssignment();
              setChapterId(id);
            }}
            className="group flex w-full items-center gap-2 rounded-lg px-3 py-1.5 data-focus:bg-gray-100"
          >
            <MdOutlineFileCopy />
            Assignment
          </button>
        </MenuItem>
      </MenuItems>
    </Menu>
  );
}
