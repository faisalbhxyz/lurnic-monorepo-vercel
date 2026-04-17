"use client";

import Checkbox from "@/components/ui/Checkbox";
import { formatDate } from "@/lib/helpers";
import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";
import StudentAction from "./StudentAction";

interface OrderListProps {
  data: IStudent[];
  studentPrefix: string | null;
}

export default function StudentList({ data, studentPrefix }: OrderListProps) {
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

  const filteredData = data.filter((std) =>
    [std.first_name, std.last_name, std.email].some((field) =>
      field?.toLowerCase().includes(search.toLowerCase())
    )
  );

  const isAllSelected =
    selected.length === filteredData.length && filteredData.length > 0;

  return (
    <>
      <div className="flex items-center justify-between mb-5">
        <button
          className="text-primary text-sm border border-primary px-3 py-1.5 rounded-md hover:bg-primary hover:text-white"
          disabled={selected.length === 0}
        >
          Apply
        </button>
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
                  <span>ID</span>
                </div>
              </th>
              <th className="p-3 font-medium">Name</th>
              <th className="p-3 font-medium">Email</th>
              <th className="p-3 font-medium">Registration Date</th>
              <th className="p-3 font-medium">Course Taken</th>
              <th className="p-3 font-medium">Action</th>
            </tr>
          </thead>
          <tbody>
            {filteredData.map((std) => (
              <tr
                key={std.id}
                className="border-t border-gray-300 hover:bg-gray-100"
              >
                <td className="p-3">
                  <div className="flex items-center gap-5">
                    <Checkbox
                      checked={selected.includes(std.id)}
                      onChange={() => toggleSelectOne(std.id)}
                    />
                    {`${studentPrefix ?? ""}`}{std.id}
                  </div>
                </td>
                <td className="p-3">
                  <div className="flex items-center gap-5">
                    <div className="flex items-center gap-3">
                      <div className="flex items-center gap-2">
                        <span className="bg-primary text-white w-8 h-8 rounded-full flex items-center justify-center font-medium">
                          {std.first_name.slice(0, 1)}
                        </span>
                        <p className="font-medium">
                          {std.first_name} {std.last_name ?? ""}
                        </p>
                      </div>
                    </div>
                  </div>
                </td>
                <td className="p-3">{std.email}</td>
                <td className="p-3">
                  <div className="flex items-center gap-2">
                    {formatDate(std.created_at)}
                  </div>
                </td>
                <td className="p-3">{std.enrollments?.length || 0}</td>
                <td className="p-3">
                  <StudentAction id={std.id} />
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
