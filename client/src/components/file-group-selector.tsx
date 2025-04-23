import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";

import { cn } from "@/utils/cn";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useQuery } from "@connectrpc/connect-query";
import { useWorkspace } from "@/hooks/use-workspace";
import { getFileGroups } from "@/types/api/file-FileService_connectquery";
import { FileTransport } from "@/lib/api-client";

export function FileGroupSelector({
  specificationId,
  selectedFileGroupId,
  setSelectedFileGroupId,
}: {
  specificationId: string;
  selectedFileGroupId: string | null;
  setSelectedFileGroupId: React.Dispatch<React.SetStateAction<string | null>>;
}) {
  const { activeWorkspace } = useWorkspace();

  const { data: fileGroupsData } = useQuery(
    getFileGroups,
    {
      workspaceId: activeWorkspace?.id ?? "",
      specificationId,
    },
    {
      enabled: !!activeWorkspace,
      transport: FileTransport,
    },
  );

  const fileGroups = fileGroupsData?.fileGroups || [];

  const [open, setOpen] = React.useState(false);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
        >
          {selectedFileGroupId
            ? fileGroups.find(
                (fileGroup) => fileGroup.id === selectedFileGroupId,
              )?.name
            : `Select file group...`}
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput placeholder="Search file group..." className="h-9" />
          <CommandList>
            <CommandEmpty>No file group found.</CommandEmpty>
            <CommandGroup>
              {fileGroups.map((fileGroup) => (
                <CommandItem
                  key={fileGroup.id}
                  value={fileGroup.name}
                  onSelect={() => {
                    setSelectedFileGroupId(fileGroup.id);
                    setOpen(false);
                  }}
                >
                  {fileGroup.name}
                  <Check
                    className={cn(
                      "ml-auto",
                      selectedFileGroupId === fileGroup.id
                        ? "opacity-100"
                        : "opacity-0",
                    )}
                  />
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
