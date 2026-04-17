"use server";

import { signIn, signOut } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import { redirect } from "next/navigation";

export const doCretendentialLogin = async (email: string, password: string) => {
  try {
    await signIn("credentials", {
      email,
      password,
      redirect: false,
    });
  } catch (error: any) {
    return {
      error: error.cause?.err.response.data.message || "Something went wrong.",
    };
  }
};

export const doCretendentialLogout = async () => {
  await signOut();
};

export const getAllCategories = async (
  session: Session
): Promise<ICategory[] | null> => {
  try {
    const res = await axiosInstance.get("/private/category", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    console.log("[ERROR]", error);

    return null;
  }
};

export const getAllSubCategories = async (
  session: Session
): Promise<ISubCategory[] | null> => {
  try {
    const res = await axiosInstance.get("/private/subcategory", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch (error) {
    console.log("[ERROR]", error);

    return null;
  }
};

export const getAllInstructorsLite = async (session: Session) => {
  try {
    const res = await axiosInstance.get("/private/instructor/lite", {
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

export const getStudentsLite = async (
  session: Session
): Promise<Pick<IStudent, "id" | "first_name" | "last_name">[]> => {
  try {
    const res = await axiosInstance.get("/private/student/lite", {
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

export const getGeneralSettings = async (
  session: Session
): Promise<GeneralSettings> => {
  try {
    const res = await axiosInstance.get("/private/general-settings", {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${session?.accessToken}`,
      },
    });
    return res.data.data;
  } catch {
    return {} as GeneralSettings;
  }
};

export const getPaymentMethods = async (
  session: Session
): Promise<IPaymentMethods[]> => {
  try {
    const res = await axiosInstance.get("/private/payment-method", {
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

export const getPaymentMethodByID = async (
  session: Session,
  id: number
): Promise<IPaymentMethods | null> => {
  try {
    const res = await axiosInstance.get(`/private/payment-method/${id}`, {
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
