"use client";

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
import { useWorkspace } from "@/hooks/use-workspace";
import { Workspace } from "@/types/api/tenantManagement_pb";

export function WorkspaceSelector({
  selectedWorspace,
  setSelectedWorkspace,
}: {
  selectedWorspace: Workspace | null;
  setSelectedWorkspace: React.Dispatch<React.SetStateAction<Workspace | null>>;
}) {
  const { workspaces } = useWorkspace();
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
          {selectedWorspace
            ? workspaces.find(
                (workspace) => workspace.id === selectedWorspace?.id,
              )?.name
            : "Select workspace..."}
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput placeholder="Search workspace..." className="h-9" />
          <CommandList>
            <CommandEmpty>No workspace found.</CommandEmpty>
            <CommandGroup>
              {workspaces.map((workspace) => (
                <CommandItem
                  key={workspace.id}
                  value={workspace.name}
                  onSelect={() => {
                    setSelectedWorkspace(workspace);
                    setOpen(false);
                  }}
                >
                  {workspace.name}
                  <Check
                    className={cn(
                      "ml-auto",
                      selectedWorspace?.id === workspace.id
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
