"use client";

import { cn } from "@/lib/cn";
import RichTextEditor from "./RichTextEditor";

export default function TextEditor({
  value,
  className,
  onChange,
  errMessage,
}: {
  value?: string;
  className?: string;
  onChange: (html: string) => void;
  errMessage?: string;
}) {
  return (
    <div className={cn("w-full mx-auto", className)}>
      <RichTextEditor content={value || ""} onChange={onChange} />
      {errMessage && <p className="text-red-500">{errMessage}</p>}
    </div>
  );
}
