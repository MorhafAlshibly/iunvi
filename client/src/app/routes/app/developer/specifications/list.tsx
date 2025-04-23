import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { paths } from "@/config/paths";
import { useWorkspace } from "@/hooks/use-workspace";
import { FileTransport } from "@/lib/api-client";
import { getSpecifications } from "@/types/api/file-FileService_connectquery";
import { useQuery } from "@connectrpc/connect-query";
import { Info } from "lucide-react";
import { useNavigate } from "react-router-dom";

const SpecificationsListRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();
  const { data: specificationsData } = useQuery(
    getSpecifications,
    {
      workspaceId: activeWorkspace?.id || "",
    },
    {
      enabled: !!activeWorkspace,
      transport: FileTransport,
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
            navigate(paths.app.developer.specifications.create.getHref());
          }}
        >
          Create Specification
        </Button>
      </div>
      <div className="grid grid-cols-1 col-span-1">
        {specificationsData?.specifications.map((specification, index) => (
          <div
            key={index}
            className="grid grid-cols-2 col-span-1 justify-items-between p-2"
          >
            <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
              <Label className="font-normal">{specification.name}</Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end">
              <Button
                size="sm"
                variant="ghost"
                onClick={() => {
                  navigate(
                    paths.app.developer.specifications.view.getHref(
                      specification.id,
                    ),
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

export default SpecificationsListRoute;
