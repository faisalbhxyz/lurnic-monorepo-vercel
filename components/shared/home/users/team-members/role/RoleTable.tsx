"use client";

import React from "react";
import RoleActions from "./RoleActions";

const roles = [
  {
    id: 1,
    name: "Administrator",
    permissions: ["Manage Users", "Edit Settings", "Access All Data"],
  },
  {
    id: 2,
    name: "Editor",
    permissions: ["Edit Content", "Publish Posts"],
  },
  {
    id: 3,
    name: "Viewer",
    permissions: ["View Reports", "Read Content"],
  },
];

export default function RoleTable({ roles }: { roles: IRole[] }) {
  return (
    <div className="border rounded-xl overflow-hidden">
      <table className="w-full text-sm">
        <thead className="bg-gray-100">
          <tr className="text-left">
            <th className="px-5 py-2 font-medium">Role</th>
            {/* <th className="px-5 py-2 font-medium">Permissions</th> */}
            <th className="px-5 py-2 font-medium text-end">Action</th>
          </tr>
        </thead>
        <tbody>
          {roles && roles.length > 0 && roles.map((role) => (
            <tr key={role.id} className="border-t border-gray-300">
              <td className="px-5 py-3">{role.name}</td>
              {/* <td className="px-5 py-3">{role.permissions &&role.permissions.join(", ")}</td> */}
              <td className="px-5 py-3 text-end">
                <RoleActions id={role.id} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
