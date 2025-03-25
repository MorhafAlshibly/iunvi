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
import { DataMode, ModelName } from "@/types/api/tenantManagement_pb";
import { useQuery } from "@connectrpc/connect-query";
import { getModels } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useWorkspace } from "@/hooks/use-workspace";

export function ModelSelector({
  selectedModelId,
  setSelectedModelId,
}: {
  selectedModelId: string | null;
  setSelectedModelId: React.Dispatch<React.SetStateAction<string | null>>;
}) {
  const { activeWorkspace } = useWorkspace();

  const { data: modelsData } = useQuery(
    getModels,
    {
      workspaceId: activeWorkspace?.id ?? "",
    },
    {
      enabled: !!activeWorkspace,
    },
  );

  const models = modelsData?.models || [];

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
          {selectedModelId
            ? models.find((model) => model.id === selectedModelId)?.name
            : `Select model...`}
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput placeholder="Search model..." className="h-9" />
          <CommandList>
            <CommandEmpty>No model found.</CommandEmpty>
            <CommandGroup>
              {models.map((model) => (
                <CommandItem
                  key={model.id}
                  value={model.name}
                  onSelect={() => {
                    setSelectedModelId(model.id);
                    setOpen(false);
                  }}
                >
                  {model.name}
                  <Check
                    className={cn(
                      "ml-auto",
                      selectedModelId === model.id
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
