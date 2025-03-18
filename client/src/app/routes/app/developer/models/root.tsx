import { Outlet } from "react-router";
import { ContentLayout } from "@/components/layouts/content";

export const ErrorBoundary = () => {
  return <div>Something went wrong!</div>;
};

const AppDeveloperModelsRoot = () => {
  return (
    <ContentLayout title="Models">
      <Outlet />
    </ContentLayout>
  );
};

export default AppDeveloperModelsRoot;
