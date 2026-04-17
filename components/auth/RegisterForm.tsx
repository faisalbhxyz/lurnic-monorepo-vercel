"use client";

import { useRouter } from "next/navigation";
import React, { useState, useTransition } from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";
import ValidationErrorMsg from "../ValidationErrorMsg";

const LoginSchema = z.object({
  name: z
    .string({ required_error: "Name is required." })
    .trim()
    .min(1, { message: "Name is required." })
    .max(100, { message: "Name should not exceed 100 characters" }),
  email: z
    .string({ required_error: "Email is required" })
    .trim()
    .email({ message: "Invalid email address" }),
  phone: z
    .string({ required_error: "Phone number is required" })
    .trim()
    .min(11, { message: "Phone number must be at least 11 characters" })
    .max(11, { message: "Phone number must be at most 11 characters" })
    .startsWith("01", { message: "Invalid phone number" }),
  password: z
    .string()
    .trim()
    .min(1, { message: "Password is required." })
    .min(6, { message: "Password must be at least 6 characters" }),
});

type TLoginSchema = z.infer<typeof LoginSchema>;

const RegisterForm = () => {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<TLoginSchema>({
    resolver: zodResolver(LoginSchema),
  });

  const handleOnSubmit = async (data: TLoginSchema) => {
    startTransition(() => {
      axiosInstance
        .post("/user/register", data, {
          headers: {
            "Content-Type": "application/json",
          },
        })
        .then((res) => {
          toast.success(res.data.message);
          router.push("/login");
          reset();
        })
        .catch((error) => {
          console.log("[ERROR]", error);

          toast.error(error.response.data.error || "Something went wrong.");
        });
    });
  };

  return (
    <form onSubmit={handleSubmit(handleOnSubmit)}>
      <div>
        <label className="block text-sm font-medium text-gray-700">Name</label>
        <input
          type="text"
          disabled={isPending}
          placeholder="Enter your name"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-[#00828a] focus:border-transparent"
          {...register("name")}
        />
        {errors.name && <ValidationErrorMsg error={errors.name.message} />}
      </div>
      <div className="mt-4">
        <label className="block text-sm font-medium text-gray-700">Email</label>
        <input
          type="email"
          disabled={isPending}
          placeholder="Enter your email"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-[#00828a] focus:border-transparent"
          {...register("email")}
        />
        {errors.email && <ValidationErrorMsg error={errors.email.message} />}
      </div>
      <div className="mt-4">
        <label className="block text-sm font-medium text-gray-700">Phone</label>
        <input
          type="text"
          disabled={isPending}
          placeholder="Enter your phone number"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-[#00828a] focus:border-transparent"
          {...register("phone")}
        />
        {errors.phone && <ValidationErrorMsg error={errors.phone.message} />}
      </div>
      <div className="mt-4">
        <div className="flex items-center justify-between">
          <label className="block text-sm font-medium text-gray-700">
            Password
          </label>
          {/* <Link href={"/forgot-password"} className="text-sm text-[#00828a]">
            Forgot Password
          </Link> */}
        </div>
        <input
          type="password"
          disabled={isPending}
          placeholder="Enter your password"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-[#00828a] focus:border-transparent"
          {...register("password")}
        />
        {errors.password && (
          <ValidationErrorMsg error={errors.password.message} />
        )}
      </div>
      <button
        type="submit"
        className="w-full cursor-pointer bg-[#00828a] text-white py-2 rounded-lg hover:bg-[#00828a] transition mt-3"
        disabled={isSubmitting}
      >
        Register
      </button>
    </form>
  );
};

export default RegisterForm;
