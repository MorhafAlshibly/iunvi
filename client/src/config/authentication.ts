export const msalConfig = {
  auth: {
    clientId: import.meta.env.VITE_AZUREADCLIENTID as string,
    authority: ("https://login.microsoftonline.com/" +
      import.meta.env.VITE_AZUREADTENANTID) as string,
    redirectUri: import.meta.env.VITE_AZUREADREDIRECTURI as string,
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
  scopes: [import.meta.env.VITE_AZUREADSCOPE as string],
};
