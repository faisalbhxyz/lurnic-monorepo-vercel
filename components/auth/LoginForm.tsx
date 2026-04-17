"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LuLoaderCircle } from "react-icons/lu";
import { doCretendentialLogin } from "@/app/actions/actions";
import { toast } from "sonner";
import ValidationErrorMsg from "../ValidationErrorMsg";

const LoginSchema = z.object({
  email: z
    .string({ required_error: "Email is required" })
    .email({ message: "Invalid email address" })
    .trim(),
  password: z
    .string()
    .trim()
    .min(1, { message: "Password is required." })
    .min(6, { message: "Password must be at least 6 characters" }),
});

type TLoginSchema = z.infer<typeof LoginSchema>;

const LoginForm = () => {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<TLoginSchema>({
    resolver: zodResolver(LoginSchema),
  });

  const handleOnSubmit = async (data: TLoginSchema) => {
    setLoading(true);
    const result = await doCretendentialLogin(data.email, data.password);

    if (result?.error) {
      toast.error(result.error);
      setLoading(false);
    } else {
      toast.success("Logged in successfully. Redirecting...");
      router.push("/");
    }
  };

  return (
    <form onSubmit={handleSubmit(handleOnSubmit)}>
      <div>
        <label className="block text-sm font-medium text-gray-700">Email</label>
        <input
          type="email"
          placeholder="Enter your email"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-[#00828a] focus:border-transparent"
          {...register("email")}
        />
        {errors.email && <ValidationErrorMsg error={errors.email.message} />}
      </div>
      <div className="mt-5">
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
        className="w-full cursor-pointer bg-[#00828a] text-white py-2 rounded-lg hover:bg-[#00828a] transition mt-8 flex items-center justify-center"
        disabled={isSubmitting}
      >
        {loading ? (
          <LuLoaderCircle size={24} className="animate-spin" />
        ) : (
          "Login"
        )}
      </button>
    </form>
  );
};

export default LoginForm;
