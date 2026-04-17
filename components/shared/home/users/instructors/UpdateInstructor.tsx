"use client";

import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import Modal from "@/components/ui/Modal";
import axiosInstance from "@/lib/axiosInstance";
import { zodResolver } from "@hookform/resolvers/zod";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { FaRegCircleCheck } from "react-icons/fa6";
import { GoPlus } from "react-icons/go";
import { RxCross2 } from "react-icons/rx";
import { toast } from "sonner";
import { z } from "zod";
import { useEditStore } from "@/hooks/useEditStore";
import FeaturedImage from "@/components/shared/FeaturedImage";
import TagInputField from "@/components/ui/TagsInputField";

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
  // phone: z
  //   .string({ required_error: "Phone number is required" })
  //   .trim()
  //   .min(11, { message: "Phone number must be at least 11 characters" })
  //   .max(11, { message: "Phone number must be at most 11 characters" })
  //   .startsWith("01", { message: "Invalid phone number" }),
  role: z
    .string()
    .trim()
    .max(100, { message: "Role should not exceed 100 characters" })
    .nullable()
    .or(z.literal("")),
  designation: z
    .string()
    .trim()
    .max(100, { message: "Role should not exceed 100 characters" })
    .nullable()
    .or(z.literal("")),
  image: z
    .any()
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;
      return file.size <= 2 * 1024 * 1024; // Check size
    }, "Max image size is 2MB.")
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;
      return ["image/png", "image/jpg", "image/jpeg"].includes(file.type); // Check file type
    }, "Only .png, .jpg & .jpeg formats are supported.")
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;
      return new Promise<boolean>((resolve) => {
        const img = document.createElement("img") as HTMLImageElement;
        img.src = URL.createObjectURL(file);
        img.onload = () => {
          const isValid = img.width <= 1000 && img.height <= 1000; // example dimensions
          resolve(isValid);
        };
        img.onerror = () => resolve(false); // in case of an error loading the image
      });
    }, "Image must be 1000x1000 pixels or smaller."),
});

type TRegisterSchema = z.infer<typeof RegisterSchema>;

export default function UpdateInstructor() {
  const { data: session } = useSession();
  const router = useRouter();
  const {
    editID,
    isInstructorEditOpen: isOpen,
    toggleInstructorEdit: setIsOpen,
  } = useEditStore();

  const {
    register,
    handleSubmit,
    reset,
    control,
    trigger,
    setValue,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<TRegisterSchema>({
    resolver: zodResolver(RegisterSchema),
  });

  useEffect(() => {
    if (editID && isOpen) {
      axiosInstance
        .get(`/private/instructor/details/${editID}`, {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${session?.accessToken}`,
          },
        })
        .then((res) => {
          const mappedImg = {
            name: res.data.data.image,
            size: 1339,
            type: "image/jpeg",
            isDBImg: true,
          };
          reset({ ...res.data.data, image: mappedImg });
        })
        .catch((error) => {
          toast.error(error.response.data.error || "Something went wrong.");
        });
    }
  }, [editID, isOpen]);

  const handleOnSubmit = async (data: TRegisterSchema) => {
    const fd = new FormData();
    fd.append("first_name", data.first_name);
    fd.append("last_name", data.last_name || "");
    fd.append("role", data.role || "");
    fd.append("designation", data.designation || "");
    if (data.image && !data.image.isDBImg) {
      fd.append("image", data.image);
    }

    axiosInstance
      .put(`/private/instructor/update/${editID}`, fd, {
        headers: {
          "Content-Type": "multipart/form-data",
          Authorization: `Bearer ${session?.accessToken}`,
        },
      })
      .then((res) => {
        setIsOpen(false);
        toast.success(res.data.message);
        router.refresh();
      })
      .catch((error) => {
        toast.error(error.response.data.error || "Something went wrong.");
      });
  };

  return (
    <>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <div className="flex items-center justify-between py-3 px-4 border-b border-gray-300">
          <p className="font-medium text-lg">Update Instructor</p>
          <button type="button" onClick={() => setIsOpen(false)}>
            <RxCross2 />
          </button>
        </div>
        <form onSubmit={handleSubmit(handleOnSubmit)} className="p-3">
          <div className="border rounded-xl">
            <div className="p-5">
              <div className="mb-4 w-[150px]">
                <label className="text-sm font-medium mb-1 block">
                  Profile Picture
                </label>
                <FeaturedImage
                  label="Upload"
                  dbImage={watch("image")}
                  onFileSelected={(file) => {
                    setValue("image", file, { shouldValidate: true });
                    trigger("image");
                  }}
                />
                {errors.image?.message && (
                  <p className="text-sm text-red-500 mt-1">
                    {String(errors.image.message)}
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
              <div className="flex items-center gap-5 mb-5">
                <div className="w-full">
                  <label className="text-sm font-medium mb-1 block">Role</label>
                  <InputField
                    placeholder="Role"
                    className="w-full"
                    {...register("role")}
                    error={errors.role?.message}
                  />
                </div>
                <div className="w-full">
                  <label className="text-sm font-medium mb-1 block">
                    Designation
                  </label>
                  <Controller
                    control={control}
                    name="designation"
                    render={({ field: { onChange, value } }) => (
                      <TagInputField onChange={onChange} value={value || ""} />
                    )}
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
