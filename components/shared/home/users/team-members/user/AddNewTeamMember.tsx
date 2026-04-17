"use client";

import { getAllRoles } from "@/app/actions/team_actions";
import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import Modal from "@/components/ui/Modal";
import SelectList from "@/components/ui/SelectList";
import axiosInstance from "@/lib/axiosInstance";
import { zodResolver } from "@hookform/resolvers/zod";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import React, { useEffect, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";
import { toast } from "sonner";
import { z } from "zod";

const formSchema = z.object({
  user_id: z
    .string({ required_error: "User id is required" })
    .min(1, "User id is required")
    .max(10, "User id must be at most 10 characters"),
  name: z
    .string({ required_error: "Name is required" })
    .min(1, "Name is required"),
  phone: z.string().trim().optional().nullable(),
  email: z
    .string()
    .trim()
    .email("Invalid email address")
    .min(1, "Email is required"),
  password: z
    .string({ required_error: "Password is required" })
    .trim()
    .min(1, "Password is required"),
  role: z.number({ required_error: "Role is required" }),
});

type IFormSchema = z.infer<typeof formSchema>;

export default function AddNewTeamMember() {
  const router = useRouter();
  const { data: session } = useSession();
  const [isOpen, setIsOpen] = useState(false);
  const [roles, setRoles] = useState<IRole[]>([]);

  const {
    control,
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<IFormSchema>({
    resolver: zodResolver(formSchema),
  });

  useEffect(() => {
    if (session && isOpen) {
      getAllRoles(session)
        .then((data) => setRoles(data))
        .catch(() => setRoles([]));
    }
  }, [session, isOpen]);

  const handleSave = async (data: IFormSchema) => {
    try {
      await axiosInstance.post("/team-member/create", data, {
        headers: {
          Authorization: `Bearer ${session?.accessToken}`,
        },
      });
      toast.success("Team member created successfully");
      setIsOpen(false);
      router.refresh();
    } catch (error) {
      console.log("[ERROR]", error);
      toast.error("Failed to create team member");
    }
  };

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>
        <LuPlus />
        Add New Team Member
      </Button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <div className="flex items-center justify-between py-3 px-4 border-b border-gray-300">
          <p className="font-medium text-lg">New Team Member</p>
          <button onClick={() => setIsOpen(false)}>
            <RxCross2 />
          </button>
        </div>
        <form className="p-3" onSubmit={handleSubmit(handleSave)}>
          <div className="p-5">
            <div className="flex items-center gap-5 mb-5">
              <div className="w-full">
                <label className="text-sm font-medium mb-1 block">
                  User ID <span className="text-red-500">*</span>
                </label>
                <InputField
                  placeholder="User ID"
                  className="w-full"
                  {...register("user_id")}
                  error={errors.user_id?.message}
                />
              </div>
              <div className="w-full">
                <label className="text-sm font-medium mb-1 block">
                  Name <span className="text-red-500">*</span>
                </label>
                <InputField
                  placeholder="Full Name"
                  className="w-full"
                  {...register("name")}
                  error={errors.name?.message}
                />
              </div>
            </div>
            <div className="flex items-center gap-5 mb-5">
              <div className="w-full">
                <label className="text-sm font-medium mb-1 block">
                  Add email address <span className="text-red-500">*</span>
                </label>
                <InputField
                  placeholder="Email"
                  className="w-full"
                  {...register("email")}
                  error={errors.email?.message}
                />
              </div>
              <div className="w-full">
                <label className="text-sm font-medium mb-1 block">
                  Add phone number
                </label>
                <InputField
                  placeholder="Phone number"
                  className="w-full"
                  {...register("phone")}
                  error={errors.phone?.message}
                />
              </div>
            </div>
            <div className="flex items-center gap-5 mb-5">
              <div className="w-full">
                <label className="text-sm font-medium mb-1 block">
                  Password
                </label>
                <InputField
                  type="password"
                  placeholder="Password"
                  className="w-full"
                  {...register("password")}
                  error={errors.password?.message}
                />
              </div>
              <div className="w-full">
                <label className="text-sm font-medium mb-1 block">Role</label>
                <Controller
                  control={control}
                  name="role"
                  render={({ field: { onChange, value } }) => (
                    <SelectList
                      options={
                        roles?.map((role) => ({
                          id: role.id,
                          name: role.name,
                          value: String(role.id),
                        })) || []
                      }
                      value={
                        roles
                          ? roles
                              .map((role) => ({
                                id: role.id,
                                name: role.name,
                                value: String(role.id),
                              }))
                              .find((role) => role.id === value) || null
                          : null
                      }
                      onChange={(d) => onChange(d.id)}
                      className="w-full"
                    />
                  )}
                />
                {errors.role && (
                  <p className="text-sm text-red-500 mt-1">
                    {errors.role.message}
                  </p>
                )}
              </div>
            </div>
          </div>
          <div className="flex items-center justify-end gap-3 mt-5">
            <button
              type="button"
              onClick={() => setIsOpen(false)}
              className="border text-sm font-medium px-4 py-2 rounded-full"
            >
              Cancel
            </button>
            <Button type="submit">Create</Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
