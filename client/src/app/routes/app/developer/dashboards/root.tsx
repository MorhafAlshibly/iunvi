import { Outlet } from "react-router-dom";
import { ContentLayout } from "@/components/layouts/content";

export const ErrorBoundary = () => {
  return <div>Something went wrong!</div>;
};

const AppDeveloperDashboardsRoot = () => {
  return (
    <ContentLayout title="Dashboards">
      <Outlet />
    </ContentLayout>
  );
};

export default AppDeveloperDashboardsRoot;
