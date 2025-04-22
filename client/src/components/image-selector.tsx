import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";

import { cn } from "@/utils/cn";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useQuery } from "@connectrpc/connect-query";
import { useWorkspace } from "@/hooks/use-workspace";
import { getImages } from "@/types/api/model-ModelService_connectquery";

export function ImageSelector({
  selectedImageName,
  setSelectedImageName,
}: {
  selectedImageName: string | null;
  setSelectedImageName: React.Dispatch<React.SetStateAction<string | null>>;
}) {
  const { activeWorkspace } = useWorkspace();

  const { data: imagesData } = useQuery(
    getImages,
    {
      workspaceId: activeWorkspace?.id ?? "",
    },
    {
      enabled: !!activeWorkspace,
    },
  );

  const images = imagesData?.images || [];

  const [open, setOpen] = React.useState(false);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
        >
          {selectedImageName ? selectedImageName : `Select image...`}
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput placeholder="Search image..." className="h-9" />
          <CommandList>
            <CommandEmpty>No image found.</CommandEmpty>
            <CommandGroup>
              {images.map((image) => (
                <CommandItem
                  key={image.name}
                  value={image.name}
                  onSelect={() => {
                    setSelectedImageName(image.name);
                    setOpen(false);
                  }}
                >
                  {image.name}
                  <Check
                    className={cn(
                      "ml-auto",
                      selectedImageName === image.name
                        ? "opacity-100"
                        : "opacity-0",
                    )}
                  />
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
