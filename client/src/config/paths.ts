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
  },
} as const;
