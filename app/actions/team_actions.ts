"use server";

import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";

export const getAllUsers = async (session: Session): Promise<IUser[]> => {
  try {
    const res = await axiosInstance.get("/team-member/collection", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });

    return res.data.data;
  } catch (error) {
    console.log("[ERROR]", error);
    return [];
  }
};

export const getUserByID = async (
  session: Session,
  id: number
): Promise<IUser | null> => {
  try {
    const res = await axiosInstance.get(`/team-member/details/${id}`, {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });

    console.log(res);

    return res.data.data;
  } catch (error) {
    console.log("[ERROR]", error);
    return null;
  }
};

export const getAllRoles = async (session: Session): Promise<IRole[]> => {
  try {
    const res = await axiosInstance.get("/role/collection", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });

    return res.data.data;
  } catch (error) {
    console.log("[ERROR]", error);
    return [];
  }
};
