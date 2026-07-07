"use client";

import React from "react";
import Image from "next/image";
import Link from "next/link";
import { formatDate, formatDateTime } from "@/lib/helpers";
import StatCard from "@/components/shared/home/dashboard/StatCard";
import {
  HiOutlineAcademicCap,
  HiOutlineBookOpen,
  HiOutlineCash,
  HiOutlineClipboardList,
  HiOutlineMail,
  HiOutlinePhone,
  HiOutlineDeviceMobile,
  HiOutlineShoppingCart,
} from "react-icons/hi";
import { BiEditAlt } from "react-icons/bi";
import StudentDetailsActions from "./StudentDetailsActions";
import { useCurrency } from "@/context/CurrencyContext";

function fullName(student: IStudentDetails) {
  return `${student.first_name}${student.last_name ? ` ${student.last_name}` : ""}`;
}

function initials(student: IStudentDetails) {
  const first = student.first_name?.charAt(0) ?? "";
  const last = student.last_name?.charAt(0) ?? "";
  return (first + last).toUpperCase() || "?";
}

function PaymentStatusBadge({ status }: { status: "paid" | "unpaid" }) {
  const isPaid = status === "paid";
  return (
    <span
      className={`inline-flex text-xs font-medium px-2 py-0.5 rounded-full capitalize ${
        isPaid ? "bg-green-100 text-green-700" : "bg-amber-100 text-amber-700"
      }`}
    >
      {status}
    </span>
  );
}

