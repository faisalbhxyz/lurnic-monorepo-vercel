"use client";

import FeaturedImage from "@/components/shared/FeaturedImage";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import Modal from "@/components/ui/Modal";
import { useEditStore } from "@/hooks/useEditStore";
import axiosInstance from "@/lib/axiosInstance";
import { zodResolver } from "@hookform/resolvers/zod";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { FaRegCircleCheck } from "react-icons/fa6";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";
import { toast } from "sonner";
import { z } from "zod";
import { profileImageSchema } from "./profileImageSchema";

const RegisterSchema = z.object({
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
  profile_image: profileImageSchema.optional(),
});

type TRegisterSchema = z.infer<typeof RegisterSchema>;

export default function UpdateStudent() {
  const { data: session } = useSession();
  const router = useRouter();
  const {
    editID,
    isStudentEditOpen: isOpen,
    toggleStudentEdit: setIsOpen,
  } = useEditStore();

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    trigger,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<TRegisterSchema>({
    resolver: zodResolver(RegisterSchema),
  });

  useEffect(() => {

    if (editID) {
      axiosInstance
        .get(`/private/student/details/${editID}`, {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${session?.accessToken}`,
          },
        })
        .then((res) => {
          const data = res.data.data;
          reset({
            first_name: data.first_name,
            last_name: data.last_name ?? "",
            email: data.email,
            phone: data.phone ?? "",
            profile_image: data.profile_image
              ? {
                  isDBImg: true,
                  name: data.profile_image,
                  type: "image/jpeg",
                }
              : undefined,
          });
        })
        .catch((error) => {
          toast.error(error.response.data.error || "Something went wrong.");
        });
    }
  }, [editID, reset]);

  const handleOnSubmit = async (data: TRegisterSchema) => {
    const fd = new FormData();
    fd.append("first_name", data.first_name);
    if (data.last_name) fd.append("last_name", data.last_name);
    if (data.phone) fd.append("phone", data.phone);
    if (data.profile_image && !data.profile_image.isDBImg) {
      fd.append("profile_image", data.profile_image);
    }

    axiosInstance
      .put(`/private/student/update/${editID}`, fd, {
        headers: {
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
        toast.error(error.response.data.error || "Something went wrong.");
      });
  };

  return (
    <>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <div className="flex items-center justify-between py-3 px-4 border-b border-gray-300">
          <p className="font-medium text-lg">Update Student</p>
          <button type="button" onClick={() => setIsOpen(false)}>
            <RxCross2 />
          </button>
        </div>
        <form onSubmit={handleSubmit(handleOnSubmit)} className="p-3">
          <div className="border rounded-xl">
            <div className="p-5">
              <div className="mb-5 max-w-xs">
                <label className="text-sm font-medium mb-1 block">
                  Profile Image
                </label>
                <FeaturedImage
                  label="Upload Profile Image"
                  dbImage={watch("profile_image")}
                  onFileSelected={(file) => {
                    setValue("profile_image", file, { shouldValidate: true });
                    trigger("profile_image");
                  }}
                />
                {errors.profile_image?.message && (
                  <p className="text-red-500 text-xs mt-1">
                    {String(errors.profile_image.message)}
                  </p>
                )}
              </div>
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
                  disabled
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
              Update
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
