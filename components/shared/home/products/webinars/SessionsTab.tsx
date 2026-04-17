import Image from "next/image";
import React from "react";

type DigitalDownloadsList = {
  data: {
    id: number;
    title: string;
  }[];
};

export default function SessionsTab({ data }: DigitalDownloadsList) {
  return (
    <div className="grid grid-cols-3 gap-3">
      {data.map((item) => (
        <div
          key={item.id}
          className="border rounded-lg overflow-hidden bg-white"
        >
          <Image
            src={"/images/no-image.png"}
            alt={"image"}
            width={500}
            height={500}
            className="w-full h-48 object-cover"
          />
          <div className="p-4 text-sm space-y-2">
            <p className="font-semibold">{item.title}</p>
            <p className="text-gray-500">test summery</p>
            <p>
              <span className="text-gray-500">Start At:</span> 12 Jun 2025 -
              19:24
            </p>
            <p>
              <span className="text-gray-500">Registration ends:</span> 11 Jun
              2025 - 23:27
            </p>
          </div>
          <div className="p-4 flex items-center justify-between border-t border-gray-300">
            <p className="text-sm text-gray-500">2 hours ago</p>
            <span className="uppercase font-medium">Free</span>
          </div>
        </div>
      ))}
    </div>
  );
}
