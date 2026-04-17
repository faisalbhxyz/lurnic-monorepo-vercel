"use client";

import React, { useState } from "react";
import GeneralTab from "@/components/shared/home/products/courses/create/settings/GeneralTab";
import CertificatesTab from "@/components/shared/home/products/courses/create/settings/CertificatesTab";
import SeoTab from "@/components/shared/home/products/courses/create/settings/SeoTab";
import FaqTab from "@/components/shared/home/products/courses/create/settings/FaqTab";

const TABS = ["General"]; //"Certificates", "SEO", "FAQ"

export default function SettingTabs({
  categories,
  subcategories,
  instructors,
}: {
  categories: ICategory[] | null;
  subcategories: ISubCategory[] | null;
  instructors: IInstructor[] | null;
}) {
  const [activeTab, setActiveTab] = useState(TABS[0]);

  const renderContent = () => {
    switch (activeTab) {
      case "General":
        return (
          <GeneralTab
            categories={categories}
            subcategories={subcategories}
            instructors={instructors}
          />
        );
      case "Certificates":
        return <CertificatesTab />;
      case "SEO":
        return <SeoTab />;
      case "FAQ":
        return <FaqTab />;
      default:
        return null;
    }
  };

  return (
    <>
      <div className="mb-4 flex space-x-1 border-b border-gray-200 text-sm font-medium">
        {TABS.map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`px-4 py-2 border-b transition-colors ${
              activeTab === tab
                ? "border-primary text-primary"
                : "border-transparent text-gray-500 hover:text-primary"
            }`}
          >
            {tab}
          </button>
        ))}
      </div>
      <div className="mt-5">{renderContent()}</div>
    </>
  );
}
