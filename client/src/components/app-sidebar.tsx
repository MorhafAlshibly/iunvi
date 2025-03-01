import * as React from "react";
import {
  AudioWaveform,
  Blocks,
  BookOpen,
  Bot,
  Command,
  Component,
  Frame,
  GalleryVerticalEnd,
  Home,
  Map,
  Newspaper,
  PieChart,
  Settings2,
  SquareTerminal,
  User as UserIcon,
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
import { ActiveUser } from "@/types/user";
import { getWorkspaces } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { GetWorkspacesRequest } from "@/types/api/tenantManagement_pb";
import { useQuery } from "@connectrpc/connect-query";
import { paths } from "@/config/paths";
import { useWorkspace } from "@/hooks/use-workspace";
import { Button } from "./ui/button";
import { useNavigate } from "react-router";

const navViewer = [
  {
    title: "Run history",
    url: "#",
    icon: Map,
    isActive: true,
  },
  {
    title: "Results",
    url: "#",
    icon: BookOpen,
  },
];
const navUser = [
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
];

const navDeveloper = [
  {
    title: "Specifications",
    url: paths.app.developer.specifications.list.getHref(),
    icon: Newspaper,
  },
  {
    title: "Registry",
    url: paths.app.developer.registry.getHref(),
    icon: Blocks,
  },
  {
    title: "Models",
    url: paths.app.developer.models.getHref(),
    icon: Component,
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
];
const navAdmin = [
  {
    title: "Workspaces",
    url: paths.app.admin.workspaces.getHref(),
    icon: Settings2,
  },
  {
    title: "Users",
    url: paths.app.admin.users.getHref(),
    icon: UserIcon,
  },
];

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const navigate = useNavigate();
  const user = useUser().data as ActiveUser;
  const { activeWorkspaceRole, activeWorkspace } = useWorkspace();
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <WorkspaceSwitcher />
      </SidebarHeader>
      <SidebarContent>
        <Authorization
          policyCheck={POLICIES["viewer:access"](
            user,
            activeWorkspace,
            activeWorkspaceRole,
          )}
        >
          <NavList title="Viewer" items={navViewer} />
        </Authorization>
        <Authorization
          policyCheck={POLICIES["user:access"](
            user,
            activeWorkspace,
            activeWorkspaceRole,
          )}
        >
          <NavList title="User" items={navUser} />
        </Authorization>
        <Authorization
          policyCheck={POLICIES["developer:access"](
            user,
            activeWorkspace,
            activeWorkspaceRole,
          )}
        >
          <NavList title="Developer" items={navDeveloper} />
        </Authorization>
        <Authorization policyCheck={POLICIES["admin:access"](user)}>
          <NavList title="Admin" items={navAdmin} />
        </Authorization>
      </SidebarContent>
      <SidebarFooter>
        <NavUser />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
