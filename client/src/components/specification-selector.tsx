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
import { getSpecifications } from "@/types/api/file-FileService_connectquery";
import { DataMode } from "@/types/api/file_pb";
import { FileTransport } from "@/lib/api-client";

export function SpecificationSelector({
  mode,
  selectedSpecificationId,
  setSelectedSpecificationId,
}: {
  mode: DataMode;
  selectedSpecificationId: string | null;
  setSelectedSpecificationId: React.Dispatch<
    React.SetStateAction<string | null>
  >;
}) {
  const { activeWorkspace } = useWorkspace();

  const { data: specificationsData } = useQuery(
    getSpecifications,
    {
      workspaceId: activeWorkspace?.id ?? "",
      mode,
    },
    {
      enabled: !!activeWorkspace,
      transport: FileTransport,
    },
  );

  const specifications = specificationsData?.specifications || [];

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
          {selectedSpecificationId
            ? specifications.find(
                (specification) => specification.id === selectedSpecificationId,
              )?.name
            : `Select ${mode == DataMode.INPUT ? "input" : "output"} specification...`}
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput placeholder="Search specification..." className="h-9" />
          <CommandList>
            <CommandEmpty>No specification found.</CommandEmpty>
            <CommandGroup>
              {specifications.map((specification) => (
                <CommandItem
                  key={specification.id}
                  value={specification.name}
                  onSelect={() => {
                    setSelectedSpecificationId(specification.id);
                    setOpen(false);
                  }}
                >
                  {specification.name}
                  <Check
                    className={cn(
                      "ml-auto",
                      selectedSpecificationId === specification.id
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
