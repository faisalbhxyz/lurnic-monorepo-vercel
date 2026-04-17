import React from "react";
import TeamMemberList from "@/components/shared/home/users/team-members/user/TeamMemberList";
import Link from "next/link";
import AddNewTeamMember from "@/components/shared/home/users/team-members/user/AddNewTeamMember";
import { getAllUsers } from "@/app/actions/team_actions";
import { auth } from "@/lib/auth";
import UpdateTeamMember from "@/components/shared/home/users/team-members/user/UpdateTeamMember";

export default async function page() {
  const session = await auth();
  if (!session) return null;

  const users = await getAllUsers(session);

  return (
    <>
      <div className="flex items-center text-sm gap-1">
        <Link href="/users" className="text-gray-500">
          Users
        </Link>
        /
        <Link href="/users/team-members" className="text-gray-700">
          Team Members
        </Link>
      </div>
      <div className="flex-between my-5">
        <h3 className="font-medium text-2xl">Team Members</h3>
        <AddNewTeamMember />
      </div>
      <div className="flex gap-3 mb-5">
        <Link
          href="/team-members"
          className="px-3 py-3 text-sm border-b border-primary text-primary"
        >
          Team Members
        </Link>
        <Link href="/team-members/role" className="px-3 py-3 text-sm">
          Role
        </Link>
      </div>
      <TeamMemberList users={users} />
      <UpdateTeamMember />
    </>
  );
}
