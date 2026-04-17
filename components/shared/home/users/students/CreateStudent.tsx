"use client";

import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import Modal from "@/components/ui/Modal";
import axiosInstance from "@/lib/axiosInstance";
import { zodResolver } from "@hookform/resolvers/zod";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { FaRegCircleCheck } from "react-icons/fa6";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";
import { toast } from "sonner";
import { z } from "zod";

const RegisterSchema = z
  .object({
    first_name: z
      .string({ required_error: "Name is required." })
      .trim()
      .min(1, { message: "Name is required." })
      .max(100, { message: "Name should not exceed 100 characters" }),
    last_name: z
      .string({ required_error: "Name is required." })
      .trim()
      .max(100, { message: "Name should not exceed 100 characters" })
      .optional(),
    email: z
      .string({ required_error: "Email is required" })
      .trim()
      .email({ message: "Invalid email address" }),
    phone: z
      .string({ required_error: "Phone number is required" })
      .trim()
      .min(11, { message: "Phone number must be at least 11 characters" })
      .max(11, { message: "Phone number must be at most 11 characters" })
      .startsWith("01", { message: "Invalid phone number" })
      .optional()
      .nullable()
      .or(z.literal("")),
    password: z
      .string()
      .trim()
      .min(1, { message: "Password is required." })
      .min(6, { message: "Password must be at least 6 characters" }),
    confirm_password: z
      .string()
      .trim()
      .min(1, { message: "Password is required." })
      .min(6, { message: "Password must be at least 6 characters" }),
  })
  .refine((data) => data.password === data.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });

type TRegisterSchema = z.infer<typeof RegisterSchema>;

export default function CreateStudent() {
  const { data: session } = useSession();
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);
  const [hideAlert, setHideAlert] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<TRegisterSchema>({
    resolver: zodResolver(RegisterSchema),
  });

  const handleOnSubmit = async (data: TRegisterSchema) => {


    axiosInstance
      .post("/private/student/register", data, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${session?.accessToken}`,
        },
      })
      .then((res) => {
        setIsOpen(false);
        toast.success(res.data.message);
        router.refresh();
        reset();
      })
      .catch((error) => {
        console.log("[ERROR]", error);
        
        toast.error(error.response.data.error || "Something went wrong.");
      });
  };

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>
        <LuPlus /> Add Student
      </Button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <div className="flex items-center justify-between py-3 px-4 border-b border-gray-300">
          <p className="font-medium text-lg">Create Student</p>
          <button type="button" onClick={() => setIsOpen(false)}>
            <RxCross2 />
          </button>
        </div>
        <form onSubmit={handleSubmit(handleOnSubmit)} className="p-3">
          <div className="border rounded-xl">
            <div className="border-b p-5 border-gray-300">
              <p className="font-medium">Create Student</p>
              {!hideAlert && (
                <div className="bg-green-50 p-4 border border-green-400 rounded-xl flex items-start mt-3">
                  <span className="p-1.5 text-green-800">
                    <FaRegCircleCheck />
                  </span>
                  <p className="text-sm font-medium ml-3 text-green-800">
                    When you create a new student, they will receive an email
                    notification.
                  </p>
                  <button
                    type="button"
                    onClick={() => setHideAlert(true)}
                    className="p-1 text-green-500"
                  >
                    <RxCross2 />
                  </button>
                </div>
              )}
            </div>
            <div className="p-5">
              <div className="flex items-center gap-5 mb-5">
                <div className="w-full">
                  <label className="text-sm font-medium mb-1 block">
                    First Name <span className="text-red-500">*</span>
                  </label>
                  <InputField
                    placeholder="First Name"
                    className="w-full"
                    {...register("first_name")}
                    error={errors.first_name?.message}
                  />
                </div>
                <div className="w-full">
                  <label className="text-sm font-medium mb-1 block">
                    Last name
                  </label>
                  <InputField
                    placeholder="Last Name"
                    className="w-full"
                    {...register("last_name")}
                    error={errors.last_name?.message}
                  />
                </div>
              </div>
              <div className="w-full mb-5">
                <label className="text-sm font-medium mb-1 block">
                  Email <span className="text-red-500">*</span>
                </label>
                <InputField
                  placeholder="Email"
                  className="w-full"
                  {...register("email")}
                  error={errors.email?.message}
                />
              </div>
              <div className="w-full mb-5">
                <label className="text-sm font-medium mb-1 block">Phone</label>
                <InputField
                  placeholder="Phone"
                  className="w-full"
                  {...register("phone")}
                  error={errors.phone?.message}
                />
              </div>
              <div className="flex items-center gap-5 mb-5">
                <div className="w-full">
                  <label className="text-sm font-medium mb-1 block">
                    Password <span className="text-red-500">*</span>
                  </label>
                  <InputField
                    placeholder="+++++++"
                    type="password"
                    className="w-full"
                    {...register("password")}
                    error={errors.password?.message}
                  />
                </div>
                <div className="w-full">
                  <label className="text-sm font-medium mb-1 block">
                    Confirm Password
                  </label>
                  <InputField
                    placeholder="+++++++"
                    type="password"
                    className="w-full"
                    {...register("confirm_password")}
                    error={errors.confirm_password?.message}
                  />
                </div>
              </div>
            </div>
          </div>
          <div className="flex items-center justify-end gap-3 mt-5">
            <button
              type="button"
              onClick={() => setIsOpen(false)}
              className="border text-sm font-medium px-4 py-2 rounded-full"
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <Button type="submit" disabled={isSubmitting}>
              Create
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
