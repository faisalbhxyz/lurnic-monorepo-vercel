"use client";

import React, { ChangeEvent, useEffect, useState } from "react";
import Image from "next/image";
import InputField from "@/components/ui/InputField";
import Label from "@/components/ui/Label";
import SelectList from "@/components/ui/SelectList";
import Button from "@/components/ui/Button";
import { Session } from "next-auth";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import FeaturedImage from "@/components/shared/FeaturedImage";

const options = [
  {
    id: 1,
    name: "courses 1",
  },
  {
    id: 2,
    name: "courses 2",
  },
];

const bannerSchema = z.object({
  title: z
    .string()
    .trim()
    .max(100, { message: "Title should not exceed 100 characters" })
    .optional()
    .nullable(),
  url: z
    .string()
    .trim()
    .url({ message: "Invalid URL format" })
    .or(z.literal("")),
  image: z
    .any()
    .refine((file) => {
      if (file && file.isDBImg) return true;
      if (!(file instanceof File)) return false;
      return file.size <= 2 * 1024 * 1024;
    }, "Max image size is 2MB.")
    .refine((file) => {
      // if (!file) return true; // Allow empty
      if (file && file.isDBImg) return true;
      return [
        "image/png",
        "image/jpg",
        "image/jpeg",
        "image/webp",
        "image/gif",
      ].includes(file.type); // Check file type
    }, "Only .png, .jpg & .jpeg formats are supported.")
    .refine((file) => {
      // if (!file) return true;
      if (file && file.isDBImg) return true;
      if (!(file instanceof File)) return false;

      return new Promise<boolean>((resolve) => {
        const img = document.createElement("img");
        img.src = URL.createObjectURL(file);
        img.onload = () => {
          const isValid = img.width <= 1920 && img.height <= 1080;
          resolve(isValid);
          URL.revokeObjectURL(img.src); // cleanup
        };
        img.onerror = () => {
          resolve(false);
          URL.revokeObjectURL(img.src); // cleanup
        };
      });
    }, "Image must be 1920x1080 pixels or smaller."),
});

type TCategorySchema = z.infer<typeof bannerSchema>;

export default function CreateBanner({
  session,
  isEdit = false,
  banner = null,
}: {
  session: Session;
  isEdit?: boolean;
  banner?: IBanner | null;
}) {
  const router = useRouter();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    trigger,
    watch,
    setValue,
  } = useForm<TCategorySchema>({
    resolver: zodResolver(bannerSchema),
  });

  useEffect(() => {
    if (isEdit && banner) {
      reset({
        title: banner.title || "",
        url: banner.url || "",
        image: {
          name: banner.image,
          size: 1339,
          type: "image/jpeg",
          isDBImg: true,
        },
      });
    }
  }, [isEdit, banner]);

  const handleSave = (data: TCategorySchema) => {
    const fd = new FormData();
    if (data.image && !data.image.isDBImg) {
      fd.append("image", data.image);
    }
    fd.append("title", data.title || "");
    fd.append("url", data.url || "");

    if (isEdit) {
      axiosInstance
        .put(`/private/banner/update/${banner?.id}`, fd, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${session?.accessToken}`,
          },
        })
        .then((res) => {
          toast.success(res.data.message);
          router.refresh();
          router.push("/banner");
        })
        .catch((error) => {
          console.log("[ERROR]", error);

          toast.error(error.response.data.error || "Something went wrong.");
        });
    } else {
      axiosInstance
        .post("/private/banner/create", fd, {
          headers: {
            "Content-Type": "multipart/form-data",
            Authorization: `Bearer ${session?.accessToken}`,
          },
        })
        .then((res) => {
          toast.success(res.data.message);
          reset();
          router.refresh();
          router.push("/banner");
        })
        .catch((error) => {
          toast.error(error.response.data.error || "Something went wrong.");
        });
    }
  };

  return (
    <form onSubmit={handleSubmit(handleSave)}>
      <div className="flex-between px-5 py-3 bg-white border-b border-gray-300">
        <div className="w-full flex items-center justify-between gap-5">
          <h3 className="font-medium">{isEdit ? "Edit" : "Create"} Banner</h3>
          <div className="flex items-center gap-3">
            <Button
              type="button"
              variant="secondary"
              className="px-4"
              onClick={() => {
                router.push("/banner");
                reset();
              }}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button type="submit" className="px-4" disabled={isSubmitting}>
              {isEdit ? "Update" : "Save"}
            </Button>
          </div>
        </div>
      </div>
      <div className="flex items-start">
        <div className="w-full p-5">
          <div>
            <p className="text-sm mb-1">Thumbnail</p>
            <FeaturedImage
              dbImage={watch("image")}
              onFileSelected={(file) => {
                setValue("image", file);
                trigger("image");
              }}
            />
            {errors.image && (
              <p className="text-red-500 text-sm mt-1">
                {String(errors.image.message)}
              </p>
            )}
          </div>
          <div className="mb-5 mt-5">
            <Label htmlFor={"name"}>Title</Label>
            <InputField
              id="name"
              className="w-full"
              {...register("title")}
              error={errors.title?.message}
            />
          </div>
          <div className="mb-5 mt-5">
            <Label htmlFor={"name"}>URL</Label>
            <InputField
              id="name"
              className="w-full"
              {...register("url")}
              error={errors.url?.message}
            />
          </div>
        </div>
        <div className="p-5 w-96 min-w-96"></div>
      </div>
    </form>
  );
}
