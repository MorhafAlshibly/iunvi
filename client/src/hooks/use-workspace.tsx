import React, { createContext, useState, useEffect } from "react";
import { Workspace, WorkspaceRole } from "@/types/api/tenantManagement_pb";
import {
  getUserWorkspaceAssignment,
  getWorkspaces,
} from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useQuery } from "@connectrpc/connect-query";
import { useUser } from "@/lib/authentication";

const WorkspaceContext = createContext({
  workspaces: [] as Workspace[],
  activeWorkspaceRole: null as WorkspaceRole | null,
  activeWorkspace: null as Workspace | null,
  changeWorkspace: (workspace: Workspace) => {},
  workspacesRefetch: () => {},
  workspaceRoleRefetch: () => {},
});

export const WorkspaceProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const user = useUser();

  const [activeWorkspace, setActiveWorkspace] = useState(
    null as Workspace | null,
  );

  const [activeWorkspaceRole, setActiveWorkspaceRole] = useState(
    null as WorkspaceRole | null,
  );

  const [workspaces, setWorkspaces] = useState([] as Workspace[]);

  const { data: workspacesData, refetch: workspacesRefetch } =
    useQuery(getWorkspaces);
  const { data: workspaceRoleData, refetch: workspaceRoleRefetch } = useQuery(
    getUserWorkspaceAssignment,
    {
      userObjectId: user.data?.objectId,
      workspaceId: activeWorkspace?.id,
    },
    { enabled: activeWorkspace !== null && user.data !== null },
  );

  useEffect(() => {
    if (workspacesData) {
      setWorkspaces(workspacesData.workspaces);
    }

    if (
      activeWorkspace &&
      !workspacesData?.workspaces.includes(activeWorkspace)
    ) {
      setActiveWorkspace(null);
    }
  }, [workspacesData]);

  useEffect(() => {
    setActiveWorkspaceRole(workspaceRoleData?.role || null);
  }, [workspaceRoleData]);

  const changeWorkspace = (workspace: Workspace) => {
    setActiveWorkspace(workspace);
  };

  return (
    <WorkspaceContext.Provider
      value={{
        workspaces,
        activeWorkspace,
        changeWorkspace,
        activeWorkspaceRole,
        workspacesRefetch,
        workspaceRoleRefetch,
      }}
    >
      {children}
    </WorkspaceContext.Provider>
  );
};

export const useWorkspace = () => React.useContext(WorkspaceContext);
