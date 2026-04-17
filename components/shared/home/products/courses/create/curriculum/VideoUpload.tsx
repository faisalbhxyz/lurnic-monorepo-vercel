import React, { useRef, useState } from "react";
import { ImFilePlay } from "react-icons/im";
import { RiUploadCloud2Line } from "react-icons/ri";
import { toast } from "sonner";

const MAX_FILE_SIZE_MB = 1024;
const ACCEPTED_FORMATS = ["video/mp4", "video/webm", "video/ogg"];

const VideoUpload: React.FC = () => {
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploading, setUploading] = useState(false);
  const [videoURL, setVideoURL] = useState<string | null>(null);

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (!ACCEPTED_FORMATS.includes(file.type)) {
      toast.error("Unsupported file format. Please upload MP4, WEBM, or OGG.");
      return;
    }

    if (file.size > MAX_FILE_SIZE_MB * 1024 * 1024) {
      toast.error("File is too large. Max allowed size is 1GB.");
      return;
    }

    toast.success(`Selected: ${file.name}`);
    simulateUpload(file);
  };

  const simulateUpload = (file: File) => {
    setUploading(true);
    setUploadProgress(0);
    setVideoURL(null);

    let progress = 0;
    const interval = setInterval(() => {
      progress += 10;
      setUploadProgress(progress);

      if (progress >= 100) {
        clearInterval(interval);
        setUploading(false);
        toast.success("Upload completed!");
        const url = URL.createObjectURL(file);
        setVideoURL(url);
      }
    }, 200);
  };

  return (
    <div className="bg-gray-50 border rounded-xl p-4 flex flex-col items-center justify-center w-full">
      <input
        type="file"
        accept={ACCEPTED_FORMATS.join(",")}
        onChange={handleFileChange}
        ref={fileInputRef}
        hidden
      />

      <span className="p-2 bg-white text-primary border rounded-md">
        <ImFilePlay size={22} />
      </span>
      <p className="text-sm text-gray-600 font-medium my-3 text-center">
        We accept MP4, WEBM, and OGG files that are less than 1GB
      </p>

      <button
        onClick={handleClick}
        disabled={uploading}
        className={`bg-white flex items-center gap-2 border text-sm px-3 py-1.5 rounded-md font-medium ${
          uploading ? "text-gray-400 cursor-not-allowed" : "text-primary"
        }`}
      >
        <RiUploadCloud2Line />
        {uploading ? "Uploading..." : "Upload Video"}
      </button>

      {uploading && (
        <div className="w-full mt-4">
          <div className="w-full bg-gray-200 rounded-full h-2.5 overflow-hidden">
            <div
              className="bg-primary h-2.5 transition-all duration-200"
              style={{ width: `${uploadProgress}%` }}
            ></div>
          </div>
          <p className="text-xs text-gray-500 text-center mt-1">
            Uploading: {uploadProgress}%
          </p>
        </div>
      )}

      {videoURL && (
        <div className="mt-6 w-full">
          <video
            src={videoURL}
            controls
            className="rounded-lg w-full max-h-64 object-contain"
          />
        </div>
      )}
    </div>
  );
};

export default VideoUpload;
