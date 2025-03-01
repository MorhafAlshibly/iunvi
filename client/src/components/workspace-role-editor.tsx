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
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";

export function WorkspaceRoleEditor({
  user,
  workspace,
  currentRole,
}: {
  user: User;
  workspace: Workspace | null;
  currentRole: WorkspaceRole;
}) {
  const [selectEnabled, setSelectEnabled] = React.useState(false);
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
    setSelectEnabled(false);
  };

  const workspaceRoleList = Object.keys(WorkspaceRole).filter(
    (key) => !isNaN(Number(WorkspaceRole[key as any])),
  );

  return (
    <>
      <Select
        onValueChange={(value) => {
          setSelectedWorkspaceRole(
            WorkspaceRole[value as keyof typeof WorkspaceRole],
          );
        }}
        defaultValue={
          selectedWorkspaceRole
            ? WorkspaceRole[selectedWorkspaceRole]
            : WorkspaceRole[WorkspaceRole.UNASSIGNED]
        }
        disabled={!selectEnabled}
      >
        <SelectTrigger>
          <SelectValue placeholder="Select a role" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            {workspaceRoleList.map((workspaceRole) => (
              <SelectItem key={workspaceRole} value={workspaceRole}>
                {workspaceRole}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
      {selectEnabled ? (
        <Button variant="ghost" size="sm" onClick={handleAssignUserToWorkspace}>
          <Check />
        </Button>
      ) : (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setSelectEnabled(true)}
        >
          <Edit />
        </Button>
      )}
    </>
  );
}
