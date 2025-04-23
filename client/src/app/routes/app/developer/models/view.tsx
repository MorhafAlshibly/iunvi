import { Label } from "@radix-ui/react-dropdown-menu";
import { useMatch } from "react-router-dom";
import { useQuery } from "@connectrpc/connect-query";
import { paths } from "@/config/paths";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { ArrowBigLeft, CircleArrowLeft, Info } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { getModel } from "@/types/api/model-ModelService_connectquery";
import { TableFieldType } from "@/types/api/file_pb";
import { FileTransport, ModelTransport } from "@/lib/api-client";
import { getSpecification } from "@/types/api/file-FileService_connectquery";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { useDarkMode } from "usehooks-ts";

const ModelsViewRoute = () => {
  const navigate = useNavigate();
  const id = useMatch(paths.app.developer.models.root.getHref() + "/:id")
    ?.params.id;
  const { data: modelData } = useQuery(
    getModel,
    {
      id: id || "",
    },
    {
      enabled: !!id,
      transport: ModelTransport,
    },
  );

  const model = modelData?.model;

  const { data: inputSpecificationData } = useQuery(
    getSpecification,
    {
      id: modelData?.model?.inputSpecificationId,
    },
    {
      enabled: !!modelData?.model?.inputSpecificationId,
      transport: FileTransport,
    },
  );

  const inputSpecification = inputSpecificationData?.specification;

  const { data: outputSpecificationData } = useQuery(
    getSpecification,
    {
      id: modelData?.model?.outputSpecificationId,
    },
    {
      enabled: !!modelData?.model?.outputSpecificationId,
      transport: FileTransport,
    },
  );

  const outputSpecification = outputSpecificationData?.specification;

  const darkMode = useDarkMode();

  return (
    <div className="grid grid-cols-1 gap-4">
      {model && inputSpecification && outputSpecification ? (
        <>
          <div className="grid grid-cols-2 col-span-1 justify-items-between">
            <div className="grid grid-cols-1 col-span-1 justify-items-start">
              <Label className="col-span-1 content-center text-lg font-medium">
                {model.name}
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
              <Label className="font-normal">Image</Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end content-center">
              <Label className="font-medium">{model.imageName}</Label>
            </div>
            <div className="grid grid-cols-1 col-span-2">
              <Separator className="mt-2" />
            </div>
          </div>
          <div className="grid grid-cols-2 col-span-1 justify-items-between p-2">
            <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
              <Label className="font-normal">Input specification</Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end content-center">
              <Button
                size="sm"
                variant="outline"
                onClick={() => {
                  navigate(
                    paths.app.developer.specifications.view.getHref(
                      inputSpecification.id,
                    ),
                  );
                }}
              >
                <Info />
                {inputSpecification.name}
              </Button>
            </div>
            <div className="grid grid-cols-1 col-span-2">
              <Separator className="mt-2" />
            </div>
          </div>
          <div className="grid grid-cols-2 col-span-1 justify-items-between p-2">
            <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
              <Label className="font-normal">Output specification</Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end content-center">
              <Button
                size="sm"
                variant="outline"
                onClick={() => {
                  navigate(
                    paths.app.developer.specifications.view.getHref(
                      outputSpecification.id,
                    ),
                  );
                }}
              >
                <Info />
                {outputSpecification.name}
              </Button>
            </div>
            <div className="grid grid-cols-1 col-span-2">
              <Separator className="mt-2" />
            </div>
          </div>
          {model.parametersSchema ? (
            <>
              <div className="grid grid-cols-1 col-span-1 content-center mt-4">
                <Label className="font-medium">Parameters Schema</Label>
              </div>
              <div className="grid grid-cols-1 col-span-1">
                <CodeMirror
                  value={model.parametersSchema}
                  height="auto"
                  extensions={[json()]}
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

export default ModelsViewRoute;
