import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";
import { useMutation } from "@connectrpc/connect-query";
import { useState } from "react";
import CodeMirror from "@uiw/react-codemirror";
import { markdown } from "@codemirror/lang-markdown";
import { Button } from "@/components/ui/button";
import { paths } from "@/config/paths";
import { useNavigate } from "react-router-dom";
import { ModelSelector } from "@/components/model-selector";
import { createDashboard } from "@/types/api/dashboard-DashboardService_connectquery";
import { CreateDashboardRequest } from "@/types/api/dashboard_pb";
import { DashboardTransport } from "@/lib/api-client";
import { useDarkMode } from "usehooks-ts";

const DashboardsCreateRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();

  const createDashboardMutation = useMutation(createDashboard, {
    transport: DashboardTransport,
  });
  const [dashboard, setDashboard] = useState<CreateDashboardRequest>({
    $typeName: "dashboard.CreateDashboardRequest",
    modelId: "",
    name: "",
    definition: "",
  });

  const validateDashboard = () => {
    if (!dashboard.modelId) return false;
    if (!dashboard.name) return false;
    if (!dashboard.definition) return false;
    return true;
  };

  const darkMode = useDarkMode();

  return (
    <div className="grid grid-cols-1 gap-4">
      <div className="grid grid-cols-2 col-span-1 gap-4 justify-items-between mb-4">
        <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
          Use this page to create dashboards using the Evidence Markdown format.
        </div>
        <div className="grid grid-cols-1 col-span-1 justify-items-end">
          <Button
            size="lg"
            variant="default"
            disabled={!validateDashboard()}
            onClick={() => {
              createDashboardMutation.mutate(dashboard);
              navigate(paths.app.developer.dashboards.list.getHref());
            }}
          >
            Create
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 col-span-1 mt-4">
        <ModelSelector
          selectedModelId={dashboard.modelId}
          setSelectedModelId={(action) => {
            setDashboard({
              ...dashboard,
              modelId:
                (typeof action == "function"
                  ? action(dashboard.modelId)
                  : action) ?? "",
            });
          }}
        />
      </div>
      <div className="grid grid-cols-1 col-span-1">
        <Input
          placeholder="Name"
          className="col-span-1"
          onChange={(e) => {
            setDashboard({
              ...dashboard,
              name: e.target.value,
            });
          }}
        />
      </div>
      <div className="grid grid-cols-1 col-span-1 justify-items-start content-center mt-4">
        <Label className="font-medium text-sm">Dashboard Markdown</Label>
      </div>
      <div className="grid grid-cols-1 col-span-1 border">
        <CodeMirror
          value={dashboard.definition}
          onChange={(value) =>
            setDashboard({
              ...dashboard,
              definition: value,
            })
          }
          theme={darkMode.isDarkMode ? "dark" : "light"}
          extensions={[markdown()]}
        />
      </div>
    </div>
  );
};

export default DashboardsCreateRoute;
