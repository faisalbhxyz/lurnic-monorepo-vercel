"use client";

import React, { useId, useLayoutEffect, useRef, useState } from "react";
import "@/styles/minar-certificate.css";

const DESIGN_WIDTH_MM = 297;

export interface MinarAcademyCertificateProps {
  studentName?: string | null;
  courseName?: string | null;
  certificateId?: string | null;
  issuedDate?: string | null;
  pricingModel?: "free" | "paid" | null;
  organizationName?: string | null;
  completionLine?: string | null;
  brandLogoSrc?: string | null;
  watermarkImageSrc?: string | null;
  watermarkOpacity?: number | null;
  signatureImageSrc?: string | null;
  signerName?: string | null;
  signerRole?: string | null;
  signerOrg?: string | null;
  dualSigners?: boolean;
  signatureImageSrc2?: string | null;
  signer2Name?: string | null;
  signer2Role?: string | null;
  signer2Org?: string | null;
  showDownloadButton?: boolean;
  onDownload?: () => void;
  className?: string;
  previewMode?: boolean;
}

function pricingWord(pricingModel?: "free" | "paid" | null) {
  return pricingModel === "paid" ? "PAID" : "FREE";
}

function mmToPx(mm: number) {
  return mm * (96 / 25.4);
}

function previewField(
  value: string | null | undefined,
  placeholder: string,
  previewMode: boolean,
  fallback = ""
) {
  const trimmed = value?.trim();
  if (trimmed) return trimmed;
  return previewMode ? placeholder : fallback;
}

