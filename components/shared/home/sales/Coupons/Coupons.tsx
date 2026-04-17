"use client";

import React, { useState } from "react";
import Button from "@/components/ui/Button";
import { LuPlus } from "react-icons/lu";
import CouponsList from "./CouponsList";
import Link from "next/link";
import { FiSearch } from "react-icons/fi";

const sampleCoupons = [
  {
    id: 1,
    name: "Spring Sale",
    discount: "15%",
    type: "Percentage",
    code: "SPRING15",
    status: "Active",
    uses: 10,
  },
  {
    id: 2,
    name: "Welcome Discount",
    discount: "$10",
    type: "Fixed",
    code: "WELCOME10",
    status: "Active",
    uses: 25,
  },
  {
    id: 3,
    name: "Black Friday Deal",
    discount: "50%",
    type: "Percentage",
    code: "BLACKFRIDAY50",
    status: "Expired",
    uses: 100,
  },
  {
    id: 4,
    name: "Holiday Special",
    discount: "$20",
    type: "Fixed",
    code: "HOLIDAY20",
    status: "Scheduled",
    uses: 0,
  },
  {
    id: 5,
    name: "Flash Sale",
    discount: "30%",
    type: "Percentage",
    code: "FLASH30",
    status: "Active",
    uses: 5,
  },
];

const statuses = ["All", "Active", "Inactive", "Trash"];

export default function Coupons() {
  const [activeTab, setActiveTab] = useState("All");
  const [search, setSearch] = useState("");

  const filteredData = sampleCoupons.filter((coupon) => {
    const matchesStatus =
      activeTab === "All" ||
      coupon.status.toLowerCase() === activeTab.toLowerCase();
    const matchesSearch = coupon.name
      .toLowerCase()
      .includes(search.toLowerCase());
    return matchesStatus && matchesSearch;
  });

  return (
    <>
      <div className="flex items-center text-sm gap-1 mb-5">
        <Link href={""} className="text-gray-500">
          Sales
        </Link>
        /<Link href={""}>Coupons</Link>
      </div>
      <div className="flex-between">
        <h3 className="font-medium text-2xl">Coupons</h3>
        <Button link src="/coupons/create">
          <LuPlus /> Add New Coupon
        </Button>
      </div>
      <div className="flex-between my-5">
        <div className="flex gap-3">
          {statuses.map((status) => {
            const count =
              status === "All"
                ? sampleCoupons.length
                : sampleCoupons.filter((c) => c.status === status).length;
            return (
              <button
                key={status}
                onClick={() => setActiveTab(status)}
                className={`px-3 py-3 border-b text-sm ${
                  activeTab === status
                    ? "border-primary text-primary"
                    : "border-transparent text-gray-700 hover:text-primary"
                }`}
              >
                {status}({count})
              </button>
            );
          })}
        </div>
        <div className="relative w-full max-w-60">
          <FiSearch className="absolute top-1/2 left-3 transform -translate-y-1/2 text-gray-500" />
          <input
            type="text"
            placeholder="Search..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full border text-sm px-3 py-1.5 rounded-md pl-8"
          />
        </div>
      </div>
      <CouponsList data={filteredData} />
    </>
  );
}
