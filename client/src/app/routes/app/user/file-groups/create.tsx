import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useWorkspace } from "@/hooks/use-workspace";
import {
  createFileGroup,
  createSpecification,
  getLandingZoneFiles,
  getSpecification,
  getSpecifications,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import {
  CreateFileGroupRequest,
  CreateSpecificationRequest,
  CreateSpecificationRequestSchema,
  DataMode,
  LandingZoneFile,
  SpecificationName,
  TableSchema,
} from "@/types/api/tenantManagement_pb";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
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
  Plus,
  PlusCircle,
} from "lucide-react";
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
import { FileSelector } from "@/components/file-selector";

const FileGroupsCreateRoute = () => {
  const navigate = useNavigate();
  const { activeWorkspace } = useWorkspace();
  const [createFileGroupInput, setCreateFileGroupInput] =
    useState<CreateFileGroupRequest>({
      $typeName: "api.CreateFileGroupRequest",
      specificationId: "",
      name: "",
      schemaFileMappings: [],
    });

  const { data: specificationData } = useQuery(
    getSpecification,
    {
      id: createFileGroupInput.specificationId,
    },
    {
      enabled: !!createFileGroupInput.specificationId,
    },
  );

  const [searchFilter, setSearchFilter] = useState<string>("");
  const {
    data: landingZoneFilesData,
    refetch: landingZoneFilesRefetch,
    hasNextPage: landingZoneFilesHasNextPage,
    fetchNextPage: landingZoneFilesFetchNextPage,
  } = useInfiniteQuery(
    getLandingZoneFiles,
    {
      workspaceId: activeWorkspace?.id ?? "",
      marker: "",
      prefix: searchFilter,
    },
    {
      pageParamKey: "marker",
      getNextPageParam: (lastPage, allPages, lastPageParam, allPageParams) =>
        lastPage.nextMarker,
      enabled: !!activeWorkspace?.id,
    },
  );

  const validateFileGroups = () => {
    if (!createFileGroupInput.specificationId) return false;
    if (!createFileGroupInput.name) return false;
    if (
      createFileGroupInput.schemaFileMappings.length !==
      specificationData?.specification?.tables.length
    )
      return false;
    try {
      createFileGroupInput.schemaFileMappings.forEach((mapping) => {
        if (!mapping.landingZoneFileName) throw new Error("File is required");
      });
    } catch (e) {
      return false;
    }
    return true;
  };

  const handleFileAssignment = (
    action: React.SetStateAction<string | null>,
    table: TableSchema,
  ) => {
    setCreateFileGroupInput((prev) => {
      let newSelectedFileName = null;
      if (typeof action == "function") {
        newSelectedFileName = action(
          prev.schemaFileMappings.find(
            (mapping) => mapping.schemaName === table.name,
          )?.landingZoneFileName ?? null,
        );
      } else {
        newSelectedFileName = action;
      }
      const mappingIndex = prev.schemaFileMappings.findIndex(
        (mapping) => mapping.schemaName === table.name,
      );
      // If the file is already mapped, update the mapping, otherwise add a new mapping
      const newMappings = [...prev.schemaFileMappings];
      if (mappingIndex !== -1) {
        newMappings[mappingIndex].landingZoneFileName =
          newSelectedFileName ?? "";
      } else {
        newMappings.push({
          $typeName: "api.SchemaFileMapping",
          schemaName: table.name,
          landingZoneFileName: newSelectedFileName ?? "",
        });
      }
      return {
        ...prev,
        schemaFileMappings: newMappings,
      };
    });
  };

  const createFileGroupMutation = useMutation(createFileGroup);

  const handleCreateFileGroup = async () => {
    if (!validateFileGroups()) {
      return;
    }
    await createFileGroupMutation.mutateAsync(createFileGroupInput);
    navigate(paths.app.user.fileGroups.list.getHref());
  };

  return (
    <div className="grid grid-cols-1">
      <div className="grid grid-cols-1 col-span-1 gap-4 mt-4 justify-items-between">
        <div className="grid grid-cols-1 col-span-1">
          <SpecificationSelector
            mode={DataMode.INPUT}
            selectedSpecificationId={createFileGroupInput.specificationId}
            setSelectedSpecificationId={(action) => {
              setCreateFileGroupInput((prev) => ({
                ...prev,
                specificationId:
                  (typeof action == "function"
                    ? action(prev.specificationId)
                    : action) ?? "",
              }));
            }}
          />
        </div>
        <div className="grid grid-cols-1 col-span-1 gap-4">
          <Input
            placeholder="Name"
            value={createFileGroupInput.name}
            onChange={(e) =>
              setCreateFileGroupInput((prev) => ({
                ...prev,
                name: e.target.value,
              }))
            }
          />
        </div>
        {specificationData?.specification ? (
          <div className="grid grid-cols-1 col-span-1 gap-4 mt-2">
            {specificationData.specification.tables.map((table) => (
              <div
                className="grid grid-cols-2 gap-4 border p-4"
                key={table.name}
              >
                <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
                  <Label className="text-sm font-normal">{table.name}</Label>
                </div>
                <div className="grid grid-cols-1 col-span-1 justify-items-end">
                  <FileSelector
                    files={landingZoneFilesData?.pages[0]?.files ?? []}
                    selectedFileName={
                      createFileGroupInput.schemaFileMappings.find(
                        (mapping) => mapping.schemaName === table.name,
                      )?.landingZoneFileName ?? null
                    }
                    setSelectedFileName={(action) =>
                      handleFileAssignment(action, table)
                    }
                    searchFilter={searchFilter}
                    setSearchFilter={setSearchFilter}
                  />
                </div>
              </div>
            ))}
            <div className="grid grid-cols-1 col-span-1 mt-4 justify-items-end">
              <Button
                onClick={handleCreateFileGroup}
                disabled={!validateFileGroups()}
              >
                Create
              </Button>
            </div>
          </div>
        ) : null}
      </div>
    </div>
  );
};

export default FileGroupsCreateRoute;
