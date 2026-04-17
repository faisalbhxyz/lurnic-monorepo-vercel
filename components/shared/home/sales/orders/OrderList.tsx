import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";
import OrderAction from "./OrderAction";
import Checkbox from "@/components/ui/Checkbox";
import { formatDate } from "@/lib/helpers";

interface OrderListProps {
  data: IOrder[];
}

export default function OrderList({ data }: OrderListProps) {
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
                  <span>ID</span>
                </div>
              </th>
              <th className="p-3 font-medium">Course</th>
              <th className="p-3 font-medium">Student</th>
              <th className="p-3 font-medium">Date</th>
              <th className="p-3 font-medium">Payment Status</th>
              <th className="p-3 font-medium">Payment Method</th>
              <th className="p-3 font-medium">Transaction ID</th>
              <th className="p-3 font-medium">Total</th>
              <th className="p-3 font-medium">Action</th>
            </tr>
          </thead>
          <tbody>
            {data.map((order) => (
              <tr
                key={order.id}
                className="border-t border-gray-300 hover:bg-gray-100"
              >
                <td className="p-3">
                  <div className="flex items-center gap-3">
                    <Checkbox
                      checked={selected.includes(order.id)}
                      onChange={() => toggleSelectOne(order.id)}
                    />
                    <div className="flex items-center gap-3">
                      <div>
                        <p className="font-medium">{order.invoice_id}</p>
                      </div>
                    </div>
                  </div>
                </td>
                <td className="p-3">{order.course.title}</td>
                <td className="p-3">{order.student.email}</td>
                <td className="p-3">
                  <div className="flex items-center gap-2">
                    {formatDate(order.created_at)}
                  </div>
                </td>
                <td className="p-3">{order.payment_status}</td>
                <td className="p-3">
                  {order.payment_method ? order.payment_method : "--"}
                </td>
                <td className="p-3">
                  {order.transaction_id ? order.transaction_id : "--"}
                </td>
                <td className="p-3">
                  <p>&#2547; {order.total}</p>
                </td>
                <td className="p-3">
                  <div className="flex items-center gap-2">
                    <OrderAction
                      isPaid={order.payment_status === "paid"}
                      id={order.id}
                    />
                  </div>
                </td>
              </tr>
            ))}
            {data.length === 0 && (
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
