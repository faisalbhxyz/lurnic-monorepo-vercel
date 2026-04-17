"use client";

import React, { useState } from "react";
import InputField from "@/components/ui/InputField";
import Label from "@/components/ui/Label";
import RadioButton from "@/components/ui/RadioButton";
import SelectList from "@/components/ui/SelectList";
import Checkbox from "@/components/ui/Checkbox";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import DatePicker from "react-datepicker";
import { HiOutlineCalendar } from "react-icons/hi";
import { LuClock } from "react-icons/lu";

const discountTypes = [
  { id: 1, name: "Percent" },
  { id: 2, name: "Amount" },
];

const appliesOptions = [
  { id: 1, name: "All courses" },
  { id: 2, name: "Specific courses" },
  { id: 3, name: "Specific category" },
];

export default function CreateCoupons() {
  const [method, setMethod] = useState("code");
  const [startDate, setStartDate] = useState<Date | null>(null);
  const [selectedTime, setSelectedTime] = useState<string | null>(null);
  const [isEndDate, setIsEndDate] = useState(false);
  const [purchaseRequirements, setPurchaseRequirements] = useState("noMinimum");
  const [limitTotal, setLimitTotal] = useState(false);
  const [limitCustomer, setLimitCustomer] = useState(false);
  const [selectDiscountType, setSelectDiscountType] = useState(
    discountTypes[0]
  );
  const [selectAppliesOption, setSelectAppliesOption] = useState(
    appliesOptions[0]
  );

  const formattedDate = startDate
    ? startDate.toLocaleDateString(undefined, {
        year: "numeric",
        month: "short",
        day: "numeric",
      })
    : "Select date";

  const formatTime = (hour: number): string => {
    const period = hour < 12 ? "AM" : "PM";
    const hour12 = hour % 12 === 0 ? 12 : hour % 12;
    return `${hour12.toString().padStart(2, "0")}:00 ${period}`;
  };

  // Generate times array from 0 to 23
  const times = Array.from({ length: 24 }, (_, i) => formatTime(i));

  return (
    <div className="flex items-start">
      <div className="w-full p-5">
        <div className="border border-gray-400/60 p-4 rounded-md bg-white">
          <p>Coupon Info</p>
          <p className="text-sm text-gray-500">
            Create a coupon code or set up automatic discounts.
          </p>
          <div className="mt-5">
            <p className="text-sm mb-1">Method</p>
            <div className="flex items-center gap-5 mb-3">
              <div className="flex items-center gap-2">
                <RadioButton
                  name="method"
                  value="code"
                  checked={method === "code"}
                  onChange={(e) => setMethod(e.target.value)}
                  id="code"
                />
                <label htmlFor="code" className="text-sm">
                  Code
                </label>
              </div>
              <div className="flex items-center gap-2">
                <RadioButton
                  name="method"
                  value="automatic"
                  checked={method === "automatic"}
                  onChange={(e) => setMethod(e.target.value)}
                  id="automatic"
                />
                <label htmlFor="automatic" className="text-sm">
                  Automatic
                </label>
              </div>
            </div>
            <div className="mb-3">
              <Label htmlFor={""}>Title</Label>
              <InputField
                placeholder="e.g. Summer Sale 2025"
                className="w-full"
              />
            </div>
            {method === "code" && (
              <div>
                <div className="flex items-center justify-between mb-1">
                  <label
                    htmlFor={""}
                    className="text-sm text-gray-700 font-medium"
                  >
                    Coupon Code
                  </label>
                  <button className="text-primary text-sm">
                    Generate Code
                  </button>
                </div>
                <InputField placeholder="e.g. SUMMER20" className="w-full" />
              </div>
            )}
          </div>
        </div>
        <div className="border border-gray-400/60 p-4 rounded-md bg-white mt-5">
          <p className="mb-3">Discount</p>
          <div className="flex items-center gap-5 mb-5">
            <div className="w-full">
              <Label htmlFor={""}>Discount Type</Label>
              <SelectList
                options={discountTypes}
                value={selectDiscountType}
                onChange={setSelectDiscountType}
              />
            </div>
            <div className="w-full">
              <Label htmlFor={""}>Discount Value</Label>
              <div className="flex items-center border focus-within:ring-1 focus-within:ring-primary w-full h-9 rounded-md">
                <span className="px-4 border-r border-gray-300 text-gray-500 h-full flex items-center">
                  %
                </span>
                <input
                  type="text"
                  placeholder="0"
                  className="w-full h-full px-3 outline-none text-sm"
                />
              </div>
            </div>
          </div>
          <div className="w-full">
            <Label htmlFor={""}>Applies to</Label>
            <SelectList
              options={appliesOptions}
              value={selectAppliesOption}
              onChange={setSelectAppliesOption}
            />
          </div>
        </div>
        <div className="border border-gray-400/60 p-4 rounded-md bg-white mt-5">
          <p className="mb-3">Usage Limitation</p>
          <div className="flex items-center gap-2 mb-2">
            <Checkbox
              checked={limitTotal}
              onChange={() => setLimitTotal((prev) => !prev)}
            />
            <label htmlFor="code" className="text-sm">
              Limit number of times this coupon can be used in total
            </label>
          </div>
          {limitTotal && (
            <InputField
              placeholder="0"
              className="w-60 outline-none text-sm mb-3"
            />
          )}
          <div className="flex items-center gap-2">
            <Checkbox
              checked={limitCustomer}
              onChange={() => setLimitCustomer((prev) => !prev)}
            />
            <label htmlFor="code" className="text-sm">
              Limit number of times this coupon can be used by a customer
            </label>
          </div>
          {limitCustomer && (
            <InputField
              placeholder="0"
              className="w-60 outline-none text-sm mt-2"
            />
          )}
        </div>
        <div className="border border-gray-400/60 p-4 rounded-md bg-white mt-5">
          <p className="mb-3">Minimum Purchase Requirements</p>
          <div className="flex items-center gap-2 mb-2">
            <RadioButton
              name="purchaseRequirements"
              value="noMinimum"
              checked={purchaseRequirements === "noMinimum"}
              onChange={(e) => setPurchaseRequirements(e.target.value)}
              id="noMinimum"
            />
            <label htmlFor="noMinimum" className="text-sm">
              No minimum requirements
            </label>
          </div>
          <div className="flex items-center gap-2 mb-2">
            <RadioButton
              name="purchaseRequirements"
              value="minimum"
              checked={purchaseRequirements === "minimum"}
              onChange={(e) => setPurchaseRequirements(e.target.value)}
              id="minimum"
            />
            <label htmlFor="minimum" className="text-sm">
              Minimum purchase amount ($)
            </label>
          </div>
          {purchaseRequirements === "minimum" && (
            <div className="flex items-center border focus-within:ring-1 mb-3 focus-within:ring-primary w-60 h-9 rounded-md">
              <span className="px-4 border-r border-gray-300 text-gray-500 h-full flex items-center">
                $
              </span>
              <input
                type="text"
                placeholder="0"
                className="w-full h-full px-3 outline-none text-sm"
              />
            </div>
          )}

          <div className="flex items-center gap-2 mb-2">
            <RadioButton
              name="purchaseRequirements"
              value="minimumQuantity"
              checked={purchaseRequirements === "minimumQuantity"}
              onChange={(e) => setPurchaseRequirements(e.target.value)}
              id="minimumQuantity"
            />
            <label htmlFor="minimumQuantity" className="text-sm">
              Minimum quantity of courses
            </label>
          </div>
          {purchaseRequirements === "minimumQuantity" && (
            <InputField placeholder="0" className="w-60 outline-none text-sm" />
          )}
        </div>
        <div className="border border-gray-400/60 p-4 rounded-md bg-white mt-5">
          <p className="mb-3">Validity</p>
          <p className="text-sm mb-1 text-gray-700">Starts from</p>
          <div className="flex items-center gap-5">
            <Menu>
              <MenuButton className="w-60 border rounded-md flex items-center px-3 gap-2 text-sm py-2">
                <HiOutlineCalendar size={18} className="text-gray-500" />
                {formattedDate}
              </MenuButton>

              <MenuItems
                transition
                anchor="bottom start"
                className="origin-top-left rounded-xl p-1 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
              >
                <DatePicker
                  selected={startDate}
                  onChange={(date) => setStartDate(date)}
                  inline
                />
              </MenuItems>
            </Menu>
            <Menu>
              <MenuButton className="w-60 border rounded-md flex items-center px-3 gap-2 text-sm py-2">
                <LuClock size={18} className="text-gray-500" />
                {selectedTime || "Select Time"}
              </MenuButton>

              <MenuItems
                transition
                anchor="bottom end"
                className="origin-top-right border bg-white h-80 rounded-xl p-1 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
              >
                {times.map((time) => (
                  <MenuItem key={time}>
                    <button
                      onClick={() => setSelectedTime(time)}
                      className="group flex w-full items-center gap-2 text-sm rounded-lg py-1.5 px-3 data-[focus]:bg-gray-100"
                    >
                      {time}
                    </button>
                  </MenuItem>
                ))}
              </MenuItems>
            </Menu>
          </div>
          <div className="flex items-center gap-2 mt-2">
            <Checkbox
              checked={isEndDate}
              onChange={() => setIsEndDate((prev) => !prev)}
            />
            <label htmlFor="code" className="text-sm">
              Set end date
            </label>
          </div>
          <p className="text-xs text-gray-500 ml-6">
            Leaving the end date blank will make the coupon valid indefinitely.
          </p>
          {isEndDate && (
            <div className="mt-3">
              <p className="text-sm mb-1 text-gray-700">Ends in</p>
              <div className="flex items-center gap-5">
                <Menu>
                  <MenuButton className="w-60 border rounded-md flex items-center px-3 gap-2 text-sm py-2">
                    <HiOutlineCalendar size={18} className="text-gray-500" />
                    {formattedDate}
                  </MenuButton>

                  <MenuItems
                    transition
                    anchor="bottom start"
                    className="origin-top-left rounded-xl p-1 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
                  >
                    <DatePicker
                      selected={startDate}
                      onChange={(date) => setStartDate(date)}
                      inline
                    />
                  </MenuItems>
                </Menu>
                <Menu>
                  <MenuButton className="w-60 border rounded-md flex items-center px-3 gap-2 text-sm py-2">
                    <LuClock size={18} className="text-gray-500" />
                    {selectedTime || "Select Time"}
                  </MenuButton>

                  <MenuItems
                    transition
                    anchor="bottom end"
                    className="origin-top-right border bg-white h-80 rounded-xl p-1 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0"
                  >
                    {times.map((time) => (
                      <MenuItem key={time}>
                        <button
                          onClick={() => setSelectedTime(time)}
                          className="group flex w-full items-center gap-2 text-sm rounded-lg py-1.5 px-3 data-[focus]:bg-gray-100"
                        >
                          {time}
                        </button>
                      </MenuItem>
                    ))}
                  </MenuItems>
                </Menu>
              </div>
            </div>
          )}
        </div>
      </div>
      <div className="p-5 w-96 min-w-96">Right</div>
    </div>
  );
}
