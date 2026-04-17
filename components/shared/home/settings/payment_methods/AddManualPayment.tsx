import Button from "@/components/ui/Button";
import InputField from "@/components/ui/InputField";
import Label from "@/components/ui/Label";
import Modal from "@/components/ui/Modal";
import React, { useRef, useState } from "react";
import { LuPlus } from "react-icons/lu";
import { RxCross2 } from "react-icons/rx";
import { LuImagePlus } from "react-icons/lu";
import { z } from "zod";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { cn } from "@/lib/cn";
import axiosInstance from "@/lib/axiosInstance";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useSession } from "next-auth/react";
import FeaturedImage from "@/components/shared/FeaturedImage";
import { useEditStore } from "@/hooks/useEditStore";

const FormSchema = z.object({
  title: z
    .string({ required_error: "Title is required." })
    .trim()
    .min(1, { message: "Title is required." })
    .max(100, { message: "Title should not exceed 100 characters" }),
  instruction: z
    .string({ required_error: "Instruction is required." })
    .trim()
    .min(1, { message: "Instruction is required." }),
  image: z
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
      if (!(file instanceof File)) return false;
      return ["image/png", "image/jpg", "image/jpeg"].includes(file.type); // Check file type
    }, "Only .png, .jpg & .jpeg formats are supported.")
    .refine((file) => {
      if (!file) return true;
      if (file.isDBImg) return true;
      if (!(file instanceof File)) return false;

      return new Promise<boolean>((resolve) => {
        const img = document.createElement("img");
        img.src = URL.createObjectURL(file);
        img.onload = () => {
          const isValid = img.width <= 50 && img.height <= 50;
          resolve(isValid);
          URL.revokeObjectURL(img.src); // cleanup
        };
        img.onerror = () => {
          resolve(false);
          URL.revokeObjectURL(img.src); // cleanup
        };
      });
    }, "Image must be 48x48 pixels or smaller."),
});

type TFormSchema = z.infer<typeof FormSchema>;

export default function AddManualPayment() {
  const router = useRouter();
  const [isOpen, setIsOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const { data: session } = useSession();
  const { refreshPaymentMethods } = useEditStore();

  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<TFormSchema>({
    resolver: zodResolver(FormSchema),
  });

  const handleSave = (data: TFormSchema) => {
    if (loading) return;
    setLoading(true);
    const fd = new FormData();
    fd.append("title", data.title);
    fd.append("instruction", data.instruction);
    if (data.image) {
      fd.append("image", data.image);
    }

    axiosInstance
      .post("/private/payment-method/create", fd, {
        headers: {
          "Content-Type": "multipart/form-data",
          Authorization: `Bearer ${session?.accessToken}`,
        },
      })
      .then((res) => {
        toast.success(res.data.message);
        setIsOpen(false);
        setLoading(false);
        refreshPaymentMethods();
        router.refresh();
      })
      .catch((error) => {
        setLoading(false);
        toast.error(error.response.data.error || "Something went wrong.");
      });
  };

  return (
    <>
      <button
        type="button"
        onClick={() => {
          setIsOpen(true);
          setLoading(false);
          reset();
        }}
        className="text-sm hover:text-primary font-medium mt-3 flex items-center gap-1"
      >
        <LuPlus /> Add Manual Payment
      </button>
      <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} className="p-0">
        <form onSubmit={handleSubmit(handleSave)}>
          <div className="flex items-center justify-between px-4 py-3 border-b border-gray-300">
            <p className="font-medium text-gray-500">
              Set up manual payment method
            </p>
            <button type="button" onClick={() => setIsOpen(false)}>
              <RxCross2 />
            </button>
          </div>
          <div className="p-4 overflow-y-auto space-y-4">
            <Label htmlFor={""}>
              Title <span className="text-red-500">*</span>
            </Label>
            <InputField
              placeholder="e.g. Bank Transfer"
              className="w-full"
              {...register("title")}
              error={errors.title?.message}
            />
            <Label htmlFor={""}>Image</Label>
            <Controller
              control={control}
              name="image"
              render={({ field: { onChange, value } }) => (
                <FeaturedImage
                  onFileSelected={(file) => onChange(file)}
                  desc="Recommended size: 48x48"
                />
              )}
            />
            {errors.image && (
              <p className="text-red-500">{String(errors.image.message)}</p>
            )}

            <div>
              <Label htmlFor={""}>
                Payment Instructions <span className="text-red-500">*</span>
              </Label>
              <textarea
                className={cn(
                  "border w-full rounded-md min-h-32 p-2 focus:outline-none",
                  {
                    "border-red-500": errors.instruction,
                  }
                )}
                {...register("instruction")}
              />
              {errors.instruction && (
                <p className="text-red-500 mb-0 text-sm font-medium">
                  {String(errors.instruction.message)}
                </p>
              )}
            </div>
          </div>
          <div className=" bg-white rounded-b-2xl flex items-center justify-between px-4 py-3 border-t border-gray-300">
            <Button
              type="button"
              onClick={() => setIsOpen(false)}
              variant="secondary"
              disabled={loading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? "Saving..." : "Save"}
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
