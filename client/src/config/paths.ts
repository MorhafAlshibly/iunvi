export const paths = {
  home: {
    path: "/",
    getHref: () => "/",
  },

  auth: {
    register: {
      path: "/auth/register",
      getHref: (redirectTo?: string | null) =>
        `/auth/register${redirectTo ? `?redirectTo=${encodeURIComponent(redirectTo)}` : ""}`,
    },
    login: {
      path: "/auth/login",
      getHref: (redirectTo?: string | null) =>
        `/auth/login${redirectTo ? `?redirectTo=${encodeURIComponent(redirectTo)}` : ""}`,
    },
  },

  app: {
    root: {
      path: "/app",
      getHref: () => "/app",
    },
    dashboard: {
      path: "",
      getHref: () => "/app",
    },
    admin: {
      root: {
        path: "admin",
        getHref: () => "/app/admin",
      },
      workspaces: {
        path: "workspaces",
        getHref: () => "/app/admin/workspaces",
      },
      users: {
        path: "users",
        getHref: () => "/app/admin/users",
      },
    },
    developer: {
      root: {
        path: "developer",
        getHref: () => "/app/developer",
      },
      registry: {
        path: "registry",
        getHref: () => "/app/developer/registry",
      },
      models: {
        path: "models",
        getHref: () => "/app/developer/models",
      },
    },
  },
} as const;
