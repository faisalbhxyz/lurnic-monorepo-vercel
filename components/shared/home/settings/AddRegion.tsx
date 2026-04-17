import React, { useState } from "react";
import Modal from "@/components/ui/Modal";
import { RxCross2 } from "react-icons/rx";
import Checkbox from "@/components/ui/Checkbox";
import Button from "@/components/ui/Button";

const countries = [
  {
    id: 1,
    name: "Afghanistan",
  },
  {
    id: 2,
    name: "Bangladesh",
  },
];

export default function AddRegion() {
  const [isOpen, setIsOpen] = useState(false);
  return (
    <>
      <button
        onClick={() => setIsOpen(true)}
        className="bg-primary/10 text-primary px-3 py-1.5 rounded-md"
      >
        Add Region
      </button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <div className="flex items-center justify-between px-4 py-3 border-b border-gray-300">
          <p className="text-xl font-medium">Add tax region</p>
          <button onClick={() => setIsOpen(false)}>
            <RxCross2 />
          </button>
        </div>
        <div className="p-4 overflow-y-auto space-y-4">
          {countries.map((item) => (
            <div key={item.id} className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Checkbox />{" "}
                <label htmlFor="" className="text-sm font-medium">
                  {item.name}
                </label>
              </div>
              <div></div>
            </div>
          ))}
        </div>
        <div className=" bg-white rounded-b-2xl flex items-center justify-between px-4 py-3 border-t border-gray-300">
          <Button onClick={() => setIsOpen(false)} variant="secondary">
            Cancel
          </Button>
          <Button>Apply</Button>
        </div>
      </Modal>
    </>
  );
}
