import { Label } from "@radix-ui/react-dropdown-menu";
import { useMatch } from "react-router-dom";
import { useQuery } from "@connectrpc/connect-query";
import { paths } from "@/config/paths";
import CodeMirror from "@uiw/react-codemirror";
import { json } from "@codemirror/lang-json";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { ArrowBigLeft, CircleArrowLeft } from "lucide-react";
import {
  getFileGroup,
  getSpecification,
} from "@/types/api/file-FileService_connectquery";
import { DataMode } from "@/types/api/file_pb";
import { FileTransport } from "@/lib/api-client";
import { useDarkMode } from "usehooks-ts";
import { Separator } from "@/components/ui/separator";

const FileGroupsViewRoute = () => {
  const navigate = useNavigate();
  const id = useMatch(paths.app.user.fileGroups.root.getHref() + "/:id")?.params
    .id;
  const { data: fileGroupsData } = useQuery(
    getFileGroup,
    {
      id: id || "",
    },
    {
      enabled: !!id,
      transport: FileTransport,
    },
  );

  const fileGroup = fileGroupsData?.fileGroup;

  const darkMode = useDarkMode();

  return (
    <div className="grid grid-cols-1 gap-4">
      {fileGroup ? (
        <>
          <div className="grid grid-cols-2 col-span-1 justify-items-between">
            <div className="grid grid-cols-1 col-span-1 justify-items-start">
              <Label className="col-span-1 content-center font-medium text-lg">
                {fileGroup?.name}
              </Label>
            </div>
            <div className="grid grid-cols-1 col-span-1 justify-items-end">
              <Button
                size="sm"
                variant="outline"
                onClick={() => {
                  navigate(-1);
                }}
              >
                <CircleArrowLeft />
                Back
              </Button>
            </div>
          </div>
          <div className="grid grid-cols-1 col-span-1 mt-4">
            <div className="grid grid-cols-2 col-span-1">
              <div className="grid grid-cols-1 col-span-1 content-center justify-items-start">
                <Label className="font-medium">File name</Label>
              </div>
              <div className="grid grid-cols-1 col-span-1 content-center justify-items-end">
                <Label className="font-medium">Schema name</Label>
              </div>
              <div className="grid grid-cols-1 col-span-2">
                <Separator className="my-2" />
              </div>
            </div>
            {fileGroup?.schemaFileMappings.map((file, index) => (
              <div key={index} className="grid grid-cols-2 col-span-1">
                <div className="grid grid-cols-1 col-span-1 content-center justify-items-start">
                  <Label className="font-normal">
                    {file.landingZoneFileName}
                  </Label>
                </div>
                <div className="grid grid-cols-1 col-span-1 content-center justify-items-end">
                  <Label className="font-light">{file.schemaName}</Label>
                </div>
                <div className="grid grid-cols-1 col-span-2">
                  <Separator className="my-2" />
                </div>
              </div>
            ))}
          </div>
        </>
      ) : null}
    </div>
  );
};

export default FileGroupsViewRoute;
