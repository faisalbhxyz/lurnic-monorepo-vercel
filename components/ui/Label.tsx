import { cn } from "@/lib/cn";
import React from "react";

interface LabelProps extends React.LabelHTMLAttributes<HTMLLabelElement> {
  htmlFor?: string;
  children: React.ReactNode;
  className?: string;
}

export default function Label({
  htmlFor,
  children,
  className,
  ...props
}: LabelProps) {
  return (
    <label
      htmlFor={htmlFor}
      {...props}
      className={cn("block mb-1 text-sm font-medium text-gray-700", className)}
    >
      {children}
    </label>
  );
}
