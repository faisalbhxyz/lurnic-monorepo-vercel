import Checkbox from "@/components/ui/Checkbox";
import InputField from "@/components/ui/InputField";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import { cn } from "@/lib/cn";
import { pageFeatures, presetColors } from "@/lib/constants";
import Image from "next/image";
import React, { useState } from "react";
import { IoMdRefresh } from "react-icons/io";
import { MdCheck } from "react-icons/md";

const COLUMN_OPTIONS = [
  { id: 1, label: "One", count: 1 },
  { id: 2, label: "Two", count: 2 },
  { id: 3, label: "Three", count: 3 },
  { id: 4, label: "Four", count: 4 },
];

type ProfileLayout = "private" | "modern" | "minimal" | "classic";

const profileLayoutOptions: {
  key: ProfileLayout;
  label: string;
  imageSrc: string;
}[] = [
  { key: "private", label: "Private", imageSrc: "/images/profile-private.svg" },
  { key: "modern", label: "Modern", imageSrc: "/images/profile-modern.svg" },
  { key: "minimal", label: "Minimal", imageSrc: "/images/profile-minimal.svg" },
  { key: "classic", label: "Classic", imageSrc: "/images/profile-classic.svg" },
];

export default function Design() {
  const [selectedColumn, setSelectedColumn] = useState(3);
  const [selectedListLayout, setSelectedListLayout] = useState("portrait");
  const [selectedInstructorProfileLayout, setSelectedInstructorProfileLayout] =
    useState("modern");
  const [selectedStudentProfileLayout, setSelectedStudentProfileLayout] =
    useState("modern");
  const [selectedPresetColor, setSelectedPresetColor] = useState(
    presetColors[0]
  );

  const renderColumnPreview = (count: number, isActive: boolean) => (
    <div className="w-20 h-7 flex items-center gap-2">
      {Array.from({ length: count }).map((_, i) => (
        <div
          key={i}
          className={cn(
            "w-full h-full border",
            isActive ? "bg-primary" : "bg-gray-100"
          )}
        />
      ))}
    </div>
  );

  return (
    <>
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">Design</p>
        <button className="text-sm font-medium text-gray-500 flex items-center gap-1">
          <IoMdRefresh size={18} />
          Reset to Default
        </button>
      </div>

      <p className="text-gray-600 mt-5 mb-2">Course</p>

      <div className="border rounded-md bg-white px-4 mb-5">
        {/* Enable Coupon Code */}
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Enable Coupon Code</p>
            <p className="text-gray-700 mt-1">
              Allow users to apply the coupon code during checkout.
            </p>
          </div>

          <div className="flex items-center gap-5">
            {COLUMN_OPTIONS.map(({ id, label, count }) => {
              const isActive = selectedColumn === id;
              return (
                <div
                  key={id}
                  onClick={() => setSelectedColumn(id)}
                  className="flex flex-col items-center gap-2 cursor-pointer"
                >
                  <div
                    className={cn(
                      "relative hover:bg-gray-100 p-2 rounded-sm",
                      isActive && "bg-gray-100"
                    )}
                  >
                    {renderColumnPreview(count, isActive)}
                    {isActive && (
                      <div className="absolute -top-1 -right-1 bg-primary text-white rounded-full p-px">
                        <MdCheck size={10} />
                      </div>
                    )}
                  </div>
                  <span className={cn(isActive && "text-primary")}>
                    {label}
                  </span>
                </div>
              );
            })}
          </div>
        </div>

        {/* Course Filter */}
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Course Filter</p>
            <p className="text-gray-700 mt-1">
              Show sorting and filtering options on course archive page
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Courses Per Page</p>
            <p className="text-gray-700 mt-1">
              Set the number of courses to display per page on the Course List
              page.
            </p>
          </div>
          <InputField type="number" className="w-20" />
        </div>
        <div className="py-4 border-b border-gray-300 text-sm gap-3">
          <p className="font-medium text-gray-700">Preferred Course Filters</p>
          <p className="text-gray-700 mt-1">
            Choose preferred filter options you&apos;d like to show on the
            course archive page.
          </p>
          <div className="mt-2 flex items-center gap-5">
            <div className="flex items-center gap-2">
              <Checkbox />
              <label htmlFor="" className="font-medium">
                Keyword Search
              </label>
            </div>
            <div className="flex items-center gap-2">
              <Checkbox />
              <label htmlFor="" className="font-medium">
                Category
              </label>
            </div>
            <div className="flex items-center gap-2">
              <Checkbox />
              <label htmlFor="" className="font-medium">
                Tag
              </label>
            </div>
            <div className="flex items-center gap-2">
              <Checkbox />
              <label htmlFor="" className="font-medium">
                Difficulty Level
              </label>
            </div>
            <div className="flex items-center gap-2">
              <Checkbox />
              <label htmlFor="" className="font-medium">
                Price Type
              </label>
            </div>
          </div>
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Course Sorting</p>
            <p className="text-gray-700 mt-1">
              If enabled, the courses will be sortable by Course Name or
              Creation Date in either Ascending or Descending order
            </p>
          </div>
          <ToggleSwitch />
        </div>
      </div>

      <p className="text-gray-600 mt-5 mb-2">Layout</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 text-sm gap-3">
          <p className="font-medium text-gray-700">Instructor List Layout</p>
          <p className="text-gray-700 mt-1">
            Choose a layout for the list of instructors inside a course page.
            You can change this at any time.
          </p>
          <div className="mt-5 flex items-center gap-8">
            <div className="flex flex-col items-center gap-2">
              <div
                onClick={() => setSelectedListLayout("portrait")}
                className={cn(
                  "relative w-[9.5rem] border rounded-md cursor-pointer",
                  selectedListLayout === "portrait" && "ring-2 ring-primary"
                )}
              >
                <Image
                  src={"/images/instructor-portrait.svg"}
                  alt={""}
                  width={300}
                  height={400}
                />
                {selectedListLayout === "portrait" && (
                  <div className="absolute -top-2 -right-2 bg-primary text-white rounded-full p-px">
                    <MdCheck />
                  </div>
                )}
              </div>
              <span>Portrait</span>
            </div>
            <div className="flex flex-col items-center gap-2">
              <div
                onClick={() => setSelectedListLayout("cover")}
                className={cn(
                  "relative w-[9.5rem] border rounded-md cursor-pointer",
                  selectedListLayout === "cover" && "ring-2 ring-primary"
                )}
              >
                <Image
                  src={"/images/instructor-cover.svg"}
                  alt={""}
                  width={300}
                  height={400}
                />
                {selectedListLayout === "cover" && (
                  <div className="absolute -top-2 -right-2 bg-primary text-white rounded-full p-px">
                    <MdCheck />
                  </div>
                )}
              </div>
              <span>Cover</span>
            </div>
            <div className="flex flex-col items-center gap-2">
              <div
                onClick={() => setSelectedListLayout("minimal")}
                className={cn(
                  "relative w-[12rem] h-[12rem] border rounded-md cursor-pointer",
                  selectedListLayout === "minimal" && "ring-2 ring-primary"
                )}
              >
                <Image
                  src={"/images/instructor-minimal.svg"}
                  alt={""}
                  width={300}
                  height={400}
                  className="w-full h-full object-cover"
                />
                {selectedListLayout === "minimal" && (
                  <div className="absolute -top-2 -right-2 bg-primary text-white rounded-full p-px">
                    <MdCheck />
                  </div>
                )}
              </div>
              <span>Minimal</span>
            </div>
          </div>
          <div className="mt-5 flex items-center gap-8">
            <div className="flex flex-col items-center gap-2">
              <div
                onClick={() => setSelectedListLayout("horizontal-portrait")}
                className={cn(
                  "relative h-[6rem] border rounded-md cursor-pointer",
                  selectedListLayout === "horizontal-portrait" &&
                    "ring-2 ring-primary"
                )}
              >
                <Image
                  src={"/images/instructor-horizontal-portrait.svg"}
                  alt={""}
                  width={300}
                  height={400}
                  className="w-full h-full object-cover"
                />
                {selectedListLayout === "horizontal-portrait" && (
                  <div className="absolute -top-2 -right-2 bg-primary text-white rounded-full p-px">
                    <MdCheck />
                  </div>
                )}
              </div>
              <span>Portrait Horizontal</span>
            </div>
            <div className="flex flex-col items-center gap-2">
              <div
                onClick={() => setSelectedListLayout("horizontal-minimal")}
                className={cn(
                  "relative h-[6rem] border rounded-md cursor-pointer",
                  selectedListLayout === "horizontal-minimal" &&
                    "ring-2 ring-primary"
                )}
              >
                <Image
                  src={"/images/instructor-horizontal-minimal.svg"}
                  alt={""}
                  width={300}
                  height={400}
                  className="w-full h-full object-cover"
                />
                {selectedListLayout === "horizontal-minimal" && (
                  <div className="absolute -top-2 -right-2 bg-primary text-white rounded-full p-px">
                    <MdCheck />
                  </div>
                )}
              </div>
              <span>Horizontal Minimal</span>
            </div>
          </div>
        </div>
        <div className="py-4 border-b border-gray-300 text-sm gap-3">
          <p className="font-medium text-gray-700">
            Instructor Public Profile Layout
          </p>
          <p className="text-gray-700 mt-1">
            Choose a layout design for a instructor’s public profile
          </p>
          <div className="mt-5 flex items-center gap-5">
            {profileLayoutOptions.map(({ key, label, imageSrc }) => (
              <div key={key} className="flex flex-col items-center gap-2">
                <div
                  onClick={() => setSelectedInstructorProfileLayout(key)}
                  className={cn(
                    "relative w-[9.5rem] border rounded-md cursor-pointer",
                    selectedInstructorProfileLayout === key &&
                      "ring-2 ring-primary"
                  )}
                >
                  <Image src={imageSrc} alt={label} width={300} height={400} />
                  {selectedInstructorProfileLayout === key && (
                    <div className="absolute -top-2 -right-2 bg-primary text-white rounded-full p-px">
                      <MdCheck />
                    </div>
                  )}
                </div>
                <span>{label}</span>
              </div>
            ))}
          </div>
        </div>
        <div className="py-4 text-sm gap-3">
          <p className="font-medium text-gray-700">
            Student Public Profile Layout
          </p>
          <p className="text-gray-700 mt-1">
            Choose a layout design for a student’s public profile
          </p>
          <div className="mt-5 flex items-center gap-5">
            {profileLayoutOptions.map(({ key, label, imageSrc }) => (
              <div key={key} className="flex flex-col items-center gap-2">
                <div
                  onClick={() => setSelectedStudentProfileLayout(key)}
                  className={cn(
                    "relative w-[9.5rem] border rounded-md cursor-pointer",
                    selectedStudentProfileLayout === key &&
                      "ring-2 ring-primary"
                  )}
                >
                  <Image src={imageSrc} alt={label} width={300} height={400} />
                  {selectedStudentProfileLayout === key && (
                    <div className="absolute -top-2 -right-2 bg-primary text-white rounded-full p-px">
                      <MdCheck />
                    </div>
                  )}
                </div>
                <span>{label}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      <p className="text-gray-600 mt-5 mb-2">Course Details</p>
      <div className="border rounded-md bg-white p-4 mb-5 border-b border-gray-300 text-sm gap-3">
        {/* Enable Coupon Code */}
        <p className="font-medium text-gray-700">Page Features</p>
        <p className="text-gray-700 mt-1">
          You can keep the following features active or inactive as per the need
          of your business model
        </p>
        <div className="grid grid-cols-3 gap-5 mt-5">
          {pageFeatures.map((item) => (
            <div
              key={item.id}
              className="flex items-center gap-2 border border-gray-400 p-3 rounded-md"
            >
              <ToggleSwitch checked={item.enable} />
              <p>{item.name}</p>
            </div>
          ))}
        </div>
      </div>

      <p className="text-gray-600 mt-5 mb-2">Colors</p>
      <div className="border rounded-md bg-white p-4 mb-5 border-b border-gray-300 text-sm">
        {/* Enable Coupon Code */}
        <p className="font-medium text-gray-700">Preset Colors</p>
        <p className="text-gray-700 mt-1">
          These colors will be used throughout your website. Choose between
          these presets or create your own custom palette.
        </p>
        <div className="grid grid-cols-4 gap-5 my-5">
          {presetColors.map((item) => (
            <label
              key={item.name}
              htmlFor={item.name}
              className="border-2 rounded-md border-primary cursor-pointer"
            >
              <div className="flex">
                <div
                  className="w-full h-8"
                  style={{ backgroundColor: item.colors.primaryColor }}
                />
                <div
                  className="w-full h-8"
                  style={{ backgroundColor: item.colors.hoverColor }}
                />
                <div
                  className="w-full h-8"
                  style={{ backgroundColor: item.colors.textColor }}
                />
                <div
                  className="w-full h-8"
                  style={{ backgroundColor: item.colors.gray }}
                />
              </div>
              <div className="flex items-center justify-between p-2">
                <p className="font-medium">{item.name}</p>
                <Checkbox
                  id={item.name}
                  checked={selectedPresetColor === item}
                  onChange={() => setSelectedPresetColor(item)}
                />
              </div>
            </label>
          ))}
        </div>
        <div className="mt-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Primary Color</p>
            <p className="text-gray-700 mt-1">Choose a primary color</p>
          </div>
          <div className="border h-8 flex items-center rounded-md overflow-hidden">
            <div
              className="w-8 h-full"
              style={{
                backgroundColor: selectedPresetColor.colors.primaryColor,
              }}
            />
            <input type="text" className="w-16 h-full outline-none px-2" />
          </div>
        </div>
        <div className="mt-5 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Primary Hover Color</p>
            <p className="text-gray-700 mt-1">Choose a primary hover color</p>
          </div>
          <div className="border h-8 flex items-center rounded-md overflow-hidden">
            <div
              className="w-8 h-full"
              style={{
                backgroundColor: selectedPresetColor.colors.hoverColor,
              }}
            />
            <input type="text" className="w-16 h-full outline-none px-2" />
          </div>
        </div>
        <div className="mt-5 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Text Color</p>
            <p className="text-gray-700 mt-1">
              Choose a text color for your website
            </p>
          </div>
          <div className="border h-8 flex items-center rounded-md overflow-hidden">
            <div
              className="w-8 h-full"
              style={{
                backgroundColor: selectedPresetColor.colors.textColor,
              }}
            />
            <input type="text" className="w-16 h-full outline-none px-2" />
          </div>
        </div>
        <div className="mt-5 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Gray</p>
            <p className="text-gray-700 mt-1">
              Choose a color for elements like table, card etc
            </p>
          </div>
          <div className="border h-8 flex items-center rounded-md overflow-hidden">
            <div
              className="w-8 h-full"
              style={{
                backgroundColor: selectedPresetColor.colors.gray,
              }}
            />
            <input type="text" className="w-16 h-full outline-none px-2" />
          </div>
        </div>
        <div className="mt-5 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Border</p>
            <p className="text-gray-700 mt-1">
              Choose a border color for your website
            </p>
          </div>
          <div className="border h-8 flex items-center rounded-md overflow-hidden">
            <div className="bg-gray-500 w-8 h-full" />
            <input type="text" className="w-16 h-full outline-none px-2" />
          </div>
        </div>
      </div>
      <p className="text-gray-600 mt-5 mb-2">Video Player</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Use Tutor Player for YouTube
            </p>
            <p className="text-gray-700 mt-1">
              Enable this option to use Tutor LMS video player for YouTube.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Use Tutor Player for Vimeo
            </p>
            <p className="text-gray-700 mt-1">
              Enable this option to use Tutor LMS video player for Vimeo.
            </p>
          </div>
          <ToggleSwitch />
        </div>
      </div>
    </>
  );
}