export default function StudentDetailsView({
  student,
  studentPrefix,
}: {
  student: IStudentDetails;
  studentPrefix?: string | null;
}) {
  const { formatAmount } = useCurrency();
  const name = fullName(student);
  const orders = student.orders ?? [];

  return (
    <div className="space-y-6 my-5">
      <div className="flex-between flex-wrap gap-4">
        <h3 className="font-medium text-2xl">Student Details</h3>
        <StudentDetailsActions studentId={student.id} />
      </div>

      <div className="border rounded-xl bg-white p-6">
        <div className="flex flex-col sm:flex-row sm:items-center gap-5">
          {student.profile_image ? (
            <Image
              src={student.profile_image}
              alt={name}
              width={80}
              height={80}
              className="w-20 h-20 rounded-full object-cover border shrink-0"
            />
          ) : (
            <div className="w-20 h-20 rounded-full bg-primary text-white flex items-center justify-center text-2xl font-semibold shrink-0">
              {initials(student)}
            </div>
          )}
          <div className="flex-1 min-w-0">
            <div className="flex flex-wrap items-center gap-2">
              <h4 className="text-xl font-semibold text-gray-900">{name}</h4>
              <span
                className={`text-xs px-2 py-0.5 rounded-full ${
                  student.status
                    ? "bg-green-100 text-green-700"
                    : "bg-gray-100 text-gray-600"
                }`}
              >
                {student.status ? "Active" : "Inactive"}
              </span>
            </div>
            <p className="text-sm text-gray-500 mt-1">
              {studentPrefix ?? ""}
              {student.id} · Joined {formatDateTime(student.created_at)}
            </p>
            <div className="flex flex-wrap gap-4 mt-3 text-sm text-gray-600">
              <span className="inline-flex items-center gap-1.5">
                <HiOutlineMail className="text-gray-400" />
                {student.email}
              </span>
              {student.phone ? (
                <span className="inline-flex items-center gap-1.5">
                  <HiOutlinePhone className="text-gray-400" />
                  {student.phone}
                </span>
              ) : null}
            </div>
          </div>
        </div>
      </div>

      <div className="border rounded-xl bg-white overflow-hidden">
        <div className="px-5 py-4 border-b">
          <h4 className="font-medium text-lg">Active Device</h4>
          <p className="text-sm text-gray-500 mt-0.5">
            Only one device can access classes at a time. A new login replaces the previous session.
          </p>
        </div>
        {student.active_device ? (
          <div className="p-5 grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
            <div>
              <p className="text-gray-500">Device</p>
              <p className="font-medium text-gray-900 mt-1 inline-flex items-center gap-1.5">
                <HiOutlineDeviceMobile className="text-gray-400" />
                {student.active_device.device_name}
              </p>
            </div>
            <div>
              <p className="text-gray-500">IP address</p>
              <p className="font-medium text-gray-900 mt-1">
                {student.active_device.ip_address || "—"}
              </p>
            </div>
            <div>
              <p className="text-gray-500">Logged in</p>
              <p className="font-medium text-gray-900 mt-1">
                {formatDateTime(student.active_device.logged_in_at)}
              </p>
            </div>
            <div>
              <p className="text-gray-500">Last active</p>
              <p className="font-medium text-gray-900 mt-1">
                {formatDateTime(student.active_device.last_seen_at)}
              </p>
            </div>
            <div className="sm:col-span-2">
              <p className="text-gray-500">Device ID</p>
              <p className="font-mono text-xs text-gray-700 mt-1 break-all">
                {student.active_device.device_id}
              </p>
            </div>
          </div>
        ) : (
          <div className="p-8 text-center text-gray-500 text-sm">
            Not logged in on any device.
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
        <StatCard
          label="Enrolled Courses"
          value={student.stats.total_enrollments}
          icon={<HiOutlineBookOpen size={20} />}
          iconClassName="bg-blue-50 text-blue-600"
        />
        <StatCard
          label="Total Orders"
          value={student.stats.total_orders ?? 0}
          hint={`${student.stats.paid_orders ?? 0} paid · ${student.stats.unpaid_orders ?? 0} unpaid`}
          icon={<HiOutlineShoppingCart size={20} />}
          iconClassName="bg-slate-100 text-slate-600"
        />
        <StatCard
          label="Total Spent"
          value={formatAmount(student.stats.total_spent ?? 0)}
          hint="Paid orders only"
          icon={<HiOutlineCash size={20} />}
          iconClassName="bg-emerald-50 text-emerald-600"
        />
        <StatCard
          label="Average Progress"
          value={`${student.stats.average_progress}%`}
          hint="Lessons, quizzes & assignments"
          icon={<HiOutlineAcademicCap size={20} />}
          iconClassName="bg-violet-50 text-violet-600"
        />
      </div>

      <div className="border rounded-xl bg-white overflow-hidden">
        <div className="px-5 py-4 border-b">
          <h4 className="font-medium text-lg">Enrolled Courses</h4>
          <p className="text-sm text-gray-500 mt-0.5">
            {student.enrollments.length === 0
              ? "This student is not enrolled in any course yet."
              : `${student.enrollments.length} course${student.enrollments.length === 1 ? "" : "s"} enrolled`}
          </p>
        </div>

        {student.enrollments.length === 0 ? (
          <div className="p-8 text-center text-gray-500 text-sm">
            No enrollments found.
          </div>
        ) : (
          <div className="divide-y">
            {student.enrollments.map((enrollment) => (
              <div
                key={enrollment.id}
                className="p-5 flex flex-col lg:flex-row lg:items-center gap-4 hover:bg-gray-50"
              >
                <div className="flex items-center gap-4 flex-1 min-w-0">
                  <Image
                    src={enrollment.featured_image || "/images/placeholder.svg"}
                    alt={enrollment.title}
                    width={96}
                    height={64}
                    className="w-24 h-16 object-cover rounded-lg border shrink-0"
                  />
                  <div className="min-w-0">
                    <Link
                      href={`/courses/${enrollment.course_id}/update`}
                      className="font-medium text-gray-900 hover:text-primary truncate block"
                    >
                      {enrollment.title}
                    </Link>
                    <p className="text-sm text-gray-500 mt-1">
                      Enrolled on {formatDateTime(enrollment.enrolled_at)}
                    </p>
                    <p className="text-xs text-gray-400 mt-1">
                      Lessons: {enrollment.lessons_completed}/{enrollment.lessons_total}
                      {" · "}
                      Quizzes: {enrollment.quizzes_completed}/{enrollment.quizzes_total}
                      {" · "}
                      Assignments: {enrollment.assignments_completed}/
                      {enrollment.assignments_total}
                    </p>
                  </div>
                </div>

                <div className="w-full lg:w-48 shrink-0">
                  <div className="flex items-center justify-between text-sm mb-1.5">
                    <span className="text-gray-500">Progress</span>
                    <span className="font-medium">{enrollment.progress_percent}%</span>
                  </div>
                  <div className="h-2 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      className="h-full bg-primary rounded-full transition-all"
                      style={{ width: `${enrollment.progress_percent}%` }}
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="border rounded-xl bg-white overflow-hidden">
        <div className="px-5 py-4 border-b flex-between flex-wrap gap-2">
          <div>
            <h4 className="font-medium text-lg">Orders & Payments</h4>
            <p className="text-sm text-gray-500 mt-0.5">
              {orders.length === 0
                ? "No orders placed yet."
                : `${orders.length} order${orders.length === 1 ? "" : "s"} placed`}
            </p>
          </div>
          <Link
            href="/orders"
            className="text-sm text-primary hover:underline"
          >
            View all orders
          </Link>
        </div>

        {orders.length === 0 ? (
          <div className="p-8 text-center text-gray-500 text-sm">
            No orders found.
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 text-left">
                <tr>
                  <th className="p-3 font-medium">Invoice</th>
                  <th className="p-3 font-medium">Course</th>
                  <th className="p-3 font-medium">Ordered At</th>
                  <th className="p-3 font-medium">Payment</th>
                  <th className="p-3 font-medium">Method</th>
                  <th className="p-3 font-medium">Transaction ID</th>
                  <th className="p-3 font-medium">Total</th>
                </tr>
              </thead>
              <tbody>
                {orders.map((order) => (
                  <tr
                    key={order.id}
                    className="border-t border-gray-200 hover:bg-gray-50"
                  >
                    <td className="p-3 font-medium">#{order.invoice_id}</td>
                    <td className="p-3">
                      <div className="flex items-center gap-3 min-w-[200px]">
                        <Image
                          src={order.featured_image || "/images/placeholder.svg"}
                          alt={order.course_title}
                          width={64}
                          height={40}
                          className="w-16 h-10 object-cover rounded border shrink-0"
                        />
                        <div className="min-w-0">
                          <Link
                            href={`/courses/${order.course_id}/update`}
                            className="font-medium text-gray-900 hover:text-primary truncate block"
                          >
                            {order.course_title}
                          </Link>
                          {order.discount > 0 ? (
                            <p className="text-xs text-gray-400 mt-0.5">
                              Discount: {order.discount_type} ({order.discount})
                            </p>
                          ) : null}
                        </div>
                      </div>
                    </td>
                    <td className="p-3 whitespace-nowrap text-gray-600">
                      {formatDateTime(order.ordered_at)}
                    </td>
                    <td className="p-3">
                      <div className="space-y-1">
                        <PaymentStatusBadge status={order.payment_status} />
                        {order.payment_status === "paid" ? (
                          <p className="text-xs text-gray-400">
                            Paid {formatDateTime(order.updated_at)}
                          </p>
                        ) : (
                          <p className="text-xs text-gray-400">Awaiting payment</p>
                        )}
                      </div>
                    </td>
                    <td className="p-3 text-gray-600">
                      {order.payment_method || order.payment_type || "—"}
                    </td>
                    <td className="p-3 text-gray-600 font-mono text-xs">
                      {order.transaction_id || "—"}
                    </td>
                    <td className="p-3 font-medium whitespace-nowrap">
                      {formatAmount(order.total)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <StatCard
          label="Quiz Submissions"
          value={student.stats.quizzes_submitted}
          icon={<HiOutlineClipboardList size={20} />}
          iconClassName="bg-amber-50 text-amber-600"
        />
        <StatCard
          label="Assignment Submissions"
          value={student.stats.assignments_submitted}
          icon={<BiEditAlt size={20} />}
          iconClassName="bg-emerald-50 text-emerald-600"
        />
      </div>
    </div>
  );
}
