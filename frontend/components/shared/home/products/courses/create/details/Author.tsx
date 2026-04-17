import { useSession } from "next-auth/react";
import Image from "next/image";
import React from "react";
import { IoIosArrowDown } from "react-icons/io";

export default function Author() {
  const { data: session } = useSession();
  return (
    <div className="border p-3 rounded-md flex items-center justify-between bg-white">
      <div className="flex items-center gap-2">
        <Image
          src={"/images/profile.png"}
          alt={"image"}
          width={30}
          height={30}
        />
        <div className="text-sm">
          <p>{session?.user?.name}</p>
          <p className="text-gray-600 text-xs">{session?.user?.email}</p>
        </div>
      </div>
      {/* <IoIosArrowDown /> */}
    </div>
  );
}
