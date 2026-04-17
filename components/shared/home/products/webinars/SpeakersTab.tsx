import React from "react";

const sampleSpeakers = [
  {
    name: "Alice Johnson",
    email: "alice.johnson@example.com",
    designation: "Lead Developer",
    company: "TechNova Inc.",
  },
  {
    name: "Brian Lee",
    email: "brian.lee@example.com",
    designation: "UX Designer",
    company: "Creative Minds",
  },
];

export default function SpeakersTab() {
  return (
    <div className="overflow-x-auto rounded-lg border">
      <table className="min-w-full divide-y divide-gray-200 text-sm text-left">
        <thead className="bg-gray-100 text-gray-700 tracking-wider">
          <tr>
            <th className="px-6 py-3 text-sm font-medium">Name</th>
            <th className="px-6 py-3 text-sm font-medium">Email</th>
            <th className="px-6 py-3 text-sm font-medium">Designation</th>
            <th className="px-6 py-3 text-sm font-medium">Company</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200">
          {sampleSpeakers.map((speaker, index) => (
            <tr key={index} className="hover:bg-gray-50">
              <td className="px-6 py-4">{speaker.name}</td>
              <td className="px-6 py-4">{speaker.email}</td>
              <td className="px-6 py-4">{speaker.designation}</td>
              <td className="px-6 py-4">{speaker.company}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
