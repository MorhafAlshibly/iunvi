export const msalConfig = {
  auth: {
    clientId: import.meta.env.VITE_CLIENTID as string,
    authority: ("https://login.microsoftonline.com/" +
      import.meta.env.VITE_TENANTID) as string,
    redirectUri: import.meta.env.VITE_REDIRECTURI as string,
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
  scopes: [import.meta.env.VITE_SCOPE as string],
  prompt: "consent",
};
