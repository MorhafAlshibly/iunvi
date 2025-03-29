import { Outlet } from "react-router-dom";
import { ContentLayout } from "@/components/layouts/content";

export const ErrorBoundary = () => {
  return <div>Something went wrong!</div>;
};

const AppUserFileGroupsRoot = () => {
  return (
    <ContentLayout title="File groups">
      <Outlet />
    </ContentLayout>
  );
};

export default AppUserFileGroupsRoot;
