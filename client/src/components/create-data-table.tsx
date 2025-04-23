import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { ChevronDown, CircleX, Plus, PlusCircle, X } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";
import { TableFieldType, TableSchema } from "@/types/api/file_pb";

const CreateDataTable = ({
  dataTables,
  setDataTables,
}: {
  dataTables: TableSchema[];
  setDataTables: React.Dispatch<React.SetStateAction<TableSchema[]>>;
}) => {
  const tableFieldTypeList = Object.keys(TableFieldType).filter(
    (key) => !isNaN(Number(TableFieldType[key as any])),
  );

  return (
    <div className="grid grid-cols-1 gap-4">
      {dataTables.map((table, index) => (
        <div
          key={index}
          className="grid grid-cols-1 col-span-1 border p-4 gap-4"
        >
          <div className="grid grid-cols-2 col-span-1 justify-items-between">
            <Label className="col-span-1 content-center">
              Table {index + 1}
            </Label>
            <div className="grid grid-cols-1 col-span-1 justify-items-end">
              <Button
                size="sm"
                variant="ghost"
                className="col-span-1"
                disabled={dataTables.length == 1}
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
            </div>
          </div>
          <Input
            placeholder="Name"
            className="col-span-1"
            value={table.name}
            onChange={(e) => {
              setDataTables((prev) => {
                const newTables = [...prev];
                newTables[index].name = e.target.value;
                return newTables;
              });
            }}
          />
          <div className="grid grid-cols-1 col-span-1 gap-4">
            {table.fields.map((field, fieldIndex) => (
              <div className="grid grid-cols-2 col-span-1 justify-items-between">
                <div className="grid grid-cols-1 col-span-1 justify-items-start">
                  <Input
                    placeholder={`Field ${fieldIndex + 1}`}
                    className="col-span-1 w-[400px]"
                    value={field.name}
                    onChange={(e) => {
                      setDataTables((prev) => {
                        const newFields = [...prev];
                        newFields[index].fields[fieldIndex].name =
                          e.target.value;
                        return newFields;
                      });
                    }}
                  />
                </div>
                <div className="grid grid-flow-col justify-self-end gap-2">
                  <Select
                    onValueChange={(value) => {
                      setDataTables((prev) => {
                        const newFields = [...prev];
                        newFields[index].fields[fieldIndex].type =
                          TableFieldType[value as keyof typeof TableFieldType];
                        return newFields;
                      });
                    }}
                    defaultValue={TableFieldType[field.type]}
                  >
                    <SelectTrigger className="w-[180px]">
                      <SelectValue placeholder="Type" />
                    </SelectTrigger>
                    <SelectContent>
                      {tableFieldTypeList.map((key) => (
                        <SelectItem key={key} value={key}>
                          {key}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>

                  <Button
                    size="icon"
                    variant="ghost"
                    disabled={table.fields.length == 1}
                    onClick={() => {
                      setDataTables((prev) => {
                        const newFields = [...prev];
                        newFields[index].fields.splice(fieldIndex, 1);
                        return newFields;
                      });
                    }}
                  >
                    <X />
                  </Button>
                </div>
              </div>
            ))}
            <div className="grid grid-cols-1 col-span-1 justify-items-start">
              <Button
                size="icon"
                variant="ghost"
                className="col-span-1"
                onClick={() => {
                  setDataTables((prev) => {
                    const newFields = [...prev];
                    newFields[index].fields.push({
                      $typeName: "file.TableField",
                      name: "",
                      type: TableFieldType.BIGINT,
                    });
                    return newFields;
                  });
                }}
              >
                <ChevronDown />
              </Button>
            </div>
          </div>
          {/* <CodeMirror
            value={table.schema}
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
          /> */}
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
