import Button from "@/components/ui/Button";
import SelectList from "@/components/ui/SelectList";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import React, { useState } from "react";

const dripType = [
  {
    id: 1,
    name: "Specific Date",
  },
  {
    id: 2,
    name: "Day after Enrollment",
  },
];

export default function LessonDrip() {
  const [selectDripType, setSelectDripType] = useState(dripType[0]);

  return (
    <>
      <p className="font-semibold">Drip</p>
      <p className="text-sm text-gray-500">
        Release lessons on specific dates or a specific number of days after a
        enrollment
      </p>
      <p className="font-medium mt-3">Lesson</p>
      <p className="text-sm">Lesson Title</p>
      <div className="my-3">
        <label htmlFor="" className="block font-medium text-sm mb-1">
          Drip Type
        </label>
        <SelectList
          options={dripType}
          value={selectDripType}
          onChange={setSelectDripType}
          className="w-full"
        />
      </div>
      <div className="flex items-center gap-2">
        <ToggleSwitch />
        <label htmlFor="">Active Drip</label>
      </div>
      <div className="flex items-center justify-between mt-5">
        <button className="text-sm border-gray-300 px-3 border py-1.5 rounded-full">
          Cancel
        </button>
        <Button>Save</Button>
      </div>
    </>
  );
}
