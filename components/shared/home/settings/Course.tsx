import Checkbox from "@/components/ui/Checkbox";
import InputField from "@/components/ui/InputField";
import RadioButton from "@/components/ui/RadioButton";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import React from "react";
import { IoMdRefresh } from "react-icons/io";

export default function Course() {
  return (
    <>
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">Course</p>
        <button className="text-sm font-medium text-gray-500 flex items-center gap-1">
          <IoMdRefresh size={18} />
          Reset to Default
        </button>
      </div>

      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Course Visibility</p>
            <p className="text-gray-700 mt-1">
              Students must be logged in to view course
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Course Content Access</p>
            <p className="text-gray-700 mt-1">
              Allow instructors and admins to view the course content without
              enrolling
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Content Summary</p>
            <p className="text-gray-700 mt-1">
              Enabling this feature will show a course content summary on the
              Course Details page.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Auto Redirect to Courses
            </p>
            <p className="text-gray-700 mt-1">
              When a user&apos;s WooCommerce order is auto-completed, they will
              be redirected to enrolled courses
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Spotlight Mode</p>
            <p className="text-gray-700 mt-1">
              This will hide the header and the footer and enable spotlight
              (full screen) mode when students view lessons.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Auto Complete Course on All Lesson Completion
            </p>
            <p className="text-gray-700 mt-1">
              If enabled, an Enrolled Course will be automatically completed if
              all its Lessons, Quizzes, and Assignments are already completed by
              the Student
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 text-sm gap-3">
          <p className="font-medium text-gray-700">Course Completion Process</p>
          <p className="text-gray-700 mt-1">
            Choose when a user can click on the “Complete Course” button
          </p>
          <div className="flex items-center gap-3 mt-3 mb-2">
            <RadioButton />
            <div>
              <p className="font-medium">Flexible</p>
              <p className="text-xs">
                Students can complete courses anytime in the Flexible mode
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3 mt-3 mb-2">
            <RadioButton />
            <div>
              <p className="font-medium">Strict</p>
              <p className="text-xs">
                Students have to complete, pass all the lessons and quizzes (if
                any) to mark a course as complete.
              </p>
            </div>
          </div>
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Course Retake</p>
            <p className="text-gray-700 mt-1">
              Enabling this feature will allow students to reset course progress
              and start over.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Publish Course Review on Admin&apos;s Approval
            </p>
            <p className="text-gray-700 mt-1">
              Enable to publish/re-publish Course Review after the approval of
              Site Admin
            </p>
          </div>
          <ToggleSwitch />
        </div>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Lesson</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">WP Editor for Lesson</p>
            <p className="text-gray-700 mt-1">
              Enable classic editor to edit lesson.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Automatically Load Next Course Content.
            </p>
            <p className="text-gray-700 mt-1">
              Enable this feature to automatically load the next course content
              after the current one is finished.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Enable Lesson Comment</p>
            <p className="text-gray-700 mt-1">
              Enable this feature to allow students to post comments on lessons.
            </p>
          </div>
          <ToggleSwitch />
        </div>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Quiz</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 text-sm gap-3">
          <p className="font-medium text-gray-700">When time expires</p>
          <p className="text-gray-700 mt-1">
            Choose which action to follow when the quiz time expires.
          </p>
          <div className="flex items-center gap-3 mt-3 mb-2">
            <RadioButton />
            <div>
              <p className="font-medium">Auto Submit</p>
              <p className="text-xs">
                The current quiz answers are submitted automatically.
              </p>
            </div>
          </div>
          <div className="flex items-center gap-3 mt-3 mb-2">
            <RadioButton />
            <div>
              <p className="font-medium">Auto Abandon</p>
              <p className="text-xs">
                Attempts must be submitted before time expires, otherwise they
                will not be counted
              </p>
            </div>
          </div>
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Correct Answer Display Time (When Reveal Mode is enabled)
            </p>
            <p className="text-gray-700 mt-1">
              Put the answer display time in seconds
            </p>
          </div>
          <InputField className="w-20" />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Default Quiz Attempt limit (When Retry Mode is enabled)
            </p>
            <p className="text-gray-700 mt-1">
              The highest number of attempts allowed for students to participate
              a quiz. 0 means unlimited. This will work as the default Quiz
              Attempt limit in case of Quiz Retry Mode.
            </p>
          </div>
          <InputField className="w-20" />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Show Quiz Previous Button
            </p>
            <p className="text-gray-700 mt-1">
              Choose whether to show or hide the previous button for each
              question.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 text-sm gap-3">
          <p className="font-medium text-gray-700">Final Grade Calculation</p>
          <p className="text-gray-700 mt-1">
            When multiple attempts are allowed, select which method should be
            used to calculate a student&apos;s final grade for the quiz.
          </p>
          <div className="mt-3 flex items-center gap-3 text-gray-800">
            <div className="flex items-center gap-3">
              <RadioButton />
              <label className="font-medium text-sm">Highest Grade</label>
            </div>
            <div className="flex items-center gap-3">
              <RadioButton />
              <label className="font-medium text-sm">Average Grade</label>
            </div>
            <div className="flex items-center gap-3">
              <RadioButton />
              <label className="font-medium text-sm">First Attempt</label>
            </div>
            <div className="flex items-center gap-3">
              <RadioButton />
              <label className="font-medium text-sm">Last Attempt</label>
            </div>
          </div>
        </div>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Lesson</p>
      <div className="border rounded-md bg-white px-4 py-2 text-sm mb-5">
        <p className="font-medium text-gray-700">Preferred Video Source</p>
        <p className="text-gray-700 mt-1">
          Select the video hosting platform(s) you want to enable.
        </p>
        <div className="flex items-center gap-2 mt-2">
          <Checkbox />
          <label htmlFor="" className="font-medium">
            HTML 5 (mp4)
          </label>
        </div>
        <div className="flex items-center gap-2 mt-2">
          <Checkbox />
          <label htmlFor="" className="font-medium">
            External URL
          </label>
        </div>
        <div className="flex items-center gap-2 mt-2">
          <Checkbox />
          <label htmlFor="" className="font-medium">
            YouTube
          </label>
        </div>
        <div className="flex items-center gap-2 mt-2">
          <Checkbox />
          <label htmlFor="" className="font-medium">
            Vimeo
          </label>
        </div>
        <div className="flex items-center gap-2 mt-2">
          <Checkbox />
          <label htmlFor="" className="font-medium">
            Embedded
          </label>
        </div>
        <div className="flex items-center gap-2 mt-2">
          <Checkbox />
          <label htmlFor="" className="font-medium">
            Shortcode
          </label>
        </div>
      </div>
    </>
  );
}
