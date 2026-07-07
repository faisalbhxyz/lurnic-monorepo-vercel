"use client";

import React, { useRef } from "react";
import { printCertificatePreview } from "@/lib/printCertificate";
import { isMinarCertificateTemplate } from "./certificate-templates";
import MinarAcademyCertificate from "../home/products/courses/create/settings/MinarAcademyCertificate";
import CertificatePreview from "../home/products/courses/create/settings/CertificatePreview";

export interface StudentCertificateData {
  id: number;
  course_title: string;
  certificate_number: string;
  student_name: string;
  template_path: string;
  title?: string | null;
  subtitle_one?: string | null;
  subtitle_two?: string | null;
  brand_logo?: string | null;
  watermark_image?: string | null;
  watermark_opacity?: number | null;
  organization_name?: string | null;
  signer_name?: string | null;
  signer_role?: string | null;
  signer_org?: string | null;
  dual_signers_enabled?: boolean;
  signer2_name?: string | null;
  signer2_role?: string | null;
  signer2_org?: string | null;
  pricing_model?: "free" | "paid" | null;
  owner_signature?: string | null;
  instructor_signature?: string | null;
  issued_at: string;
  download_url?: string;
}

import { formatCertificateIssuedAt } from "@/lib/certificate-format";

export default function StudentCertificateView({
  certificate,
  showDownloadButton = true,
}: {
  certificate: StudentCertificateData;
  showDownloadButton?: boolean;
}) {
  const certificateRef = useRef<HTMLDivElement>(null);
  const issuedDate = formatCertificateIssuedAt(certificate.issued_at);
  const dualSigners = Boolean(certificate.dual_signers_enabled);
  const handleDownload = () => {
    if (!certificateRef.current) return;
    printCertificatePreview(certificateRef.current);
  };
  const signatureImageSrc =
    certificate.instructor_signature ||
    (!dualSigners ? certificate.owner_signature : null) ||
    null;
  const signatureImageSrc2 = dualSigners ? certificate.owner_signature || null : null;

  if (isMinarCertificateTemplate(certificate.template_path)) {
    return (
      <div ref={certificateRef}>
      <MinarAcademyCertificate
        studentName={certificate.student_name}
        courseName={certificate.course_title}
        certificateId={certificate.certificate_number}
        issuedDate={issuedDate}
        pricingModel={certificate.pricing_model}
        organizationName={certificate.organization_name}
        completionLine={certificate.subtitle_one}
        brandLogoSrc={certificate.brand_logo}
        watermarkImageSrc={certificate.watermark_image}
        watermarkOpacity={certificate.watermark_opacity}
        signatureImageSrc={signatureImageSrc}
        signerName={certificate.signer_name}
        signerRole={certificate.signer_role}
        signerOrg={certificate.signer_org}
        dualSigners={dualSigners}
        signatureImageSrc2={signatureImageSrc2}
        signer2Name={certificate.signer2_name}
        signer2Role={certificate.signer2_role}
        signer2Org={certificate.signer2_org}
        showDownloadButton={showDownloadButton}
        onDownload={handleDownload}
      />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <CertificatePreview
        templatePath={certificate.template_path}
        title={certificate.title}
        subtitleOne={certificate.subtitle_one}
        subtitleTwo={certificate.subtitle_two}
        courseTitle={certificate.course_title}
        ownerSignatureSrc={certificate.owner_signature}
        instructorSignatureSrc={certificate.instructor_signature}
      />
      {showDownloadButton ? (
        <p className="text-sm text-gray-600 text-center">
          Use your browser print dialog to save this certificate as PDF.
        </p>
      ) : null}
    </div>
  );
}
