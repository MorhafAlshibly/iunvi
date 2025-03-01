import { Button } from "@/components/ui/button";
import { paths } from "@/config/paths";
import { useWorkspace } from "@/hooks/use-workspace";
import { getSpecifications } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useQuery } from "@connectrpc/connect-query";
import { Info } from "lucide-react";
import { useNavigate } from "react-router";

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
      <div className="grid grid-cols-1 col-span-1 gap-4">
        {specificationsData?.specifications.map((specification, index) => (
          <div
            key={index}
            className="grid grid-cols-2 col-span-1 justify-items-between border p-4"
          >
            <span className="grid grid-cols-1 col-span-1 justify-items-start content-center">
              {specification.name}
            </span>
            <span className="grid grid-cols-1 col-span-1 justify-items-end">
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
            </span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SpecificationsListRoute;
