import { Label } from "@radix-ui/react-dropdown-menu";
import { useMatch } from "react-router-dom";
import { useQuery } from "@connectrpc/connect-query";
import { paths } from "@/config/paths";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { ArrowBigLeft, CircleArrowLeft } from "lucide-react";
import { Separator } from "@/components/ui/separator";
import { getSpecification } from "@/types/api/file-FileService_connectquery";
import { DataMode, TableFieldType } from "@/types/api/file_pb";
import { FileTransport } from "@/lib/api-client";

const SpecificationsViewRoute = () => {
  const navigate = useNavigate();
  const id = useMatch(
    paths.app.developer.specifications.root.getHref() + "/:id",
  )?.params.id;
  const { data: specificationData } = useQuery(
    getSpecification,
    {
      id: id || "",
    },
    {
      enabled: !!id,
      transport: FileTransport,
    },
  );

  const specification = specificationData?.specification;

  return (
    <div className="grid grid-cols-1 gap-4">
      <div className="grid grid-cols-2 col-span-1 justify-items-between">
        <div className="grid grid-cols-1 col-span-1 justify-items-start">
          <Label className="col-span-1 content-center text-lg font-medium">
            {specification?.name}
          </Label>
        </div>
        <div className="grid grid-cols-1 col-span-1 justify-items-end">
          <Button
            size="lg"
            variant="outline"
            onClick={() => {
              navigate(paths.app.developer.specifications.list.getHref());
            }}
          >
            <CircleArrowLeft />
            Back
          </Button>
        </div>
      </div>
      {specification ? (
        <>
          <Label className="col-span-1 content-center mt-4 text-md font-normal">
            Data tables -{" "}
            {specificationData?.mode == DataMode.INPUT ? "CSV" : "Parquet"}
          </Label>
          {specification.tables.map((table, index) => (
            <div
              key={index}
              className="grid grid-cols-1 col-span-1 border p-4 gap-4"
            >
              <Label className="col-span-1 content-center font-normal">
                {table.name}
              </Label>
              <div className="grid grid-cols-1 col-span-1 gap-2">
                {table.fields.map((field, index) => (
                  <div
                    key={index}
                    className="grid grid-cols-2 col-span-1 justify-items-between"
                  >
                    <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
                      <Label className="text-sm font-light">{field.name}</Label>
                    </div>
                    <div className="grid grid-cols-1 col-span-1 justify-items-end content-center">
                      <Label className="text-sm font-medium">
                        {TableFieldType[field.type]}
                      </Label>
                    </div>
                    <div className="col-span-2">
                      <Separator className="mt-2" />
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </>
      ) : null}
    </div>
  );
};

export default SpecificationsViewRoute;
