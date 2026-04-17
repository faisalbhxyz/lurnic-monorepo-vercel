import { Editor } from "@tiptap/react";
import React from "react";
import { BsTypeBold } from "react-icons/bs";
import { MdOutlineFormatListBulleted } from "react-icons/md";
import { LuItalic, LuUnderline } from "react-icons/lu";
// import { LuStrikethrough } from "react-icons/lu";
import { LuListOrdered } from "react-icons/lu";
import { IoCodeSharp } from "react-icons/io5";
import { LuHighlighter } from "react-icons/lu";
import { Tooltip } from "react-tooltip";
import Alignment from "./Alignment";
import MoreOptions from "./MoreOptions";
import Formatting from "./Formatting";
import FileUpload from "./FileUpload";
import InsertLink from "./InsertLink";
import InsertTable from "./InsertTable";
import ColorPicker from "./ColorPicker";

interface ToolBarProps {
  editor: Editor | null;
}

export default function ToolBar({ editor }: ToolBarProps) {
  if (!editor) return null;

  const Options = [
    {
      icon: <BsTypeBold className="size-5" />,
      onClick: () => editor.chain().focus().toggleBold().run(),
      isActive: editor.isActive("bold"),
      name: "Bold",
    },
    {
      icon: <LuItalic className="size-4" />,
      onClick: () => editor.chain().focus().toggleItalic().run(),
      isActive: editor.isActive("italic"),
      name: "Italic",
    },
    {
      icon: <LuUnderline className="size-4" />,
      onClick: () => editor.chain().focus().toggleUnderline().run(),
      isActive: editor.isActive("underline"),
      name: "Underline",
    },
    // {
    //   icon: <LuStrikethrough className="size-4" />,
    //   onClick: () => editor.chain().focus().toggleStrike().run(),
    //   isActive: editor.isActive("strike"),
    // },
  ];

  const Options2 = [
    {
      icon: <MdOutlineFormatListBulleted className="size-5" />,
      onClick: () => editor.chain().focus().toggleBulletList().run(),
      isActive: editor.isActive("bulletList"),
      name: "Bulleted List",
    },
    {
      icon: <LuListOrdered className="size-5" />,
      onClick: () => editor.chain().focus().toggleOrderedList().run(),
      isActive: editor.isActive("orderedList"),
      name: "Numbered List",
    },
    {
      icon: <IoCodeSharp className="size-4" />,
      onClick: () => editor.chain().focus().toggleCodeBlock().run(),
      isActive: editor.isActive("code"),
      name: "Code Block",
    },
    {
      icon: <LuHighlighter className="size-4" />,
      onClick: () => editor.chain().focus().toggleHighlight().run(),
      isActive: editor.isActive("highlight"),
      name: "Highlight",
    },
  ];

  return (
    <div className="rounded-md p-1.5 bg-slate-50 flex items-center gap-1">
      <Formatting editor={editor} />
      <div className="bg-gray-300 h-5 w-px" />
      {Options.map((option, i) => (
        <button
          type="button"
          key={i}
          data-tooltip-id={option.name}
          data-tooltip-content={option.name}
          onClick={option.onClick}
          className={`w-6 h-6 flex items-center justify-center text-gray-700 ${
            option.isActive ? "bg-gray-300" : "hover:bg-gray-200"
          } rounded`}
        >
          {option.icon}
          <Tooltip
            id={option.name}
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
        </button>
      ))}
      <ColorPicker editor={editor} />
      <div className="bg-gray-300 h-5 w-px" />
      <Alignment editor={editor} />
      <div className="bg-gray-300 h-5 w-px" />
      <div className="hidden md:flex items-center gap-1">
        <InsertLink editor={editor} />
        <FileUpload editor={editor} />
        <InsertTable editor={editor} />
        <div className="bg-gray-300 h-5 w-px" />
      </div>
      {Options2.map((option, i) => (
        <button
          type="button"
          key={i}
          onClick={option.onClick}
          data-tooltip-id={option.name}
          data-tooltip-content={option.name}
          className={`w-6 h-6 hidden md:flex items-center justify-center text-gray-700 ${
            option.isActive ? "bg-gray-300" : "hover:bg-gray-200"
          } rounded`}
        >
          {option.icon}
          <Tooltip
            id={option.name}
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
        </button>
      ))}
      <div className="md:hidden">
        <MoreOptions editor={editor} />
      </div>
    </div>
  );
}
