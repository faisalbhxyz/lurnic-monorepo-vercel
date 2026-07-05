"use client";

import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";
import ClassNoteAction from "./ClassNoteAction";

export default function ClassNoteList({ data }: { data: IAcademicNoteClass[] }) {
  const [search, setSearch] = useState("");

  const filteredData = data.filter((item) =>
    item.title.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <>
      <div className="flex items-center justify-end mb-5">
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
              <th className="p-3 font-medium">Class</th>
              <th className="p-3 font-medium">Slug</th>
              <th className="p-3 font-medium">Notes</th>
              <th className="p-3 font-medium">Status</th>
              <th className="p-3 font-medium">Action</th>
            </tr>
          </thead>
          <tbody>
            {filteredData.map((item) => (
              <tr
                key={item.id}
                className="border-t border-gray-300 hover:bg-gray-100"
              >
                <td className="p-3">
                  <div className="flex items-center gap-3">
                    {item.icon_color && (
                      <div
                        className="w-8 h-8 rounded-md flex items-center justify-center text-white text-xs font-bold"
                        style={{ backgroundColor: item.icon_color }}
                      >
                        {item.icon_label || item.title.charAt(0)}
                      </div>
                    )}
                    <span className="font-medium">{item.title}</span>
                  </div>
                </td>
                <td className="p-3 text-gray-600">{item.slug}</td>
                <td className="p-3">{item.note_count ?? 0}</td>
                <td className="p-3">
                  <span
                    className={`px-2 py-0.5 rounded text-xs ${
                      item.is_published
                        ? "bg-green-100 text-green-700"
                        : "bg-gray-100 text-gray-600"
                    }`}
                  >
                    {item.is_published ? "Published" : "Draft"}
                  </span>
                </td>
                <td className="p-3">
                  <ClassNoteAction id={item.id} />
                </td>
              </tr>
            ))}
            {filteredData.length === 0 && (
              <tr>
                <td colSpan={5} className="p-5 text-center text-gray-500">
                  No classes found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </>
  );
}
