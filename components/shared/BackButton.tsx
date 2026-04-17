"use client";

import { useRouter } from "next/navigation";
import React from "react";
import { IoArrowBackOutline } from "react-icons/io5";

export default function BackButton({ buttonText }: { buttonText?: string }) {
  const router = useRouter();
  return (
    <button
      onClick={() => router.back()}
      className="flex items-center gap-2 text-gray-600"
    >
      <IoArrowBackOutline />
      {buttonText || "Back"}
    </button>
  );
}
