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
import { getDashboards } from "@/types/api/dashboard-DashboardService_connectquery";

export function DashboardSelector({
  modelRunId,
  selectedDashboardId,
  setSelectedDashboardId,
}: {
  modelRunId: string;
  selectedDashboardId: string | null;
  setSelectedDashboardId: React.Dispatch<React.SetStateAction<string | null>>;
}) {
  const { activeWorkspace } = useWorkspace();

  const { data: dashboardsData } = useQuery(
    getDashboards,
    {
      workspaceId: activeWorkspace?.id ?? "",
      modelRunId,
    },
    {
      enabled: !!activeWorkspace,
    },
  );

  const dashboards = dashboardsData?.dashboards || [];

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
          {selectedDashboardId
            ? dashboards.find(
                (dashboard) => dashboard.id === selectedDashboardId,
              )?.name
            : `Select dashboard...`}
          <ChevronsUpDown className="opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="p-0">
        <Command>
          <CommandInput placeholder="Search dashboard..." className="h-9" />
          <CommandList>
            <CommandEmpty>No dashboard found.</CommandEmpty>
            <CommandGroup>
              {dashboards.map((dashboard) => (
                <CommandItem
                  key={dashboard.id}
                  value={dashboard.name}
                  onSelect={() => {
                    setSelectedDashboardId(dashboard.id);
                    setOpen(false);
                  }}
                >
                  {dashboard.name}
                  <Check
                    className={cn(
                      "ml-auto",
                      selectedDashboardId === dashboard.id
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
