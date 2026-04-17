"use client";

import Checkbox from "@/components/ui/Checkbox";
import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";

interface Order {
  id: number;
  name: string;
  email: string;
  webinar: string;
  status: string;
  date: string;
}

interface OrderListProps {
  data: Order[];
}

export default function DigitalAccessList({ data }: OrderListProps) {
  const [selected, setSelected] = useState<number[]>([]);
  const [search, setSearch] = useState("");

  const toggleSelectAll = () => {
    if (selected.length === filteredData.length) {
      setSelected([]);
    } else {
      setSelected(filteredData.map((course) => course.id));
    }
  };

  const toggleSelectOne = (id: number) => {
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const filteredData = data.filter((coupon) =>
    coupon.name.toLowerCase().includes(search.toLowerCase())
  );

  const isAllSelected =
    selected.length === filteredData.length && filteredData.length > 0;
  return (
    <>
      <div className="flex items-center justify-between mb-5">
        <div></div>
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
      <div className="border rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-100">
            <tr className="text-left">
              <th className="p-3">
                <div className="flex items-center gap-5 font-medium">
                  <Checkbox
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                  />
                  <span>Name</span>
                </div>
              </th>
              <th className="p-3 font-medium">Email</th>
              <th className="p-3 font-medium">Webinar</th>
              <th className="p-3 font-medium">Status</th>
              <th className="p-3 font-medium text-end">Registered At</th>
            </tr>
          </thead>
          <tbody>
            {filteredData.map((order) => (
              <tr
                key={order.id}
                className="border-t border-gray-300 hover:bg-gray-100"
              >
                <td className="p-3">
                  <div className="flex items-center gap-5">
                    <Checkbox
                      checked={selected.includes(order.id)}
                      onChange={() => toggleSelectOne(order.id)}
                      className="w-4 h-4"
                    />
                    <div className="flex items-center gap-3">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">{order.name}</p>
                      </div>
                    </div>
                  </div>
                </td>
                <td className="p-3">{order.email}</td>
                <td className="p-3">{order.webinar}</td>
                <td className="p-3">
                  <div className="flex items-center gap-2">{order.status}</div>
                </td>
                <td className="p-3 text-end">
                  <div className="flex items-center justify-end gap-2 ">
                    {order.date}
                  </div>
                </td>
              </tr>
            ))}
            {filteredData.length === 0 && (
              <tr>
                <td colSpan={6} className="p-5 text-center text-gray-500">
                  No courses found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </>
  );
}
