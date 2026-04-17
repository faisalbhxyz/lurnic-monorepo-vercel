import ToggleSwitch from "@/components/ui/ToggleSwitch";
import React from "react";
import { IoMdRefresh } from "react-icons/io";
import SelectPage from "./SelectPage";
import InputField from "@/components/ui/InputField";

export default function Monetization() {
  return (
    <>
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">Course</p>
        <button className="text-sm font-medium text-gray-500 flex items-center gap-1">
          <IoMdRefresh size={18} />
          Reset to Default
        </button>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Options</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Select eCommerce Engine</p>
            <p className="text-gray-700 mt-1">
              Select a monetization option to generate revenue by selling
              courses.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Cart Page</p>
            <p className="text-gray-700 mt-1">
              Select the page you wish to set as the cart page.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Checkout Page</p>
            <p className="text-gray-700 mt-1">
              Select the page to be used as the checkout page.
            </p>
          </div>
          <SelectPage />
        </div>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Currency</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Currency</p>
            <p className="text-gray-700 mt-1">
              Choose the currency for transactions.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Currency Position</p>
            <p className="text-gray-700 mt-1">
              Set the position of the currency symbol.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Thousand Separator</p>
            <p className="text-gray-700 mt-1">
              Specify the thousand separator.
            </p>
          </div>
          <InputField className="w-20 text-center" />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Decimal Separator</p>
            <p className="text-gray-700 mt-1">Specify the decimal separator.</p>
          </div>
          <InputField className="w-20 text-center" />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Number of Decimals</p>
            <p className="text-gray-700 mt-1">
              Set the number of decimal places.
            </p>
          </div>
          <InputField className="w-20 text-center" />
        </div>
      </div>
      <div className="p-4 bg-white border rounded-md flex items-center justify-between text-sm gap-3">
        <div>
          <p className="font-medium text-gray-700">Enable Revenue Sharing</p>
          <p className="text-gray-700 mt-1">
            Allow revenue generated from selling courses to be shared with
            course creators.
          </p>
        </div>
        <ToggleSwitch />
      </div>
    </>
  );
}
