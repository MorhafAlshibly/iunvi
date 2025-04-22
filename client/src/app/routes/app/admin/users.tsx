import { UserWorkspacePanel } from "@/components/user-workspace-panel";
import { ContentLayout } from "@/components/layouts/content";
import { Separator } from "@/components/ui/separator";
import { useQuery } from "@connectrpc/connect-query";
import { getUsers } from "@/types/api/tenant-TenantService_connectquery";

const UsersRoute = () => {
  const { data } = useQuery(getUsers);
  const users = data?.users || [];
  return (
    <ContentLayout title="Users">
      <div>
        <div className="p-4">
          {users.map((user) => (
            <div key={user.id} className="flex text-sm border p-2 mb-4">
              <span className="flex-1 content-center">{user.displayName}</span>
              <span className="flex-1 text-right">
                <UserWorkspacePanel user={user} />
              </span>
            </div>
          ))}
        </div>
      </div>
    </ContentLayout>
  );
};

export default UsersRoute;
