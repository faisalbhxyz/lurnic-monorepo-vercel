import React, { ReactNode } from "react";

export default function layout({ children }: { children: ReactNode }) {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gradient-to-b from-gray-100 to-primary/40">
      {children}
    </div>
  );
}
