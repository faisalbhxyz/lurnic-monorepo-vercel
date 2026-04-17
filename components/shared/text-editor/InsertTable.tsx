import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { MdKeyboardArrowDown, MdOutlineTableChart } from "react-icons/md";
import { Editor } from "@tiptap/react";
import { Tooltip } from "react-tooltip";

const InsertTable = ({ editor }: { editor: Editor | null }) => {
  if (!editor) return null;

  const isTablePresent = editor.can().deleteTable(); // Checks if a table exists

  const options = [
    {
      text: "Insert table",
      onClick: () =>
        editor
          .chain()
          .focus()
          .insertTable({ rows: 3, cols: 3, withHeaderRow: true })
          .run(),
      disabled: false,
    },
    {
      text: "Add column before",
      onClick: () => editor.chain().focus().addColumnBefore().run(),
      disabled: !isTablePresent,
    },
    {
      text: "Add column after",
      onClick: () => editor.chain().focus().addColumnAfter().run(),
      disabled: !isTablePresent,
    },
    {
      text: "Add row before",
      onClick: () => editor.chain().focus().addRowBefore().run(),
      disabled: !isTablePresent,
    },
    {
      text: "Add row after",
      onClick: () => editor.chain().focus().addRowAfter().run(),
      disabled: !isTablePresent,
    },
    {
      text: "Delete column",
      onClick: () => editor.chain().focus().deleteColumn().run(),
      disabled: !isTablePresent,
    },
    {
      text: "Delete row",
      onClick: () => editor.chain().focus().deleteRow().run(),
      disabled: !isTablePresent,
    },
    {
      text: "Delete table",
      onClick: () => editor.chain().focus().deleteTable().run(),
      disabled: !isTablePresent,
    },
  ];

  return (
    <Menu as="div" className="relative">
      <MenuButton
        data-tooltip-id="insert-table"
        data-tooltip-content="Insert table"
        className="inline-flex items-center rounded pl-1 h-6 justify-center  text-sm/6 font-semibold focus:outline-none data-[hover]:bg-gray-200 data-[open]:bg-gray-200 data-[focus]:outline-1 data-[focus]:outline-white"
      >
        <MdOutlineTableChart className="size-4 fill-gray-700" />
        <MdKeyboardArrowDown size={17} className="fill-gray-600" />
      </MenuButton>
      <Tooltip
        id="insert-table"
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
        className="md:origin-top-right z-10 rounded-md mt-1 border bg-white p-1 text-sm/6 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
      >
        {options.map((option, i) => (
          <div key={i}>
            <MenuItem>
              <button
                onClick={option.onClick}
                disabled={option.disabled}
                className={`group flex w-full items-center gap-2 rounded-md px-2 py-1 ${
                  option.disabled ? "opacity-50" : "hover:bg-gray-100"
                }`}
              >
                {option.text}
              </button>
            </MenuItem>

            {i === 0 && <hr className="my-1 border-gray-300" />}

            {i === options.length - 4 && (
              <hr className="my-1 border-gray-300" />
            )}
          </div>
        ))}
      </MenuItems>
    </Menu>
  );
};

export default InsertTable;
