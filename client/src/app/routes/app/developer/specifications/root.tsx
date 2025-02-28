import { Outlet } from "react-router";
import { ContentLayout } from "@/components/layouts/content";

export const ErrorBoundary = () => {
  return <div>Something went wrong!</div>;
};

const AppDeveloperSpecificationsRoot = () => {
  return (
    <ContentLayout title="Specifications">
      <Outlet />
    </ContentLayout>
  );
};

export default AppDeveloperSpecificationsRoot;
