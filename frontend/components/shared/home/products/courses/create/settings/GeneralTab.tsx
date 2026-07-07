"use client";

import React, { useState } from "react";
import { IoIosArrowDown, IoIosArrowUp } from "react-icons/io";

export default function GeneralTab() {
  const [isOpen, setIsOpen] = useState(true);

  return (
    <div className="max-w-md">
      <div className="border p-5 rounded-md">
        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className={`flex items-center justify-between w-full ${
            isOpen ? "text-primary" : ""
          }`}
        >
          <p className="font-semibold text-start">
            Why should I mention the Language
          </p>
          {isOpen ? <IoIosArrowUp /> : <IoIosArrowDown />}
        </button>

        {isOpen && (
          <p className="text-sm mt-5 text-gray-600">
            It is important to mention the language of the course content so
            that the students who understand the language can enroll in the
            course. You can set the language in the Details tab.
          </p>
        )}
      </div>
    </div>
  );
}
