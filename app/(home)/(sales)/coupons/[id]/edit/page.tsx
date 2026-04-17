import CreateCoupons from "@/components/shared/home/sales/Coupons/CreateCoupons";
import Button from "@/components/ui/Button";
import React from "react";

export default function page() {
  return (
    <>
      <div className="flex-between px-5 py-3 bg-white border-b border-gray-300">
        <div className="w-full flex items-center justify-between gap-5">
          <h3 className="font-medium">Update Coupon</h3>
          <div className="flex items-center gap-3">
            <Button variant="secondary" className="px-4">
              Cancel
            </Button>
            <Button className="px-4">Save</Button>
          </div>
        </div>
      </div>
      <CreateCoupons />
    </>
  );
}
