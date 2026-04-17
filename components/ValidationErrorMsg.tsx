import { cn } from "@/lib/cn";
import React from "react";

export default function ValidationErrorMsg({
  error,
  className,
}: {
  error: string | undefined;
  className?: string;
}) {
  return <p className={cn("text-red-500 text-sm", className)}>{error}</p>;
}
