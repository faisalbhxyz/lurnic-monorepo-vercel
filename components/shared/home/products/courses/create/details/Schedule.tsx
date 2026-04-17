import ToggleSwitch from "@/components/ui/ToggleSwitch";
import React from "react";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import DatePicker from "react-datepicker";
import { HiOutlineCalendar } from "react-icons/hi";
import { LuClock } from "react-icons/lu";
import Checkbox from "@/components/ui/Checkbox";
import { Controller, useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";

export default function Schedule() {
  const { control, watch, setValue } = useFormContext<TCourseSchema>();

  const isSchedule = watch("is_scheduled");
  const date = watch("schedule_date");
  const time = watch("schedule_time");

  const toLocalDateString = (d: Date) => {
    const year = d.getFullYear();
    const month = (d.getMonth() + 1).toString().padStart(2, "0");
    const day = d.getDate().toString().padStart(2, "0");
    return `${year}-${month}-${day}`;
  };

  // Format for UI
  const formattedDate = date
    ? new Date(date).toLocaleDateString(undefined, {
        year: "numeric",
        month: "short",
        day: "numeric",
      })
    : "Select date";

  // Format times
  const formatTime = (hour: number): string => {
    const period = hour < 12 ? "AM" : "PM";
    const hour12 = hour % 12 === 0 ? 12 : hour % 12;
    return `${hour12.toString().padStart(2, "0")}:00 ${period}`;
  };

  const times = Array.from({ length: 24 }, (_, i) => formatTime(i));

  // Convert Date to ISO (fixes Go parsing)
  const toISO = (d: Date) => d.toISOString();

  return (
    <div className="border border-gray-400 rounded-md bg-white">
      <div className="flex items-center justify-between p-2.5">
        <span className="text-sm">Schedule</span>
        <Controller
          control={control}
          name="is_scheduled"
          render={({ field: { value, onChange } }) => (
            <ToggleSwitch checked={value} onChange={onChange} />
          )}
        />
      </div>

      {isSchedule && (
        <div className="p-2.5">
          <div className="border flex items-center rounded-md">
            {/* DATE PICKER */}
            <Menu>
              <MenuButton className="w-full flex items-center justify-center gap-2 text-sm py-2">
                <HiOutlineCalendar size={18} className="text-gray-500" />
                {formattedDate}
              </MenuButton>

              <MenuItems
                transition
                anchor="bottom start"
                className="origin-top-left rounded-xl p-1 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
              >
                <Controller
                  control={control}
                  name="schedule_date"
                  render={({ field: { onChange, value } }) => (
                    <DatePicker
                      selected={value ? new Date(value) : null}
                      inline
                      onChange={(d) => {
                        if (!d) return;
                        onChange(toLocalDateString(d)); // Send ISO to backend
                      }}
                    />
                  )}
                />
              </MenuItems>
            </Menu>

            <div className="w-px h-9 bg-gray-300" />

            {/* TIME PICKER */}
            <Menu>
              <MenuButton className="w-full flex items-center justify-center gap-2 text-sm py-2">
                <LuClock size={18} className="text-gray-500" />
                {typeof time === "string" ? time : "Select Time"}
              </MenuButton>

              <MenuItems
                transition
                anchor="bottom end"
                className="origin-top-right border bg-white h-80 rounded-xl p-1 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
              >
                {times.map((t) => (
                  <MenuItem key={t}>
                    <button
                      onClick={() => setValue("schedule_time", t as any)}
                      className="group flex w-full items-center gap-2 text-sm rounded-lg py-1.5 px-3 data-[focus]:bg-gray-100"
                    >
                      {t}
                    </button>
                  </MenuItem>
                ))}
              </MenuItems>
            </Menu>
          </div>

          {/* SHOW COMING SOON */}
          <Controller
            control={control}
            name="show_comming_soon"
            render={({ field: { value, onChange } }) => (
              <div className="flex items-start gap-2 mt-2">
                <span className="pt-0.5">
                  <Checkbox id="show" checked={value} onChange={onChange} />
                </span>
                <label
                  htmlFor="show"
                  className="text-sm font-medium text-gray-500"
                >
                  Show coming soon in course list & details page
                </label>
              </div>
            )}
          />
        </div>
      )}
    </div>
  );
}
