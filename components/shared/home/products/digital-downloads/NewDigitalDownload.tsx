"use client";

import Button from "@/components/ui/Button";
import { GenericSelectWithSearch } from "@/components/ui/GenericSelectWithSearch";
import InputField from "@/components/ui/InputField";
import Label from "@/components/ui/Label";
import Modal from "@/components/ui/Modal";
import React, { useState } from "react";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";

type Category = {
  id: number;
  name: string;
};

const category: Category[] = [
  { id: 1, name: "Arifin" },
  { id: 2, name: "Borhan" },
];

export default function NewDigitalDownload() {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(
    null
  );

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>
        <LuPlus /> New Digital Download
      </Button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <div className="flex items-start justify-between py-4 px-5 border-b border-gray-300">
          <div>
            <p className="font-medium text-lg">
              Describe your digital download
            </p>
            <p className="text-sm text-gray-500">
              Provide a title and summary for your digital download. This will
              help users understand what they will download.
            </p>
          </div>
          <button onClick={() => setIsOpen(false)} className="p-1">
            <RxCross2 />
          </button>
        </div>
        <div className="p-5">
          <div className="w-full mb-5">
            <Label>
              Title <span className="text-red-500">*</span>
            </Label>
            <InputField placeholder="e.g. Ebook Download" className="w-full" />
            <p className="text-sm text-gray-600">
              A short, concise title that describes the digital download.
            </p>
          </div>
          <div className="mb-5">
            <div className="flex items-center justify-between">
              <Label htmlFor="title">Summary</Label>
              <span className="text-gray-500 text-sm">50 characters</span>
            </div>
            <textarea
              id="course-summary"
              rows={4}
              className="bg-white border w-full min-h-20 rounded-md px-3 py-2 outline-none focus:border-primary"
              placeholder="Enter a summary for the digital download."
            />
            <p className="text-[13px] text-gray-500">
              A paragraph regarding your digital download.
            </p>
          </div>
          <div>
            <Label>
              Category <span className="text-red-500">*</span>
            </Label>
            <GenericSelectWithSearch
              items={category}
              selectedItem={selectedCategory}
              onSelect={setSelectedCategory}
              getLabel={(student) => student.name}
              className="w-full"
            />
          </div>
          <div className="flex items-center justify-end gap-3 mt-5">
            <button
              onClick={() => setIsOpen(false)}
              className="border text-sm font-medium px-4 py-2 rounded-full"
            >
              Cancel
            </button>
            <Button>Create</Button>
          </div>
        </div>
      </Modal>
    </>
  );
}
