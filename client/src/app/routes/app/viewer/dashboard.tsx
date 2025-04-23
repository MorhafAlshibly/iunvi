import { useMatch } from "react-router-dom";
import { useQuery } from "@connectrpc/connect-query";
import { paths } from "@/config/paths";
import { useEffect, useRef, useState } from "react";
import { DashboardSelector } from "@/components/dashboard-selector";
import { ContentLayout } from "@/components/layouts/content";
import { getModelRunDashboard } from "@/types/api/dashboard-DashboardService_connectquery";
import { DashboardTransport } from "@/lib/api-client";

const DashboardRoute = () => {
  const id = useMatch(paths.app.viewer.dashboard.getHref(":id"))?.params.id;
  const [dashboardId, setDashboardId] = useState<string | null>(null);
  const [iframeSrc, setIframeSrc] = useState<string | undefined>(undefined);

  const { data: dashboardData } = useQuery(
    getModelRunDashboard,
    {
      modelRunId: id || "",
      dashboardId: dashboardId || "",
    },
    {
      enabled: !!id && !!dashboardId,
      transport: DashboardTransport,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
    },
  );

  useEffect(() => {
    if (!dashboardData?.dashboardSasUrl) return;
    const dashboardUrlParts = dashboardData.dashboardSasUrl.split("?");
    const baseUrl = dashboardUrlParts[0];
    const sasToken = dashboardUrlParts[1];
    const entryUrl = `${baseUrl}/entry.html?${sasToken}`;
    setIframeSrc(entryUrl);
  }, [dashboardData?.dashboardSasUrl]);

  return (
    <ContentLayout title="Dashboard">
      <div className="grid grid-cols-1 gap-4">
        <DashboardSelector
          modelRunId={id || ""}
          selectedDashboardId={dashboardId}
          setSelectedDashboardId={setDashboardId}
        />
        <div className="grid grid-cols-1 col-span-1">
          <iframe
            src={iframeSrc}
            width="100%"
            height="1080px"
            style={{ border: "none" }}
            title="Dashboard"
          />
        </div>
      </div>
    </ContentLayout>
  );
};

export default DashboardRoute;
