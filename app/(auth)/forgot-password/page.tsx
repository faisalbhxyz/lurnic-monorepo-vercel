import React from "react";
import Image from "next/image";
import { Metadata } from "next";
import Link from "next/link";
import ForgotPassword from "@/components/auth/ForgotPassword";

export const metadata: Metadata = {
  title: "Forgot Password | Perfact",
  description: "Forgot Password to Perfact",
};
const page = () => {
  return (
    <div className="bg-white p-10 rounded-2xl shadow-lg max-w-lg w-full">
      <Image
        src={"/logo/perfact-logo.svg"}
        alt={"logo"}
        width={135}
        height={100}
        className="w-[120px] h-auto mb-3"
      />
      <p className="text-2xl font-medium mb-5">Forgot Password</p>
      <p className="text-lg mt-2 text-gray-500 mb-5">
        Type in the email you&apos;ve used at the sign-up and we&apos;ll send
        you a password reset link.
      </p>

      <ForgotPassword />
      <p className="text-center mt-3">
        <Link href={"/login"} className="text-[#00828a] underline">
          Return to login
        </Link>
      </p>
    </div>
  );
};

export default page;
