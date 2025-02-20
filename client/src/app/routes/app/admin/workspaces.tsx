import { ContentLayout } from "@/components/layouts/content";
import { CreateWorkspace } from "@/components/create-workspace";
import { useUser } from "@/lib/authentication";
import {
  createWorkspace,
  editWorkspace,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { ROLES } from "@/types/user";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { useWorkspace } from "@/hooks/use-workspace";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { Edit } from "lucide-react";
import { EditWorkspace } from "@/components/edit-workspace";

const WorkspacesRoute = () => {
  const { workspaces, refetch } = useWorkspace();
  const createWorkspaceHandler = useMutation(createWorkspace);
  const editWorkspaceHandler = useMutation(editWorkspace);

  const createWorkspaceFn = async (workspaceName: string) => {
    await createWorkspaceHandler.mutateAsync({
      name: workspaceName,
    });
    refetch();
  };

  const editWorkspaceFn = async (
    id: Uint8Array<ArrayBufferLike>,
    name: string,
  ) => {
    await editWorkspaceHandler.mutateAsync({
      id: id,
      name: name,
    });
    refetch();
  };

  return (
    <ContentLayout title="Workspaces">
      <div className="flex justify-end">
        <CreateWorkspace onSubmit={createWorkspaceFn} />
      </div>
      <div>
        <div className="p-4">
          {workspaces.map((workspace) => (
            <>
              <div key={workspace.id.toString()} className="flex text-sm">
                <span className="flex-1 content-center">{workspace.name}</span>
                <span className="flex-1 text-right">
                  <EditWorkspace
                    onSubmit={(name) => editWorkspaceFn(workspace.id, name)}
                  />
                </span>
              </div>
              <Separator className="my-2" />
            </>
          ))}
        </div>
      </div>
    </ContentLayout>
  );
};

export default WorkspacesRoute;
