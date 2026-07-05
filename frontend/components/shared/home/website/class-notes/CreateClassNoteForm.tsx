"use client";

import React, { useEffect } from "react";
import InputField from "@/components/ui/InputField";
import Label from "@/components/ui/Label";
import Button from "@/components/ui/Button";
import { Session } from "next-auth";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";
import { useRouter } from "next/navigation";

const classSchema = z.object({
  title: z.string().trim().min(1, "Title is required").max(150),
  slug: z.string().trim().max(180).optional(),
  icon_label: z.string().trim().max(10).optional(),
  icon_color: z.string().trim().max(20).optional(),
  position: z.coerce.number().gte(0),
  is_published: z.boolean(),
});

type ClassFormData = z.infer<typeof classSchema>;

export default function CreateClassNoteForm({
  session,
  isEdit = false,
  classData = null,
}: {
  session: Session;
  isEdit?: boolean;
  classData?: IAcademicNoteClass | null;
}) {
  const router = useRouter();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<ClassFormData>({
    resolver: zodResolver(classSchema),
    defaultValues: {
      is_published: true,
      position: 0,
    },
  });

  useEffect(() => {
    if (isEdit && classData) {
      reset({
        title: classData.title,
        slug: classData.slug,
        icon_label: classData.icon_label || "",
        icon_color: classData.icon_color || "",
        position: classData.position,
        is_published: classData.is_published,
      });
    }
  }, [isEdit, classData, reset]);

  const onSubmit = (data: ClassFormData) => {
    if (!session?.accessToken) {
      toast.error("Session expired. Please sign in again.");
      return;
    }

    const fd = new FormData();
    fd.append("title", data.title);
    fd.append("slug", data.slug || "");
    fd.append("icon_label", data.icon_label || "");
    fd.append("icon_color", data.icon_color || "");
    fd.append("position", String(data.position));
    fd.append("is_published", String(data.is_published));

    const headers = {
      Authorization: `Bearer ${session.accessToken}`,
    };

    if (isEdit && classData) {
      axiosInstance
        .put(`/private/academic-notes/classes/update/${classData.id}`, fd, {
          headers,
        })
        .then((res) => {
          toast.success(res.data.message);
          router.refresh();
        })
        .catch((error) => {
          toast.error(error.response?.data?.error || "Something went wrong.");
        });
    } else {
      axiosInstance
        .post("/private/academic-notes/classes/create", fd, { headers })
        .then((res) => {
          toast.success(res.data.message);
          router.push("/class-notes");
        })
        .catch((error) => {
          toast.error(error.response?.data?.error || "Something went wrong.");
        });
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 max-w-xl">
      <div>
        <Label htmlFor="title">Class Title</Label>
        <InputField
          id="title"
          placeholder="e.g. HSC, ৮ম শ্রেণি"
          {...register("title")}
          error={errors.title?.message}
        />
      </div>
      <div>
        <Label htmlFor="slug">Slug (optional)</Label>
        <InputField
          id="slug"
          placeholder="hsc"
          {...register("slug")}
          error={errors.slug?.message}
        />
      </div>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <Label htmlFor="icon_label">Icon Label</Label>
          <InputField
            id="icon_label"
            placeholder="H, ৮, ১ম"
            {...register("icon_label")}
            error={errors.icon_label?.message}
          />
        </div>
        <div>
          <Label htmlFor="icon_color">Icon Color</Label>
          <InputField
            id="icon_color"
            placeholder="#E91E63"
            {...register("icon_color")}
            error={errors.icon_color?.message}
          />
        </div>
      </div>
      <div>
        <Label htmlFor="position">Position</Label>
        <InputField
          id="position"
          type="number"
          {...register("position")}
          error={errors.position?.message}
        />
      </div>
      <label className="flex items-center gap-2 text-sm">
        <input type="checkbox" {...register("is_published")} />
        Published
      </label>
      <div className="flex gap-3">
        <Button type="submit" disabled={isSubmitting}>
          {isEdit ? "Update Class" : "Create Class"}
        </Button>
        <Button
          type="button"
          variant="secondary"
          onClick={() => router.push("/class-notes")}
        >
          Cancel
        </Button>
      </div>
    </form>
  );
}
