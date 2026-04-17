import { getGeneralSettings } from "@/app/actions/actions";
import SettingTabs from "@/components/shared/home/settings/SettingTabs";
import Button from "@/components/ui/Button";
import { auth } from "@/lib/auth";
import axiosInstance from "@/lib/axiosInstance";
import { Session } from "next-auth";
import React from "react";



export default async function page() {
  const session = await auth();
  if (!session) return null;

  const generalSettings = await getGeneralSettings(session);

  // console.log("generalSettings", generalSettings);

  return (
    <>
      <SettingTabs generalsettings={generalSettings} />
    </>
  );
}
