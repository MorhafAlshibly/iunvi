import { Interceptor } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { instance } from "./msal";
import { loginRequest } from "@/config/authentication";

const authInterceptor: Interceptor = (next) => async (request) => {
  const token = await instance.acquireTokenSilent(loginRequest);
  request.header.set("Authorization", `Bearer ${token.accessToken}`);
  return await next(request);
};

export const transport = createConnectTransport({
  baseUrl: import.meta.env.VITE_BACKENDURL as string,
  interceptors: [authInterceptor],
});
