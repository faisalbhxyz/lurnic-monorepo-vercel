"use client";

import React, { useRef, useState } from "react";
import { TiImage } from "react-icons/ti";
import { Editor } from "@tiptap/react";
import Button from "../../ui/Button";
import { Tooltip } from "react-tooltip";
import Modal from "../../ui/Modal";
import Image from "next/image";

export default function FileUpload({ editor }: { editor: Editor | null }) {
  const [isOpen, setIsOpen] = useState(false);
  const [imageUrl, setImageUrl] = useState<string>("");
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [isDragging, setIsDragging] = useState(false);

  if (!editor) return null;

  const handleUploadFile = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click();
    }
  };

  const addImageFromFile = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];

    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        const src = reader.result as string;
        setImageUrl(src);
        editor.chain().focus().setImage({ src }).run(); // <-- Insert into editor
      };
      reader.readAsDataURL(file);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const file = e.dataTransfer.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        const src = reader.result as string;
        setImageUrl(src);
        editor.chain().focus().setImage({ src }).run(); // <-- Insert into editor
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <>
      <button
        type="button"
        onClick={() => setIsOpen(true)}
        data-tooltip-id="insert-image"
        data-tooltip-content="Insert image"
        className="w-6 h-6 flex items-center justify-center text-gray-700 hover:bg-gray-200 rounded"
      >
        <TiImage size={18} />
      </button>
      <Tooltip
        id="insert-image"
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
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <div className="px-4 flex-1 p-2 overflow-auto flex items-center flex-col">
          {imageUrl && (
            <Image
              src={imageUrl}
              alt={"file"}
              width={200}
              height={200}
              className="rounded-xl w-full h-56 object-cover"
            />
          )}
          <div
            onClick={handleUploadFile}
            className={`flex flex-col items-center border border-dashed border-gray-600 hover:bg-gray-50 cursor-pointer rounded-lg w-full mt-5 p-10 ${
              isDragging ? "border-primary bg-primary/10" : "border-gray-400"
            }`}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
          >
            <div className="flex items-center gap-5">
              <Button
                type="button"
                className="font-medium bg-blue-400 text-white"
              >
                Upload image
              </Button>
              <input
                type="file"
                ref={fileInputRef}
                accept="image/*"
                hidden
                onChange={addImageFromFile}
              />
            </div>
            <p className="text-gray-600 text-sm font-medium mt-2">
              Drag and drop images here
            </p>
          </div>
        </div>
      </Modal>
    </>
  );
}
