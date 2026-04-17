"use client";

import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";
import Image from "next/image";
import PaymentAction from "./PaymentAction";

interface CategoryListProps {
  data: IPaymentMethods[] | null | undefined;
}

export default function PaymentList({ data }: CategoryListProps) {
  if (!Array.isArray(data)) {
    return (
      <div className="p-5 text-center text-gray-500">
        No payment methods found
      </div>
    );
  }

  /* Now we are 100% sure `data` is an array */
  const [selected, setSelected] = useState<number[]>([]);
  const [search, setSearch] = useState("");

  const filteredData = data.filter((item) =>
    item.title?.toLowerCase().includes(search.toLowerCase())
  );

  const toggleSelectAll = () => {
    if (selected.length === filteredData.length) {
      setSelected([]);
    } else {
      setSelected(filteredData.map((i) => i.id));
    }
  };

  const toggleSelectOne = (id: number) => {
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const isAllSelected =
    filteredData.length > 0 && selected.length === filteredData.length;

  return (
    <>
      {/* Top Bar */}
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

      {/* Table */}
      <div className="border rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-100">
            <tr className="text-left">
              <th className="p-3">
                <div className="flex items-center gap-3 font-medium">
                  <input
                    type="checkbox"
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                    className="w-4 h-4"
                  />
                  <span>Image</span>
                </div>
              </th>
              <th className="p-3 font-medium">Title</th>
              <th className="p-3 font-medium">Instruction</th>
              <th className="p-3 font-medium">Status</th>
              <th className="p-3 font-medium">Action</th>
            </tr>
          </thead>

          <tbody>
            {filteredData.length > 0 ? (
              filteredData.map((item) => (
                <tr
                  key={item.id}
                  className="border-t border-gray-300 hover:bg-gray-100"
                >
                  {/* Image + Checkbox */}
                  <td className="p-3">
                    <div className="flex items-center gap-3">
                      <input
                        type="checkbox"
                        checked={selected.includes(item.id)}
                        onChange={() => toggleSelectOne(item.id)}
                        className="w-4 h-4"
                      />

                      <div>
                        {item.image ? (
                          <Image
                            src={item.image}
                            alt="image"
                            width={64}
                            height={64}
                            className="size-8 object-contain rounded-md"
                          />
                        ) : (
                          "--"
                        )}
                      </div>
                    </div>
                  </td>

                  {/* Title */}
                  <td className="p-3">{item.title}</td>

                  {/* Instruction */}
                  <td className="p-3">{item.instruction}</td>
                  <td className="p-3">{item.status ? "Active" : "Inactive"}</td>

                  {/* Actions */}
                  <td className="p-3">
                    <PaymentAction id={item.id} />
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={7} className="p-5 text-center text-gray-500">
                  No payment method found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </>
  );
}
