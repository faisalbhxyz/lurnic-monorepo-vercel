import React from "react";
import AddNewWebinarAccess from "@/components/shared/home/sales/webinar-access/AddNewWebinarAccess";
import WebinarAccessList from "@/components/shared/home/sales/webinar-access/WebinarAccessList";
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
      <div className="flex items-center text-sm gap-1">
        <Link href={""} className="text-gray-500">
          Sales
        </Link>
        /<Link href={""}>Webinar Access</Link>
      </div>
      <div className="flex-between my-5">
        <h3 className="font-medium text-2xl">Webinar Access</h3>
        <AddNewWebinarAccess />
      </div>
      <WebinarAccessList data={sampleData} />
    </>
  );
}
