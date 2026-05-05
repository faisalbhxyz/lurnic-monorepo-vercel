import { cn } from "@/lib/cn";
import { signOut, useSession } from "next-auth/react";
import Image from "next/image";
import Link from "next/link";
import React, { useEffect, useRef, useState } from "react";
import { FiUser } from "react-icons/fi";
import { IoIosArrowDown } from "react-icons/io";
import { LuLogOut, LuWallet } from "react-icons/lu";

export default function Profile() {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const handleClickOutside = (event: MouseEvent) => {
    if (
      dropdownRef.current &&
      !dropdownRef.current.contains(event.target as Node)
    ) {
      setIsOpen(false);
    }
  };

  useEffect(() => {
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const User = () => {
    const { data: session } = useSession();

    return (
      <>
        <div className="w-10 min-w-10 h-10 rounded-md">
          <Image
            src={"/images/profile.png"}
            alt={"image"}
            width={200}
            height={200}
          />
        </div>
        <div className="text-sm truncate">
          <p className="font-semibold">{session?.user?.name}</p>
          <p className="truncate text-gray-600">{session?.user?.email}</p>
        </div>
      </>
    );
  };

  return (
    <div ref={dropdownRef} className="relative">
      <div
        onClick={() => setIsOpen((prev) => !prev)}
        className="flex items-center gap-3 hover:bg-slate-200 p-1 rounded-md cursor-pointer"
      >
        <User />
        <div>
          <IoIosArrowDown />
        </div>
      </div>
      <div
        className={cn(
          "absolute bottom-full mb-3 left-0 -right-3 bg-white border rounded-lg duration-200 z-50",
          isOpen ? "visible opacity-100" : "invisible opacity-0"
        )}
      >
        <div className="p-3 border-b border-gray-200 flex items-center gap-3">
          <User />
        </div>
        <div className="p-3">
          <Link
            href={""}
            className="hover:bg-gray-100 flex items-center gap-2 px-2 py-1 w-full rounded-md"
          >
            <LuWallet /> Billing
          </Link>
          <Link
            href={""}
            className="hover:bg-gray-100 flex items-center gap-2 px-2 py-1 w-full rounded-md"
          >
            <FiUser /> My Account
          </Link>
          <div className="border-t border-gray-200 my-3" />
          <button
            className="hover:bg-gray-100 flex items-center gap-2 px-2 py-1 w-full rounded-md"
            onClick={async () => {
              await signOut({ callbackUrl: "/login" });
            }}
          >
            <LuLogOut />
            Logout
          </button>
        </div>
      </div>
    </div>
  );
}
