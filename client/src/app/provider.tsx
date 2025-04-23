import { MsalProvider } from "@azure/msal-react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { instance } from "@/lib/msal";
import { WorkspaceProvider } from "@/hooks/use-workspace";
import { ThemeProvider } from "@/lib/theme";

interface AppProviderProps {
  children: React.ReactNode;
}

export const AppProvider = ({ children }: AppProviderProps) => {
  return (
    <QueryClientProvider client={new QueryClient()}>
      <ThemeProvider>
        <MsalProvider instance={instance}>
          <WorkspaceProvider>{children}</WorkspaceProvider>
        </MsalProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
};
