import { MsalProvider } from "@azure/msal-react";
import { TransportProvider } from "@connectrpc/connect-query";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { transport } from "@/lib/backend";
import { instance } from "@/lib/msal";

interface AppProviderProps {
  children: React.ReactNode;
}

export const AppProvider = ({ children }: AppProviderProps) => {
  return (
    <TransportProvider transport={transport}>
      <QueryClientProvider client={new QueryClient()}>
        <MsalProvider instance={instance}>{children}</MsalProvider>
      </QueryClientProvider>
    </TransportProvider>
  );
};
