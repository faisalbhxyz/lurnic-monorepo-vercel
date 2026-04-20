"use server";

import { signIn, signOut } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import axios from "axios";
import { Session } from "next-auth";
import { redirect } from "next/navigation";

export const doCretendentialLogin = async (email: string, password: string) => {
  try {
    const result = await signIn("credentials", {
      email,
      password,
      redirect: false,
    });

    // Auth.js/NextAuth signIn() can return an error payload without throwing.
    if (result && typeof result === "object" && "error" in result) {
      const signInError = (result as { error?: unknown }).error;
      if (typeof signInError === "string" && signInError.trim()) {
        return {
          error:
            signInError === "CredentialsSignin"
              ? "Invalid email or password."
              : signInError,
        };
      }
    }
  } catch (error: unknown) {
    const maybeAxiosErr = (error as any)?.cause?.err ?? error;
    const apiMessage = axios.isAxiosError(maybeAxiosErr)
      ? ((maybeAxiosErr.response?.data as any)?.message ??
          (maybeAxiosErr.response?.data as any)?.error)
      : undefined;

    const errCode =
      (axios.isAxiosError(maybeAxiosErr) ? maybeAxiosErr.code : undefined) ??
      (typeof (maybeAxiosErr as any)?.code === "string"
        ? (maybeAxiosErr as any).code
        : undefined);

    return {
      error:
        (typeof apiMessage === "string" && apiMessage.trim()) ||
        (errCode === "ECONNREFUSED" &&
          "API unreachable. Check API_INTERNAL_URL / NEXT_PUBLIC_API_URL and that the API is running.") ||
        (errCode === "ENOTFOUND" &&
          "API host not found. Check API_INTERNAL_URL / NEXT_PUBLIC_API_URL.") ||
        (errCode === "ERR_INVALID_URL" &&
          "API base URL invalid/missing. Check API_INTERNAL_URL / NEXT_PUBLIC_API_URL.") ||
        "Something went wrong.",
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
