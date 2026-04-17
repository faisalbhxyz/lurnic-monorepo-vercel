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

const categorySchema = z.object({
  name: z
    .string({ required_error: "Name is required." })
    .min(1, { message: "Name is required." })
    .max(100, { message: "Name should not exceed 100 characters" }),
  slug: z
    .string({
      required_error: "Slug is required.",
      invalid_type_error:
        "Slug is required and only lowercase english letters, numbers, and hyphens (-) are allowed",
    })
    .min(1, { message: "Slug is required." })
    .max(100, { message: "Slug should not exceed 100 characters" })
    .regex(/^[a-z0-9-]+$/, {
      message:
        "Slug can only contain lowercase english letters, numbers, and hyphens (-)",
    }),
  description: z.string().optional(),
});

type TCategorySchema = z.infer<typeof categorySchema>;

export default function CreateCategory({
  session,
  isEdit = false,
  category = null,
}: {
  session: Session;
  isEdit?: boolean;
  category?: ICategory | null;
}) {
  const router = useRouter();
  const [base64Image, setBase64Image] = useState<string | null>(null);
  const [selectOption, setSelectOption] = useState(options[0]);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
    watch,
    setValue,
  } = useForm<TCategorySchema>({
    resolver: zodResolver(categorySchema),
  });

  useEffect(() => {
    if (isEdit && category) {
      reset({
        name: category.name,
        slug: category.slug,
        description: category.description || "",
      });
    }
  }, [isEdit, category]);

  useEffect(() => {
    if (watch("name")) {
      setValue(
        "slug",
        watch("name")
          .toLowerCase()
          .replace(/\s+/g, "-") // replace spaces with hyphens
          .replace(/[^a-z0-9-]/g, "") // remove all special characters except hyphens
      );
    }
  }, [watch("name")]);

  const handleImageUpload = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setBase64Image(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSave = (data: TCategorySchema) => {
    if (isEdit) {
      axiosInstance
        .put(`/private/category/update/${category?.id}`, data, {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${session?.accessToken}`,
          },
        })
        .then((res) => {
          toast.success(res.data.message);
          reset();
          router.refresh();
          router.push("/categories");
        })
        .catch((error) => {
          console.log("[ERROR]", error);

          toast.error(error.response.data.error || "Something went wrong.");
        });
    } else {
      axiosInstance
        .post("/private/category/create", data, {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${session?.accessToken}`,
          },
        })
        .then((res) => {
          toast.success(res.data.message);
          reset();
          router.refresh();
          router.push("/categories");
        })
        .catch((error) => {
          console.log("[ERROR]", error);

          toast.error(error.response.data.error || "Something went wrong.");
        });
    }
  };

  return (
    <form onSubmit={handleSubmit(handleSave)}>
      <div className="flex-between px-5 py-3 bg-white border-b border-gray-300">
        <div className="w-full flex items-center justify-between gap-5">
          <h3 className="font-medium">{isEdit ? "Edit" : "Create"} Category</h3>
          <div className="flex items-center gap-3">
            <Button
              type="button"
              variant="secondary"
              className="px-4"
              onClick={() => {
                router.push("/categories");
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
          <div className="mb-5">
            <Label htmlFor={"name"}>Name</Label>
            <InputField
              id="name"
              className="w-full"
              {...register("name")}
              error={errors.name?.message}
            />
            <p className="text-xs mt-2 text-gray-700">
              The name is how it appears on your site.
            </p>
          </div>
          <div className="mb-5">
            <Label htmlFor={"slug"}>Slug</Label>
            <InputField
              id="slug"
              className="w-full"
              {...register("slug")}
              error={errors.slug?.message}
            />
            <p className="text-xs mt-2 text-gray-700">
              The “slug” is the URL-friendly version of the name. It is usually
              all lowercase and contains only letters, numbers, and hyphens.
            </p>
          </div>
          {/* <SelectList
            options={options}
            value={selectOption}
            onChange={setSelectOption}
          /> */}
          {/* <p className="text-xs mt-2 text-gray-700">
            Assign a parent term to create a hierarchy. The term Jazz, for
            example, would be the parent of Bebop and Big Band.
          </p> */}
          <div className="my-5">
            <Label htmlFor={"desc"}>Description</Label>
            <textarea
              id="desc"
              placeholder="Description"
              className="border outline-none rounded-md w-full h-24 min-h-24 bg-white text-sm px-3 py-2"
              {...register("description")}
            />
            {errors.description && (
              <p className="text-red-500 text-xs mt-1">
                {errors.description.message}
              </p>
            )}
            <p className="text-xs mt-2 text-gray-700">
              The description is not prominent by default; however, some themes
              may show it.
            </p>
          </div>
          {/* <div>
            <p className="text-sm mb-1">Thumbnail</p>
            <div className="flex items-center gap-3">
              <Image
                src={base64Image || "/images/placeholder.svg"}
                alt="Uploaded"
                width={100}
                height={50}
                className="w-10 h-10 object-cover rounded-md"
              />
              <label
                htmlFor="image-upload"
                className="border border-primary text-primary text-sm px-5 py-2 rounded-md cursor-pointer"
              >
                Upload image
                <input
                  type="file"
                  id="image-upload"
                  accept="image/*"
                  hidden
                  onChange={handleImageUpload}
                />
              </label>
            </div>
          </div> */}
        </div>
        <div className="p-5 w-96 min-w-96"></div>
      </div>
    </form>
  );
}
