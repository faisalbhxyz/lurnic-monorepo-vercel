import Link from "next/link";
import { formatDate } from "@/lib/helpers";

function formatMoney(amount: number) {
  return `৳ ${amount.toLocaleString("en-BD")}`;
}

export default function RecentOrders({ orders }: { orders: IOrder[] }) {
  return (
    <div className="border rounded-xl overflow-hidden bg-white">
      <div className="flex-between px-4 py-3 border-b border-gray-200">
        <h4 className="font-medium text-base">Recent Payments</h4>
        <Link href="/orders" className="text-sm text-primary hover:underline">
          View all
        </Link>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-gray-50">
            <tr className="text-left text-gray-500">
              <th className="p-3 font-medium">Invoice</th>
              <th className="p-3 font-medium">Student</th>
              <th className="p-3 font-medium">Amount</th>
              <th className="p-3 font-medium">Status</th>
              <th className="p-3 font-medium text-end">Date</th>
            </tr>
          </thead>
          <tbody>
            {orders.map((order) => (
              <tr key={order.id} className="border-t border-gray-200">
                <td className="p-3 font-medium text-gray-900">{order.invoice_id}</td>
                <td className="p-3">
                  <p className="text-gray-700">{order.student.email}</p>
                  <p className="text-xs text-gray-500">{order.course.title}</p>
                </td>
                <td className="p-3 text-gray-700">{formatMoney(order.total)}</td>
                <td className="p-3">
                  <span
                    className={`inline-flex px-2 py-0.5 rounded-full text-xs font-medium capitalize ${
                      order.payment_status === "paid"
                        ? "bg-green-100 text-green-700"
                        : "bg-amber-100 text-amber-700"
                    }`}
                  >
                    {order.payment_status}
                  </span>
                </td>
                <td className="p-3 text-end text-gray-600">
                  {formatDate(order.created_at)}
                </td>
              </tr>
            ))}
            {orders.length === 0 && (
              <tr>
                <td colSpan={5} className="p-6 text-center text-gray-500">
                  No orders yet.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
