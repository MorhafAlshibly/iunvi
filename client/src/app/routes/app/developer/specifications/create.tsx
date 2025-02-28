import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useWorkspace } from "@/hooks/use-workspace";
import { createInputSpecification } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { CreateInputSpecificationRequest } from "@/types/api/tenantManagement_pb";
import { useMutation } from "@connectrpc/connect-query";
import { useState } from "react";

const SpecificationsCreateRoute = () => {
  const createInputSpecificationMutation = useMutation(
    createInputSpecification,
  );
  const { activeWorkspace } = useWorkspace();
  const [specification, setSpecification] =
    useState<CreateInputSpecificationRequest>({
      $typeName: "api.CreateInputSpecificationRequest",
      workspaceId: activeWorkspace?.id || "",
      name: "",
      parametersSchema: "",
      csvs: [],
    });

  return (
    <div className="flex flex-col space-y-4">
      <Label className="flex-1 items-center font-semibold text-lg">
        Create specification
      </Label>
      <div className="flex-1">
        <div>
          <Label className="flex items-center">Name</Label>
          <Input
            placeholder="Name"
            onChange={(e) => {
              setSpecification({
                ...specification,
                name: e.target.value,
              });
            }}
          />
        </div>
        <div>
          <Label className="flex items-center">Parameters</Label>
          <Input
            placeholder="Parameters"
            onChange={(e) => {
              setSpecification({
                ...specification,
                parametersSchema: e.target.value,
              });
            }}
          />
        </div>
        <div>
          <Label className="flex items-center">CSVs</Label>
          <Input placeholder="CSVs" />
        </div>
      </div>
    </div>
  );
};

export default SpecificationsCreateRoute;
