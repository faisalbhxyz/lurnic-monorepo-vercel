import React, { useState } from "react";
import { MdKeyboardArrowRight } from "react-icons/md";

interface AccordionItem {
  title: string;
  description: string;
}

export default function CourseInfo() {
  const accordingData: AccordionItem[] = [
    {
      title: "Title",
      description:
        "Wireframing outlines the basic structure and layout of a design, serving as a visual guide before detailed development.",
    },
    {
      title: "Summary",
      description:
        "User-centered design ensures products meet the needs and preferences of the end-users, enhancing usability and satisfaction.",
    },
  ];

  const [isPlusAccording, setIsPlusAccording] = useState<number | null>(null);

  const handleBorderClick = (index: number) => {
    setIsPlusAccording((prevIndex) => (prevIndex === index ? null : index));
  };

  return (
    <div>
      <p className="mb-3 font-medium">Course Content</p>

      <div className="flex gap-3 flex-col w-full">
        {accordingData.map((according, index) => (
          <article
            key={index}
            className="border border-[#e5eaf2] bg-white rounded p-3"
          >
            <div
              className="flex gap-2 cursor-pointer items-center justify-between w-full"
              onClick={() => handleBorderClick(index)}
            >
              <h2 className="text-gray-700 font-medium text-sm">
                {according.title}
              </h2>
              <p>
                <MdKeyboardArrowRight
                  className={`text-[1.3rem] text-text transition-all duration-300 ${
                    isPlusAccording === index &&
                    "rotate-[90deg] !text-[#3B9DF8]"
                  }`}
                />
              </p>
            </div>
            <div
              className={`grid transition-all duration-300 overflow-hidden ease-in-out ${
                isPlusAccording === index
                  ? "grid-rows-[1fr] opacity-100 mt-4"
                  : "grid-rows-[0fr] opacity-0"
              }`}
            >
              <p className="text-[#424242] text-[0.9rem] overflow-hidden">
                {according.description}
              </p>
            </div>
          </article>
        ))}
      </div>
    </div>
  );
}
