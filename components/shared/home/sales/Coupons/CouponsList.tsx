import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";
// import CouponsAction from "./CouponsAction";
import Link from "next/link";
import Checkbox from "@/components/ui/Checkbox";

interface Coupons {
  id: number;
  name: string;
  discount: string;
  type: string;
  code: string;
  status: string;
  uses: number;
}

interface CouponsListProps {
  data: Coupons[];
}

export default function CouponsList({ data }: CouponsListProps) {
  const [selected, setSelected] = useState<number[]>([]);

  const toggleSelectAll = () => {
    if (selected.length === data.length) {
      setSelected([]);
    } else {
      setSelected(data.map((course) => course.id));
    }
  };

  const toggleSelectOne = (id: number) => {
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const isAllSelected = selected.length === data.length && data.length > 0;

  return (
    <>
      <div className="border rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-100">
            <tr className="text-left">
              <th className="p-3">
                <div className="flex items-center gap-3 font-medium">
                  <Checkbox
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                  />
                  <span>Name</span>
                </div>
              </th>
              <th className="p-3 font-medium">Discount</th>
              <th className="p-3 font-medium">Type</th>
              <th className="p-3 font-medium">Code</th>
              <th className="p-3 font-medium">Status</th>
              <th className="p-3 font-medium">Uses</th>
              <th className="p-3 font-medium">Action</th>
            </tr>
          </thead>
          <tbody>
            {data.map((coupon) => (
              <tr
                key={coupon.id}
                className="border-t border-gray-300 hover:bg-gray-100"
              >
                <td className="p-3">
                  <div className="flex items-center gap-3">
                    <Checkbox
                      checked={selected.includes(coupon.id)}
                      onChange={() => toggleSelectOne(coupon.id)}
                    />
                    <div className="flex items-center gap-3">
                      <div>
                        <p className="font-medium">{coupon.name}</p>
                      </div>
                    </div>
                  </div>
                </td>
                <td className="p-3">{coupon.discount}</td>
                <td className="p-3">{coupon.type}</td>
                <td className="p-3">{coupon.code}</td>
                <td className="p-3">{coupon.status}</td>
                <td className="p-3">{coupon.uses}</td>
                <td className="p-3">
                  <Link
                    href={`/coupons/${coupon.id}/edit`}
                    className="text-primary border border-primary px-3 py-1 rounded-md hover:bg-primary hover:text-white"
                  >
                    Edit
                  </Link>
                  {/* <CouponsAction /> */}
                </td>
              </tr>
            ))}
            {data.length === 0 && (
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
