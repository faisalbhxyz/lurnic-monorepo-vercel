import React from "react";
import AddNewDigitalAccess from "@/components/shared/home/sales/digital-access/AddNewDigitalAccess";
import DigitalAccessList from "@/components/shared/home/sales/digital-access/DigitalAccessList";
import Link from "next/link";

const sampleData = [
  {
    id: 1,
    name: "Amarrajjo",
    email: "amerrajjonowga.dev@gmail.com",
    webinar: "Title webniar",
    status: "Registered",
    date: "August 21, 2024",
  },
];

export default function page() {
  return (
    <>
      <div className="flex items-center text-sm gap-1 mb-5">
        <Link href={""} className="text-gray-500">
          Sales
        </Link>
        /<Link href={""}>Digital Access</Link>
      </div>
      <div className="flex-between mb-5">
        <h3 className="font-medium text-2xl">Digital Access</h3>
        <AddNewDigitalAccess />
      </div>

      <DigitalAccessList data={sampleData} />
    </>
  );
}
