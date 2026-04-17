"use client";

import { formatDate } from "@/lib/helpers";
import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";
import InstructorAction from "./InstructorAction";

interface OrderListProps {
  data: IInstructor[];
  teacherPrefix: string | null;
}

export default function InstructorList({ data, teacherPrefix }: OrderListProps) {
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

  const filteredData = data.filter((ins) => {
    const searchTerm = search.toLowerCase();
    return (
      ins.first_name.toLowerCase().includes(searchTerm) ||
      ins.last_name?.toLowerCase().includes(searchTerm) || // optional chaining in case last_name is null
      ins.email.toLowerCase().includes(searchTerm)
    );
  });

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
                  <input
                    type="checkbox"
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                    className="w-4 h-4"
                  />
                  <span>ID</span>
                </div>
              </th>
              <th className="p-3 font-medium">Name</th>
              <th className="p-3 font-medium">Email</th>
              {/* <th className="p-3 font-medium">Phone Number</th> */}
              {/* <th className="p-3 font-medium">Last Logged In</th> */}
              <th className="p-3 font-medium">Joined</th>
              <th className="p-3 font-medium text-end">Action</th>
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
                    <input
                      type="checkbox"
                      checked={selected.includes(order.id)}
                      onChange={() => toggleSelectOne(order.id)}
                      className="w-4 h-4"
                    />
                    <p className="font-medium">{teacherPrefix ?? ""}{order.id}</p>
                  </div>
                </td>
                <td className="p-3">
                  <div className="flex items-center gap-5">
                    <div className="flex items-center gap-3">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">
                          {order.first_name} {order.last_name ?? ""}
                        </p>
                      </div>
                    </div>
                  </div>
                </td>
                <td className="p-3">{order.email}</td>
                {/* <td className="p-3">+8801234567890</td> */}
                {/* <td className="p-3">N/A</td> */}
                <td className="p-3">
                  <div className="flex items-center gap-2">
                    {formatDate(order.created_at)}
                  </div>
                </td>
                <td className="p-3 text-end">
                  <InstructorAction id={order.id} />
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
