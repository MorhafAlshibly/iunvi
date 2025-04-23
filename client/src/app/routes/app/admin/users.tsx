import { UserWorkspacePanel } from "@/components/user-workspace-panel";
import { ContentLayout } from "@/components/layouts/content";
import { Separator } from "@/components/ui/separator";
import { useQuery } from "@connectrpc/connect-query";
import { getUsers } from "@/types/api/tenant-TenantService_connectquery";
import { TenantTransport } from "@/lib/api-client";
import { Label } from "@/components/ui/label";

const UsersRoute = () => {
  const { data } = useQuery(getUsers, undefined, {
    transport: TenantTransport,
  });
  const users = data?.users || [];
  return (
    <ContentLayout title="Users">
      <div className="grid grid-cols-1 col-span-1 gap-4">
        {users.map((user) => (
          <div key={user.id} className="grid grid-cols-2 col-span-1 mt-2">
            <div className="grid grid-cols-1 content-center col-span-1 justify-items-start">
              <Label className="font-normal">{user.displayName}</Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end">
              <UserWorkspacePanel user={user} />
            </div>
            <div className="grid grid-cols-1 col-span-2">
              <Separator className="mt-2" />
            </div>
          </div>
        ))}
      </div>
    </ContentLayout>
  );
};

export default UsersRoute;
