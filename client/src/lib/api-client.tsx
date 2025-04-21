import { Interceptor, Transport } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { TransportProvider } from "@connectrpc/connect-query";
import { instance } from "./msal";
import { loginRequest } from "@/config/authentication";

const authInterceptor: Interceptor = (next) => async (request) => {
  const token = await instance.acquireTokenSilent(loginRequest);
  request.header.set("Authorization", `Bearer ${token.accessToken}`);
  return await next(request);
};

export const transports = {
  tenant: createConnectTransport({
    baseUrl: import.meta.env.VITE_TENANTURL as string,
    interceptors: [authInterceptor],
  }),
  file: createConnectTransport({
    baseUrl: import.meta.env.VITE_FILEURL as string,
    interceptors: [authInterceptor],
  }),
  model: createConnectTransport({
    baseUrl: import.meta.env.VITE_MODELURL as string,
    interceptors: [authInterceptor],
  }),
  dashboard: createConnectTransport({
    baseUrl: import.meta.env.VITE_DASHBOARDURL as string,
    interceptors: [authInterceptor],
  }),
};

export const TransportsProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  // Loop over the transports and create a TransportProvider for each and nest them
  let tail = children;
  for (const [key, transport] of Object.entries(transports)) {
    tail = (
      <TransportProvider key={key} transport={transport}>
        {tail}
      </TransportProvider>
    );
  }
  return tail;
};
