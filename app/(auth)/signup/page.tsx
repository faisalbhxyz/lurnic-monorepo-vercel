import React from "react";
import Image from "next/image";
import { Metadata } from "next";
import Link from "next/link";
import RegisterForm from "@/components/auth/RegisterForm";
import logo from "@/public/logo/perfact-logo.svg";

export const metadata: Metadata = {
  title: "Register | Perfact",
  description: "Register to Perfact",
};
const page = () => {
  return (
    <div className="bg-white p-10 rounded-2xl shadow-lg max-w-lg w-full">
      <div className="flex flex-col items-center">
        {/* <Image
          src={logo}
          alt={"logo"}
          width={135}
          height={100}
          className="w-[150px] h-auto"
        /> */}
        <p className="text-lg mt-2 text-gray-500">Create an account.</p>
      </div>
      <RegisterForm />
      <div className="flex items-center mt-4">
        <div className="flex-grow border-t border-gray-300"></div>
        <span className="mx-4 text-gray-500">OR</span>
        <div className="flex-grow border-t border-gray-300"></div>
      </div>
      {/* <div className="mt-6 text-center">
          <button className="w-full flex items-center justify-center space-x-2 border border-gray-300 py-3 rounded-lg hover:bg-gray-100 transition">
            <Image
              src="https://www.svgrepo.com/show/355037/google.svg"
              alt="Google"
              width={24}
              height={24}
              className="w-6 h-6"
            />
            <span className="text-sm font-medium text-gray-700">
              Continue with Google
            </span>
          </button>
        </div> */}
      <p className="text-center mt-3">
        Already have an account?{" "}
        <Link href={"/login"} className="text-[#00828a]">
          Login
        </Link>
      </p>
    </div>
  );
};

export default page;
