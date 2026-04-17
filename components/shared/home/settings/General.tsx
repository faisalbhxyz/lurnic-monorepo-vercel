"use client";

import React, { useEffect, useState } from "react";
import SelectPage from "./SelectPage";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import InputField from "@/components/ui/InputField";
import { IoMdRefresh } from "react-icons/io";
import { FiImage, FiTrash2 } from "react-icons/fi";
import { LuReplace } from "react-icons/lu";
import { z } from "zod";
import Image from "next/image";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import Button from "@/components/ui/Button";
import ValidationErrorMsg from "@/components/ValidationErrorMsg";
import axiosInstance from "@/lib/axiosInstance";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

const GeneralSettingsSchema = z.object({
  org_name: z
    .string({ required_error: "Name is required." })
    .trim()
    .min(1, { message: "Name is required." })
    .max(100, { message: "Name should not exceed 100 characters" }),
  student_prefix: z
    .string({ required_error: "Student prefix is required." })
    .trim()
    .min(1, { message: "Student prefix is required." })
    .max(10, { message: "Student prefix should not exceed 10 characters" }),
  teacher_prefix: z
    .string({ required_error: "Teacher prefix is required." })
    .trim()
    .min(1, { message: "Teacher prefix is required." })
    .max(10, { message: "Teacher prefix should not exceed 10 characters" }),
  logo: z
    .any()
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;
      return file.size <= 2 * 1024 * 1024;
    }, "Max image size is 2MB.")
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      return ["image/png", "image/jpg", "image/jpeg", "image/svg+xml"].includes(
        file.type
      ); // Check file type
    }, "Only .png, .jpg, .jpeg, .svg formats are supported.")
    .refine((file) => {
      if (!file) return true;
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;

      return new Promise<boolean>((resolve) => {
        const img = document.createElement("img");
        img.src = URL.createObjectURL(file);
        img.onload = () => {
          const isValid = img.width <= 400 && img.height <= 100;
          resolve(isValid);
          URL.revokeObjectURL(img.src); // cleanup
        };
        img.onerror = () => {
          resolve(false);
          URL.revokeObjectURL(img.src); // cleanup
        };
      });
    }, "Image must be 350x75 pixels or smaller."),
  favicon: z
    .any()
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;
      return file.size <= 2 * 1024 * 1024;
    }, "Max image size is 2MB.")
    .refine((file) => {
      if (!file) return true; // Allow empty
      if (file.isDBImg) return true;
      return ["image/png", "image/x-icon", "image/svg+xml"].includes(file.type); // Check file type
    }, "Only .png, .ico, .svg formats are supported.")
    .refine((file) => {
      if (!file) return true;
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;

      return new Promise<boolean>((resolve) => {
        const img = document.createElement("img");
        img.src = URL.createObjectURL(file);
        img.onload = () => {
          
          const isValid = img.width <= 75 && img.height <= 75;
          resolve(isValid);
          URL.revokeObjectURL(img.src); // cleanup
        };
        img.onerror = () => {
          resolve(false);
          URL.revokeObjectURL(img.src); // cleanup
        };
      });
    }, "Image must be 64x64 pixels or smaller."),
});

type TGeneralSettings = z.infer<typeof GeneralSettingsSchema>;