export default function MinarAcademyCertificate({
  studentName,
  courseName,
  certificateId,
  issuedDate,
  pricingModel = "free",
  organizationName,
  completionLine,
  brandLogoSrc,
  watermarkImageSrc,
  watermarkOpacity = 30,
  signatureImageSrc,
  signerName,
  signerRole,
  signerOrg,
  dualSigners = false,
  signatureImageSrc2,
  signer2Name,
  signer2Role,
  signer2Org,
  showDownloadButton = false,
  onDownload,
  className = "",
  previewMode = false,
}: MinarAcademyCertificateProps) {
  const uid = useId().replace(/:/g, "");
  const scalerRef = useRef<HTMLDivElement>(null);
  const [scale, setScale] = useState(1);

  const displayStudent = previewField(studentName, "[STUDENT NAME]", previewMode, "[STUDENT NAME]");
  const displayCourse = previewField(courseName, "[COURSE NAME]", previewMode, "[COURSE NAME]");
  const displayCertId = previewField(certificateId, "[CERT-ID]", previewMode);
  const displayIssued = previewField(issuedDate, "[ISSUED-DATE]", previewMode);
  const displayOrg = previewField(organizationName, "[Organization]", previewMode, "your organization");
  const displayCompletion = previewField(
    completionLine,
    "[Completion Line]",
    previewMode,
    "has successfully completed"
  );
  const word = pricingWord(pricingModel);

  const showSignerName = previewField(signerName, "[Signer-Name]", previewMode);
  const showSignerRole = previewField(signerRole, "[Signer-Role]", previewMode);
  const showSignerOrg = previewField(signerOrg, "[Signer-Organization]", previewMode);
  const showSigner2Name = previewField(signer2Name, "[Signer-Name]", previewMode);
  const showSigner2Role = previewField(signer2Role, "[Signer-Role]", previewMode);
  const showSigner2Org = previewField(signer2Org, "[Signer-Organization]", previewMode);

  const renderSignerBlock = (
    imageSrc: string | null | undefined,
    name: string,
    role: string,
    org: string,
    key: string
  ) => (
    <div className="minar-cert-sig" key={key}>
      {imageSrc ? (
        <img
          className="minar-cert-sig-image"
          src={imageSrc}
          alt="Signature"
          crossOrigin="anonymous"
        />
      ) : null}
      <div className="minar-cert-sig-line" />
      {name ? <div className="minar-cert-sig-name">{name}</div> : null}
      {role ? <div className="minar-cert-sig-role">{role}</div> : null}
      {org ? <div className="minar-cert-sig-org">{org}</div> : null}
    </div>
  );

  const resolvedWatermarkOpacity = Math.min(
    100,
    Math.max(0, watermarkOpacity ?? 30)
  );

  useLayoutEffect(() => {
    const scaler = scalerRef.current;
    if (!scaler) return;

    const measure = () => {
      const width = scaler.clientWidth;
      if (width <= 0) return;
      setScale(width / mmToPx(DESIGN_WIDTH_MM));
    };

    measure();
    const observer = new ResizeObserver(measure);
    observer.observe(scaler);
    return () => observer.disconnect();
  }, []);

  return (
    <div className={`minar-cert-root ${className}`.trim()}>
      {showDownloadButton ? (
        <button
          type="button"
          className="minar-cert-download-btn"
          onClick={onDownload ?? (() => window.print())}
        >
          Download PDF
        </button>
      ) : null}

      <div ref={scalerRef} className="minar-cert-scaler">
        <div className="minar-cert-page" style={{ transform: `scale(${scale})` }}>
          <svg
            className="minar-cert-header-graphic"
            viewBox="0 0 1600 520"
            preserveAspectRatio="none"
            xmlns="http://www.w3.org/2000/svg"
            aria-hidden
          >
            <defs>
              <linearGradient id={`${uid}-tealLeft`} x1="0" y1="0" x2="800" y2="463" gradientUnits="userSpaceOnUse">
                <stop offset="0%" stopColor="#072e2a" />
                <stop offset="55%" stopColor="#0d5850" />
                <stop offset="100%" stopColor="#1c8577" />
              </linearGradient>
              <linearGradient id={`${uid}-tealRight`} x1="1600" y1="0" x2="800" y2="463" gradientUnits="userSpaceOnUse">
                <stop offset="0%" stopColor="#072e2a" />
                <stop offset="55%" stopColor="#0d5850" />
                <stop offset="100%" stopColor="#1c8577" />
              </linearGradient>
              <linearGradient id={`${uid}-goldLeft`} x1="0" y1="463" x2="800" y2="0" gradientUnits="userSpaceOnUse">
                <stop offset="0%" stopColor="#A87433" />
                <stop offset="45%" stopColor="#E6C88F" />
                <stop offset="60%" stopColor="#F3DEB4" />
                <stop offset="100%" stopColor="#C28F46" />
              </linearGradient>
              <linearGradient id={`${uid}-goldRight`} x1="1600" y1="463" x2="800" y2="0" gradientUnits="userSpaceOnUse">
                <stop offset="0%" stopColor="#A87433" />
                <stop offset="45%" stopColor="#E6C88F" />
                <stop offset="60%" stopColor="#F3DEB4" />
                <stop offset="100%" stopColor="#C28F46" />
              </linearGradient>
              <radialGradient id={`${uid}-medalGrad`} cx="50%" cy="50%" r="60%">
                <stop offset="0%" stopColor="#fff2c2" />
                <stop offset="60%" stopColor="#e0b04a" />
                <stop offset="100%" stopColor="#a97a1f" />
              </radialGradient>
              <linearGradient id={`${uid}-tailGrad`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#f0c95c" />
                <stop offset="100%" stopColor="#b8860b" />
              </linearGradient>
              <clipPath id={`${uid}-leftTriClip`}>
                <polygon points="0,0 800,0 0,463" />
              </clipPath>
              <clipPath id={`${uid}-rightTriClip`}>
                <polygon points="1600,0 800,0 1600,463" />
              </clipPath>
            </defs>
            <polygon points="0,0 800,0 0,463" fill={`url(#${uid}-tealLeft)`} />
            <polygon points="1600,0 800,0 1600,463" fill={`url(#${uid}-tealRight)`} />
            <polygon
              points="0,463 37.56,527.9 837.56,64.9 800,0"
              fill={`url(#${uid}-goldLeft)`}
              stroke="#C28F46"
              strokeWidth="60"
              strokeLinejoin="round"
              clipPath={`url(#${uid}-leftTriClip)`}
            />
            <polygon
              points="1600,463 1562.44,527.9 762.44,64.9 800,0"
              fill={`url(#${uid}-goldRight)`}
              stroke="#C28F46"
              strokeWidth="60"
              strokeLinejoin="round"
              clipPath={`url(#${uid}-rightTriClip)`}
            />
            <rect x="754" y="148" width="92" height="10" fill="#e0b04a" />
            <polygon points="754,152 754,244 800,224 846,244 846,152" fill={`url(#${uid}-tailGrad)`} />
            <circle cx="800" cy="70" r="88" fill={`url(#${uid}-medalGrad)`} stroke="#8a6316" strokeWidth="4" />
            <circle cx="800" cy="70" r="67" fill="none" stroke="#8a6316" strokeWidth="2" opacity="0.55" />
          </svg>

          <svg
            className="minar-cert-footer-bar"
            viewBox="0 0 2970 48"
            preserveAspectRatio="none"
            xmlns="http://www.w3.org/2000/svg"
            aria-hidden
          >
            <defs>
              <linearGradient id={`${uid}-footerGrad`} x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" stopColor="#0b463f" />
                <stop offset="50%" stopColor="#14685e" />
                <stop offset="100%" stopColor="#0b463f" />
              </linearGradient>
            </defs>
            <rect x="0" y="0" width="2970" height="48" fill={`url(#${uid}-footerGrad)`} />
          </svg>

          {watermarkImageSrc ? (
            <img
              className="minar-cert-watermark-bg"
              src={watermarkImageSrc}
              alt=""
              crossOrigin="anonymous"
              style={{ opacity: resolvedWatermarkOpacity / 100 }}
            />
          ) : null}

          <div className="minar-cert-content">
            {brandLogoSrc ? (
              <img className="minar-cert-brand-logo" src={brandLogoSrc} alt="Brand Logo" crossOrigin="anonymous" />
            ) : (
              <div className="minar-cert-brand-placeholder">Your Logo</div>
            )}

            <div className="minar-cert-recipient-name">{displayStudent}</div>
            <div className="minar-cert-completion-line">{displayCompletion}</div>
            <div className="minar-cert-course-name">{displayCourse}</div>
            <div className="minar-cert-course-subline">
              a <span className="minar-cert-pricing-word">{word}</span> online course offered by {displayOrg}.
            </div>

            <div className="minar-cert-meta-row">
              <div>
                CERTIFICATE ID: <strong>{displayCertId}</strong>
              </div>
              <div>
                Issued <span className="minar-cert-issued-at">{displayIssued}</span>
              </div>
            </div>
          </div>

          <div
            className={`minar-cert-signatures ${dualSigners ? "minar-cert-signatures--dual" : ""}`}
          >
            {dualSigners ? (
              <>
                {renderSignerBlock(
                  signatureImageSrc,
                  showSignerName,
                  showSignerRole,
                  showSignerOrg,
                  "signer-1"
                )}
                {renderSignerBlock(
                  signatureImageSrc2,
                  showSigner2Name,
                  showSigner2Role,
                  showSigner2Org,
                  "signer-2"
                )}
              </>
            ) : (
              renderSignerBlock(
                signatureImageSrc,
                showSignerName,
                showSignerRole,
                showSignerOrg,
                "signer-single"
              )
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
