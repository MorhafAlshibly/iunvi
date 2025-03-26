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
        root: {
          path: "models",
          getHref: () => "/app/developer/models",
        },
        list: {
          path: "",
          getHref: () => "/app/developer/models",
        },
        view: {
          path: ":id",
          getHref: (id: string) => `/app/developer/models/${id}`,
        },
        create: {
          path: "create",
          getHref: () => "/app/developer/models/create",
        },
      },
      dashboards: {
        root: {
          path: "dashboards",
          getHref: () => "/app/developer/dashboards",
        },
        list: {
          path: "",
          getHref: () => "/app/developer/dashboards",
        },
        view: {
          path: ":id",
          getHref: (id: string) => `/app/developer/dashboards/${id}`,
        },
        create: {
          path: "create",
          getHref: () => "/app/developer/dashboards/create",
        },
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
      runModels: {
        path: "run-models",
        getHref: () => "/app/user/run-models",
      },
    },
    viewer: {
      root: {
        path: "viewer",
        getHref: () => "/app/viewer",
      },
      runHistory: {
        path: "run-history",
        getHref: () => "/app/viewer/run-history",
      },
      dashboard: {
        path: "dashboard/:id",
        getHref: (id: string) => `/app/viewer/dashboard/${id}`,
      },
    },
  },
} as const;
