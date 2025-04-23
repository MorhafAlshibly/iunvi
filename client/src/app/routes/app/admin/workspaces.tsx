import { ContentLayout } from "@/components/layouts/content";
import { CreateWorkspace } from "@/components/create-workspace";
import { useMutation } from "@connectrpc/connect-query";
import { useWorkspace } from "@/hooks/use-workspace";
import { Separator } from "@/components/ui/separator";
import { EditWorkspace } from "@/components/edit-workspace";
import {
  createWorkspace,
  editWorkspace,
} from "@/types/api/tenant-TenantService_connectquery";
import { TenantTransport } from "@/lib/api-client";

const WorkspacesRoute = () => {
  const { workspaces, workspacesRefetch } = useWorkspace();
  const createWorkspaceHandler = useMutation(createWorkspace, {
    transport: TenantTransport,
  });
  const editWorkspaceHandler = useMutation(editWorkspace, {
    transport: TenantTransport,
  });

  const createWorkspaceFn = async (workspaceName: string) => {
    await createWorkspaceHandler.mutateAsync({
      name: workspaceName,
    });
    workspacesRefetch();
  };

  const editWorkspaceFn = async (id: string, name: string) => {
    await editWorkspaceHandler.mutateAsync({
      id: id,
      name: name,
    });
    workspacesRefetch();
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
              <div key={workspace.id} className="flex text-sm">
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
