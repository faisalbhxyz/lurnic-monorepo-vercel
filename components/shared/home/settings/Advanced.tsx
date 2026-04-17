import ToggleSwitch from "@/components/ui/ToggleSwitch";
import React from "react";
import { IoMdRefresh } from "react-icons/io";
import SelectPage from "./SelectPage";
import InputField from "@/components/ui/InputField";

export default function Advanced() {
  return (
    <>
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">Advanced</p>
        <button className="text-sm font-medium text-gray-500 flex items-center gap-1">
          <IoMdRefresh size={18} />
          Reset to Default
        </button>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Course</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Gutenberg Editor</p>
            <p className="text-gray-700 mt-1">
              Enable this to create courses using the Gutenberg Editor.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Hide Course Products on Shop Page
            </p>
            <p className="text-gray-700 mt-1">
              Enable to hide course products on shop page.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Course Archive Page</p>
            <p className="text-gray-700 mt-1">
              This page will be used to list all the published courses.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Instructor Registration Page
            </p>
            <p className="text-gray-700 mt-1">
              Choose the page for instructor registration.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Student Registration Page
            </p>
            <p className="text-gray-700 mt-1">
              Choose the page for student registration.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">YouTube API Key</p>
            <p className="text-gray-700 mt-1">
              To host live videos on your platform using YouTube, enter your
              YouTube API key.
            </p>
          </div>
          <InputField />
        </div>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Base Permalink</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Course Permalink</p>
            <p className="text-gray-700 mt-1">
              https://amerrajjonowga.com/courses/sample-course
            </p>
          </div>
          <InputField />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Lesson Permalink</p>
            <p className="text-gray-700 mt-1">
              https://amerrajjonowga.com/courses/sample-course/lessons/sample-lesson/
            </p>
          </div>
          <InputField />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Quiz Permalink</p>
            <p className="text-gray-700 mt-1">
              https://amerrajjonowga.com/courses/sample-course/quizzes/sample-quiz/
            </p>
          </div>
          <InputField />
        </div>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Options</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Profile Completion</p>
            <p className="text-gray-700 mt-1">
              Enabling this feature will show a notification bar to students and
              instructors to complete their profile information
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Enable Tutor Login</p>
            <p className="text-gray-700 mt-1">
              Enable to use the tutor login modal instead of the default
              WordPress login page
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Erase upon uninstallation
            </p>
            <p className="text-gray-700 mt-1">
              Delete all data during uninstallation
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Maintenance Mode</p>
            <p className="text-gray-700 mt-1">
              Enabling maintenance mode will display a custom message on the
              frontend. During maintenance mode, visitors cannot access site
              content, but the wp-admin dashboard remains accessible.
            </p>
          </div>
          <ToggleSwitch />
        </div>
      </div>
    </>
  );
}
