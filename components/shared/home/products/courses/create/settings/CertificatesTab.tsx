import InputField from "@/components/ui/InputField";
import Image from "next/image";
import React, { useRef, useState } from "react";
import { FiUpload } from "react-icons/fi";
import { IoIosArrowDown, IoIosArrowUp } from "react-icons/io";

const certificates = [
  { id: 1, name: "/images/Certificat-14.jpg" },
  { id: 2, name: "/images/Certificat-15.jpg" },
  { id: 3, name: "/images/Certificat-16.jpg" },
  { id: 4, name: "/images/Certificat-17.jpg" },
];

export default function CertificatesTab() {
  const [isOpen, setIsOpen] = useState(true);
  const [selectedCertificate, setSelectedCertificate] = useState(
    certificates[0].name
  );
  const [ownerSignatureImage, setOwnerSignatureImage] = useState<string | null>(
    null
  );
  const [instructorSignatureImage, setInstructorSignatureImage] = useState<
    string | null
  >(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const fileInputRef2 = useRef<HTMLInputElement>(null);

  const handleOwnerFileChange = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setOwnerSignatureImage(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const triggerOwnerFileInput = () => {
    fileInputRef.current?.click();
  };
  const handleInstructorFileChange = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setInstructorSignatureImage(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const triggerInstructorFileInput = () => {
    fileInputRef2.current?.click();
  };

  return (
    <div className="flex items-start gap-10">
      <div className="w-full">
        <p className="font-medium mb-2">Choose a design</p>
        <div className="border rounded-md overflow-hidden mb-4">
          <Image
            src={selectedCertificate}
            alt="Selected Certificate"
            width={500}
            height={500}
            className="w-full h-auto object-cover"
          />
        </div>

        <div className="flex items-center justify-center gap-3">
          {certificates.map((image) => (
            <button
              key={image.id}
              onClick={() => setSelectedCertificate(image.name)}
              className={`w-12 h-10 border-2 rounded-md overflow-hidden p-0 ${
                selectedCertificate === image.name
                  ? "border-primary"
                  : "border-transparent"
              }`}
            >
              <Image
                src={image.name}
                alt={`Certificate ${image.id}`}
                width={100}
                height={100}
                className="w-full h-full"
              />
            </button>
          ))}
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium mb-1">Certificate Title</label>
          <InputField className="w-full" />
        </div>
        <div className="mt-5">
          <label className="text-sm font-medium mb-1">
            Certificate Subtitle One
          </label>
          <InputField className="w-full" />
        </div>
        <div className="mt-5">
          <label className="text-sm font-medium mb-1">
            Certificate Subtitle Two
          </label>
          <InputField className="w-full" />
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium block mb-2">
            School Owner Signature (150x250 px)
          </label>
          <div className="flex items-center gap-4">
            {ownerSignatureImage && (
              <Image
                src={ownerSignatureImage}
                alt="Uploaded"
                width={100}
                height={100}
                className="w-20 h-20 object-cover rounded"
              />
            )}
            <button
              type="button"
              onClick={triggerOwnerFileInput}
              className="border px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline"
            >
              <FiUpload /> Upload
            </button>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              onChange={handleOwnerFileChange}
              hidden
            />
          </div>
        </div>
        <div className="mt-5">
          <label className="text-sm font-medium block mb-2">
            Instructor Signature (150x250 px)
          </label>
          <div className="flex items-center gap-4">
            {instructorSignatureImage && (
              <Image
                src={instructorSignatureImage}
                alt="Uploaded"
                width={100}
                height={100}
                className="w-20 h-20 object-cover rounded"
              />
            )}
            <button
              type="button"
              onClick={triggerInstructorFileInput}
              className="border px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline"
            >
              <FiUpload /> Upload
            </button>
            <input
              ref={fileInputRef2}
              type="file"
              accept="image/*"
              onChange={handleInstructorFileChange}
              hidden
            />
          </div>
        </div>
      </div>

      <div className="w-80 min-w-80">
        <div className="border p-5 rounded-md">
          <button
            onClick={() => setIsOpen(!isOpen)}
            className={`flex items-center justify-between w-full ${
              isOpen ? "text-primary" : "text-gray-700"
            }`}
          >
            <p className="font-semibold text-start">
              Why should I mention the Language
            </p>
            {isOpen ? <IoIosArrowUp /> : <IoIosArrowDown />}
          </button>

          {isOpen && (
            <p className="text-sm mt-5 text-gray-600">
              It is important to mention the language of the course content so
              that the students who understand the language can enroll in the
              course.
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
