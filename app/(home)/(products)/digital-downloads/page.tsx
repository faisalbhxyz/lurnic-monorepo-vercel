import React from "react";
import DigitalDownloadsList from "@/components/shared/home/products/digital-downloads/DigitalDownloadsList";
import NewDigitalDownload from "@/components/shared/home/products/digital-downloads/NewDigitalDownload";
import Link from "next/link";

const sampleData = [
  {
    id: 1,
    title: "Digital Download",
  },
];

export default function page() {
  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href={""} className="text-gray-500">
          Products
        </Link>
        /<Link href={""}>Digital Downloads</Link>
      </div>
      <div className="flex-between mt-5">
        <h3 className="font-medium text-2xl">Digital Downloads</h3>
        <NewDigitalDownload />
      </div>
      <div className="mt-5 grid grid-cols-3 gap-3">
        {sampleData.map((item) => (
          <DigitalDownloadsList key={item.id} data={item} />
        ))}
      </div>
    </>
  );
}
