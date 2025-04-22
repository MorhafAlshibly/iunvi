import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Check, ChevronsUpDown, Command, Edit, Info, Plus } from "lucide-react";
import { useEffect, useState } from "react";
import { WorkspaceSelector } from "./workspace-selector";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { WorkspaceRoleEditor } from "./workspace-role-editor";
import { getUserWorkspaceAssignment } from "@/types/api/tenant-TenantService_connectquery";
import { User, Workspace, WorkspaceRole } from "@/types/api/tenant_pb";

export function UserWorkspacePanel({ user }: { user: User }) {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedWorkspace, setSelectedWorkspace] = useState<Workspace | null>(
    null,
  );

  const { data, refetch } = useQuery(
    getUserWorkspaceAssignment,
    { userObjectId: user.id, workspaceId: selectedWorkspace?.id },
    { enabled: selectedWorkspace !== null },
  );

  useEffect(() => {
    refetch();
  }, [selectedWorkspace]);

  const workspaceRole = data?.role || WorkspaceRole.UNASSIGNED;

  return (
    <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost" size="sm">
          <Info />
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{user.displayName}</DialogTitle>
          <DialogDescription></DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <div className="col-span-4">
              <WorkspaceSelector
                selectedWorkspace={selectedWorkspace}
                setSelectedWorkspace={setSelectedWorkspace}
              />
            </div>
            <div className="grid grid-cols-4 col-span-4 p-2">
              <Label className="col-span-2 flex items-center">
                Workspace role
              </Label>
              <div className="col-span-2 flex gap-2 justify-end">
                <WorkspaceRoleEditor
                  user={user}
                  workspace={selectedWorkspace}
                  currentRole={workspaceRole}
                />
              </div>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
