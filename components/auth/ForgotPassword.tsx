"use client";

import { useRouter } from "next/navigation";
import React from "react";

const ForgotPassword = () => {
  const router = useRouter();

  const handleLogin = () => {
    router.push("/");
  };

  return (
    <form className="space-y-4">
      <div>
        <label className="block text-sm font-medium text-gray-700">
          Email Address
        </label>
        <input
          type="email"
          placeholder="Enter your work email"
          className="mt-1 w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-[#00828a] focus:border-transparent"
          required
        />
      </div>
      <button
        type="submit"
        onClick={handleLogin}
        className="w-full cursor-pointer bg-[#00828a] text-white py-2 rounded-lg hover:bg-[#00828a] transition mt-3"
      >
        Send me link
      </button>
    </form>
  );
};

export default ForgotPassword;
