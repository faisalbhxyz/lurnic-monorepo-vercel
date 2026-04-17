"use client";

import React, { useEffect, useState } from "react";
import { Menu, MenuButton, MenuItems } from "@headlessui/react";
import { Editor } from "@tiptap/react";
import { SketchPicker, ColorResult } from "react-color";
import { RiFontColor } from "react-icons/ri";
import { MdKeyboardArrowDown } from "react-icons/md";
import { Tooltip } from "react-tooltip";

interface ColorPickerProps {
  editor: Editor | null;
}

const ColorPicker: React.FC<ColorPickerProps> = ({ editor }) => {
  const [color, setColor] = useState<string>("#000");

  useEffect(() => {
    if (editor) {
      const currentColor = editor.getAttributes("textStyle").color || "#000";
      setColor(currentColor);
    }
  }, [editor]);

  if (!editor) return null;

  const handleColorChange = (color: ColorResult) => {
    setColor(color.hex);
    editor.chain().focus().setColor(color.hex).run();
  };

  return (
    <Menu>
      <MenuButton
        data-tooltip-id="color"
        data-tooltip-content="Color"
        className="inline-flex items-center rounded pl-1 h-6 justify-center  text-sm/6 font-semibold focus:outline-none data-[hover]:bg-gray-200 data-[open]:bg-gray-300 data-[focus]:outline-1 data-[focus]:outline-white"
      >
        <RiFontColor className="size-4 fill-gray-700" />
        <MdKeyboardArrowDown size={17} className="fill-gray-600" />
      </MenuButton>
      <Tooltip
        id="color"
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
        anchor="bottom"
        className="origin-top z-10 rounded-md mt-1 border bg-white text-sm/6 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
      >
        <SketchPicker
          color={color}
          onChangeComplete={handleColorChange}
          data-testid="setColor"
        />
      </MenuItems>
    </Menu>
  );
};

export default ColorPicker;
