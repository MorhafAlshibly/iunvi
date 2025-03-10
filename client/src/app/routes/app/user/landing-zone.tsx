import { ContentLayout } from "@/components/layouts/content";
import { Button } from "@/components/ui/button";
import { useWorkspace } from "@/hooks/use-workspace";
import {
  createLandingZoneSharedAccessSignature,
  getLandingZoneFiles,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import Uppy, { Meta, UppyFile, Body } from "@uppy/core";
import { useEffect, useState } from "react";
import { Dashboard, useUppyEvent } from "@uppy/react";
import "@uppy/core/dist/style.min.css";
import "@uppy/dashboard/dist/style.min.css";
import DataLakePlugin from "@/lib/uppy-data-lake";
import { Label } from "@/components/ui/label";
import { RefreshCcw } from "lucide-react";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { formatBytes } from "@/utils/bytes";
import { Separator } from "@/components/ui/separator";

const LandingZoneRoute = () => {
  const { activeWorkspace } = useWorkspace();

  const createLandingZoneSharedAccessSignatureMutation = useMutation(
    createLandingZoneSharedAccessSignature,
  );

  const [nextPageMarker, setNextPageMarker] = useState<string | undefined>(
    undefined,
  );

  const nextPage = () => {
    if (!landingZoneFilesResponse?.nextMarker) {
      return;
    }
    setNextPageMarker(landingZoneFilesResponse.nextMarker);
    landingZoneFilesRefetch();
  };

  const { data: landingZoneFilesResponse, refetch: landingZoneFilesRefetch } =
    useQuery(
      getLandingZoneFiles,
      {
        workspaceId: activeWorkspace?.id,
        marker: nextPageMarker,
      },
      {
        enabled: !!activeWorkspace?.id,
      },
    );

  const getAzureSas = async (fileName: string | undefined) => {
    if (!activeWorkspace?.id || !fileName) {
      return "";
    }
    const res =
      await createLandingZoneSharedAccessSignatureMutation.mutateAsync({
        workspaceId: activeWorkspace?.id,
        fileName,
      });
    return res.url;
  };

  const [uppy] = useState(() =>
    new Uppy({
      restrictions: {
        allowedFileTypes: [".csv"],
      },
      autoProceed: false,
      allowMultipleUploadBatches: false,
    }).use(DataLakePlugin, {
      getSasUrl: async (file: UppyFile<Meta, Body>) => {
        return getAzureSas(file.name);
      },
    }),
  );

  return (
    <ContentLayout title="Landing zone">
      <div className="grid grid-cols-1 gap-4">
        <div className="grid grid-cols-1 col-span-1 justify-items-start">
          <Dashboard uppy={uppy} />
        </div>
        <div className="grid grid-cols-2 col-span-1 justify-items-between">
          <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
            <Label className="text-lg font-semibold">Files</Label>
          </div>
          <div className="grid grid-cols-1 col-span-1 justify-items-end">
            <Button
              size="sm"
              variant="ghost"
              onClick={() => landingZoneFilesRefetch}
            >
              <RefreshCcw />
            </Button>
          </div>
        </div>
        <div className="grid grid-cols-1 col-span-1">
          {landingZoneFilesResponse?.files.map((file) => (
            <div key={file.name} className="grid grid-cols-3 col-span-1 p-2">
              <div className="grid grid-cols-1 col-span-1 content-center justify-items-start">
                <Label className="font-normal">{file.name}</Label>
              </div>
              <div className="grid grid-cols-1 col-span-1 content-center justify-items-end">
                <Label className="font-light">
                  {formatBytes(Number(file.size))}
                </Label>
              </div>
              <div className="grid grid-cols-1 col-span-1 content-center justify-items-end">
                <Label className="font-light">
                  {file.lastModified
                    ? timestampDate(file.lastModified).toLocaleString()
                    : ""}
                </Label>
              </div>
              <div className="grid grid-cols-1 col-span-3">
                <Separator className="mt-2" />
              </div>
            </div>
          ))}
        </div>
      </div>
    </ContentLayout>
  );
};

export default LandingZoneRoute;
