"use client";

import ToggleSwitch from "@/components/ui/ToggleSwitch";
import React, { useState } from "react";
import TeamMemberActions from "./TeamMemberActions";

const sampleMembers = [
  {
    id: 1,
    name: "Alice Johnson",
    email: "alice@example.com",
    status: "Active",
  },
  {
    id: 2,
    name: "Bob Smith",
    email: "bob@example.com",
    status: "Invited",
  },
  {
    id: 3,
    name: "Charlie Brown",
    email: "charlie@example.com",
    status: "Inactive",
  },
];

export default function TeamMemberList({ users }: { users: IUser[] }) {
  return (
    <>
      <div className="border rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-100">
            <tr className="text-left">
              <th className="px-5 py-2 font-medium">User ID</th>
              <th className="px-5 py-2 font-medium">Name & Email</th>
              {/* <th className="px-5 py-2 font-medium">Status</th> */}
              <th className="px-5 py-2 font-medium text-end">Action</th>
            </tr>
          </thead>
          <tbody>
            {users &&
              users.length > 0 &&
              users.map((member) => (
                <tr key={member.id} className="border-t border-gray-300">
                  <td className="px-5 py-2">{member.user_id}</td>
                  <td className="px-5 py-2">
                    <div className="font-medium">{member.name}</div>
                    <div className="text-gray-500">{member.email}</div>
                  </td>
                  {/* <td className="px-5 py-2">
                    <span
                      className={`inline-block px-2 py-1 rounded text-xs font-medium`}
                    >
                      <ToggleSwitch />
                    </span>
                  </td> */}
                  <td className="px-5 py-2 text-end">
                    <TeamMemberActions id={member.id} />
                  </td>
                </tr>
              ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
