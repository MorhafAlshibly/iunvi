import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";

import { useMutation, useQuery } from "@connectrpc/connect-query";
import { createRef, useEffect, useState } from "react";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import { CircleX, Cross, PlusCircle } from "lucide-react";
import { paths } from "@/config/paths";
import { useNavigate } from "react-router-dom";
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
import { FileGroupSelector } from "@/components/file-group-selector";
import { ModelSelector } from "@/components/model-selector";
import { ContentLayout } from "@/components/layouts/content";
import Form from "@rjsf/core";
import { RJSFSchema } from "@rjsf/utils";
import validator from "@rjsf/validator-ajv8";
import {
  createModelRun,
  getModel,
} from "@/types/api/model-ModelService_connectquery";
import { CreateModelRunRequest } from "@/types/api/model_pb";

const RunModelsCreateRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();

  const createModelRunMutation = useMutation(createModelRun);
  const [modelRun, setModelRun] = useState<CreateModelRunRequest>({
    $typeName: "model.CreateModelRunRequest",
    modelId: "",
    inputFileGroupId: "",
    parameters: undefined,
    name: "",
  });

  const { data: modelData } = useQuery(
    getModel,
    {
      id: modelRun.modelId,
    },
    {
      enabled: !!modelRun.modelId,
    },
  );

  const model = modelData?.model;

  const validateModelRun = () => {
    if (!modelRun.modelId) return false;
    if (!modelRun.inputFileGroupId) return false;
    if (!modelRun.name) return false;
    if (model && model.parametersSchema && !modelRun.parameters) return false;
    if (model?.parametersSchema && !formRef.current?.validateForm())
      return false;
    return true;
  };

  useEffect(() => {
    // Reset the model run when the model changes
    setModelRun((modelRun) => ({
      ...modelRun,
      inputFileGroupId: "",
      parameters: undefined,
    }));
  }, [modelRun.modelId]);

  const formRef = createRef<Form>();

  return (
    <ContentLayout title="Run model">
      <div className="grid grid-cols-1 gap-4">
        <div className="grid grid-cols-2 col-span-1 gap-4 justify-items-between mb-4">
          <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
            Use this page to run a model.
          </div>
          <div className="grid grid-cols-1 col-span-1 justify-items-end">
            <Button
              size="lg"
              variant="default"
              disabled={!validateModelRun()}
              onClick={() => {
                createModelRunMutation.mutate(modelRun);
                navigate(paths.app.viewer.runHistory.getHref());
              }}
            >
              Run
            </Button>
          </div>
        </div>
        <div className="grid grid-cols-1 col-span-1 mt-4">
          <ModelSelector
            selectedModelId={modelRun.modelId}
            setSelectedModelId={(action) => {
              setModelRun({
                ...modelRun,
                modelId:
                  (typeof action == "function"
                    ? action(modelRun.modelId)
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
              setModelRun({
                ...modelRun,
                name: e.target.value,
              });
            }}
          />
        </div>
        {model ? (
          <>
            <div className="grid grid-cols-1 col-span-1 mt-4">
              <FileGroupSelector
                specificationId={model.inputSpecificationId}
                selectedFileGroupId={modelRun.inputFileGroupId}
                setSelectedFileGroupId={(action) => {
                  setModelRun({
                    ...modelRun,
                    inputFileGroupId:
                      (typeof action == "function"
                        ? action(modelRun.inputFileGroupId)
                        : action) ?? "",
                  });
                }}
              />
            </div>
            {model.parametersSchema ? (
              <>
                <div className="grid grid-cols-1 col-span-1 justify-items-start content-center mt-4">
                  <Label className="font-medium text-sm">Parameters</Label>
                </div>
                <div className="grid grid-cols-1 col-span-1 border">
                  <Form
                    schema={JSON.parse(model.parametersSchema)}
                    validator={validator}
                    onChange={(e) => {
                      setModelRun({
                        ...modelRun,
                        parameters: JSON.stringify(e.formData),
                      });
                    }}
                    ref={formRef}
                    onError={() => console.log("errors")}
                    uiSchema={{
                      "ui:submitButtonOptions": { norender: true },
                    }}
                  />
                </div>
              </>
            ) : null}
          </>
        ) : null}
      </div>
    </ContentLayout>
  );
};

export default RunModelsCreateRoute;
