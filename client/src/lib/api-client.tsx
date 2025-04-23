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

export const TenantTransport = createConnectTransport({
  baseUrl: import.meta.env.VITE_TENANTURL as string,
  interceptors: [authInterceptor],
});

export const FileTransport = createConnectTransport({
  baseUrl: import.meta.env.VITE_FILEURL as string,
  interceptors: [authInterceptor],
});

export const ModelTransport = createConnectTransport({
  baseUrl: import.meta.env.VITE_MODELURL as string,
  interceptors: [authInterceptor],
});

export const DashboardTransport = createConnectTransport({
  baseUrl: import.meta.env.VITE_DASHBOARDURL as string,
  interceptors: [authInterceptor],
});
