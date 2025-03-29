import { Outlet } from "react-router-dom";

import { DashboardLayout } from "@/components/layouts/dashboard";
import { ProtectedRoute } from "@/lib/authentication";

export const ErrorBoundary = () => {
  return <div>Something went wrong!</div>;
};

const AppRoot = () => {
  return (
    <ProtectedRoute>
      <DashboardLayout>
        <Outlet />
      </DashboardLayout>
    </ProtectedRoute>
  );
};

export default AppRoot;
