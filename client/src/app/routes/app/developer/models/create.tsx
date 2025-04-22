import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";
import { useMutation } from "@connectrpc/connect-query";
import { useState } from "react";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import { CircleX, Cross, PlusCircle } from "lucide-react";
import { paths } from "@/config/paths";
import { useNavigate } from "react-router-dom";
import CreateDataTable from "@/components/create-data-table";

import { SpecificationSelector } from "@/components/specification-selector";
import { ImageSelector } from "@/components/image-selector";
import { createModel } from "@/types/api/model-ModelService_connectquery";
import { CreateModelRequest } from "@/types/api/model_pb";
import { DataMode } from "@/types/api/file_pb";

const ModelsCreateRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();

  const createModelMutation = useMutation(createModel);
  const [model, setModel] = useState<CreateModelRequest>({
    $typeName: "model.CreateModelRequest",
    inputSpecificationId: "",
    outputSpecificationId: "",
    parametersSchema: undefined,
    name: "",
    imageName: "",
  });

  const validateModel = () => {
    if (!model.inputSpecificationId) return false;
    if (!model.outputSpecificationId) return false;
    if (!model.name) return false;
    if (!model.imageName) return false;
    if (model.parametersSchema) {
      try {
        JSON.parse(model.parametersSchema);
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
          Use this page to create models based on an input and output. Utilize
          parameters to create an input form using react-jsonschema-form.
        </div>
        <div className="grid grid-cols-1 col-span-1 justify-items-end">
          <Button
            size="lg"
            variant="default"
            disabled={!validateModel()}
            onClick={() => {
              createModelMutation.mutate(model);
              navigate(paths.app.developer.models.list.getHref());
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
            setModel({
              ...model,
              name: e.target.value,
            });
          }}
        />
      </div>
      <div className="grid grid-cols-1 col-span-1 mt-4">
        <SpecificationSelector
          mode={DataMode.INPUT}
          selectedSpecificationId={model.inputSpecificationId}
          setSelectedSpecificationId={(action) => {
            setModel({
              ...model,
              inputSpecificationId:
                (typeof action == "function"
                  ? action(model.inputSpecificationId)
                  : action) ?? "",
            });
          }}
        />
      </div>
      <div className="grid grid-cols-1 col-span-1">
        <SpecificationSelector
          mode={DataMode.OUTPUT}
          selectedSpecificationId={model.outputSpecificationId}
          setSelectedSpecificationId={(action) => {
            setModel({
              ...model,
              outputSpecificationId:
                (typeof action == "function"
                  ? action(model.outputSpecificationId)
                  : action) ?? "",
            });
          }}
        />
      </div>{" "}
      <div className="grid grid-cols-1 col-span-1 mt-4">
        <ImageSelector
          selectedImageName={model.imageName}
          setSelectedImageName={(action) => {
            setModel({
              ...model,
              imageName:
                (typeof action == "function"
                  ? action(model.imageName)
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
          value={model.parametersSchema || ""}
          onChange={(value) =>
            setModel({
              ...model,
              parametersSchema: value.trim() ? value : undefined,
            })
          }
          extensions={[json()]}
        />
      </div>
    </div>
  );
};

export default ModelsCreateRoute;
