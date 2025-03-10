import { ContentLayout } from "@/components/layouts/content";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useWorkspace } from "@/hooks/use-workspace";
import {
  createRegistryTokenPassword,
  getImages,
  getRegistryTokenPasswords,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { RegistryTokenPassword } from "@/types/registry";
import { timestampDate } from "@bufbuild/protobuf/wkt";
import { useMutation, useQuery } from "@connectrpc/connect-query";
import { Label } from "@radix-ui/react-dropdown-menu";
import { Check, Plus, RefreshCcw } from "lucide-react";
import { useEffect, useState } from "react";

const RegistryRoute = () => {
  const { activeWorkspace } = useWorkspace();
  const registryUrl = `${import.meta.env.VITE_REGISTRYNAME}.azurecr.io`;
  const registryUsername = `webapp-${activeWorkspace?.id.toLowerCase()}`;
  const scope = `scope-${activeWorkspace?.id.toLowerCase()}`;
  const [passwords, setPasswords] = useState<(RegistryTokenPassword | null)[]>([
    null,
    null,
  ]);

  const {
    data: registryTokenPasswords,
    refetch: refetchRegistryTokenPasswords,
  } = useQuery(
    getRegistryTokenPasswords,
    {
      workspaceId: activeWorkspace?.id,
    },
    { enabled: activeWorkspace != null },
  );

  useEffect(() => {
    if (!registryTokenPasswords) {
      return;
    }
    if (passwords[0] === null && registryTokenPasswords.password1) {
      const password1 = {
        password: ".........................",
        createdAt: timestampDate(registryTokenPasswords.password1),
      };
      setPasswords((prev) => [password1, prev[1]]);
    }
    if (passwords[1] === null && registryTokenPasswords.password2) {
      const password2 = {
        password: ".........................",
        createdAt: timestampDate(registryTokenPasswords.password2),
      };
      setPasswords((prev) => [prev[0], password2]);
    }
  }, [registryTokenPasswords]);

  const createRegistryTokenPasswordMutation = useMutation(
    createRegistryTokenPassword,
  );

  const handleCreatePassword = async (password2: boolean) => {
    const { password, createdAt } =
      await createRegistryTokenPasswordMutation.mutateAsync({
        workspaceId: activeWorkspace?.id,
        password2: password2,
      });
    if (!password2) {
      setPasswords((prev) => [
        {
          password,
          createdAt: createdAt ? timestampDate(createdAt) : new Date(),
        },
        prev[1],
      ]);
    } else {
      setPasswords((prev) => [
        prev[0],
        {
          password,
          createdAt: createdAt ? timestampDate(createdAt) : new Date(),
        },
      ]);
    }
    await refetchRegistryTokenPasswords();
  };

  const { data: imagesData, refetch: refetchImages } = useQuery(
    getImages,
    {
      workspaceId: activeWorkspace?.id,
    },
    { enabled: activeWorkspace != null },
  );

  const images = imagesData?.images ?? [];

  return (
    <ContentLayout title="Registry">
      <div className="grid grid-cols-1 gap-4">
        <code className="col-span-1 text-sm sm:text-base inline-flex flex-col text-left space-x-4 bg-gray-800 text-white rounded-lg p-4 pl-6">
          <span className="flex-1">
            <span>docker login -u {registryUsername} -p</span>
            <span className="text-yellow-500"> $password </span>
            <span>{registryUrl}</span>
          </span>
          <span className="flex-1">
            <span>docker tag</span>
            <span className="text-yellow-500"> $image </span>
            <span>
              {registryUrl}/{scope}/
            </span>
            <span className="text-yellow-500">$image</span>
          </span>
          <span className="flex-1">
            <span>
              docker push {registryUrl}/{scope}/
            </span>
            <span className="text-yellow-500">$image</span>
          </span>
        </code>
        {passwords.map((password, index) => (
          <div
            className="grid grid-cols-3 justify-between gap-4 col-span-1"
            key={"password" + (index + 1)}
          >
            <Label className="flex col-span-1 items-center">
              Password {index + 1}
            </Label>
            {password ? (
              <>
                <code className="flex justify-end gap-4 col-span-2 items-center text-sm sm:text-base inline-flex text-left items-center space-x-4 bg-gray-800 text-white rounded-lg p-4 pl-6">
                  <span>{password.password}</span>
                  <Label className="flex justify-end gap-4 col-span-1 items-center">
                    {password.createdAt.toLocaleString()}
                  </Label>
                  <Button
                    size="sm"
                    variant="ghost"
                    onClick={() => handleCreatePassword(index === 1)}
                  >
                    <RefreshCcw />
                  </Button>
                </code>
              </>
            ) : (
              <code className="flex justify-end gap-4 col-span-2 items-center text-sm sm:text-base inline-flex text-left items-center space-x-4 bg-gray-800 text-white rounded-lg p-4 pl-6">
                <span>Does not exist</span>
                <Button
                  size="sm"
                  variant="ghost"
                  onClick={() => handleCreatePassword(index === 1)}
                >
                  <Plus />
                </Button>
              </code>
            )}
          </div>
        ))}
        <div>
          <div className="flex mt-10">
            <Label className="flex-1 items-center font-semibold text-lg">
              Images
            </Label>
            <div className="flex-1 text-right">
              <Button size="sm" variant="ghost" onClick={() => refetchImages()}>
                <RefreshCcw />
              </Button>
            </div>
          </div>
          <Separator className="my-2" />
          {images.map((image) => (
            <div key={image.name} className="grid grid-cols-3">
              <Label className="flex col-span-3 items-center">
                {image.name}
              </Label>
              <Separator className="col-span-3 my-2" />
            </div>
          ))}
        </div>
      </div>
    </ContentLayout>
  );
};

export default RegistryRoute;
