import { ContentLayout } from "@/components/layouts/content";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { paths } from "@/config/paths";
import { useWorkspace } from "@/hooks/use-workspace";
import { getModelRuns } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useQuery } from "@connectrpc/connect-query";
import { Info } from "lucide-react";
import { useNavigate } from "react-router-dom";

const RunHistoryRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();
  const { data: modelRunsData } = useQuery(
    getModelRuns,
    {
      workspaceId: activeWorkspace?.id || "",
    },
    {
      enabled: !!activeWorkspace,
    },
  );

  return (
    <ContentLayout title="Run history">
      <div className="grid grid-cols-1 gap-4">
        <div className="grid grid-cols-1 col-span-1">
          {modelRunsData?.modelRuns.map((modelRun, index) => (
            <div
              key={index}
              className="grid grid-cols-2 col-span-1 justify-items-between p-2"
            >
              <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
                <Label className="font-normal">{modelRun.name}</Label>
              </div>
              <div className="grid grid-cols-1 col-span-1 justify-items-end">
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() => {
                    navigate(paths.app.viewer.dashboard.getHref(modelRun.id));
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
    </ContentLayout>
  );
};

export default RunHistoryRoute;
