import { ContentLayout } from "@/components/layouts/content";
import { useUser } from "@/lib/authentication";
import { ROLES } from "@/types/user";

const HomeRoute = () => {
  const user = useUser();
  return (
    <ContentLayout title="Home">
      Welcome {user.data?.displayName}!
    </ContentLayout>
  );
};

export default HomeRoute;
