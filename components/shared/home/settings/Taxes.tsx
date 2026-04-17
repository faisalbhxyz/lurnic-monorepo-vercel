import React from "react";
import CountryActions from "./CountryActions";
import RadioButton from "@/components/ui/RadioButton";
import Checkbox from "@/components/ui/Checkbox";
import AddRegion from "./AddRegion";

export default function Taxes() {
  return (
    <>
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">Taxes</p>
      </div>
      <div className="border rounded-lg bg-white">
        <div className="border-b border-gray-300 px-4 py-3">
          <p className="text-lg font-medium">Tax Regions & Rates</p>
          <p className="text-sm text-gray-500 font-medium">
            Specify regions and their applicable tax rates.
          </p>
        </div>
        <div className="p-4">
          <ul className="border border-gray-200 rounded-md text-sm">
            <li className="flex items-center px-3 py-2 font-medium">
              <div className="w-full">Country</div>
              <div className="w-36 min-w-36">Tax rate</div>
            </li>
            <li className="flex items-center px-3 py-2 font-medium border-t border-gray-200">
              <div className="w-full">Bangladesh</div>
              <div className="w-36 min-w-36 flex items-center justify-between">
                <span>0%</span>
                <CountryActions />
              </div>
            </li>
            <li className="flex items-center px-3 py-2 font-medium border-t border-gray-200">
              <AddRegion />
            </li>
          </ul>
        </div>
      </div>
      <div className="border rounded-lg bg-white mt-5">
        <div className="border-b border-gray-300 px-4 py-3">
          <p className="text-lg font-medium">Global Tax Settingss</p>
          <p className="text-sm text-gray-500 font-medium">
            Set how taxes are displayed and applied to your courses.
          </p>
        </div>
        <div className="p-4">
          <div className="flex items-center gap-2">
            <RadioButton />
            <label htmlFor="" className="text-sm">
              Tax is already included in my prices
            </label>
          </div>
          <div className="flex items-center gap-2 mt-2">
            <RadioButton />
            <label htmlFor="" className="text-sm">
              Tax should be calculated and displayed on the checkout page
            </label>
          </div>
          <div className="flex items-center gap-2 mt-2">
            <Checkbox />
            <label htmlFor="" className="text-sm">
              Display prices inclusive tax
            </label>
          </div>
          <p className="text-sm ml-6 text-gray-500">
            Show prices with tax included, so customers see the final amount
            they’ll pay upfront.
          </p>
        </div>
      </div>
    </>
  );
}
