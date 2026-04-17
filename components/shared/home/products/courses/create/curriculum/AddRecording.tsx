import React, { useRef, useState } from "react";
import { LuFile } from "react-icons/lu";
import { RiUploadCloud2Line } from "react-icons/ri";
import { toast } from "sonner";
import { IoClose } from "react-icons/io5"; // Cross icon

const MAX_FILE_SIZE = 1024 * 1024 * 1024; // 1GB
const ACCEPTED_TYPES = ["video/mp4", "video/webm", "video/ogg"];

type AddRecordingProps = {
  onUploadComplete?: (file: File | null) => void;
};

const AddRecording: React.FC<AddRecordingProps> = ({ onUploadComplete }) => {
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!ACCEPTED_TYPES.includes(file.type)) {
      toast.error("Unsupported file type. Please upload MP4, WEBM, or OGG.");
      onUploadComplete?.(null);
      return;
    }

    if (file.size > MAX_FILE_SIZE) {
      toast.error("File size exceeds 1GB limit.");
      onUploadComplete?.(null);
      return;
    }

    setSelectedFile(file);
    onUploadComplete?.(file);
  };

  const handleCancel = () => {
    setSelectedFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = ""; // clear input value
    }
    onUploadComplete?.(null);
  };

  return (
    <div className="bg-gray-50 border rounded-xl min-h-40 p-4 flex flex-col items-center justify-center relative">
      <span className="p-2 bg-white text-primary border rounded-md mb-2">
        <LuFile size={22} />
      </span>

      {selectedFile ? (
        <div className="w-full flex flex-col items-center">
          <div className="flex items-center gap-2 bg-white border px-3 py-1.5 rounded-md text-gray-700 text-sm">
            <span className="truncate max-w-[200px]">{selectedFile.name}</span>
            <button
              type="button"
              onClick={handleCancel}
              className="text-red-500 hover:text-red-700"
            >
              <IoClose size={18} />
            </button>
          </div>
        </div>
      ) : (
        <>
          <p className="text-sm text-gray-600 font-medium my-3 text-center">
            We accept MP4, WEBM, and OGG files that are less than 1GB
          </p>
          <button
            type="button"
            onClick={handleButtonClick}
            className="bg-white flex items-center gap-2 border text-sm px-3 py-1.5 rounded-md text-primary font-medium"
          >
            <RiUploadCloud2Line /> Add Recording
          </button>
        </>
      )}

      <input
        type="file"
        accept=".mp4,.webm,.ogg"
        onChange={handleFileChange}
        ref={fileInputRef}
        hidden
      />
    </div>
  );
};

export default AddRecording;
