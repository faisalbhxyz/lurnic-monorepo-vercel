"use client";

import React, { useEffect, useState } from "react";
import { IoMdRefresh } from "react-icons/io";
import SelectPage from "./SelectPage";
import InputField from "@/components/ui/InputField";
import ToggleSwitch from "@/components/ui/ToggleSwitch";
import Button from "@/components/ui/Button";
import { Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react";
import { MdKeyboardArrowDown } from "react-icons/md";
import { CURRENCIES, CurrencyCode, DEFAULT_CURRENCY } from "@/lib/currency";
import axiosInstance from "@/lib/axiosInstance";
import { useSession } from "next-auth/react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";
import { useCurrency } from "@/context/CurrencyContext";

export default function Monetization({
  generalSettings,
}: {
  generalSettings: GeneralSettings;
}) {
  const { data: session } = useSession();
  const router = useRouter();
  const { currency: activeCurrency } = useCurrency();
  const [selectedCurrency, setSelectedCurrency] = useState<CurrencyCode>(
    (generalSettings?.currency as CurrencyCode) || DEFAULT_CURRENCY
  );
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    setSelectedCurrency(
      (generalSettings?.currency as CurrencyCode) || DEFAULT_CURRENCY
    );
  }, [generalSettings?.currency]);

  const selectedOption =
    CURRENCIES.find((item) => item.code === selectedCurrency) ??
    CURRENCIES[0];

  const handleSaveCurrency = async () => {
    const accessToken = session?.accessToken;
    if (!accessToken) {
      toast.error("Session expired. Please sign in again.");
      return;
    }

    setIsSaving(true);
    try {
      const res = await axiosInstance.put(
        "/private/general-settings/currency",
        { currency: selectedCurrency },
        {
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        }
      );
      toast.success(res.data.message || "Currency updated successfully");
      router.refresh();
    } catch (error: unknown) {
      const message =
        error &&
        typeof error === "object" &&
        "response" in error &&
        error.response &&
        typeof error.response === "object" &&
        "data" in error.response &&
        error.response.data &&
        typeof error.response.data === "object" &&
        "error" in error.response.data
          ? String(error.response.data.error)
          : "Failed to update currency.";
      toast.error(message);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <>
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">Course</p>
        <button className="text-sm font-medium text-gray-500 flex items-center gap-1">
          <IoMdRefresh size={18} />
          Reset to Default
        </button>
      </div>
      <p className="text-gray-600 mt-5 mb-1">Options</p>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Select eCommerce Engine</p>
            <p className="text-gray-700 mt-1">
              Select a monetization option to generate revenue by selling
              courses.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Cart Page</p>
            <p className="text-gray-700 mt-1">
              Select the page you wish to set as the cart page.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Checkout Page</p>
            <p className="text-gray-700 mt-1">
              Select the page to be used as the checkout page.
            </p>
          </div>
          <SelectPage />
        </div>
      </div>
      <div className="flex items-center justify-between mb-1">
        <p className="text-gray-600">Currency</p>
        <Button
          type="button"
          onClick={handleSaveCurrency}
          disabled={isSaving || selectedCurrency === activeCurrency}
        >
          Save Currency
        </Button>
      </div>
      <div className="border rounded-md bg-white px-4 mb-5">
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Currency</p>
            <p className="text-gray-700 mt-1">
              Choose the currency for transactions.
            </p>
          </div>
          <Menu>
            <MenuButton className="inline-flex items-center justify-between text-gray-600 border gap-2 rounded-md py-1.5 px-3 text-sm min-w-52 font-medium focus:outline-none data-[focus]:outline-1 data-[focus]:outline-white">
              <span>
                {selectedOption.symbol} {selectedOption.label}
              </span>
              <MdKeyboardArrowDown className="size-5" />
            </MenuButton>
            <MenuItems
              transition
              anchor="bottom end"
              className="w-64 origin-top-right rounded-xl border bg-white p-1 text-sm transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0 z-20"
            >
              {CURRENCIES.map((item) => (
                <MenuItem key={item.code}>
                  <button
                    type="button"
                    onClick={() => setSelectedCurrency(item.code)}
                    className="group flex w-full items-center gap-2 rounded-lg py-1.5 px-3 data-[focus]:bg-gray-100"
                  >
                    <span className="w-5 text-center">{item.symbol}</span>
                    <span>{item.label}</span>
                  </button>
                </MenuItem>
              ))}
            </MenuItems>
          </Menu>
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Currency Position</p>
            <p className="text-gray-700 mt-1">
              Set the position of the currency symbol.
            </p>
          </div>
          <SelectPage />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Thousand Separator</p>
            <p className="text-gray-700 mt-1">
              Specify the thousand separator.
            </p>
          </div>
          <InputField className="w-20 text-center" />
        </div>
        <div className="py-4 border-b border-gray-300 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Decimal Separator</p>
            <p className="text-gray-700 mt-1">Specify the decimal separator.</p>
          </div>
          <InputField className="w-20 text-center" />
        </div>
        <div className="py-4 flex items-center justify-between text-sm gap-3">
          <div>
            <p className="font-medium text-gray-700">Number of Decimals</p>
            <p className="text-gray-700 mt-1">
              Set the number of decimal places.
            </p>
          </div>
          <InputField className="w-20 text-center" />
        </div>
      </div>
      <div className="p-4 bg-white border rounded-md flex items-center justify-between text-sm gap-3">
        <div>
          <p className="font-medium text-gray-700">Enable Revenue Sharing</p>
          <p className="text-gray-700 mt-1">
            Allow revenue generated from selling courses to be shared with
            course creators.
          </p>
        </div>
        <ToggleSwitch />
      </div>
    </>
  );
}
