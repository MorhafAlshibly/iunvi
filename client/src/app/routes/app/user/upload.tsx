import { ContentLayout } from "@/components/layouts/content";
import { Button } from "@/components/ui/button";
import { useWorkspace } from "@/hooks/use-workspace";
import { createLandingZoneSharedAccessSignature } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useMutation } from "@connectrpc/connect-query";
import Uppy, { Meta, UppyFile, Body } from "@uppy/core";
import { useEffect, useState } from "react";
import { Dashboard, useUppyEvent } from "@uppy/react";
import AwsS3 from "@uppy/aws-s3";
import AzureBlob from "uppy-azure-blob";
import Tus from "@uppy/tus";

import "@uppy/core/dist/style.min.css";
import "@uppy/dashboard/dist/style.min.css";
import AzureBlobPlugin from "@/lib/uppy-data-lake";

const UploadRoute = () => {
  const { activeWorkspace } = useWorkspace();

  const createLandingZoneSharedAccessSignatureMutation = useMutation(
    createLandingZoneSharedAccessSignature,
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
    }).use(AzureBlobPlugin, {
      getSasUrl: async (file: UppyFile<Meta, Body>) => {
        return getAzureSas(file.name);
      },
      chunkSize: 10 * 1024 * 1024,
      maxRetries: 3,
    }),
  );

  // const [uppy] = useState(() =>
  //   new Uppy({
  //     restrictions: {
  //       allowedFileTypes: [".csv"],
  //     },
  //     autoProceed: false,
  //     allowMultipleUploadBatches: false,
  //   }).use(Tus, {
  //     endpoint:
  //       "https://app-tus-iunvi-dev-eastus-001-f5czeub6c8h8dea3.eastus-01.azurewebsites.net",

  //   }),
  // );

  // const [uppy] = useState(() =>
  //   new Uppy({
  //     restrictions: {
  //       allowedFileTypes: [".csv"],
  //     },
  //     autoProceed: false,
  //     allowMultipleUploadBatches: false,
  //     // @ts-ignore
  //   }).use(AzureBlob, {
  //     endpoint: "https://saiunvideveastus001.blob.core.windows.net",
  //     containerName: "landing-zone",
  //     sas: "?sv=2021-08-06",
  //   }),
  // );

  // useEffect(() => {
  //   getAzureSas("test").then((url) => {
  //     const sas = "?" + url.split("?")[1];
  //     uppy?.getPlugin("AzureBlob")?.setOptions({
  //       sas,
  //     });
  //   });
  // }, []);

  // const [uppy] = useState(() =>
  //   new Uppy({
  //     restrictions: {
  //       allowedFileTypes: [".csv"],
  //     },
  //     autoProceed: false,
  //     allowMultipleUploadBatches: false,
  //   }).use(AwsS3, {
  //     getUploadParameters: async (file: UppyFile<Meta, Body>) => {
  //       const url = await getAzureSas(file.name);
  //       return {
  //         method: "PUT",
  //         url,
  //         fields: {},
  //         headers: {
  //           "x-ms-blob-type": "BlockBlob",
  //         },
  //       };
  //     },
  //   } as any),
  // );

  // useEffect(() => {
  //   uppy?.getPlugin("AwsS3")?.setOptions({
  //     getUploadParameters: () => {
  //       return getAzureSas().then((url) => {
  //         return {
  //           method: "POST",
  //           url: url,
  //           fields: {},
  //           headers: {
  //             "x-ms-blob-type": "PageBlob",
  //           },
  //         };
  //       });
  //     },
  //   });
  // }, []);

  // useEffect(async () => {
  //   uppy?.getPlugin("AwsS3")?.setOptions({
  //     getUploadParameters: (file) => {
  //       return getAzureSas({
  //         filename: file.name,
  //         contentType: file.type,
  //         extension: file.extension,
  //       }).then((data) => {
  //         return {
  //           method: "PUT",
  //           url: data.url,
  //           fields: {},
  //           headers: {
  //             "x-ms-blob-type": "BlockBlob",
  //           },
  //         };
  //       });
  //     },
  //   });
  // }, []);

  return (
    <ContentLayout title="Upload">
      <div className="grid grid-cols-1 gap-4">
        <div className="grid grid-cols-1 col-span-1 justify-items-start">
          <Dashboard uppy={uppy} />
        </div>
      </div>
    </ContentLayout>
  );
};

export default UploadRoute;
