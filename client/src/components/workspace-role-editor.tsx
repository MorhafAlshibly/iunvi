"use client";

import * as React from "react";
import { Check, ChevronsUpDown, Edit } from "lucide-react";

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
import {
  User,
  Workspace,
  WorkspaceRole,
} from "@/types/api/tenantManagement_pb";
import { assignUserToWorkspace } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useMutation } from "@connectrpc/connect-query";

export function WorkspaceRoleEditor({
  user,
  workspace,
  currentRole,
}: {
  user: User;
  workspace: Workspace | null;
  currentRole: WorkspaceRole;
}) {
  const [open, setOpen] = React.useState(false);
  const [comboboxEnabled, setComboboxEnabled] = React.useState(false);
  const [selectedWorkspaceRole, setSelectedWorkspaceRole] =
    React.useState(currentRole);

  React.useEffect(() => {
    setSelectedWorkspaceRole(currentRole);
  }, [currentRole]);

  const assignUserToWorkspaceMutation = useMutation(assignUserToWorkspace);

  const handleAssignUserToWorkspace = async () => {
    if (workspace && selectedWorkspaceRole !== currentRole) {
      await assignUserToWorkspaceMutation.mutateAsync({
        userObjectId: user.id,
        workspaceId: workspace.id,
        role: selectedWorkspaceRole,
      });
    }
    setComboboxEnabled(false);
  };

  const workspaceRoleList = Object.keys(WorkspaceRole).filter(
    (key) => !isNaN(Number(WorkspaceRole[key as any])),
  );

  return (
    <>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between"
            disabled={!comboboxEnabled}
          >
            {selectedWorkspaceRole
              ? WorkspaceRole[selectedWorkspaceRole]
              : WorkspaceRole[WorkspaceRole.UNASSIGNED]}
            <ChevronsUpDown className="opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="p-0">
          <Command>
            <CommandInput placeholder="Search role..." className="h-9" />
            <CommandList>
              <CommandEmpty>No role found.</CommandEmpty>
              <CommandGroup>
                {workspaceRoleList.map((workspaceRole) => (
                  <CommandItem
                    key={workspaceRole}
                    value={workspaceRole}
                    onSelect={() => {
                      setSelectedWorkspaceRole(
                        WorkspaceRole[
                          workspaceRole as keyof typeof WorkspaceRole
                        ],
                      );
                      setOpen(false);
                    }}
                  >
                    {workspaceRole}
                    <Check
                      className={cn(
                        "ml-auto",
                        WorkspaceRole[selectedWorkspaceRole] === workspaceRole
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
      {comboboxEnabled ? (
        <Button variant="ghost" size="sm" onClick={handleAssignUserToWorkspace}>
          <Check />
        </Button>
      ) : (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setComboboxEnabled(true)}
        >
          <Edit />
        </Button>
      )}
    </>
  );
}
