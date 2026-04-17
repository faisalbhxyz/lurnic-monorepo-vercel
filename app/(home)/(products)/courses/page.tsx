import Courses from "@/components/shared/home/products/courses/Courses";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import React from "react";

const getCourses = async (session: Session): Promise<CourseDetails[]> => {
  try {
    const response = await axiosInstance.get("/private/course", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });

    return response.data.data;
  } catch (error) {
    console.log(error);
    
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
