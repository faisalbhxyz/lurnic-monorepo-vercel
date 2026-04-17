"use client";

import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";
// import CouponsAction from "./CouponsAction";
import Link from "next/link";
import Image from "next/image";
import CategoryAction from "./CategoryAction";

interface CategoryListProps {
  data: ICategory[];
}

export default function CategoryList({ data }: CategoryListProps) {
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
                <div className="flex items-center gap-3 font-medium">
                  <input
                    type="checkbox"
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                    className="w-4 h-4"
                  />
                  <span>Name</span>
                </div>
              </th>
              {/* <th className="p-3 font-medium">Name</th> */}
              <th className="p-3 font-medium">Slug</th>
              <th className="p-3 font-medium">Description</th>
              <th className="p-3 font-medium">Action</th>
            </tr>
          </thead>
          <tbody>
            {filteredData.map((coupon) => (
              <tr
                key={coupon.id}
                className="border-t border-gray-300 hover:bg-gray-100"
              >
                <td className="p-3">
                  <div className="flex items-center gap-3">
                    <input
                      type="checkbox"
                      checked={selected.includes(coupon.id)}
                      onChange={() => toggleSelectOne(coupon.id)}
                      className="w-4 h-4"
                    />
                    <div className="flex items-center gap-3">
                      <div>
                        {coupon.name}
                      </div>
                    </div>
                  </div>
                </td>
                {/* <td className="p-3">{coupon.name}</td> */}
                <td className="p-3">{coupon.slug}</td>
                <td className="p-3">{coupon.description}</td>
                <td className="p-3">
                  <CategoryAction id={coupon.id} />
                </td>
              </tr>
            ))}
            {filteredData.length === 0 && (
              <tr>
                <td colSpan={7} className="p-5 text-center text-gray-500">
                  No coupons found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </>
  );
}
