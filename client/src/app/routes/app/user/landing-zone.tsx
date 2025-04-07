import { ContentLayout } from "@/components/layouts/content";
import { Button } from "@/components/ui/button";
import { useWorkspace } from "@/hooks/use-workspace";
import {
  createLandingZoneSharedAccessSignature,
  getLandingZoneFiles,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
  useSuspenseInfiniteQuery,
} from "@connectrpc/connect-query";
import Uppy, { Meta, UppyFile, Body } from "@uppy/core";
import { useEffect, useState } from "react";
import { Dashboard, useUppyEvent } from "@uppy/react";
import "@uppy/core/dist/style.min.css";
import "@uppy/dashboard/dist/style.min.css";
import "./uppy.css";
import DataLakePlugin from "@/lib/uppy-blob";
import { Label } from "@/components/ui/label";
import { ChevronDown, RefreshCcw } from "lucide-react";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { formatBytes } from "@/utils/bytes";
import { Separator } from "@/components/ui/separator";
import { LandingZoneFile } from "@/types/api/tenantManagement_pb";
import { Input } from "@/components/ui/input";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";

const LandingZoneRoute = () => {
  const { activeWorkspace } = useWorkspace();

  const [currentPage, setCurrentPage] = useState<number>(1);
  const [searchFilter, setSearchFilter] = useState<string>("");

  const createLandingZoneSharedAccessSignatureMutation = useMutation(
    createLandingZoneSharedAccessSignature,
  );

  const {
    data: landingZoneFilesData,
    refetch: landingZoneFilesRefetch,
    hasNextPage: landingZoneFilesHasNextPage,
    fetchNextPage: landingZoneFilesFetchNextPage,
  } = useInfiniteQuery(
    getLandingZoneFiles,
    {
      workspaceId: activeWorkspace?.id || "",
      marker: "",
      prefix: searchFilter,
    },
    {
      pageParamKey: "marker",
      getNextPageParam: (lastPage, allPages, lastPageParam, allPageParams) =>
        lastPage.nextMarker,
      enabled: !!activeWorkspace?.id,
    },
  );

  const refreshFiles = () => {
    setCurrentPage(1);
    setSearchFilter("");
    landingZoneFilesRefetch();
  };

  const handlePageChange = (page: number) => {
    if (page < 1) return;
    if (page === (landingZoneFilesData?.pages.length ?? 0) + 1) {
      if (!landingZoneFilesHasNextPage) return;
      landingZoneFilesFetchNextPage();
    }
    setCurrentPage(page);
  };

  const handleSearchFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCurrentPage(1);
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
          {landingZoneFilesData?.pages[currentPage - 1]?.files.map(
            (file: LandingZoneFile) => (
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
            ),
          )}
        </div>
        <div className="grid grid-cols-1 col-span-1 justify-items-center">
          <Pagination>
            <PaginationContent>
              <PaginationItem>
                <PaginationPrevious
                  onClick={() => handlePageChange(currentPage - 1)}
                  // disabled={currentPage === 1}
                />
              </PaginationItem>
              {Array.from(
                { length: landingZoneFilesData?.pages.length || 0 },
                (_, page) => page,
              ).map((page: number) => (
                <PaginationItem key={page}>
                  <PaginationLink
                    onClick={() => handlePageChange(page + 1)}
                    isActive={currentPage === page + 1}
                  >
                    {page + 1}
                  </PaginationLink>
                </PaginationItem>
              ))}
              <PaginationItem>
                <PaginationEllipsis />
              </PaginationItem>
              <PaginationItem>
                <PaginationNext
                  onClick={() => handlePageChange(currentPage + 1)}
                  // disabled={
                  //   currentPage === landingZoneFilesData?.pages.length &&
                  //   !landingZoneFilesHasNextPage
                  // }
                />
              </PaginationItem>
            </PaginationContent>
          </Pagination>
        </div>
      </div>
    </ContentLayout>
  );
};

export default LandingZoneRoute;
