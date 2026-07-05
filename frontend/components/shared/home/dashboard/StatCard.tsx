import React from "react";

interface StatCardProps {
  label: string;
  value: string | number;
  hint?: string;
  icon: React.ReactNode;
  iconClassName?: string;
}

export default function StatCard({
  label,
  value,
  hint,
  icon,
  iconClassName = "bg-primary/10 text-primary",
}: StatCardProps) {
  return (
    <div className="border rounded-xl p-4 bg-white">
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <p className="text-sm text-gray-500">{label}</p>
          <p className="text-2xl font-semibold text-gray-900 mt-1 truncate">
            {value}
          </p>
          {hint ? <p className="text-xs text-gray-400 mt-1">{hint}</p> : null}
        </div>
        <div
          className={`shrink-0 w-10 h-10 rounded-lg flex items-center justify-center ${iconClassName}`}
        >
          {icon}
        </div>
      </div>
    </div>
  );
}
