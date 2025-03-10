export const paths = {
  landing: {
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
    home: {
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
      specifications: {
        root: {
          path: "specifications",
          getHref: () => "/app/developer/specifications",
        },
        list: {
          path: "",
          getHref: () => "/app/developer/specifications",
        },
        view: {
          path: ":id",
          getHref: (id: string) => `/app/developer/specifications/${id}`,
        },
        create: {
          path: "create",
          getHref: () => "/app/developer/specifications/create",
        },
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
    user: {
      root: {
        path: "user",
        getHref: () => "/app/user",
      },
      landingZone: {
        path: "landing-zone",
        getHref: () => "/app/user/landing-zone",
      },
      fileGroups: {
        root: {
          path: "file-groups",
          getHref: () => "/app/user/file-groups",
        },
        list: {
          path: "",
          getHref: () => "/app/user/file-groups",
        },
        view: {
          path: ":id",
          getHref: (id: string) => `/app/user/file-groups/${id}`,
        },
        create: {
          path: "create",
          getHref: () => "/app/user/file-groups/create",
        },
      },
    },
  },
} as const;
