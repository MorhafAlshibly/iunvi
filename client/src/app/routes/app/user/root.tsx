import { Outlet } from "react-router-dom";

import { Authorization, POLICIES } from "@/lib/authorization";
import { useUser } from "@/lib/authentication";
import { useWorkspace } from "@/hooks/use-workspace";
import { ActiveUser } from "@/types/user";

export const ErrorBoundary = () => {
  return <div>Something went wrong!</div>;
};

const AppUserRoot = () => {
  const user = useUser().data as ActiveUser;
  const { activeWorkspace, activeWorkspaceRole } = useWorkspace();
  return (
    <Authorization
      policyCheck={POLICIES["user:access"](
        user,
        activeWorkspace,
        activeWorkspaceRole,
      )}
    >
      <Outlet />
    </Authorization>
  );
};

export default AppUserRoot;
