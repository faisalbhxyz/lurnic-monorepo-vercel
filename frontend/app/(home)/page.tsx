import DashboardOverview from "@/components/shared/home/dashboard/DashboardOverview";
import { getDashboardData } from "@/app/actions/dashboard_actions";
import { auth } from "@/lib/auth";

export default async function Home() {
  const session = await auth();
  if (!session) return null;

  const data = await getDashboardData(session);
  const userName =
    (session.user as { name?: string } | undefined)?.name ?? "Admin";

  return <DashboardOverview data={data} userName={userName} />;
}
