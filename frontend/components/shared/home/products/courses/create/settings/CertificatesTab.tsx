"use client";

import InputField from "@/components/ui/InputField";
import Image from "next/image";
import React, { useRef } from "react";
import { Controller, useFormContext } from "react-hook-form";
import { FiDownload, FiPlus, FiUpload, FiX } from "react-icons/fi";
import { IoIosArrowDown, IoIosArrowUp } from "react-icons/io";
import { TCourseSchema } from "@/schema/course.schema";
import { CERTIFICATE_TEMPLATE_MINAR } from "@/components/shared/certificates/certificate-templates";
import CertificatePreview from "./CertificatePreview";
import { printCertificatePreview } from "@/lib/printCertificate";

export default function CertificatesTab() {
  const [isOpen, setIsOpen] = React.useState(true);
  const {
    register,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useFormContext<TCourseSchema>();

  React.useEffect(() => {
    setValue("certificate_settings.template_path", CERTIFICATE_TEMPLATE_MINAR, {
      shouldDirty: false,
    });
  }, [setValue]);

  const brandLogoFileInputRef = useRef<HTMLInputElement>(null);
  const watermarkFileInputRef = useRef<HTMLInputElement>(null);
  const ownerFileInputRef = useRef<HTMLInputElement>(null);
  const instructorFileInputRef = useRef<HTMLInputElement>(null);
  const certificatePreviewRef = useRef<HTMLDivElement>(null);

  const isEnabled = watch("certificate_settings.is_enabled");
  const subtitleOne = watch("certificate_settings.subtitle_one");
  const courseTitle = watch("title");
  const pricingModel = watch("pricing_model");
  const organizationName = watch("certificate_settings.organization_name");
  const brandLogo = watch("certificate_settings.brand_logo");
  const watermarkImage = watch("certificate_settings.watermark_image");
  const watermarkOpacity = watch("certificate_settings.watermark_opacity");
  const ownerSignature = watch("certificate_settings.owner_signature");
  const instructorSignature = watch("certificate_settings.instructor_signature");
  const signerName = watch("certificate_settings.signer_name");
  const signerRole = watch("certificate_settings.signer_role");
  const signerOrg = watch("certificate_settings.signer_org");
  const dualSignersEnabled = watch("certificate_settings.dual_signers_enabled");
  const signer2Name = watch("certificate_settings.signer2_name");
  const signer2Role = watch("certificate_settings.signer2_role");
  const signer2Org = watch("certificate_settings.signer2_org");

  const brandLogoPreview =
    brandLogo && typeof brandLogo === "object" && "isDBImg" in brandLogo
      ? (brandLogo.name as string)
      : brandLogo instanceof File
        ? URL.createObjectURL(brandLogo)
        : null;

  const watermarkPreview =
    watermarkImage &&
    typeof watermarkImage === "object" &&
    "isDBImg" in watermarkImage
      ? (watermarkImage.name as string)
      : watermarkImage instanceof File
        ? URL.createObjectURL(watermarkImage)
        : null;

  const ownerPreview =
    ownerSignature && typeof ownerSignature === "object" && "isDBImg" in ownerSignature
      ? (ownerSignature.name as string)
      : ownerSignature instanceof File
        ? URL.createObjectURL(ownerSignature)
        : null;

  const instructorPreview =
    instructorSignature &&
    typeof instructorSignature === "object" &&
    "isDBImg" in instructorSignature
      ? (instructorSignature.name as string)
      : instructorSignature instanceof File
        ? URL.createObjectURL(instructorSignature)
        : null;

  const handleDownloadPreview = () => {
    if (!certificatePreviewRef.current) return;
    printCertificatePreview(certificatePreviewRef.current);
  };

  return (
    <div className="flex items-start gap-10">
      <div className="w-full">
        <div className="mb-6 flex items-center justify-between rounded-md border p-4">
          <div>
            <p className="font-medium">Enable certificate</p>
            <p className="text-sm text-gray-500 mt-1">
              Auto-issue when a student reaches the completion threshold.
            </p>
          </div>
          <label className="inline-flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
              {...register("certificate_settings.is_enabled")}
            />
            <span className="text-sm font-medium">Enabled</span>
          </label>
        </div>

        <div className="mb-5">
          <label className="text-sm font-medium mb-1 block">
            Minimum completion (%)
          </label>
          <InputField
            type="number"
            min={1}
            max={100}
            className="w-full max-w-xs"
            disabled={!isEnabled}
            {...register("certificate_settings.completion_percent")}
          />
          {errors.certificate_settings?.completion_percent && (
            <p className="text-sm text-red-500 mt-1">
              {errors.certificate_settings.completion_percent.message}
            </p>
          )}
          <p className="text-xs text-gray-500 mt-1">
            Example: 80 means the student must complete 80% of selected items
            below.
          </p>
        </div>

        <div className="mb-6 rounded-md border p-4">
          <p className="font-medium mb-1">Count toward progress</p>
          <p className="text-sm text-gray-500 mb-3">
            Choose what counts for this course. Different courses can use
            different rules.
          </p>
          <div className="space-y-2">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                {...register("certificate_settings.count_lessons")}
              />
              Lessons (published)
            </label>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                {...register("certificate_settings.count_quizzes")}
              />
              Quizzes (submitted)
            </label>
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                {...register("certificate_settings.count_assignments")}
              />
              Assignments (submitted)
            </label>
          </div>
          {errors.certificate_settings?.count_lessons && (
            <p className="text-sm text-red-500 mt-2">
              {errors.certificate_settings.count_lessons.message}
            </p>
          )}
        </div>

        <div className="mb-2 flex max-w-3xl items-center justify-between gap-3">
          <p className="font-medium">Certificate preview</p>
          <button
            type="button"
            onClick={handleDownloadPreview}
            aria-label="Download certificate PDF"
            title="Download PDF"
            className="inline-flex shrink-0 items-center justify-center rounded-md border border-gray-200 bg-white p-1.5 text-gray-600 transition-colors hover:border-[#0b463f] hover:bg-[#0b463f] hover:text-white"
          >
            <FiDownload size={14} />
          </button>
        </div>
        <div ref={certificatePreviewRef} className="mb-4 max-w-3xl">
          <CertificatePreview
            templatePath={CERTIFICATE_TEMPLATE_MINAR}
            subtitleOne={subtitleOne}
            courseTitle={courseTitle}
            pricingModel={pricingModel}
            organizationName={organizationName}
            brandLogoSrc={brandLogoPreview}
            watermarkImageSrc={watermarkPreview}
            watermarkOpacity={watermarkOpacity}
            ownerSignatureSrc={ownerPreview}
            instructorSignatureSrc={instructorPreview}
            signerName={signerName}
            signerRole={signerRole}
            signerOrg={signerOrg}
            dualSigners={Boolean(dualSignersEnabled)}
            signer2Name={signer2Name}
            signer2Role={signer2Role}
            signer2Org={signer2Org}
          />
        </div>

        <input
          type="hidden"
          {...register("certificate_settings.template_path")}
          value={CERTIFICATE_TEMPLATE_MINAR}
        />

        <div className="mt-5">
          <label className="text-sm font-medium block mb-2">
            Certificate Logo
          </label>
          <div className="flex items-center gap-4">
            {brandLogoPreview && (
              <Image
                src={brandLogoPreview}
                alt="Certificate logo"
                width={100}
                height={100}
                className="h-20 w-20 rounded object-contain border"
              />
            )}
            <button
              type="button"
              disabled={!isEnabled}
              onClick={() => brandLogoFileInputRef.current?.click()}
              className="border px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline disabled:opacity-50"
            >
              <FiUpload /> Upload Logo
            </button>
            <Controller
              control={control}
              name="certificate_settings.brand_logo"
              render={({ field: { onChange } }) => (
                <input
                  ref={brandLogoFileInputRef}
                  type="file"
                  accept="image/*"
                  hidden
                  disabled={!isEnabled}
                  onChange={(event) => {
                    const file = event.target.files?.[0];
                    if (file) onChange(file);
                  }}
                />
              )}
            />
          </div>
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium block mb-2">
            Watermark Image
          </label>
          <p className="text-xs text-gray-500 mb-2">
            Optional faded background image shown behind certificate text.
          </p>
          <div className="flex items-center gap-4">
            {watermarkPreview && (
              <Image
                src={watermarkPreview}
                alt="Certificate watermark"
                width={100}
                height={100}
                className="h-20 w-20 rounded object-contain border opacity-60"
              />
            )}
            <button
              type="button"
              disabled={!isEnabled}
              onClick={() => watermarkFileInputRef.current?.click()}
              className="border px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline disabled:opacity-50"
            >
              <FiUpload /> Upload Watermark
            </button>
            <Controller
              control={control}
              name="certificate_settings.watermark_image"
              render={({ field: { onChange } }) => (
                <input
                  ref={watermarkFileInputRef}
                  type="file"
                  accept="image/*"
                  hidden
                  disabled={!isEnabled}
                  onChange={(event) => {
                    const file = event.target.files?.[0];
                    if (file) onChange(file);
                  }}
                />
              )}
            />
          </div>
          <div className="mt-4 max-w-md">
            <div className="mb-2 flex items-center justify-between gap-3">
              <label className="text-sm font-medium" htmlFor="watermark-opacity">
                Watermark opacity
              </label>
              <span className="text-sm text-gray-600">
                {Math.min(100, Math.max(0, Number(watermarkOpacity ?? 30)))}%
              </span>
            </div>
            <input
              id="watermark-opacity"
              type="range"
              min={0}
              max={100}
              step={1}
              disabled={!isEnabled || !watermarkPreview}
              className="w-full accent-[#0b463f] disabled:opacity-50"
              {...register("certificate_settings.watermark_opacity", {
                valueAsNumber: true,
              })}
            />
            <p className="text-xs text-gray-500 mt-1">
              Default is 30%. Upload a watermark image to adjust transparency.
            </p>
          </div>
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium mb-1">
            Organization Name (offered by)
          </label>
          <InputField
            className="w-full"
            disabled={!isEnabled}
            placeholder="e.g. 10 Minute School"
            {...register("certificate_settings.organization_name")}
          />
        </div>

        <div className="mt-5">
          <label className="text-sm font-medium mb-1">Completion Line</label>
          <InputField
            className="w-full"
            disabled={!isEnabled}
            placeholder="has successfully completed"
            {...register("certificate_settings.subtitle_one")}
          />
        </div>

        <div className="mt-5 rounded-md border border-gray-200 p-4">
          <p className="text-sm font-semibold text-gray-900 mb-4">Signer 1</p>

          <div className="space-y-5">
            <div>
              <label className="text-sm font-medium mb-1 block">Signer Name</label>
              <InputField
                className="w-full"
                disabled={!isEnabled}
                placeholder="e.g. Ayman Sadiq"
                {...register("certificate_settings.signer_name")}
              />
            </div>

            <div>
              <label className="text-sm font-medium mb-1 block">Signer Role</label>
              <InputField
                className="w-full"
                disabled={!isEnabled}
                placeholder="e.g. Course Instructor"
                {...register("certificate_settings.signer_role")}
              />
            </div>

            <div>
              <label className="text-sm font-medium mb-1 block">Signer Organization</label>
              <InputField
                className="w-full"
                disabled={!isEnabled}
                placeholder="e.g. CEO, 10 Minute School"
                {...register("certificate_settings.signer_org")}
              />
            </div>

            <div>
              <label className="text-sm font-medium block mb-2">
                Signature Image (150x250 px)
              </label>
              <div className="flex items-center gap-4">
                {instructorPreview && (
                  <Image
                    src={instructorPreview}
                    alt="Signer 1 signature"
                    width={100}
                    height={100}
                    className="w-20 h-20 object-cover rounded"
                  />
                )}
                <button
                  type="button"
                  disabled={!isEnabled}
                  onClick={() => instructorFileInputRef.current?.click()}
                  className="border px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline disabled:opacity-50"
                >
                  <FiUpload /> Upload
                </button>
                <Controller
                  control={control}
                  name="certificate_settings.instructor_signature"
                  render={({ field: { onChange } }) => (
                    <input
                      ref={instructorFileInputRef}
                      type="file"
                      accept="image/*"
                      hidden
                      disabled={!isEnabled}
                      onChange={(event) => {
                        const file = event.target.files?.[0];
                        if (file) onChange(file);
                      }}
                    />
                  )}
                />
              </div>
            </div>
          </div>

          {!dualSignersEnabled ? (
            <button
              type="button"
              disabled={!isEnabled}
              onClick={() =>
                setValue("certificate_settings.dual_signers_enabled", true, {
                  shouldDirty: true,
                })
              }
              className="mt-5 inline-flex items-center gap-2 rounded-md border border-[#0b463f] bg-[#0b463f]/5 px-4 py-2.5 text-sm font-medium text-[#0b463f] transition-colors hover:bg-[#0b463f] hover:text-white disabled:opacity-50"
            >
              <FiPlus size={16} />
              Add Another sign
            </button>
          ) : null}

          {dualSignersEnabled ? (
            <div className="mt-6 border-t border-gray-200 pt-5 space-y-5">
              <div className="flex items-center justify-between gap-3">
                <p className="text-sm font-semibold text-gray-900">Signer 2</p>
                <button
                  type="button"
                  disabled={!isEnabled}
                  onClick={() => {
                    setValue("certificate_settings.dual_signers_enabled", false, {
                      shouldDirty: true,
                    });
                    setValue("certificate_settings.signer2_name", "", {
                      shouldDirty: true,
                    });
                    setValue("certificate_settings.signer2_role", "", {
                      shouldDirty: true,
                    });
                    setValue("certificate_settings.signer2_org", "", {
                      shouldDirty: true,
                    });
                    setValue("certificate_settings.owner_signature", null, {
                      shouldDirty: true,
                    });
                  }}
                  className="inline-flex items-center gap-1.5 text-xs font-medium text-gray-500 hover:text-red-600 disabled:opacity-50"
                >
                  <FiX size={14} />
                  Remove
                </button>
              </div>

              <div>
                <label className="text-sm font-medium mb-1 block">Signer Name</label>
                <InputField
                  className="w-full"
                  disabled={!isEnabled}
                  placeholder="e.g. Ayman Sadiq"
                  {...register("certificate_settings.signer2_name")}
                />
              </div>

              <div>
                <label className="text-sm font-medium mb-1 block">Signer Role</label>
                <InputField
                  className="w-full"
                  disabled={!isEnabled}
                  placeholder="e.g. Course Instructor"
                  {...register("certificate_settings.signer2_role")}
                />
              </div>

              <div>
                <label className="text-sm font-medium mb-1 block">Signer Organization</label>
                <InputField
                  className="w-full"
                  disabled={!isEnabled}
                  placeholder="e.g. CEO, 10 Minute School"
                  {...register("certificate_settings.signer2_org")}
                />
              </div>

              <div>
                <label className="text-sm font-medium block mb-2">
                  Signature Image (150x250 px)
                </label>
                <div className="flex items-center gap-4">
                  {ownerPreview && (
                    <Image
                      src={ownerPreview}
                      alt="Signer 2 signature"
                      width={100}
                      height={100}
                      className="w-20 h-20 object-cover rounded"
                    />
                  )}
                  <button
                    type="button"
                    disabled={!isEnabled}
                    onClick={() => ownerFileInputRef.current?.click()}
                    className="border px-4 py-2 rounded-md flex items-center gap-2 text-sm font-medium text-primary hover:underline disabled:opacity-50"
                  >
                    <FiUpload /> Upload
                  </button>
                  <Controller
                    control={control}
                    name="certificate_settings.owner_signature"
                    render={({ field: { onChange } }) => (
                      <input
                        ref={ownerFileInputRef}
                        type="file"
                        accept="image/*"
                        hidden
                        disabled={!isEnabled}
                        onChange={(event) => {
                          const file = event.target.files?.[0];
                          if (file) onChange(file);
                        }}
                      />
                    )}
                  />
                </div>
              </div>
            </div>
          ) : null}

          <p className="text-xs text-gray-500 mt-4">
            Signature image appears above the line. Signer name, role, and organization
            appear below the line in regular text. With two signers, they appear on
            the left and right of the certificate.
          </p>
        </div>
      </div>

      <div className="w-80 min-w-80">
        <div className="border p-5 rounded-md">
          <button
            type="button"
            onClick={() => setIsOpen(!isOpen)}
            className={`flex items-center justify-between w-full ${
              isOpen ? "text-primary" : "text-gray-700"
            }`}
          >
            <p className="font-semibold text-start">How certificates work</p>
            {isOpen ? <IoIosArrowUp /> : <IoIosArrowDown />}
          </button>

          {isOpen && (
            <p className="text-sm mt-5 text-gray-600">
              Progress is calculated from the items you select above. When a
              student reaches your completion percentage, a certificate is
              issued automatically. Students can open the download link and save
              as PDF from the browser print dialog.
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
