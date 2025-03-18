import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";
import { createRunModel } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import {
  CreateRunModelRequest,
  CreateRunModelRequestSchema,
  DataMode,
  TableFieldType,
} from "@/types/api/tenantManagement_pb";
import { useMutation } from "@connectrpc/connect-query";
import { useState } from "react";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import { CircleX, Cross, PlusCircle } from "lucide-react";
import { paths } from "@/config/paths";
import { useNavigate } from "react-router";
import CreateDataTable from "@/components/create-data-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { SpecificationSelector } from "@/components/specification-selector";
import { ImageSelector } from "@/components/image-selector";

const RunModelsCreateRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();

  const createRunModelMutation = useMutation(createRunModel);
  const [runmodel, setRunModel] = useState<CreateRunModelRequest>({
    $typeName: "api.CreateRunModelRequest",
    inputSpecificationId: "",
    outputSpecificationId: "",
    parametersSchema: undefined,
    name: "",
    imageName: "",
  });

  const validateRunModel = () => {
    if (!runmodel.inputSpecificationId) return false;
    if (!runmodel.outputSpecificationId) return false;
    if (!runmodel.name) return false;
    if (!runmodel.imageName) return false;
    if (runmodel.parametersSchema) {
      try {
        JSON.parse(runmodel.parametersSchema);
      } catch (e) {
        return false;
      }
    }
    return true;
  };

  return (
    <div className="grid grid-cols-1 gap-4">
      <div className="grid grid-cols-2 col-span-1 gap-4 justify-items-between mb-4">
        <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
          Use this page to create runmodels based on an input and output.
          Utilize parameters to create an input form using
          react-jsonschema-form.
        </div>
        <div className="grid grid-cols-1 col-span-1 justify-items-end">
          <Button
            size="lg"
            variant="default"
            disabled={!validateRunModel()}
            onClick={() => {
              createRunModelMutation.mutate(runmodel);
              navigate(paths.app.developer.runmodels.list.getHref());
            }}
          >
            Create
          </Button>
        </div>
      </div>
      <div className="grid grid-cols-1 col-span-1">
        <Input
          placeholder="Name"
          className="col-span-1"
          onChange={(e) => {
            setRunModel({
              ...runmodel,
              name: e.target.value,
            });
          }}
        />
      </div>
      <div className="grid grid-cols-1 col-span-1 mt-4">
        <SpecificationSelector
          mode={DataMode.INPUT}
          selectedSpecificationId={runmodel.inputSpecificationId}
          setSelectedSpecificationId={(action) => {
            setRunModel({
              ...runmodel,
              inputSpecificationId:
                (typeof action == "function"
                  ? action(runmodel.inputSpecificationId)
                  : action) ?? "",
            });
          }}
        />
      </div>
      <div className="grid grid-cols-1 col-span-1">
        <SpecificationSelector
          mode={DataMode.OUTPUT}
          selectedSpecificationId={runmodel.outputSpecificationId}
          setSelectedSpecificationId={(action) => {
            setRunModel({
              ...runmodel,
              outputSpecificationId:
                (typeof action == "function"
                  ? action(runmodel.outputSpecificationId)
                  : action) ?? "",
            });
          }}
        />
      </div>{" "}
      <div className="grid grid-cols-1 col-span-1 mt-4">
        <ImageSelector
          selectedImageName={runmodel.imageName}
          setSelectedImageName={(action) => {
            setRunModel({
              ...runmodel,
              imageName:
                (typeof action == "function"
                  ? action(runmodel.imageName)
                  : action) ?? "",
            });
          }}
        />
      </div>
      <div className="grid grid-cols-1 col-span-1 justify-items-start content-center mt-4">
        <Label className="font-medium text-sm">Parameters schema</Label>
      </div>
      <div className="grid grid-cols-1 col-span-1 border">
        <CodeMirror
          value={runmodel.parametersSchema || ""}
          onChange={(value) =>
            setRunModel({
              ...runmodel,
              parametersSchema: value.trim() ? value : undefined,
            })
          }
          extensions={[json()]}
        />
      </div>
    </div>
  );
};

export default RunModelsCreateRoute;
