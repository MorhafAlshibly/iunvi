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
import "./uppy.css";
import DataLakePlugin from "@/lib/uppy-data-lake";
import { Label } from "@/components/ui/label";
import { ChevronDown, RefreshCcw } from "lucide-react";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { formatBytes } from "@/utils/bytes";
import { Separator } from "@/components/ui/separator";
import { LandingZoneFile } from "@/types/api/tenantManagement_pb";
import { Input } from "@/components/ui/input";

const LandingZoneRoute = () => {
  const { activeWorkspace } = useWorkspace();

  const [files, setFiles] = useState<Map<string, LandingZoneFile>>(new Map());
  const [nextMarker, setNextMarker] = useState<string | undefined>(undefined);
  const [searchFilter, setSearchFilter] = useState<string>("");

  const createLandingZoneSharedAccessSignatureMutation = useMutation(
    createLandingZoneSharedAccessSignature,
  );

  const { data: landingZoneFilesData, refetch: landingZoneFilesRefetch } =
    useQuery(
      getLandingZoneFiles,
      {
        workspaceId: activeWorkspace?.id,
        marker: nextMarker,
        prefix: searchFilter,
      },
      {
        enabled: !!activeWorkspace?.id,
      },
    );

  useEffect(() => {
    if (landingZoneFilesData?.files) {
      setFiles((prevFiles) => {
        const newFiles = new Map(prevFiles);
        landingZoneFilesData.files.forEach((file) => {
          newFiles.set(file.name, file);
        });
        return newFiles;
      });
    }
  }, [landingZoneFilesData]);

  const getMoreFiles = () => {
    setNextMarker(landingZoneFilesData?.nextMarker);
  };

  const refreshFiles = () => {
    setNextMarker(undefined);
    setFiles(new Map());
    setSearchFilter("");
    landingZoneFilesRefetch();
  };

  const handleSearchFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNextMarker(undefined);
    setFiles(new Map());
    setSearchFilter(e.target.value);
    landingZoneFilesRefetch();
  };

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
        <div className="grid grid-cols-1 col-span-1 justify-items">
          <Dashboard uppy={uppy} />
        </div>
        <div className="grid grid-cols-2 col-span-1 justify-items-between">
          <div className="grid grid-cols-1 col-span-1 justify-items-start content-center">
            <Label className="text-lg font-semibold">Files</Label>
          </div>
          <div className="grid grid-cols-1 col-span-1 justify-items-end">
            <Button size="sm" variant="ghost" onClick={refreshFiles}>
              <RefreshCcw />
            </Button>
          </div>
        </div>
        <div className="grid grid-cols-1 col-span-1 justify-items-start">
          <Input
            type="text"
            placeholder="Search"
            onChange={handleSearchFilterChange}
            value={searchFilter}
          />
        </div>

        <div className="grid grid-cols-1 col-span-1">
          {Array.from(files.values()).map((file) => (
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
        {landingZoneFilesData?.nextMarker ? (
          <div className="grid grid-cols-1 col-span-1 justify-items-center">
            <Button size="sm" variant="outline" onClick={getMoreFiles}>
              <ChevronDown />
            </Button>
          </div>
        ) : null}
      </div>
    </ContentLayout>
  );
};

export default LandingZoneRoute;
