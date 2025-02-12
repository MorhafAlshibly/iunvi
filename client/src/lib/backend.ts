import { createConnectTransport } from "@connectrpc/connect-web";

export const transport = createConnectTransport({
  baseUrl: import.meta.env.VITE_BACKENDURL as string,
});
