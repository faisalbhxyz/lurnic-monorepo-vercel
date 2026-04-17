import ToggleSwitch from "@/components/ui/ToggleSwitch";
import Image from "next/image";
import React, { useEffect, useState } from "react";
import { IoMdRefresh } from "react-icons/io";
import { LiaSlidersHSolid } from "react-icons/lia";
import AddManualPayment from "./AddManualPayment";
import SelectList from "@/components/ui/SelectList";
import InputField from "@/components/ui/InputField";
import { MdOutlineFileCopy } from "react-icons/md";
import Button from "@/components/ui/Button";
import PaymentList from "./PaymentList";
import axiosInstance from "@/lib/axiosInstance";
import { toast } from "sonner";
import { useSession } from "next-auth/react";
import { getPaymentMethods } from "@/app/actions/actions";
import { useEditStore } from "@/hooks/useEditStore";
import UpdateManualPayment from "./UpdateManualPayment";

const environments = [
  {
    id: 1,
    name: "Option 1",
  },
  {
    id: 1,
    name: "Option 1",
  },
];

export default function PaymentMethods() {
  const [isOpen, setIsOpen] = useState(false);
  const [selectEnvironment, setSelectEnvironment] = useState(environments[0]);
  const [paymentMethods, setPaymentMethods] = useState<IPaymentMethods[]>([]);
  const { data: session } = useSession();
  const { refreshPaymentMethodsCount } = useEditStore();

  useEffect(() => {
    if (session) {
      getPaymentMethods(session).then((res) => {
        setPaymentMethods(res);
      });
    }
  }, [session, refreshPaymentMethodsCount]);

  return (
    <>
      <div className="flex items-center justify-between mb-3">
        <p className="text-xl font-medium">Payment Methods</p>
        {/* <button className="text-sm font-medium text-gray-500 flex items-center gap-1">
          <IoMdRefresh size={18} />
          Reset to Default
        </button> */}
      </div>
      <p className="text-gray-600 mt-5 mb-2">Payment methods</p>
      <PaymentList data={paymentMethods} />
      {/* <div className="bg-white p-3 border rounded-md">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Image src={"/images/paypal.svg"} alt={""} width={20} height={20} />
            Paypal
          </div>
          <div className="flex items-center gap-2">
            <ToggleSwitch />
            <button
              onClick={() => setIsOpen((prev) => !prev)}
              className="text-primary"
            >
              <LiaSlidersHSolid size={20} />
            </button>
          </div>
        </div>
        <div
          className={`grid transition-all duration-300 overflow-hidden ease-in-out ${
            isOpen
              ? "grid-rows-[1fr] opacity-100 mt-4"
              : "grid-rows-[0fr] opacity-0"
          }`}
        >
          <div className="text-[#424242] text-[0.9rem] overflow-hidden border rounded-md">
            <div className="p-4 space-y-4">
              <div className="flex items-center justify-between">
                <p>Environment</p>
                <SelectList
                  className="min-w-56"
                  options={environments}
                  value={selectEnvironment}
                  onChange={setSelectEnvironment}
                />
              </div>
              <div className="flex items-center justify-between">
                <p>Merchant email</p>
                <InputField className="w-full max-w-56" />
              </div>
              <div className="flex items-center justify-between">
                <p>Client id</p>
                <InputField className="w-full max-w-56" />
              </div>
              <div className="flex items-center justify-between">
                <p>Secret id</p>
                <InputField className="w-full max-w-56" />
              </div>
              <div className="flex items-center justify-between">
                <p>Webhook id</p>
                <InputField className="w-full max-w-56" />
              </div>
              <div className="flex items-center justify-between">
                <p>Webhook url</p>
                <div className="flex items-center gap-3 max-w-56">
                  <p className="truncate text-xs text-blue-600 font-medium">
                    https://amerrajjonowga.com/wp-json/tutor/v
                  </p>
                  <Button
                    variant="secondary"
                    className="flex items-center gap-1"
                  >
                    <MdOutlineFileCopy />
                    Copy
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div> */}
      <AddManualPayment />
      <UpdateManualPayment />
    </>
  );
}
