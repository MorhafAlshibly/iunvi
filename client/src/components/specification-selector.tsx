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
import { SpecificationName } from "@/types/api/tenantManagement_pb";

export function SpecificationSelector({
  specifications,
  selectedSpecification,
  setSelectedSpecification,
}: {
  specifications: SpecificationName[];
  selectedSpecification: SpecificationName | null;
  setSelectedSpecification: React.Dispatch<
    React.SetStateAction<SpecificationName | null>
  >;
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
          {selectedSpecification
            ? specifications.find(
                (specification) =>
                  specification.id === selectedSpecification?.id,
              )?.name
            : "Select specification..."}
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput placeholder="Search specification..." className="h-9" />
          <CommandList>
            <CommandEmpty>No specification found.</CommandEmpty>
            <CommandGroup>
              {specifications.map((specification) => (
                <CommandItem
                  key={specification.id}
                  value={specification.name}
                  onSelect={() => {
                    setSelectedSpecification(specification);
                    setOpen(false);
                  }}
                >
                  {specification.name}
                  <Check
                    className={cn(
                      "ml-auto",
                      selectedSpecification?.id === specification.id
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
