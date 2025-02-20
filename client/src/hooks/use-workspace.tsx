import React, { createContext, useState, useEffect } from "react";
import { Workspace } from "@/types/api/tenantManagement_pb";
import { getWorkspaces } from "@/types/api/tenantManagement-TenantManagementService_connectquery";
import { useQuery } from "@connectrpc/connect-query";

const WorkspaceContext = createContext({
  workspaces: [] as Workspace[],
  activeWorkspace: null as Workspace | null,
  changeWorkspace: (workspace: Workspace) => {},
  refetch: () => {},
});

export const WorkspaceProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [activeWorkspace, setActiveWorkspace] = useState(
    null as Workspace | null,
  );

  const [workspaces, setWorkspaces] = useState([] as Workspace[]);

  const { data: workspacesData, refetch } = useQuery(getWorkspaces);

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

  const changeWorkspace = (workspace: Workspace) => {
    setActiveWorkspace(workspace);
  };

  return (
    <WorkspaceContext.Provider
      value={{ workspaces, activeWorkspace, changeWorkspace, refetch }}
    >
      {children}
    </WorkspaceContext.Provider>
  );
};

export const useWorkspace = () => React.useContext(WorkspaceContext);
