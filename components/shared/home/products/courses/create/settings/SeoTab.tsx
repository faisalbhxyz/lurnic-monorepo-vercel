import React, { useRef, useState } from "react";
import InputField from "@/components/ui/InputField";
import { FiUpload } from "react-icons/fi";
import { LuImage } from "react-icons/lu";
import Image from "next/image";

export default function SeoTab() {
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [fileName, setFileName] = useState<string | null>(null);
  const [fileSize, setFileSize] = useState<string | null>(null);

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const formatFileSize = (size: number) => {
    return size > 1024 * 1024
      ? `${(size / (1024 * 1024)).toFixed(2)} MB`
      : `${(size / 1024).toFixed(2)} KB`;
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file && file.type.startsWith("image/")) {
      const reader = new FileReader();
      reader.onloadend = () => {
        const base64 = reader.result as string;
        setImagePreview(base64);
        setFileName(file.name);
        setFileSize(formatFileSize(file.size));
      };
      reader.readAsDataURL(file);
    } else {
      alert("Please select a valid image file.");
    }
  };

  const handleDelete = () => {
    setImagePreview(null);
    setFileName(null);
    setFileSize(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  return (
    <div className="flex items-start gap-10">
      <div className="w-full">
        <div className="mt-5">
          <label className="text-sm font-medium mb-1">Meta Title</label>
          <InputField placeholder="Write meta title" className="w-full" />
          <p className="text-[13px] text-gray-500 mt-1">
            Tip: Keeping it under 60 characters helps it display nicely in
            search results.
          </p>
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium mb-1">Meta Description</label>
          <textarea
            rows={4}
            className="bg-white border w-full min-h-20 rounded-md px-3 py-2 outline-none focus:border-primary"
            placeholder="Write meta description"
          />
          <p className="text-[13px] text-gray-500 mt-1">
            Tip: Aim for under 160 characters to make it more engaging in search
            results.
          </p>
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium mb-1">OpenGraph Image</label>
          <div className="bg-gray-100 min-h-36 border rounded-xl flex flex-col items-center justify-center gap-2 p-4 relative">
            {imagePreview ? (
              <div className="w-full flex flex-col items-center">
                <Image
                  src={imagePreview}
                  alt="Preview"
                  width={400}
                  height={400}
                  className="w-80 h-auto object-contain rounded-md"
                />
                <div className="w-full flex items-center justify-between mt-5 text-sm">
                  <div>
                    <p className="font-medium">{fileName}</p>
                    <p>{fileSize}</p>
                  </div>
                  <button
                    onClick={handleDelete}
                    className="text-red-500 font-medium border px-3 py-1.5 rounded-md"
                  >
                    Delete
                  </button>
                </div>
              </div>
            ) : (
              <>
                <span className="bg-white p-2 rounded-md">
                  <LuImage />
                </span>
                <p className="text-sm text-gray-600">
                  Only image files are accepted.
                </p>
                <button
                  type="button"
                  onClick={handleUploadClick}
                  className="border mt-2 px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline"
                >
                  <FiUpload /> Upload
                </button>
                <input
                  type="file"
                  accept="image/*"
                  ref={fileInputRef}
                  onChange={handleFileChange}
                  hidden
                />
              </>
            )}
          </div>
          <p className="text-[13px] text-gray-500 mt-3">
            This image will be displayed in preview when your course is shared
            on social media, enhancing visibility and attracting potential
            students. (Recommended size: 1200x630 pixels)
          </p>
        </div>
      </div>

      <div className="w-80 min-w-80" />
    </div>
  );
}
