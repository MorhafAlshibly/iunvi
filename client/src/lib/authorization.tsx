import { ReactNode, useCallback } from "react";
import { useUser } from "@/lib/authentication";
import { ROLES, Role, ActiveUser } from "@/types/user";
import { WorkspaceRole } from "@/types/api/tenantManagement_pb";

export const POLICIES = {
  "admin:access": (user: ActiveUser) => user.role === ROLES.ADMIN,
  "developer:access": (
    user: ActiveUser,
    activeWorkspaceRole: WorkspaceRole | null,
  ) =>
    user.role === ROLES.ADMIN ||
    activeWorkspaceRole === WorkspaceRole.DEVELOPER,
  "user:access": (
    user: ActiveUser,
    activeWorkspaceRole: WorkspaceRole | null,
  ) =>
    user.role === ROLES.ADMIN ||
    activeWorkspaceRole === WorkspaceRole.USER ||
    activeWorkspaceRole === WorkspaceRole.DEVELOPER,
  "viewer:access": (
    user: ActiveUser,
    activeWorkspaceRole: WorkspaceRole | null,
  ) =>
    user.role === ROLES.ADMIN ||
    activeWorkspaceRole === WorkspaceRole.VIEWER ||
    activeWorkspaceRole === WorkspaceRole.USER ||
    activeWorkspaceRole === WorkspaceRole.DEVELOPER,
};

export const useAuthorization = () => {
  const user = useUser();
  if (!user.data) {
    throw Error("User does not exist!");
  }

  const checkAccess = useCallback(
    ({ allowedRoles }: { allowedRoles: Role[] }) => {
      if (allowedRoles && allowedRoles.length > 0 && user.data.role) {
        return allowedRoles?.includes(user.data.role);
      }

      return true;
    },
    [user.data.role],
  );

  return { checkAccess };
};

type AuthorizationProps = {
  forbiddenFallback?: ReactNode;
  children: ReactNode;
} & (
  | {
      allowedRoles: Role[];
      policyCheck?: never;
    }
  | {
      allowedRoles?: never;
      policyCheck: boolean;
    }
);

export const Authorization = ({
  policyCheck,
  allowedRoles,
  forbiddenFallback = null,
  children,
}: AuthorizationProps) => {
  const { checkAccess } = useAuthorization();

  let canAccess = false;

  if (allowedRoles) {
    canAccess = checkAccess({ allowedRoles });
  }

  if (typeof policyCheck !== "undefined") {
    canAccess = policyCheck;
  }

  return <>{canAccess ? children : forbiddenFallback}</>;
};
