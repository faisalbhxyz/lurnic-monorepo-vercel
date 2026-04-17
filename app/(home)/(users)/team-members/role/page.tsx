import { getAllRoles } from "@/app/actions/team_actions";
import AddNewRole from "@/components/shared/home/users/team-members/role/AddNewRole";
import RoleTable from "@/components/shared/home/users/team-members/role/RoleTable";
import { auth } from "@/lib/auth";
import Link from "next/link";
import React from "react";

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const roles = await getAllRoles(session);

  return (
    <>
      <div className="flex items-center text-sm gap-1 text-gray-500">
        <Link href="">Users</Link>/
        <Link href="/team-members">Team Members</Link>/
        <Link href="/team-members/role" className="text-gray-700">
          Role
        </Link>
      </div>
      <div className="flex-between my-5">
        <h3 className="font-medium text-2xl">Team Members</h3>
        <AddNewRole />
      </div>
      <div className="flex gap-3 mb-5">
        <Link href="/team-members" className="px-3 py-3 text-sm">
          Team Members
        </Link>
        <Link
          href="/team-members/role"
          className="px-3 py-3 text-sm border-b border-primary text-primary"
        >
          Role
        </Link>
      </div>

      <RoleTable roles={roles}/>
    </>
  );
}
