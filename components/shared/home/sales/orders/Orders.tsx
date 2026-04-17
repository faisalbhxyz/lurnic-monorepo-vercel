"use client";

import React, { useState } from "react";
import OrderList from "./OrderList";
import Link from "next/link";
import { FiSearch } from "react-icons/fi";

const statuses = ["All", "Paid", "Unpaid"];

export default function Orders({ orders }: { orders: IOrder[] }) {
  const [activeTab, setActiveTab] = useState("All");
  const [search, setSearch] = useState("");

  const filteredData = orders.filter((order) => {
    const paymentStatus = order.payment_status.toLowerCase();

    const matchesStatus =
      activeTab.toLowerCase() === "all"
        ? true
        : activeTab.toLowerCase() === "paid or unpaid"
        ? paymentStatus === "paid" || paymentStatus === "unpaid"
        : paymentStatus === activeTab.toLowerCase();

    const matchesSearch = order.course.title
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
        /<Link href={""}>Orders</Link>
      </div>
      <h3 className="font-medium text-2xl">Orders</h3>
      <div className="flex-between my-5">
        <div className="flex gap-3">
          {statuses.map((status) => {
            const count =
              status === "All"
                ? orders.length
                : orders.filter(
                    (c) =>
                      c.payment_status.toLowerCase() === status.toLowerCase()
                  ).length;
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
      <OrderList data={filteredData} />
    </>
  );
}
