import { Label } from "@radix-ui/react-dropdown-menu";
import { useMatch } from "react-router-dom";
import { useQuery } from "@connectrpc/connect-query";
import { paths } from "@/config/paths";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { ArrowBigLeft, CircleArrowLeft } from "lucide-react";
import { getSpecification } from "@/types/api/file-FileService_connectquery";
import { DataMode } from "@/types/api/file_pb";
import { FileTransport } from "@/lib/api-client";

const FileGroupsViewRoute = () => {
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
          <Label className="col-span-1 content-center text-lg">
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
      <Label className="col-span-1 content-center mt-4 text-lg">
        Data tables -{" "}
        {specificationData?.mode == DataMode.INPUT ? "CSV" : "Parquet"}
      </Label>
      {specification?.tables.map((table, index) => (
        <div
          key={index}
          className="grid grid-cols-1 col-span-1 border p-4 gap-4"
        >
          <Label className="col-span-1 content-center">{table.name}</Label>
          <CodeMirror
            value={table.schema}
            height="auto"
            extensions={[json()]}
            editable={false}
            className="col-span-1 border"
          />
        </div>
      ))}
    </div>
  );
};

export default FileGroupsViewRoute;
