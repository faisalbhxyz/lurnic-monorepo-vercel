import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import {
  MdOutlineFileUpload,
  MdOutlineFormatListBulleted,
} from "react-icons/md";
import { LuListOrdered } from "react-icons/lu";
import { IoCodeSharp } from "react-icons/io5";
import { LuHighlighter } from "react-icons/lu";

import { Editor } from "@tiptap/react";
import { HiOutlineDotsHorizontal } from "react-icons/hi";
import FileUpload from "./FileUpload";
import InsertTable from "./InsertTable";
import InsertLink from "./InsertLink";

const MoreOptions = ({ editor }: { editor: Editor | null }) => {
  if (!editor) return null;

  const addImage = () => {
    const url = window.prompt("URL");
    if (url) {
      editor.chain().focus().setImage({ src: url }).run();
    }
  };

  const Options = [
    {
      icon: <MdOutlineFormatListBulleted className="size-5" />,
      onClick: () => editor.chain().focus().toggleBulletList().run(),
      isActive: editor.isActive("bulletList"),
    },
    {
      icon: <LuListOrdered className="size-5" />,
      onClick: () => editor.chain().focus().toggleOrderedList().run(),
      isActive: editor.isActive("orderedList"),
    },
    {
      icon: <IoCodeSharp className="size-4" />,
      onClick: () => editor.chain().focus().toggleCodeBlock().run(),
      isActive: editor.isActive("code"),
    },
    {
      icon: <LuHighlighter className="size-4" />,
      onClick: () => editor.chain().focus().toggleHighlight().run(),
      isActive: editor.isActive("highlight"),
    },
    {
      icon: <MdOutlineFileUpload className="size-5" />,
      onClick: () => addImage(),
      isActive: editor.isActive("image"),
    },
  ];

  return (
    <Menu>
      <MenuButton className="inline-flex items-center rounded px-1 h-6 justify-center  text-sm/6 font-semibold focus:outline-none data-[hover]:bg-gray-200 data-[open]:bg-gray-300 data-[focus]:outline-1 data-[focus]:outline-white">
        <HiOutlineDotsHorizontal size={17} className="fill-gray-600" />
      </MenuButton>

      <MenuItems
        transition
        anchor="bottom end"
        className="origin-top-right z-10 flex flex-col rounded-md mt-1 border bg-white p-1 text-sm/6 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
      >
        <div className="flex items-center gap-1">
          <InsertLink editor={editor} />
          <FileUpload editor={editor} />
          <InsertTable editor={editor} />
        </div>
        <div className="flex items-center gap-1">
          {Options.map((option, i) => (
            <MenuItem key={i}>
              <button
                onClick={option.onClick}
                className={`group gap-2 rounded-md w-6 min-w-6 h-6 flex items-center justify-center text-gray-700 ${
                  option.isActive ? "bg-gray-300" : "data-[focus]:bg-gray-200"
                }`}
              >
                {option.icon}
              </button>
            </MenuItem>
          ))}
        </div>
      </MenuItems>
    </Menu>
  );
};

export default MoreOptions;
