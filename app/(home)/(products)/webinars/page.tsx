import NewWebinar from "@/components/shared/home/products/webinars/NewWebinar";
import WebinerList from "@/components/shared/home/products/webinars/WebinerList";
import Link from "next/link";
import React from "react";

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
        /<Link href={""}>Webinars</Link>
      </div>
      <div className="flex-between mt-5">
        <h3 className="font-medium text-2xl">Webinars</h3>
        <NewWebinar />
      </div>
      <div className="mt-5">
        <WebinerList data={sampleData} />
      </div>
    </>
  );
}
