import Link from "next/link";
import { formatDate } from "@/lib/helpers";

export default function RecentEnrollments({
  enrollments,
}: {
  enrollments: Enrollment[];
}) {
  return (
    <div className="border rounded-xl overflow-hidden bg-white">
      <div className="flex-between px-4 py-3 border-b border-gray-200">
        <h4 className="font-medium text-base">Recent Enrollments</h4>
        <Link href="/enrollments" className="text-sm text-primary hover:underline">
          View all
        </Link>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-gray-50">
            <tr className="text-left text-gray-500">
              <th className="p-3 font-medium">Student</th>
              <th className="p-3 font-medium">Course</th>
              <th className="p-3 font-medium text-end">Date</th>
            </tr>
          </thead>
          <tbody>
            {enrollments.map((entry) => (
              <tr key={entry.id} className="border-t border-gray-200">
                <td className="p-3">
                  <p className="font-medium text-gray-900">
                    {entry.student.first_name} {entry.student.last_name}
                  </p>
                  <p className="text-xs text-gray-500">{entry.student.email}</p>
                </td>
                <td className="p-3 text-gray-700">{entry.course.title}</td>
                <td className="p-3 text-end text-gray-600">
                  {formatDate(entry.created_at)}
                </td>
              </tr>
            ))}
            {enrollments.length === 0 && (
              <tr>
                <td colSpan={3} className="p-6 text-center text-gray-500">
                  No enrollments yet.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
