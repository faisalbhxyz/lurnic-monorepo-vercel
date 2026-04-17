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
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";
import { toast } from "sonner";
import { z } from "zod";

const formSchema = z.object({
  title: z.string().min(1, "Role Title is required"),
  permissions: z.array(z.string()),
});

type TFormSchema = z.infer<typeof formSchema>;

export default function AddNewRole() {
  const { data: session } = useSession();
  const [isOpen, setIsOpen] = useState(false);
  const router = useRouter();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<TFormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      permissions: [],
    },
  });

  const handleOnSave = async (data: TFormSchema) => {
    try {
      await axiosInstance.post("/role/create", data, {
        headers: {
          Authorization: `Bearer ${session?.accessToken}`,
        },
      });

      reset();
      router.refresh();
      toast.success("Role created successfully");
      setIsOpen(false);
    } catch (error) {
      toast.error("Failed to create role");
    }
  };

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>
        <LuPlus />
        Add New Role
      </Button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <div className="flex items-center justify-between py-3 px-4 border-b border-gray-300">
          <p className="font-medium text-lg">Create Role</p>
          <button type="button" onClick={() => setIsOpen(false)}>
            <RxCross2 />
          </button>
        </div>
        <form className="p-3" onSubmit={handleSubmit(handleOnSave)}>
          <div className="p-5">
            <div className="flex items-center gap-5 mb-5">
              <div className="w-full">
                <label className="text-sm font-medium mb-1 block">
                  Role Title <span className="text-red-500">*</span>
                </label>
                <InputField
                  placeholder="admin"
                  className="w-full"
                  {...register("title")}
                  error={errors.title?.message}
                />
              </div>
            </div>
          </div>
          <div className="flex items-center justify-end gap-3 mt-5">
            <button
              type="button"
              onClick={() => setIsOpen(false)}
              className="border text-sm font-medium px-4 py-2 rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
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
