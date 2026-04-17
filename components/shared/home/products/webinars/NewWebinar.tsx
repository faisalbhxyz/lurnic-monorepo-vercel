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

export default function NewWebinar() {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(
    null
  );

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>
        <LuPlus /> New Webinar
      </Button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <div className="flex items-start justify-between py-4 px-5 border-b border-gray-300">
          <div>
            <p className="font-medium text-lg">Describe your webinar</p>
            <p className="text-sm text-gray-500">
              Provide a title and summary for your webinar. This will help users
              understand what they will download.
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
            <InputField
              placeholder="e.g. Introduction to Web Design"
              className="w-full"
            />
            <p className="text-sm text-gray-600">
              A short, concise title that describes the webinar.
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
              placeholder="Enter a summary for the webinar."
            />
            <p className="text-[13px] text-gray-500">
              A paragraph regarding your webinar.
            </p>
          </div>
          <div className="w-full mb-5">
            <Label>
              Webinar Start Time <span className="text-red-500">*</span>
            </Label>
            <InputField type="datetime-local" className="w-full" />
            <p className="text-sm text-gray-600">
              When will this webinar start? Choose a date and time in the
              future.
            </p>
          </div>
          <div className="w-full mb-5">
            <Label>
              Registration End Time <span className="text-red-500">*</span>
            </Label>
            <InputField type="datetime-local" className="w-full" />
            <p className="text-sm text-gray-600">
              When will registration close? Choose a date and time before or
              equal to the webinar start time.
            </p>
          </div>
          <div>
            <Label>Webinar Timezone</Label>
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
