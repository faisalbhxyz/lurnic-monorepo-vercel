"use client";

import Image from "next/image";
import React from "react";
import { isMinarCertificateTemplate } from "@/components/shared/certificates/certificate-templates";
import MinarAcademyCertificate from "./MinarAcademyCertificate";

const TEMPLATE_BG: Record<string, string> = {
  "/images/Certificat-14.jpg": "#f7f5f0",
  "/images/Certificat-15.jpg": "#f7f5f0",
  "/images/Certificat-16.jpg": "#ffffff",
  "/images/Certificat-17.jpg": "#ffffff",
};

export interface CertificatePreviewProps {
  templatePath: string;
  title?: string | null;
  subtitleOne?: string | null;
  subtitleTwo?: string | null;
  courseTitle?: string | null;
  pricingModel?: "free" | "paid" | null;
  organizationName?: string | null;
  brandLogoSrc?: string | null;
  watermarkImageSrc?: string | null;
  watermarkOpacity?: number | null;
  ownerSignatureSrc?: string | null;
  instructorSignatureSrc?: string | null;
  signerName?: string | null;
  signerRole?: string | null;
  signerOrg?: string | null;
  dualSigners?: boolean;
  signer2Name?: string | null;
  signer2Role?: string | null;
  signer2Org?: string | null;
  className?: string;
}

function overlayBg(templatePath: string) {
  return TEMPLATE_BG[templatePath] ?? "#ffffff";
}

export default function CertificatePreview({
  templatePath,
  title,
  subtitleOne,
  subtitleTwo,
  courseTitle,
  pricingModel,
  organizationName,
  brandLogoSrc,
  watermarkImageSrc,
  watermarkOpacity,
  ownerSignatureSrc,
  instructorSignatureSrc,
  signerName,
  signerRole,
  signerOrg,
  dualSigners,
  signer2Name,
  signer2Role,
  signer2Org,
  className = "",
}: CertificatePreviewProps) {
  if (isMinarCertificateTemplate(templatePath)) {
    const signatureImageSrc = instructorSignatureSrc || (!dualSigners ? ownerSignatureSrc : null) || null;
    const signatureImageSrc2 = dualSigners ? ownerSignatureSrc || null : null;

    return (
      <MinarAcademyCertificate
        className={className}
        previewMode
        studentName="[STUDENT NAME]"
        courseName={courseTitle}
        pricingModel={pricingModel}
        organizationName={organizationName}
        completionLine={subtitleOne}
        brandLogoSrc={brandLogoSrc}
        watermarkImageSrc={watermarkImageSrc}
        watermarkOpacity={watermarkOpacity}
        signatureImageSrc={signatureImageSrc}
        signerName={signerName}
        signerRole={signerRole}
        signerOrg={signerOrg}
        dualSigners={dualSigners}
        signatureImageSrc2={signatureImageSrc2}
        signer2Name={signer2Name}
        signer2Role={signer2Role}
        signer2Org={signer2Org}
      />
    );
  }

  const bg = overlayBg(templatePath);
  const displayTitle = title?.trim() || "Certificate of Completion";
  const displayCourse = courseTitle?.trim() || "Course Name";

  return (
    <div
      className={`@container relative w-full overflow-hidden rounded-md border ${className}`}
      style={{ aspectRatio: "1650 / 1275", backgroundColor: bg }}
    >
      <Image
        src={templatePath}
        alt="Certificate preview"
        fill
        className="object-cover"
        sizes="(max-width: 768px) 100vw, 600px"
      />

      <div
        className="absolute left-1/2 -translate-x-1/2 text-center"
        style={{ top: "26.5%", width: "72%" }}
      >
        <span
          className="inline-block px-2 py-0.5 text-[clamp(0.5rem,2.6cqw,1rem)] font-semibold uppercase leading-tight tracking-[0.2em] text-[#b8956a]"
          style={{ backgroundColor: bg }}
        >
          {displayTitle}
        </span>
      </div>

      {subtitleOne?.trim() ? (
        <div
          className="absolute left-1/2 -translate-x-1/2 text-center"
          style={{ top: "39%", width: "88%" }}
        >
          <span
            className="inline-block px-1.5 py-0.5 text-[clamp(0.4rem,1.7cqw,0.7rem)] font-medium uppercase leading-snug tracking-wide text-gray-800"
            style={{ backgroundColor: bg }}
          >
            {subtitleOne.trim()}
          </span>
        </div>
      ) : null}

      <div
        className="absolute left-1/2 -translate-x-1/2 text-center"
        style={{ top: "46.5%", width: "62%" }}
      >
        <span
          className="block text-[clamp(1.1rem,5.2cqw,2.5rem)] leading-none text-gray-900"
          style={{ fontFamily: "'Brush Script MT', 'Segoe Script', cursive" }}
        >
          Your Name
        </span>
        <div className="mx-auto mt-1 h-px w-[70%] bg-gray-400/80" />
      </div>

      {subtitleTwo?.trim() ? (
        <div
          className="absolute left-1/2 -translate-x-1/2 text-center"
          style={{ top: "57.5%", width: "88%" }}
        >
          <span
            className="inline-block px-1.5 py-0.5 text-[clamp(0.4rem,1.7cqw,0.7rem)] font-medium uppercase leading-snug tracking-wide text-gray-800"
            style={{ backgroundColor: bg }}
          >
            {subtitleTwo.trim()}
          </span>
        </div>
      ) : null}

      {courseTitle?.trim() ? (
        <div
          className="absolute left-1/2 -translate-x-1/2 text-center"
          style={{ top: "63.5%", width: "78%" }}
        >
          <span
            className="inline-block px-2 py-0.5 text-[clamp(0.55rem,2.4cqw,0.95rem)] font-bold uppercase leading-tight tracking-[0.18em] text-[#b8956a]"
            style={{ backgroundColor: bg }}
          >
            {displayCourse}
          </span>
        </div>
      ) : null}

      {ownerSignatureSrc ? (
        <div
          className="absolute flex flex-col items-center"
          style={{ left: "17%", top: "73.5%", width: "24%" }}
        >
          <div className="relative h-[clamp(1.75rem,7.5cqw,3.5rem)] w-full">
            <Image
              src={ownerSignatureSrc}
              alt="Owner signature"
              fill
              className="object-contain object-bottom"
              unoptimized={ownerSignatureSrc.startsWith("blob:")}
            />
          </div>
        </div>
      ) : null}

      {instructorSignatureSrc ? (
        <div
          className="absolute flex flex-col items-center"
          style={{ right: "17%", top: "73.5%", width: "24%" }}
        >
          <div className="relative h-[clamp(1.75rem,7.5cqw,3.5rem)] w-full">
            <Image
              src={instructorSignatureSrc}
              alt="Instructor signature"
              fill
              className="object-contain object-bottom"
              unoptimized={instructorSignatureSrc.startsWith("blob:")}
            />
          </div>
        </div>
      ) : null}
    </div>
  );
}
