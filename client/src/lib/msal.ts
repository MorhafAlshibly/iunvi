import {
  PublicClientApplication,
  EventType,
  AccountInfo,
} from "@azure/msal-browser";

import { msalConfig } from "@/config/authentication";

export const instance = new PublicClientApplication(msalConfig);

// Default to using the first account if no account is active
if (!instance.getActiveAccount() && instance.getAllAccounts().length) {
  instance.setActiveAccount(instance.getAllAccounts()[0]);
}

instance.addEventCallback((message) => {
  if (
    (message.eventType === EventType.LOGIN_SUCCESS ||
      message.eventType === EventType.ACQUIRE_TOKEN_SUCCESS ||
      message.eventType === EventType.SSO_SILENT_SUCCESS) &&
    message.payload
  ) {
    instance.setActiveAccount(message.payload as AccountInfo);
  }
});
