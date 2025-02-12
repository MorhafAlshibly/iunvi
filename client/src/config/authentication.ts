export const msalConfig = {
  auth: {
    clientId: import.meta.env.VITE_AZUREADCLIENTID,
    authority:
      "https://login.microsoftonline.com/" +
      import.meta.env.VITE_AZUREADTENANTID,
    redirectUri: import.meta.env.VITE_AZUREADREDIRECTURI,
  },
  cache: {
    cacheLocation: "sessionStorage", // 'sessionStorage' is more secure but 'localStorage' keeps the user signed in
    storeAuthStateInCookie: false, // For IE11/Edge compatibility
  },
  system: {
    tokenRenewalOffsetSeconds: 300, // Refresh token 5 minutes before expiry
  },
};

export const loginRequest = {
  scopes: ["User.Read"], // Define permissions you need (e.g., Microsoft Graph API)
};
