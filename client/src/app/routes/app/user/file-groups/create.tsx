import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";
import {
  createSpecification,
  getSpecification,
  getSpecifications,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import {
  CreateSpecificationRequest,
  CreateSpecificationRequestSchema,
  DataMode,
  FileSchema,
  SpecificationName,
} from "@/types/api/tenantManagement_pb";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { useState } from "react";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import {
  Check,
  ChevronsUpDown,
  CircleX,
  Command,
  Cross,
  PlusCircle,
} from "lucide-react";
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
import { cn } from "@/utils/cn";
import {
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
} from "cmdk";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { SpecificationSelector } from "@/components/specification-selector";

const FileGroupsCreateRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();
  const [selectedSpecification, setSelectedSpecification] =
    useState<SpecificationName | null>(null);

  const { data: specifications } = useQuery(
    getSpecifications,
    {
      workspaceId: activeWorkspace?.id || "",
      mode: DataMode.INPUT,
    },
    {
      enabled: !!activeWorkspace,
    },
  );

  const { data: specificationData } = useQuery(
    getSpecification,
    {
      id: selectedSpecification?.id || "",
    },
    {
      enabled: !!selectedSpecification,
    },
  );

  return (
    <div className="grid grid-cols-1">
      <div className="grid grid-cols-1 col-span-1 gap-4 mt-4 justify-items-between">
        <div className="grid grid-cols-1 col-span-1">
          <SpecificationSelector
            specifications={specifications?.specifications || []}
            selectedSpecification={selectedSpecification}
            setSelectedSpecification={setSelectedSpecification}
          />
        </div>
        {specificationData?.specification ? (
          <div className="grid grid-cols-1 col-span-1 mt-4 gap-4">
            <div className="grid grid-cols-1 col-span-1 content-center">
              <Label className="text-md font-medium">
                {specificationData.specification.name}
              </Label>
            </div>
            {specificationData.specification.tables.map((table) => (
              <div className="grid grid-cols-1 gap-4 border p-4">
                <div className="grid grid-cols-1 col-span-1 justify-items-start">
                  <Label className="text-sm font-normal">{table.name}</Label>
                </div>{" "}
                <div className="grid grid-cols-1 col-span-1 justify-items-start">
                  <select>{table.name}</select>
                </div>
              </div>
            ))}
          </div>
        ) : null}
      </div>
    </div>
  );
};

export default FileGroupsCreateRoute;
