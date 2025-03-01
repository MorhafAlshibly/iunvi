import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";
import { createInputSpecification } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { FileSchema } from "@/types/api/tenantManagement_pb";
import { useMutation } from "@connectrpc/connect-query";
import { useState } from "react";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import { CircleX, PlusCircle } from "lucide-react";

const CreateDataTable = ({
  dataTables,
  setDataTables,
}: {
  dataTables: FileSchema[];
  setDataTables: React.Dispatch<React.SetStateAction<FileSchema[]>>;
}) => {
  return (
    <div className="grid grid-cols-1 gap-4">
      {dataTables.map((csv, index) => (
        <div
          key={index}
          className="grid grid-cols-1 col-span-1 border p-4 gap-4"
        >
          <div className="grid grid-cols-2 col-span-1 justify-items-between">
            <Label className="col-span-1 content-center">
              Table {index + 1}
            </Label>
            <div className="grid grid-cols-1 col-span-1 justify-items-end">
              {index > 0 ? (
                <Button
                  size="sm"
                  variant="ghost"
                  className="col-span-1"
                  onClick={() => {
                    setDataTables((prev) => {
                      const newTables = [...prev];
                      newTables.splice(index, 1);
                      return newTables;
                    });
                  }}
                >
                  <CircleX />
                </Button>
              ) : null}
            </div>
          </div>
          <Input
            placeholder="Name"
            className="col-span-1"
            value={csv.name}
            onChange={(e) => {
              setDataTables((prev) => {
                const newTables = [...prev];
                newTables[index].name = e.target.value;
                return newTables;
              });
            }}
          />
          <CodeMirror
            value={csv.schema}
            height="auto"
            placeholder={"Schema"}
            extensions={[json()]}
            onChange={(value) => {
              setDataTables((prev) => {
                const newTables = [...prev];
                newTables[index].schema = value;
                return newTables;
              });
            }}
            className="col-span-1 border"
          />
        </div>
      ))}
      <div className="grid grid-cols-1 col-span-1 justify-items-center">
        <Button
          size="lg"
          variant="outline"
          className="col-span-1"
          onClick={() => {
            setDataTables((prev) => [
              ...prev,
              {
                $typeName: "api.FileSchema",
                name: "",
                schema: "",
              },
            ]);
          }}
        >
          <PlusCircle />
        </Button>
      </div>
    </div>
  );
};

export default CreateDataTable;
