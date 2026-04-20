import Courses from "@/components/shared/home/products/courses/Courses";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import React from "react";

const getCourses = async (session: Session): Promise<CourseDetails[]> => {
  try {
    // #region agent log
    fetch("/api/debug-log", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Debug-Session-Id": "01d620",
      },
      body: JSON.stringify({
        sessionId: "01d620",
        runId: "pre-fix",
        hypothesisId: "H2",
        location: "courses/page.tsx:getCourses(entry)",
        message: "Server component fetching /private/course",
        data: { hasAccessToken: Boolean(session?.accessToken) },
        timestamp: Date.now(),
      }),
    }).catch(() => {});
    // #endregion agent log
    const response = await axiosInstance.get("/private/course", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });

    // #region agent log
    fetch("/api/debug-log", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Debug-Session-Id": "01d620",
      },
      body: JSON.stringify({
        sessionId: "01d620",
        runId: "pre-fix",
        hypothesisId: "H2",
        location: "courses/page.tsx:getCourses(success)",
        message: "Server component /private/course response",
        data: {
          status: response.status,
          count: Array.isArray(response.data?.data) ? response.data.data.length : null,
        },
        timestamp: Date.now(),
      }),
    }).catch(() => {});
    // #endregion agent log
    return response.data.data;
  } catch (error) {
    console.log(error);
    // #region agent log
    fetch("/api/debug-log", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Debug-Session-Id": "01d620",
      },
      body: JSON.stringify({
        sessionId: "01d620",
        runId: "pre-fix",
        hypothesisId: "H2",
        location: "courses/page.tsx:getCourses(error)",
        message: "Server component /private/course failed",
        data: {},
        timestamp: Date.now(),
      }),
    }).catch(() => {});
    // #endregion agent log
    return [];
  }
};

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const courseDetails = await getCourses(session);

  // console.log("courseDetails", courseDetails);
  

  return <Courses courses={courseDetails} />;
}