export default function General({
  generalSettings,
}: {
  generalSettings: GeneralSettings;
}) {
  const { data: session } = useSession();
  const router = useRouter();
  const [logoPreview, setLogoPreview] = useState<string | null>(null);
  const [faviconPreview, setFaviconPreview] = useState<string | null>(null);

  const {
    register,
    setValue,
    trigger,
    reset,
    formState: { errors, isSubmitting },
    handleSubmit,
  } = useForm<TGeneralSettings>({
    resolver: zodResolver(GeneralSettingsSchema),
  });

  const handleImageChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    type: "logo" | "favicon",
    setPreview: React.Dispatch<React.SetStateAction<string | null>>
  ) => {
    const file = e.target.files?.[0];

    console.log(file);

    if (file) {
      setValue(type, file, {
        shouldValidate: true,
      });
      trigger(type);
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreview(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSetNull = (type: "logo" | "favicon") => {
    setValue(type, null);
    type === "logo" ? setLogoPreview(null) : setFaviconPreview(null);
  };

  useEffect(() => {
    if (generalSettings) {
      if (generalSettings.logo) {
        setLogoPreview(generalSettings.logo);
      }
      if (generalSettings.favicon) {
        setFaviconPreview(generalSettings.favicon);
      }
      reset({
        org_name: generalSettings.org_name,
        student_prefix: generalSettings.student_prefix,
        teacher_prefix: generalSettings.teacher_prefix,
        logo: {
          name: generalSettings.logo,
          size: 100,
          type: "image/png",
          isDBImg: true,
        },
        favicon: {
          name: generalSettings.favicon,
          size: 100,
          type: "image/png",
          isDBImg: true,
        },
      });
    }
  }, [generalSettings]);

  const handleSave = (data: TGeneralSettings) => {
    const fd = new FormData();
    fd.append("org_name", data.org_name);
    fd.append("student_prefix", data.student_prefix);
    fd.append("teacher_prefix", data.teacher_prefix);
    if (data.logo && !data.logo.isDBImg) {
      fd.append("logo", data.logo);
    }
    if (data.favicon && !data.favicon.isDBImg) {
      fd.append("favicon", data.favicon);
    }

    axiosInstance
      .put("/private/general-settings/update", fd, {
        headers: {
          "Content-Type": "multipart/form-data",
          Authorization: `Bearer ${session?.accessToken}`,
        },
      })
      .then((res) => {
        toast.success(res.data.message);
        router.refresh();
      })
      .catch((error) => {
        toast.error(error.response.data.error || "Something went wrong.");
      });
  };

  return (
    <form onSubmit={handleSubmit(handleSave)}>
      {/* {JSON.stringify(errors, null, 2)} */}
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">General</p>
        {/* <button className="text-sm font-medium text-gray-500 flex items-center gap-1">
          <IoMdRefresh size={18} />
          Reset to Default
        </button> */}
        <Button type="submit" disabled={isSubmitting}>
          Save Changes
        </Button>
      </div>

      {/* Organization Info */}
      <div className="border bg-white p-4 rounded-md text-sm gap-3 mb-5 space-y-4">
        <div>
          <p className="font-medium text-gray-700 mb-1">Organization Name</p>
          <InputField
            placeholder="Enter your organization name"
            {...register("org_name")}
            error={errors.org_name?.message}
          />
        </div>

        <div className="flex items-center gap-6">
          {/* Logo Upload */}
          <div>
            <p className="font-medium text-gray-700 mb-2">Logo</p>
            <div className="w-[350px] h-[75px] rounded border relative group bg-gray-100 overflow-hidden">
              {logoPreview ? (
                <>
                  <Image
                    src={logoPreview}
                    alt="Logo"
                    width={350}
                    height={75}
                    className="w-full h-full object-contain"
                  />
                  {/* Hover Overlay */}
                  <div className="absolute inset-0 bg-black/70 opacity-0 group-hover:opacity-100 transition flex items-center justify-center gap-3 text-white text-sm">
                    <label className="flex items-center gap-1 cursor-pointer hover:text-blue-500">
                      <LuReplace size={16} />
                      <input
                        type="file"
                        accept="image/*"
                        className="hidden"
                        onChange={(e) =>
                          handleImageChange(e, "logo", setLogoPreview)
                        }
                      />
                    </label>
                    <button
                      onClick={() => handleSetNull("logo")}
                      className="flex items-center gap-1 hover:text-red-500"
                    >
                      <FiTrash2 size={16} />
                    </button>
                  </div>
                </>
              ) : (
                <label className="w-full h-full flex items-center justify-center cursor-pointer text-gray-400 hover:text-gray-600">
                  <FiImage size={24} />
                  <input
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={(e) =>
                      handleImageChange(e, "logo", setLogoPreview)
                    }
                  />
                </label>
              )}
            </div>
            {errors.logo && (
              <ValidationErrorMsg error={String(errors.logo?.message)} />
            )}
          </div>

          {/* Favicon Upload */}
          <div>
            <p className="font-medium text-gray-700 mb-2">Favicon</p>
            <div className="w-16 h-16 rounded border relative group bg-gray-100 overflow-hidden">
              {faviconPreview ? (
                <>
                  <Image
                    src={faviconPreview}
                    alt="Favicon"
                    width={64}
                    height={64}
                    className="w-full h-full object-contain"
                  />
                  {/* Hover Overlay */}
                  <div className="absolute inset-0 bg-black/70 opacity-0 group-hover:opacity-100 transition flex items-center justify-center gap-2 text-white text-xs">
                    <label className="flex items-center gap-1 cursor-pointer hover:text-blue-500">
                      <LuReplace size={16} />
                      <input
                        type="file"
                        accept="image/*"
                        className="hidden"
                        onChange={(e) =>
                          handleImageChange(e, "favicon", setFaviconPreview)
                        }
                      />
                    </label>
                    <button
                      onClick={() => handleSetNull("favicon")}
                      className="flex items-center gap-1 hover:text-red-500"
                    >
                      <FiTrash2 size={16} />
                    </button>
                  </div>
                </>
              ) : (
                <label className="w-full h-full flex items-center justify-center cursor-pointer text-gray-400 hover:text-gray-600">
                  <FiImage size={18} />
                  <input
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={(e) =>
                      handleImageChange(e, "favicon", setFaviconPreview)
                    }
                  />
                </label>
              )}
            </div>
            {errors.favicon && (
              <ValidationErrorMsg error={String(errors.favicon?.message)} />
            )}
          </div>
        </div>

        <div className="flex gap-4">
          <div className="flex-1">
            <p className="font-medium text-gray-700 mb-1">Student ID Prefix</p>
            <InputField
              placeholder="e.g., STD-"
              {...register("student_prefix")}
              error={errors.student_prefix?.message}
            />
          </div>
          <div className="flex-1">
            <p className="font-medium text-gray-700 mb-1">Teacher ID Prefix</p>
            <InputField
              placeholder="e.g., TCH-"
              {...register("teacher_prefix")}
              error={errors.teacher_prefix?.message}
            />
          </div>
        </div>
      </div>

      {/* Dashboard Page */}
      <div className="border bg-white p-4 rounded-md flex items-center justify-between text-sm gap-3 mb-5">
        <div>
          <p className="font-medium text-gray-700">Dashboard Page</p>
          <p className="text-gray-700 mt-1">
            This page will be used for student and instructor dashboard
          </p>
        </div>
        <SelectPage />
      </div>

      {/* Terms and Privacy Pages */}
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 flex items-center justify-between text-sm gap-3 border-b border-gray-300">
          <div>
            <p className="font-medium text-gray-700">
              Terms and Conditions Page
            </p>
            <p className="text-gray-700 mt-1">
              This page will be used as the Terms and Conditions page
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Privacy Policy</p>
            <p className="text-gray-700 mt-1">
              Choose the page for privacy policy.
            </p>
          </div>
          <SelectPage />
        </div>
      </div>

      {/* Other Settings */}
      <p className="text-gray-600 mt-5 mb-1">Others</p>
      <div className="border bg-white p-4 rounded-md flex items-center justify-between text-sm gap-3 mb-3">
        <div>
          <p className="font-medium text-gray-700">Enable Marketplace</p>
          <p className="text-gray-700 mt-1">
            Allow multiple instructors to sell their courses.
          </p>
        </div>
        <ToggleSwitch />
      </div>

      <div className="border bg-white p-4 rounded-md flex items-center justify-between text-sm gap-3 mb-5">
        <div>
          <p className="font-medium text-gray-700">Pagination</p>
          <p className="text-gray-700 mt-1">
            Set the number of rows to be displayed per page
          </p>
        </div>
        <InputField type="number" className="w-20" />
      </div>

      {/* Instructor Settings */}
      <p className="text-gray-600 mt-5 mb-1">Instructor</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Become an Instructor Button
            </p>
            <p className="text-gray-700 mt-1">
              Enable the option to display this button on the student dashboard.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Allow Instructors to Publish Courses
            </p>
            <p className="text-gray-700 mt-1">
              Enable instructors to publish the course directly. If disabled,
              admins will be able to review course content before publishing.
            </p>
          </div>
          <ToggleSwitch />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">
              Allow Instructors to Trash Courses
            </p>
            <p className="text-gray-700 mt-1">
              Enable this setting to allow instructors to delete courses.
            </p>
          </div>
          <ToggleSwitch />
        </div>
      </div>
    </form>
  );
}
