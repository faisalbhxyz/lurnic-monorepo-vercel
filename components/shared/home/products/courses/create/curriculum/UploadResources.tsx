"use client";

import React, { useRef, useState, useEffect } from "react";
import { toast } from "sonner";
import { RiUploadCloud2Line } from "react-icons/ri";
import { IoClose } from "react-icons/io5";
import axiosInstance from "@/lib/axiosInstance";
import { useSession } from "next-auth/react";

const MAX_FILE_SIZE = 50 * 1024 * 1024; // 50MB
const ACCEPTED_TYPES = [
  "application/pdf",
  "application/msword", // .doc
  "application/vnd.openxmlformats-officedocument.wordprocessingml.document", // .docx
  "application/zip",
  "application/x-zip-compressed",
];

export type UploadResourceItem =
  | File
  | {
      id: number;
      course_id: number;
      name: string;
      url: string;
      type: string;
      size: number;
      isDBImg: true;
    };

type UploadResourcesProps = {
  value?: UploadResourceItem[];
  onFilesSelected?: (files: UploadResourceItem[]) => void;
};

const UploadResources: React.FC<UploadResourcesProps> = ({
  value = [],
  onFilesSelected,
}) => {
  const { data: session } = useSession();
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [selectedFiles, setSelectedFiles] = useState<UploadResourceItem[]>([]);

  // Initialize selectedFiles from value prop
  useEffect(() => {
    if (value) {
      setSelectedFiles(value);
    }
  }, [value]);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    const validFiles: File[] = [];

    for (const file of files) {
      if (!ACCEPTED_TYPES.includes(file.type)) {
        toast.error(`${file.name} has an unsupported file type.`);
        continue;
      }

      if (file.size > MAX_FILE_SIZE) {
        toast.error(`${file.name} exceeds the 50MB limit.`);
        continue;
      }

      validFiles.push(file);
    }

    if (validFiles.length > 0) {
      const updatedFiles = [...selectedFiles, ...validFiles];
      setSelectedFiles(updatedFiles);
      onFilesSelected?.(updatedFiles);
    }

    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleDeleteFromDB = async (course_id: number, id: number) => {
    const isConfirm = confirm(
      "Are you sure you want to delete this payment method?"
    );
    if (!isConfirm) return Promise.reject("Cancelled");

    try {
      const res = await axiosInstance.delete(
        `/private/course/delete-resource/${course_id}/${id}`,
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${session?.accessToken}`,
          },
        }
      );

      toast.success(res.data.message);
      return Promise.resolve(res);
    } catch (error: any) {
      toast.error(error?.response?.data?.error || "Something went wrong.");
      return Promise.reject(error);
    }
  };

  const handleRemove = async (index: number) => {
    const fileToRemove = selectedFiles[index] as any;

    if (!fileToRemove) return;

    if (fileToRemove.isDBImg) {
      try {
        await handleDeleteFromDB(fileToRemove.course_id, fileToRemove.id);
      } catch {
        return;
      }
    }

    const updatedFiles = selectedFiles.filter((_, i) => i !== index);
    setSelectedFiles(updatedFiles);
    onFilesSelected?.(updatedFiles);
  };

  return (
    <div className="bg-gray-50 border rounded-xl min-h-40 p-4 flex flex-col items-center justify-center relative w-full">
      <p className="text-sm text-gray-600 font-medium mb-3 text-center">
        We accept PDF, Word, and ZIP files that are less than 50MB
      </p>

      <button
        type="button"
        onClick={handleButtonClick}
        className="bg-white flex items-center gap-2 border text-sm px-3 py-1.5 rounded-md text-primary font-medium mb-3"
      >
        <RiUploadCloud2Line /> Upload Resources
      </button>

      <input
        type="file"
        accept=".pdf,.doc,.docx,.zip"
        multiple
        onChange={handleFileChange}
        ref={fileInputRef}
        hidden
      />

      {selectedFiles.length > 0 && (
        <div className="w-full max-h-40 overflow-auto space-y-2">
          {selectedFiles.map((file, index) => (
            <div
              key={index}
              className="flex items-center justify-between bg-white border px-3 py-1.5 rounded-md text-gray-700 text-sm"
            >
              <span className="truncate max-w-[200px]">{file.name}</span>
              <button
                type="button"
                onClick={() => handleRemove(index)}
                className="text-red-500 hover:text-red-700"
              >
                <IoClose size={18} />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default UploadResources;
