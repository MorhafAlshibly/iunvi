import { Label } from "@radix-ui/react-dropdown-menu";
import { useMatch } from "react-router-dom";
import { useQuery } from "@connectrpc/connect-query";
import { paths } from "@/config/paths";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { ArrowBigLeft, CircleArrowLeft, Info } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { getModel } from "@/types/api/model-ModelService_connectquery";
import { DashboardTransport, ModelTransport } from "@/lib/api-client";
import CodeMirror from "@uiw/react-codemirror";
import { useDarkMode } from "usehooks-ts";
import {
  getDashboard,
  getDashboardMarkdown,
} from "@/types/api/dashboard-DashboardService_connectquery";
import { markdown } from "@codemirror/lang-markdown";

const DashboardsViewRoute = () => {
  const navigate = useNavigate();
  const id = useMatch(paths.app.developer.dashboards.root.getHref() + "/:id")
    ?.params.id;

  const { data: dashboardData } = useQuery(
    getDashboard,
    {
      id: id || "",
    },
    {
      enabled: !!id,
      transport: DashboardTransport,
    },
  );

  const dashboard = dashboardData?.dashboard;

  const { data: dashboardMarkdownData } = useQuery(
    getDashboardMarkdown,
    {
      id: dashboard?.id || "",
    },
    {
      enabled: !!id,
      transport: DashboardTransport,
    },
  );

  const dashboardMarkdown = dashboardMarkdownData?.markdown;

  const { data: modelData } = useQuery(
    getModel,
    {
      id: dashboard?.modelId || "",
    },
    {
      enabled: !!dashboard?.modelId,
      transport: ModelTransport,
    },
  );

  const model = modelData?.model;

  const darkMode = useDarkMode();

  return (
    <div className="grid grid-cols-1 gap-4">
      {dashboard && model ? (
        <>
          <div className="grid grid-cols-2 col-span-1 justify-items-between">
            <div className="grid grid-cols-1 col-span-1 justify-items-start">
              <Label className="col-span-1 content-center text-lg font-medium">
                {dashboard.name}
              </Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end">
              <Button
                size="sm"
                variant="outline"
                onClick={() => {
                  navigate(-1);
                }}
              >
                <CircleArrowLeft />
                Back
              </Button>
            </div>
          </div>
          <div className="grid grid-cols-2 col-span-1 justify-items-between p-2 mt-4">
            <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
              <Label className="font-normal">Model</Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end content-center">
              <Button
                size="sm"
                variant="outline"
                onClick={() => {
                  navigate(paths.app.developer.models.view.getHref(model.id));
                }}
              >
                <Info />
                {model.name}
              </Button>
            </div>
            <div className="grid grid-cols-1 col-span-2">
              <Separator className="mt-2" />
            </div>
          </div>
          {dashboardMarkdown ? (
            <>
              <div className="grid grid-cols-1 col-span-1 content-center mt-4">
                <Label className="font-medium">Evidence Markdown</Label>
              </div>
              <div className="grid grid-cols-1 col-span-1">
                <CodeMirror
                  value={dashboardMarkdown}
                  height="auto"
                  extensions={[markdown()]}
                  editable={false}
                  theme={darkMode.isDarkMode ? "dark" : "light"}
                  className="col-span-1 border"
                />
              </div>
            </>
          ) : null}
        </>
      ) : null}
    </div>
  );
};

export default DashboardsViewRoute;
