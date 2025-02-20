import * as React from "react";
import {
  AudioWaveform,
  BookOpen,
  Bot,
  Command,
  Frame,
  GalleryVerticalEnd,
  Map,
  PieChart,
  Settings2,
  SquareTerminal,
} from "lucide-react";

import { NavList } from "@/components/nav-list";
import { NavUser } from "@/components/nav-user";
import { WorkspaceSwitcher } from "@/components/workspace-switcher";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { useLogout, useUser } from "@/lib/authentication";
import { Authorization, POLICIES } from "@/lib/authorization";
import { User } from "@/types/user";
import { getWorkspaces } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { GetWorkspacesRequest } from "@/types/api/tenantManagement_pb";
import { useQuery } from "@connectrpc/connect-query";
import { paths } from "@/config/paths";

// This is sample data.
const data = {
  workspaces: [
    {
      name: "Acme Inc",
      logo: GalleryVerticalEnd,
      plan: "Enterprise",
    },
    {
      name: "Acme Corp.",
      logo: AudioWaveform,
      plan: "Startup",
    },
    {
      name: "Evil Corp.",
      logo: Command,
      plan: "Free",
    },
  ],
  navViewer: [
    {
      title: "Run history",
      url: "#",
      icon: Map,
      isActive: true,
    },
    {
      title: "Dashboard",
      url: "#",
      icon: BookOpen,
    },
  ],
  navUser: [
    {
      title: "Upload files",
      url: "#",
      icon: Settings2,
    },
    {
      title: "Run models",
      url: "#",
      icon: PieChart,
    },
  ],
  navDeveloper: [
    {
      title: "Specifications",
      url: "#",
      icon: SquareTerminal,
      isActive: true,
      items: [
        {
          title: "History",
          url: "#",
        },
        {
          title: "Starred",
          url: "#",
        },
        {
          title: "Settings",
          url: "#",
        },
      ],
    },
    {
      title: "Models",
      url: "#",
      icon: Bot,
      items: [
        {
          title: "Genesis",
          url: "#",
        },
        {
          title: "Explorer",
          url: "#",
        },
        {
          title: "Quantum",
          url: "#",
        },
      ],
    },
    {
      title: "Charts",
      url: "#",
      icon: BookOpen,
      items: [
        {
          title: "Introduction",
          url: "#",
        },
        {
          title: "Get Started",
          url: "#",
        },
        {
          title: "Tutorials",
          url: "#",
        },
        {
          title: "Changelog",
          url: "#",
        },
      ],
    },
  ],
  navAdmin: [
    {
      title: "Workspaces",
      url: paths.app.admin.workspaces.getHref(),
      icon: Settings2,
    },
    {
      title: "Users",
      url: "#",
      icon: PieChart,
    },
  ],
};

export function AppSidebar({
  logoutFn,
  ...props
}: React.ComponentProps<typeof Sidebar> & {
  logoutFn: () => void;
}) {
  const user = useUser().data as User;
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <WorkspaceSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <NavList title="Viewer" items={data.navViewer} />
        <NavList title="User" items={data.navUser} />
        <NavList title="Developer" items={data.navDeveloper} />
        <Authorization policyCheck={POLICIES["admin:access"](user as User)}>
          <NavList title="Admin" items={data.navAdmin} />
        </Authorization>
      </SidebarContent>
      <SidebarFooter>
        <NavUser
          user={{
            name: user.displayName,
            email: user.username,
            avatar: "",
          }}
          logoutFn={logoutFn}
        />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
