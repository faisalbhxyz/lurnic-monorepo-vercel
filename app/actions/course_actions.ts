"use server";

import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";

export const getCoursesLite = async (
  session: Session
): Promise<Pick<CourseDetails, "id" | "title">[]> => {
  try {
    const res = await axiosInstance.get("/private/course/lite", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    return [];
  }
};

export const getCourseByID = async (
  session: Session,
  id: string
): Promise<CourseDetails | null> => {
  try {
    const res = await axiosInstance.get(`/private/course/${id}`, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    return null;
  }
};
