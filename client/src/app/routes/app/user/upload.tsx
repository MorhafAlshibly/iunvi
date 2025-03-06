import { ContentLayout } from "@/components/layouts/content";
import { Button } from "@/components/ui/button";
import { useWorkspace } from "@/hooks/use-workspace";
import { createLandingZoneSharedAccessSignature } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useMutation } from "@connectrpc/connect-query";

const UploadRoute = () => {
  const { activeWorkspace } = useWorkspace();
  const createLandingZoneSharedAccessSignatureMutation = useMutation(
    createLandingZoneSharedAccessSignature,
  );

  return (
    <ContentLayout title="Upload">
      <div className="grid grid-cols-1 gap-4">
        <Button
          onClick={() => {
            createLandingZoneSharedAccessSignatureMutation.mutate({
              workspaceId: activeWorkspace?.id,
            });
          }}
        >
          Create SAS
        </Button>
      </div>
    </ContentLayout>
  );
};

export default UploadRoute;
