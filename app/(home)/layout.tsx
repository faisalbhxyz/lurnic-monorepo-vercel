import Sidebar from "@/components/shared/sidebar/Sidebar";
import { auth } from "@/lib/auth";
import React, { ReactNode } from "react";
import { getGeneralSettings } from "../actions/actions";

type LayoutProps = {
  children: ReactNode;
};

export default async function Layout({ children }: LayoutProps) {
  const session = await auth();
  if (!session) return null;

  const generalSettings = await getGeneralSettings(session);

  return (
    <div className="flex bg-slate-100">
      <Sidebar orgLogo={generalSettings?.logo} />
      <main className="flex-1 overflow-y-auto p-5">
        <div className="bg-white h-full rounded-2xl p-5">{children}</div>
      </main>
    </div>
  );
}
