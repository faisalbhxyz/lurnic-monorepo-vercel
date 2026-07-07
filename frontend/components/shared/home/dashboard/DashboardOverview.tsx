"use client";

import Link from "next/link";
import { DashboardData } from "@/app/actions/dashboard_actions";
import StatCard from "./StatCard";
import RecentEnrollments from "./RecentEnrollments";
import RecentOrders from "./RecentOrders";
import {
  HiOutlineAcademicCap,
  HiOutlineBookOpen,
  HiOutlineCash,
  HiOutlineClipboardList,
  HiOutlineUserGroup,
  HiOutlineUsers,
} from "react-icons/hi";
import { LuPlus } from "react-icons/lu";
import Button from "@/components/ui/Button";
import { useCurrency } from "@/context/CurrencyContext";

const quickLinks = [
  { label: "Create Course", href: "/courses/create" },
  { label: "Manage Students", href: "/students" },
  { label: "View Enrollments", href: "/enrollments" },
  { label: "View Orders", href: "/orders" },
];

export default function DashboardOverview({
  data,
  userName,
}: {
  data: DashboardData;
  userName: string;
}) {
  const { stats, recentEnrollments, recentOrders } = data;
  const { formatAmount } = useCurrency();

  return (
    <div className="space-y-6">
      <div className="flex-between flex-wrap gap-4">
        <div>
          <p className="text-sm text-gray-500">Dashboard</p>
          <h3 className="font-medium text-2xl mt-1">Welcome back, {userName}</h3>
          <p className="text-sm text-gray-500 mt-1">
            Here&apos;s what&apos;s happening across your LMS today.
          </p>
        </div>
        <Button link src="/courses/create">
          <LuPlus /> New Course
        </Button>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-6 gap-4">
        <StatCard
          label="Students"
          value={stats.students}
          icon={<HiOutlineUserGroup size={20} />}
        />
        <StatCard
          label="Courses"
          value={stats.courses}
          hint={`${stats.publishedCourses} published`}
          icon={<HiOutlineBookOpen size={20} />}
          iconClassName="bg-blue-50 text-blue-600"
        />
        <StatCard
          label="Enrollments"
          value={stats.enrollments}
          hint={`${stats.enrollmentsThisMonth} this month`}
          icon={<HiOutlineAcademicCap size={20} />}
          iconClassName="bg-violet-50 text-violet-600"
        />
        <StatCard
          label="Instructors"
          value={stats.instructors}
          icon={<HiOutlineUsers size={20} />}
          iconClassName="bg-slate-100 text-slate-600"
        />
        <StatCard
          label="Total Revenue"
          value={formatAmount(stats.totalRevenue)}
          hint={`${formatAmount(stats.revenueThisMonth)} this month`}
          icon={<HiOutlineCash size={20} />}
          iconClassName="bg-green-50 text-green-600"
        />
        <StatCard
          label="Pending Payments"
          value={stats.pendingPayments}
          hint={`${stats.paidOrders} paid orders`}
          icon={<HiOutlineClipboardList size={20} />}
          iconClassName="bg-amber-50 text-amber-600"
        />
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-2 gap-4">
        <RecentEnrollments enrollments={recentEnrollments} />
        <RecentOrders orders={recentOrders} />
      </div>

      <div className="border rounded-xl p-4 bg-white">
        <h4 className="font-medium text-base mb-3">Quick Links</h4>
        <div className="flex flex-wrap gap-2">
          {quickLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className="text-sm px-3 py-1.5 rounded-md border border-gray-200 text-gray-700 hover:border-primary hover:text-primary transition-colors"
            >
              {link.label}
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
