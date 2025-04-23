import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { paths } from "@/config/paths";
import { useWorkspace } from "@/hooks/use-workspace";
import { DashboardTransport } from "@/lib/api-client";
import { getDashboards } from "@/types/api/dashboard-DashboardService_connectquery";
import { useQuery } from "@connectrpc/connect-query";
import { Info } from "lucide-react";
import { useNavigate } from "react-router-dom";

const DashboardsListRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();
  const { data: dashboardsData } = useQuery(
    getDashboards,
    {
      workspaceId: activeWorkspace?.id || "",
    },
    {
      enabled: !!activeWorkspace,
      transport: DashboardTransport,
    },
  );

  return (
    <div className="grid grid-cols-1 gap-4">
      <div className="grid grid-cols-1 col-span-1 justify-items-end">
        <Button
          size="lg"
          variant="default"
          className="mb-4"
          onClick={() => {
            navigate(paths.app.developer.dashboards.create.getHref());
          }}
        >
          Create Dashboard
        </Button>
      </div>
      <div className="grid grid-cols-1 col-span-1">
        {dashboardsData?.dashboards.map((dashboard, index) => (
          <div
            key={index}
            className="grid grid-cols-2 col-span-1 justify-items-between p-2"
          >
            <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
              <Label className="font-normal">{dashboard.name}</Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end">
              <Button
                size="sm"
                variant="ghost"
                onClick={() => {
                  navigate(
                    paths.app.developer.dashboards.view.getHref(dashboard.id),
                  );
                }}
              >
                <Info />
              </Button>
            </div>
            <div className="grid grid-cols-1 col-span-2">
              <Separator className="mt-2" />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default DashboardsListRoute;
