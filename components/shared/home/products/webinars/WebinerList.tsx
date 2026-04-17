"use client";

import Image from "next/image";
import React, { useState } from "react";
import SessionsTab from "./SessionsTab";
import SpeakersTab from "./SpeakersTab";

type DigitalDownloadsList = {
  data: {
    id: number;
    title: string;
  }[];
};

const TABS: { key: "sessions" | "speakers"; label: string }[] = [
  { key: "sessions", label: "Sessions" },
  { key: "speakers", label: "Speakers" },
];

export default function WebinerList({ data }: DigitalDownloadsList) {
  const [activeTab, setActiveTab] = useState<"sessions" | "speakers">(
    "sessions"
  );

  const renderContent = () => {
    switch (activeTab) {
      case "sessions":
        return <SessionsTab data={data} />;
      case "speakers":
        return <SpeakersTab />;
      default:
        return null;
    }
  };

  return (
    <div>
      {/* Tab buttons */}
      <div className="flex space-x-1 mb-4 text-sm">
        {TABS.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            className={`px-4 py-2 border-b transition-colors ${
              activeTab === tab.key
                ? "border-primary text-primary"
                : "border-transparent text-gray-500 hover:text-primary"
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Render tab content */}
      <div>{renderContent()}</div>
    </div>
  );
}
