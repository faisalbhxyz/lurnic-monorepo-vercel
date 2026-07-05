import { z } from "zod";

export const profileImageSchema = z
  .any()
  .refine((file) => {
    if (!file) return true;
    if (file.isDBImg) return true;
    if (!(file instanceof File)) return false;
    return file.size <= 2 * 1024 * 1024;
  }, "Max image size is 2MB.")
  .refine((file) => {
    if (!file) return true;
    if (file.isDBImg) return true;
    if (!(file instanceof File)) return false;
    return ["image/png", "image/jpg", "image/jpeg"].includes(file.type);
  }, "Only .png, .jpg & .jpeg formats are supported.")
  .refine((file) => {
    if (!file) return true;
    if (file.isDBImg) return true;
    if (!(file instanceof File)) return false;
    return new Promise<boolean>((resolve) => {
      const img = document.createElement("img") as HTMLImageElement;
      img.src = URL.createObjectURL(file);
      img.onload = () => {
        resolve(img.width <= 1000 && img.height <= 1000);
      };
      img.onerror = () => resolve(false);
    });
  }, "Image must be 1000x1000 pixels or smaller.");
