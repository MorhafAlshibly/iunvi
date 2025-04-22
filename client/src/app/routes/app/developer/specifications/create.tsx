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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { createSpecification } from "@/types/api/file-FileService_connectquery";
import {
  CreateSpecificationRequest,
  DataMode,
  TableFieldType,
} from "@/types/api/file_pb";

const SpecificationsCreateRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();

  const createSpecificationMutation = useMutation(createSpecification);
  const [specification, setSpecification] =
    useState<CreateSpecificationRequest>({
      $typeName: "file.CreateSpecificationRequest",
      workspaceId: activeWorkspace?.id || "",
      name: "",
      mode: DataMode.INPUT,
      tables: [
        {
          $typeName: "file.TableSchema",
          name: "",
          fields: [
            {
              $typeName: "file.TableField",
              name: "",
              type: TableFieldType.BIGINT,
            },
          ],
        },
      ],
    });

  const validateSpecification = () => {
    // check if tables are valid JSON
    try {
      specification.tables.forEach((table) => {
        // check if name is not empty
        if (!table.name) {
          throw new Error("Name is required");
        }
        table.fields.forEach((field) => {
          // check if name is not empty
          if (!field.name) {
            throw new Error("Name is required");
          }
        });
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
        Use this page to create specifications using DuckDB types.
      </span>
      <div className="grid grid-cols-2 col-span-1 gap-4 justify-items-between mb-4">
        <div className="col-span-1">
          <Select
            defaultValue={DataMode[DataMode.INPUT]}
            onValueChange={(value) => {
              setSpecification({
                ...specification,
                mode: DataMode[value as keyof typeof DataMode],
              });
            }}
          >
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Mode" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value={DataMode[DataMode.INPUT]}>Input</SelectItem>
              <SelectItem value={DataMode[DataMode.OUTPUT]}>Output</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="grid grid-cols-1 col-span-1 justify-items-end">
          <Button
            size="lg"
            variant="default"
            disabled={!validateSpecification()}
            onClick={() => {
              createSpecificationMutation.mutate(specification);
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
            ...specification,
            name: e.target.value,
          });
        }}
      />
      <Label className="col-span-1 content-center mt-4 text-lg">
        Data tables - {specification.mode == DataMode.INPUT ? "CSV" : "Parquet"}
      </Label>
      <CreateDataTable
        dataTables={specification.tables}
        setDataTables={(action) => {
          setSpecification({
            ...specification,
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
