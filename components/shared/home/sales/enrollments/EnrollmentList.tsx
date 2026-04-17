"use client";

import Button from "@/components/ui/Button";
import Checkbox from "@/components/ui/Checkbox";
import axiosInstance from "@/lib/axiosInstance";
import { formatDate } from "@/lib/helpers";
import { Session } from "next-auth";
import { useRouter } from "next/navigation";
import React, { useState } from "react";
import { FiSearch } from "react-icons/fi";
import { toast } from "sonner";

interface EnrollmentListProps {
  data: Enrollment[];
  session: Session;
}

export default function EnrollmentList({ data, session }: EnrollmentListProps) {
  const router = useRouter();
  const [selected, setSelected] = useState<number[]>([]);
  const [search, setSearch] = useState("");

  const filteredEnrollments = data.filter((entry) => {
    const fullName = `${entry.student.first_name} ${entry.student.last_name}`;
    return (
      fullName.toLowerCase().includes(search.toLowerCase()) ||
      entry.course.title.toLowerCase().includes(search.toLowerCase())
    );
  });

  const toggleSelectAll = () => {
    if (selected.length === filteredEnrollments.length) {
      setSelected([]);
    } else {
      setSelected(filteredEnrollments.map((entry) => entry.id));
    }
  };

  const toggleSelectOne = (id: number) => {
    setSelected((prev) =>
      prev.includes(id) ? prev.filter((i) => i !== id) : [...prev, id]
    );
  };

  const isAllSelected =
    selected.length === filteredEnrollments.length &&
    filteredEnrollments.length > 0;

  const handleDelete = (id: number) => {
    if (!confirm("Are you sure you want to delete this enrollment?")) return;

    axiosInstance
      .delete(`/private/enrollment/delete/${id}`, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${session?.accessToken}`,
        },
      })
      .then((res) => {
        router.refresh();
        toast.success("Enrollment deleted successfully.");
      })
      .catch((error) => {
        toast.error(error.response.data.error || "Something went wrong.");
      });
  };

  return (
    <>
      <div className="flex items-center justify-between mb-5">
        <div></div>
        <div className="relative w-full max-w-60">
          <FiSearch className="absolute top-1/2 left-3 transform -translate-y-1/2 text-gray-500" />
          <input
            type="text"
            placeholder="Search by student or course..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full border text-sm px-3 py-1.5 rounded-md pl-8"
          />
        </div>
      </div>
      <div className="border rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-100">
            <tr className="text-left">
              <th className="p-3">
                <div className="flex items-center gap-5 font-medium">
                  <Checkbox
                    checked={isAllSelected}
                    onChange={toggleSelectAll}
                  />
                  <span>Name</span>
                </div>
              </th>
              <th className="p-3 font-medium">Email</th>
              <th className="p-3 font-medium">Course</th>
              <th className="p-3 font-medium text-end">Enrolled At</th>
              <th className="p-3 font-medium text-end">Action</th>
            </tr>
          </thead>
          <tbody>
            {filteredEnrollments.map((entry) => (
              <tr
                key={entry.id}
                className="border-t border-gray-300 hover:bg-gray-100"
              >
                <td className="p-3">
                  <div className="flex items-center gap-5">
                    <Checkbox
                      checked={selected.includes(entry.id)}
                      onChange={() => toggleSelectOne(entry.id)}
                    />
                    <div className="flex items-center gap-3">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">
                          {entry.student.first_name} {entry.student.last_name}
                        </p>
                      </div>
                    </div>
                  </div>
                </td>
                <td className="p-3">{entry.student.email}</td>
                <td className="p-3">
                  <div className="flex items-center gap-2">
                    {entry.course.title}
                  </div>
                </td>
                <td className="p-3 text-end">{formatDate(entry.created_at)}</td>
                <td className="p-3 flex items-center justify-end">
                  <Button
                    className="bg-red-500 hover:bg-red-600"
                    onClick={() => handleDelete(entry.id)}
                  >
                    Delete
                  </Button>
                </td>
              </tr>
            ))}
            {filteredEnrollments.length === 0 && (
              <tr>
                <td colSpan={6} className="p-5 text-center text-gray-500">
                  No enrollments found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </>
  );
}
