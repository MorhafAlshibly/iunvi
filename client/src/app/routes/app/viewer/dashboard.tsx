import { Label } from "@radix-ui/react-dropdown-menu";
import { useMatch } from "react-router-dom";
import { useQuery } from "@connectrpc/connect-query";
import {
  getModelRunDashboard,
  getSpecification,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { paths } from "@/config/paths";
import { DataMode, TableFieldType } from "@/types/api/tenantManagement_pb";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router";
import { ArrowBigLeft, CircleArrowLeft } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { useState } from "react";
import { DashboardSelector } from "@/components/dashboard-selector";
import { ContentLayout } from "@/components/layouts/content";

const DashboardRoute = () => {
  const navigate = useNavigate();
  const id = useMatch(paths.app.viewer.dashboard.getHref(":id"))?.params.id;

  const [dashboardId, setDashboardId] = useState<string | null>(null);

  const { data: dashboardData } = useQuery(
    getModelRunDashboard,
    {
      modelRunId: id || "",
      dashboardId: dashboardId || "",
    },
    {
      enabled: !!id && !!dashboardId,
    },
  );

  return (
    <ContentLayout title="Dashboard">
      <div className="grid grid-cols-1 gap-4">
        <div className="grid grid-cols-1 col-span-1 justify-items-start">
          <DashboardSelector
            modelRunId={id || ""}
            selectedDashboardId={dashboardId}
            setSelectedDashboardId={setDashboardId}
          />
        </div>
        <div className="grid grid-cols-1 col-span-1"></div>
      </div>
    </ContentLayout>
  );
};

export default DashboardRoute;
