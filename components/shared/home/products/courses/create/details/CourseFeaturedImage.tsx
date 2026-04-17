"use client";

import React, { useRef, useState, useEffect } from "react";
import Image from "next/image";
import { BiImageAdd } from "react-icons/bi";
import { LuReplace, LuTrash2 } from "react-icons/lu";
import { useFormContext } from "react-hook-form";
import { TCourseSchema } from "@/schema/course.schema";

interface FeaturedImageProps {
  label?: string;
  desc?: string;
  onFileSelected: (file: File | null) => void;
}

export default function CourseFeaturedImage({
  label = "Upload Thumbnail",
  desc = "Image must be under 2MB",
  onFileSelected,
}: FeaturedImageProps) {
  const [thumbnail, setThumbnail] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement | null>(null);

  const { watch, setValue } = useFormContext<TCourseSchema>();
  const featuredImage = watch("featured_image");

  useEffect(() => {
    if (!thumbnail) {
      setPreview(null);
      return;
    }

    const reader = new FileReader();
    reader.onloadend = () => setPreview(reader.result as string);
    reader.readAsDataURL(thumbnail);
  }, [thumbnail]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setThumbnail(file);
      onFileSelected(file);
    }
  };

  const removeImage = () => {
    setThumbnail(null);
    onFileSelected(null);
    setValue("featured_image", null);
    // Reset the input so the same file can be selected again
    if (inputRef.current) {
      inputRef.current.value = "";
    }
  };

  const triggerInput = () => inputRef.current?.click();

  return (
    <>
      <input
        ref={inputRef}
        type="file"
        accept="image/*"
        onChange={handleChange}
        hidden
      />

      {featuredImage?.isDBImg ? (
        <ImagePreview
          imageUrl={featuredImage.name}
          onReplace={triggerInput}
          onRemove={removeImage}
        />
      ) : preview ? (
        <ImagePreview
          imageUrl={preview}
          onReplace={triggerInput}
          onRemove={removeImage}
        />
      ) : (
        <div
          onClick={triggerInput}
          className="border border-dashed rounded-lg h-44 bg-white flex flex-col items-center justify-center cursor-pointer p-4"
        >
          <BiImageAdd size={33} className="text-gray-400" />
          <p className="text-sm my-2 bg-blue-200 text-blue-700 font-medium px-2.5 py-1 rounded-md">
            {label}
          </p>
          <p className="text-xs text-gray-500 text-center">{desc}</p>
        </div>
      )}
    </>
  );
}

function ImagePreview({
  imageUrl,
  onReplace,
  onRemove,
}: {
  imageUrl: string;
  onReplace: () => void;
  onRemove: () => void;
}) {
  return (
    <div className="relative group border rounded-lg h-44 overflow-hidden p-2">
      <Image
        src={imageUrl}
        alt="Thumbnail"
        fill
        className="object-contain rounded-lg"
      />
      <div className="absolute inset-0 bg-black/40 flex items-center justify-center gap-3 opacity-0 group-hover:opacity-100 transition-opacity duration-300 rounded-md">
        <button
          type="button"
          onClick={onReplace}
          className="bg-white text-sm px-3 py-1.5 rounded"
        >
          <LuReplace size={16} />
        </button>
        <button
          type="button"
          onClick={onRemove}
          className="bg-red-500 text-white text-sm px-3 py-1.5 rounded"
        >
          <LuTrash2 size={16} />
        </button>
      </div>
    </div>
  );
}
