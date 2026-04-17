import Image from "next/image";
import React from "react";

type DigitalDownloadsList = {
  data: {
    id: number;
    title: string;
  };
};

export default function DigitalDownloadsList({ data }: DigitalDownloadsList) {
  return (
    <div className="border rounded-lg overflow-hidden bg-white">
      <Image
        src={"/images/no-image.png"}
        alt={"image"}
        width={500}
        height={500}
        className="w-full h-48 object-cover"
      />
      <div className="p-4 text-sm">
        <p className="font-semibold">{data.title}</p>
        <p className="text-gray-500">test summery</p>
      </div>
      <div className="p-4 flex items-center justify-between border-t border-gray-300">
        <p className="text-sm text-gray-500">2 hours ago</p>
        <span className="uppercase font-medium">Free</span>
      </div>
    </div>
  );
}
