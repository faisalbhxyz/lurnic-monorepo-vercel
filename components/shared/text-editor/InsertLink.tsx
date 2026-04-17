"use client";

import React, { useState } from "react";
import { Editor } from "@tiptap/react";
import { HiLink } from "react-icons/hi";
import { Tooltip } from "react-tooltip";
import Modal from "@/components/ui/Modal";
import InputField from "@/components/ui/InputField";

const selectWindow = [
  { id: 1, name: "the same window" },
  { id: 2, name: "a new window" },
];

export default function InsertLink({ editor }: { editor: Editor | null }) {
  const [isOpen, setIsOpen] = useState(false);
  const [url, setUrl] = useState("");
  const [linkTitle, setLinkTitle] = useState("");
  const [selectedWindow, setSelectedWindow] = useState(selectWindow[1]);

  if (!editor) return null;

  const handleInsertLink = () => {
    if (!url) return;

    editor
      .chain()
      .focus()
      .extendMarkRange("link")
      .setLink({
        href: url,
        target: selectedWindow.id === 2 ? "_blank" : "_self",
      })
      .run();

    setUrl("");
    setLinkTitle("");
    setIsOpen(false);
  };

  const isTextSelected = !editor.state.selection.empty;

  return (
    <>
      <button
        type="button"
        onClick={() => setIsOpen(true)}
        disabled={!isTextSelected}
        data-tooltip-id="insert-link"
        data-tooltip-content="Insert link"
        className="w-6 h-6 disabled:opacity-50 flex items-center justify-center text-gray-700 hover:bg-gray-200 disabled:hover:bg-transparent rounded"
      >
        <HiLink size={18} />
      </button>

      {isTextSelected && (
        <Tooltip
          id="insert-link"
          style={{
            backgroundColor: "white",
            color: "black",
            fontWeight: 500,
            border: "1px solid red",
            padding: "2px 8px",
            borderRadius: "8px",
            boxShadow: "0 2px 4px rgba(0, 0, 0, 0.2)",
          }}
          className="border border-gray-300"
        />
      )}

      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)}>
        <div className="p-4 space-y-3">
          <InputField
            placeholder="http://"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            className="w-full"
          />
          <p className="text-sm text-gray-600">
            http:// is required for external links
          </p>

          <InputField
            placeholder="Link title"
            value={linkTitle}
            onChange={(e) => setLinkTitle(e.target.value)}
            className="w-full"
          />
          <p className="text-sm text-gray-600">
            Used for accessibility and SEO
          </p>

          <select
            className="w-full border rounded-md px-3 py-2 mt-2"
            value={selectedWindow.id}
            onChange={(e) =>
              setSelectedWindow(
                selectWindow.find(
                  (opt) => opt.id === parseInt(e.target.value)
                ) || selectWindow[0]
              )
            }
          >
            {selectWindow.map((option) => (
              <option key={option.id} value={option.id}>
                Open in {option.name}
              </option>
            ))}
          </select>

          <div className="mt-4 flex justify-end gap-2">
            <button
              onClick={() => setIsOpen(false)}
              className="px-4 py-2 text-sm bg-gray-100 rounded hover:bg-gray-200"
            >
              Cancel
            </button>
            <button
              onClick={handleInsertLink}
              className="px-4 py-2 text-sm bg-primary text-white rounded hover:bg-blue-700"
            >
              Insert Link
            </button>
          </div>
        </div>
      </Modal>
    </>
  );
}
