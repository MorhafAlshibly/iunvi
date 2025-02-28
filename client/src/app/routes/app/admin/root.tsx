import { Outlet } from "react-router";

import { Authorization, POLICIES } from "@/lib/authorization";
import { useUser } from "@/lib/authentication";
import { ActiveUser } from "@/types/user";

export const ErrorBoundary = () => {
  return <div>Something went wrong!</div>;
};

const AppAdminRoot = () => {
  const user = useUser().data as ActiveUser;
  return (
    <Authorization policyCheck={POLICIES["admin:access"](user)}>
      <Outlet />
    </Authorization>
  );
};

export default AppAdminRoot;
