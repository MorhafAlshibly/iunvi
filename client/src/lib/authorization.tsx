import { ReactNode, useCallback } from "react";
import { useMsal, useIsAuthenticated } from "@azure/msal-react";

export enum ROLES {
  ADMIN = "ADMIN",
  USER = "USER",
}

type RoleTypes = keyof typeof ROLES;

export const POLICIES = {};

export const useAuthorization = () => {
  const { accounts } = useMsal();
  const isAuthenticated = useIsAuthenticated();
  if (!isAuthenticated) {
    throw Error("User does not exist!");
  }

  let role = accounts[0]?.idTokenClaims?.role as RoleTypes;

  const checkAccess = useCallback(
    ({ allowedRoles }: { allowedRoles: RoleTypes[] }) => {
      if (allowedRoles && allowedRoles.length > 0 && role) {
        return allowedRoles?.includes(role);
      }

      return true;
    },
    [role],
  );

  return { checkAccess, role };
};

type AuthorizationProps = {
  forbiddenFallback?: ReactNode;
  children: ReactNode;
} & (
  | {
      allowedRoles: RoleTypes[];
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
