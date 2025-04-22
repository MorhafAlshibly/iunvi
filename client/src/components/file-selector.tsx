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
import { LandingZoneFile } from "@/types/api/file_pb";

export function FileSelector({
  files,
  selectedFileName,
  setSelectedFileName,
  searchFilter,
  setSearchFilter,
}: {
  files: LandingZoneFile[];
  selectedFileName: string | null;
  setSelectedFileName: React.Dispatch<React.SetStateAction<string | null>>;
  searchFilter: string;
  setSearchFilter: React.Dispatch<React.SetStateAction<string>>;
}) {
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
          {selectedFileName ?? "Select file..."}
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput
            placeholder="Search file..."
            className="h-9"
            onValueChange={setSearchFilter}
            value={searchFilter}
          />
          <CommandList>
            <CommandEmpty>No file found.</CommandEmpty>
            <CommandGroup>
              {files.map((file) => (
                <CommandItem
                  key={file.name}
                  value={file.name}
                  onSelect={() => {
                    setSelectedFileName(file.name);
                    setOpen(false);
                  }}
                >
                  {file.name}
                  <Check
                    className={cn(
                      "ml-auto",
                      selectedFileName === file.name
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
