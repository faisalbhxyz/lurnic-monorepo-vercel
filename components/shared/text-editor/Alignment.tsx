import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import {
  MdFormatAlignLeft,
  MdFormatAlignCenter,
  MdFormatAlignRight,
  MdKeyboardArrowDown,
} from "react-icons/md";

import { Editor } from "@tiptap/react";
import { Tooltip } from "react-tooltip";

const Alignment = ({ editor }: { editor: Editor | null }) => {
  if (!editor) return null;

  const Options = [
    {
      icon: <MdFormatAlignLeft className="size-4" />,
      onClick: () => editor.chain().focus().setTextAlign("left").run(),
      isActive: editor.isActive({ textAlign: "left" }),
    },
    {
      icon: <MdFormatAlignCenter className="size-4" />,
      onClick: () => editor.chain().focus().setTextAlign("center").run(),
      isActive: editor.isActive({ textAlign: "center" }),
    },
    {
      icon: <MdFormatAlignRight className="size-4" />,
      onClick: () => editor.chain().focus().setTextAlign("right").run(),
      isActive: editor.isActive({ textAlign: "right" }),
    },
  ];

  return (
    <Menu>
      <MenuButton
        data-tooltip-id="alignment"
        data-tooltip-content="Alignment"
        className="inline-flex items-center rounded pl-1 h-6 justify-center  text-sm/6 font-semibold focus:outline-none data-[hover]:bg-gray-200 data-[open]:bg-gray-300 data-[focus]:outline-1 data-[focus]:outline-white"
      >
        <MdFormatAlignLeft className="size-4 fill-gray-700" />
        <MdKeyboardArrowDown size={17} className="fill-gray-600" />
      </MenuButton>
      <Tooltip
        id="alignment"
        style={{
          backgroundColor: "white",
          color: "black",
          fontWeight: 500,
          border: "1px solid red",
          padding: "2px 8px ",
          borderRadius: "8px",
          boxShadow: "0 2px 4px rgba(0, 0, 0, 0.2)",
        }}
        className="border border-gray-300"
      />

      <MenuItems
        transition
        anchor="bottom end"
        className="origin-top-right z-10 rounded-md mt-1 border bg-white p-1 text-sm/6 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
      >
        {Options.map((option, i) => (
          <MenuItem key={i}>
            <button
              onClick={option.onClick}
              className={`group flex w-full items-center gap-2 rounded-md p-1.5 data-[focus]:bg-gray-100 ${
                option.isActive ? "bg-gray-300" : "hover:bg-gray-200"
              }`}
            >
              {option.icon}
            </button>
          </MenuItem>
        ))}
      </MenuItems>
    </Menu>
  );
};

export default Alignment;
