"use client";

import React, { useState } from "react";
import General from "./General";
import Course from "./Course";
import Monetization from "./Monetization";
import Checkout from "./Checkout";
import Advanced from "./Advanced";
import Design from "./Design";
import PaymentMethods from "./payment_methods/PaymentMethods";
import Taxes from "./Taxes";

export default function SettingTabs({
  generalsettings,
}: {
  generalsettings: GeneralSettings;
}) {
  const [activeTab, setActiveTab] = useState("General");

  const tabs = [
    "General",
    "Course",
    "Monetization",
    "Payment Methods",
    "Taxes",
    "Checkout",
    "Design",
    "Advanced",
  ];

  const renderContent = () => {
    switch (activeTab) {
      case "General":
        return <General generalSettings={generalsettings} />;
      case "Course":
        return <Course />;
      case "Monetization":
        return <Monetization />;
      case "Payment Methods":
        return <PaymentMethods />;
      case "Taxes":
        return <Taxes />;
      case "Checkout":
        return <Checkout />;
      case "Design":
        return <Design />;
      case "Advanced":
        return <Advanced />;
      default:
        return null;
    }
  };

  return (
    <>
      <div className="flex rounded-lg max-w-5xl mx-auto p-10">
        <aside className="w-56 min-w-56 sticky top-20 self-start bg-gray-50 border border-gray-200 rounded-lg p-3">
          {tabs.map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`block w-full text-left rounded-md border-l-4 rounded-l-none px-4 py-2 text-sm font-medium transition-colors ${
                activeTab === tab
                  ? "bg-white text-primary border-primary"
                  : "text-gray-700 hover:bg-gray-100 border-transparent"
              }`}
            >
              {tab}
            </button>
          ))}
        </aside>
        <div className="w-full px-5">{renderContent()}</div>
      </div>
    </>
  );
}
