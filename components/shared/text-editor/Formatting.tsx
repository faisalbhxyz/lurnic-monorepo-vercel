import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { MdKeyboardArrowDown } from "react-icons/md";
import { Editor } from "@tiptap/react";
import { Tooltip } from "react-tooltip";

const Formatting = ({ editor }: { editor: Editor | null }) => {
  if (!editor) return null;

  const Options = [
    {
      text: "Heading 1",
      onClick: () => editor.chain().focus().toggleHeading({ level: 1 }).run(),
      isActive: editor.isActive("heading", { level: 1 }),
      sizeClass: "text-3xl font-bold",
    },
    {
      text: "Heading 2",
      onClick: () => editor.chain().focus().toggleHeading({ level: 2 }).run(),
      isActive: editor.isActive("heading", { level: 2 }),
      sizeClass: "text-2xl font-bold",
    },
    {
      text: "Heading 3",
      onClick: () => editor.chain().focus().toggleHeading({ level: 3 }).run(),
      isActive: editor.isActive("heading", { level: 3 }),
      sizeClass: "text-xl font-bold",
    },
  ];

  // Find the active option
  const activeOption = Options.find((option) => option.isActive);

  return (
    <Menu>
      <MenuButton
        data-tooltip-id="formatting"
        data-tooltip-content="Formatting"
        className="inline-flex items-center text-gray-600 rounded pl-1.5 h-6 justify-center text-sm font-medium focus:outline-none hover:bg-gray-200 open:bg-gray-300"
      >
        {activeOption ? (
          <span>{activeOption.text}</span>
        ) : (
          <span>Paragraph</span> // Default label when no active heading
        )}
        <MdKeyboardArrowDown size={17} className="fill-gray-600 ml-1" />
      </MenuButton>
      <Tooltip
        id="formatting"
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
        className="origin-top-right z-10 rounded-md mt-1 border bg-white p-1 text-sm transition duration-100 ease-out focus:outline-none"
      >
        {Options.map((option, i) => (
          <MenuItem key={i}>
            <button
              onClick={option.onClick}
              className={`group flex w-full items-center gap-2 rounded-md p-1.5 focus:bg-gray-100 ${
                option.isActive ? "bg-gray-300" : "hover:bg-gray-200"
              }`}
            >
              <span className={option.sizeClass}>{option.text}</span>
            </button>
          </MenuItem>
        ))}
      </MenuItems>
    </Menu>
  );
};

export default Formatting;
