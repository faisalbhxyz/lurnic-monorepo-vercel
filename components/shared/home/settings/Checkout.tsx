import ToggleSwitch from "@/components/ui/ToggleSwitch";
import React from "react";
import { IoMdRefresh } from "react-icons/io";

export default function Checkout() {
  return (
    <>
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">Checkout</p>
        <button className="text-sm font-medium text-gray-500 flex items-center gap-1">
          <IoMdRefresh size={18} />
          Reset to Default
        </button>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Checkout Configuration</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Enable Coupon Code</p>
            <p className="text-gray-700 mt-1">
              Allow users to apply the coupon code during checkout.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Enable &quot;Buy Now&quot; Button
            </p>
            <p className="text-gray-700 mt-1">
              Allow users to purchase courses directly without adding them to
              the cart.
            </p>
          </div>
          <ToggleSwitch />
        </div>
      </div>
    </>
  );
}
