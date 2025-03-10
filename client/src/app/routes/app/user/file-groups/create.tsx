import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";
import {
  createSpecification,
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
    },
    {
      enabled: !!activeWorkspace,
    },
  );

  return (
    <div className="grid grid-cols-1">
      <div className="grid grid-cols-2 col-span-1 gap-4 justify-items-between">
        <div className="grid grid-cols-1 col-span-2">
          <SpecificationSelector
            specifications={specifications?.specifications || []}
            selectedSpecification={selectedSpecification}
            setSelectedSpecification={setSelectedSpecification}
          />
        </div>
      </div>
    </div>
  );
};

export default FileGroupsCreateRoute;
