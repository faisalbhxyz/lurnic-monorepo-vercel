"use server";

import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";

export interface DashboardStats {
  students: number;
  courses: number;
  publishedCourses: number;
  enrollments: number;
  enrollmentsThisMonth: number;
  instructors: number;
  totalRevenue: number;
  revenueThisMonth: number;
  paidOrders: number;
  pendingPayments: number;
}

export interface DashboardData {
  stats: DashboardStats;
  recentEnrollments: Enrollment[];
  recentOrders: IOrder[];
}

const emptyStats: DashboardStats = {
  students: 0,
  courses: 0,
  publishedCourses: 0,
  enrollments: 0,
  enrollmentsThisMonth: 0,
  instructors: 0,
  totalRevenue: 0,
  revenueThisMonth: 0,
  paidOrders: 0,
  pendingPayments: 0,
};

function authHeaders(session: Session) {
  const token = session?.accessToken;
  if (!token) return null;
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };
}

function sortByNewest<T extends { created_at: string }>(items: T[]) {
  return [...items].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );
}

function startOfCurrentMonth() {
  const now = new Date();
  return new Date(now.getFullYear(), now.getMonth(), 1);
}

export async function getDashboardData(session: Session): Promise<DashboardData> {
  const headers = authHeaders(session);
  if (!headers) {
    return { stats: emptyStats, recentEnrollments: [], recentOrders: [] };
  }

  try {
    const [studentsRes, coursesRes, enrollmentsRes, ordersRes, instructorsRes] =
      await Promise.all([
        axiosInstance.get("/private/student", { headers }),
        axiosInstance.get("/private/course", { headers }),
        axiosInstance.get("/private/enrollment", { headers }),
        axiosInstance.get("/private/order", { headers }),
        axiosInstance.get("/private/instructor", { headers }),
      ]);

    const students: IStudent[] = studentsRes.data.data ?? [];
    const courses: CourseDetails[] = coursesRes.data.data ?? [];
    const enrollments: Enrollment[] = enrollmentsRes.data.data ?? [];
    const orders: IOrder[] = ordersRes.data.data ?? [];
    const instructors: IInstructor[] = instructorsRes.data.data ?? [];

    const monthStart = startOfCurrentMonth();
    const paidOrders = orders.filter((o) => o.payment_status === "paid");
    const totalRevenue = paidOrders.reduce((sum, o) => sum + (o.total || 0), 0);
    const revenueThisMonth = paidOrders
      .filter((o) => new Date(o.created_at) >= monthStart)
      .reduce((sum, o) => sum + (o.total || 0), 0);
    const enrollmentsThisMonth = enrollments.filter(
      (e) => new Date(e.created_at) >= monthStart
    ).length;

    return {
      stats: {
        students: students.length,
        courses: courses.length,
        publishedCourses: courses.filter((c) => c.visibility === "public").length,
        enrollments: enrollments.length,
        enrollmentsThisMonth,
        instructors: instructors.length,
        totalRevenue,
        revenueThisMonth,
        paidOrders: paidOrders.length,
        pendingPayments: orders.filter((o) => o.payment_status === "unpaid").length,
      },
      recentEnrollments: sortByNewest(enrollments).slice(0, 5),
      recentOrders: sortByNewest(orders).slice(0, 5),
    };
  } catch {
    return { stats: emptyStats, recentEnrollments: [], recentOrders: [] };
  }
}
