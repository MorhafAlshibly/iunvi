import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";
import {
  createInputSpecification,
  createOutputSpecification,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import {
  CreateInputSpecificationRequest,
  CreateInputSpecificationRequestSchema,
  CreateOutputSpecificationRequest,
  FileSchema,
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

const SpecificationsCreateRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();
  const [dataMode, setDataMode] = useState<"Input" | "Output">("Input");

  const createInputSpecificationMutation = useMutation(
    createInputSpecification,
  );
  const createOutputSpecificationMutation = useMutation(
    createOutputSpecification,
  );

  const createSpecificationMutation =
    dataMode === "Input"
      ? createInputSpecificationMutation
      : createOutputSpecificationMutation;

  const [inputSpecification, setInputSpecification] =
    useState<CreateInputSpecificationRequest>({
      $typeName: "api.CreateInputSpecificationRequest",
      workspaceId: activeWorkspace?.id || "",
      name: "",
      parametersSchema: "",
      tables: [
        {
          $typeName: "api.FileSchema",
          name: "",
          schema: "",
        },
      ],
    });

  const [outputSpecification, setOutputSpecification] =
    useState<CreateOutputSpecificationRequest>({
      $typeName: "api.CreateOutputSpecificationRequest",
      workspaceId: activeWorkspace?.id || "",
      name: "",
      tables: [
        {
          $typeName: "api.FileSchema",
          name: "",
          schema: "",
        },
      ],
    });

  const specification =
    dataMode === "Input" ? inputSpecification : outputSpecification;
  const setSpecification =
    dataMode === "Input" ? setInputSpecification : setOutputSpecification;

  const validateSpecification = () => {
    if (dataMode === "Input") {
      // check if paramertersSchema is valid JSON
      try {
        JSON.parse(inputSpecification.parametersSchema);
      } catch (e) {
        return false;
      }
    }
    // check if tables are valid JSON
    try {
      specification.tables.forEach((table) => {
        // check if name is not empty
        if (!table.name) {
          throw new Error("Name is required");
        }
        JSON.parse(table.schema);
      });
    } catch (e) {
      return false;
    }
    // check if name is not empty
    if (!specification.name) {
      return false;
    }
    return true;
  };

  return (
    <div className="grid grid-cols-1 gap-4">
      <span className="col-span-1 text-md mb-4">
        Use this page to create specifications using JSON Schema 2020-12 for the
        Parameters and the Frictionless Table Schema for the Data tables.
      </span>
      <div className="grid grid-cols-2 col-span-1 gap-4 justify-items-between mb-4">
        <div className="col-span-1">
          <Select
            defaultValue={dataMode}
            onValueChange={(value) => setDataMode(value as any)}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Mode" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="Input">Input</SelectItem>
              <SelectItem value="Output">Output</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="grid grid-cols-1 col-span-1 justify-items-end">
          <Button
            size="lg"
            variant="default"
            disabled={!validateSpecification()}
            onClick={() => {
              createSpecificationMutation.mutate(specification as any);
              navigate(paths.app.developer.specifications.list.getHref());
            }}
          >
            Create
          </Button>
        </div>
      </div>
      <Input
        placeholder="Name"
        className="col-span-1"
        onChange={(e) => {
          setSpecification({
            ...(specification as any),
            name: e.target.value,
          });
        }}
      />
      {specification == inputSpecification ? (
        <CodeMirror
          value={inputSpecification.parametersSchema}
          height="auto"
          extensions={[json()]}
          placeholder={"Parameters"}
          onChange={(value) => {
            setInputSpecification({
              ...inputSpecification,
              parametersSchema: value,
            });
          }}
          className="col-span-1 border"
        />
      ) : null}
      <Label className="col-span-1 content-center mt-4 text-lg">
        Data tables - {specification == inputSpecification ? "CSV" : "Parquet"}
      </Label>
      <CreateDataTable
        dataTables={specification.tables}
        setDataTables={(action) => {
          setSpecification({
            ...(specification as any),
            tables: Array.isArray(action)
              ? action
              : action(specification.tables),
          });
        }}
      />
    </div>
  );
};

export default SpecificationsCreateRoute;
