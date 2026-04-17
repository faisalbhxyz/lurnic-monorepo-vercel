"use client";

import React, { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import Profile from "./Profile";
import { IoIosArrowDown } from "react-icons/io";
import { LuSettings } from "react-icons/lu";
import { cn } from "@/lib/cn";
import defaultLogo from "@/public/logo/dashlearn-logo.svg";

// Icons
import {
  HiOutlineViewGrid,
  HiOutlineCube,
  HiOutlineCurrencyDollar,
  HiOutlineUserGroup,
} from "react-icons/hi";
import Image from "next/image";
import { MdWebStories } from "react-icons/md";
import { IoEarthSharp } from "react-icons/io5";

type MenuItem = {
  label: string;
  href?: string;
  submenu?: MenuItem[];
  icon?: React.ReactNode;
};

const menuItems: MenuItem[] = [
  { label: "Dashboard", href: "/", icon: <HiOutlineViewGrid size={18} /> },
  {
    label: "Products",
    icon: <HiOutlineCube size={18} />,
    submenu: [
      { label: "Courses", href: "/courses" },
      // { label: "Digital Downloads", href: "/digital-downloads" },
      // { label: "Webinars", href: "/webinars" },
      { label: "Category", href: "/categories" },
      { label: "Sub Category", href: "/subcategory" },
    ],
  },
  {
    label: "Sales",
    icon: <HiOutlineCurrencyDollar size={18} />,
    submenu: [
      { label: "Orders", href: "/orders" },
      { label: "Enrollments", href: "/enrollments" },
      // { label: "Digital Access", href: "/digital-access" },
      // { label: "Webinar Access", href: "/webinar-access" },
      // { label: "Coupons", href: "/coupons" },
    ],
  },
  {
    label: "Users",
    icon: <HiOutlineUserGroup size={18} />,
    submenu: [
      { label: "Students", href: "/students" },
      { label: "Instructors", href: "/instructors" },
      { label: "Team Members", href: "/team-members" },
    ],
  },
  {
    label: "Website",
    icon: <IoEarthSharp size={18} />,
    submenu: [{ label: "Banner", href: "/banner" }],
  },
];

export default function Sidebar({ orgLogo }: { orgLogo: string | null }) {
  const pathname = usePathname();
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggleMenu = (index: number) => {
    setOpenIndex((prev) => (prev === index ? null : index));
  };

  const isActive = (href?: string) => href && pathname === href;

  return (
    <div className="w-64 h-screen bg-slate-100 flex flex-col sticky top-0">
      <div className="p-4">
        <Link href="/" className="text-2xl font-bold block">
          {orgLogo ? (
            <Image
              src={orgLogo}
              width={350}
              height={75}
              alt="image"
              className="w-[115px] h-[30px]"
            />
          ) : (
            <Image
              src={defaultLogo}
              width={350}
              height={75}
              alt="image"
              className="w-[115px] h-[30px]"
            />
          )}
        </Link>
      </div>
      <ul className="flex-1 overflow-y-auto p-2 text-sm space-y-1">
        {menuItems.map((item, index) => (
          <li key={index}>
            <Link
              href={item.href || "#"}
              className={cn(
                "flex items-center justify-between py-1.5 px-3 rounded-md cursor-pointer hover:bg-primary/10 hover:text-primary w-full",
                openIndex === index && "bg-primary/10 text-primary"
              )}
              onClick={() => toggleMenu(index)}
            >
              <div className="flex items-center gap-2 font-medium">
                {item.icon}
                <span>{item.label}</span>
              </div>
              {item.submenu && (
                <IoIosArrowDown
                  className={`transition-all duration-300 ${
                    openIndex === index && "rotate-180"
                  }`}
                />
              )}
            </Link>

            {item.submenu && (
              <div
                className={cn(
                  "grid transition-all duration-300 overflow-hidden ease-in-out",
                  openIndex === index
                    ? "grid-rows-[1fr] opacity-100 mt-1"
                    : "grid-rows-[0fr] opacity-0"
                )}
              >
                <div className="relative overflow-hidden space-y-1 ml-5">
                  <div className="absolute w-0.5 bg-gray-300 top-4 bottom-3" />
                  {item.submenu.map((submenu, subIndex) => (
                    <Link
                      key={subIndex}
                      href={submenu.href || "#"}
                      className={cn(
                        "relative block text-gray-700 font-medium px-3 py-1.5 hover:text-primary rounded-md",
                        isActive(submenu.href) && "text-primary"
                      )}
                    >
                      {isActive(submenu.href) && (
                        <div className="absolute left-0 h-6 w-0.5 bg-primary" />
                      )}
                      {submenu.label}
                    </Link>
                  ))}
                </div>
              </div>
            )}
          </li>
        ))}
      </ul>
      <div>
        <Link
          href="/settings"
          className={cn(
            "px-4 py-2 font-medium text-sm hover:bg-primary/10 flex items-center gap-2",
            isActive("/settings")
              ? "text-primary"
              : "text-gray-800 hover:text-primary"
          )}
        >
          <LuSettings size={18} /> Settings
        </Link>
        <div className="p-3 border-t border-slate-300">
          <Profile />
        </div>
      </div>
    </div>
  );
}
