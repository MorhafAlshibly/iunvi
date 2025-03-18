import { QueryClient, useQueryClient } from "@tanstack/react-query";
import { useMemo } from "react";
import { createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router/dom";

import { paths } from "@/config/paths";

import {
  default as AppRoot,
  ErrorBoundary as AppRootErrorBoundary,
} from "./routes/app/root";

import {
  default as AppDeveloperRoot,
  ErrorBoundary as AppDeveloperRootErrorBoundary,
} from "./routes/app/developer/root";

import {
  default as AppDeveloperSpecificationsRoot,
  ErrorBoundary as AppDeveloperSpecificationsRootErrorBoundary,
} from "./routes/app/developer/specifications/root";

import {
  default as AppAdminRoot,
  ErrorBoundary as AppAdminRootErrorBoundary,
} from "./routes/app/admin/root";

import {
  default as AppUserRoot,
  ErrorBoundary as AppUserRootErrorBoundary,
} from "./routes/app/user/root";

import {
  default as AppUserFileGroupsRoot,
  ErrorBoundary as AppUserFileGroupsRootErrorBoundary,
} from "./routes/app/user/file-groups/root";

import {
  default as AppDeveloperModelsRoot,
  ErrorBoundary as AppDeveloperModelsRootErrorBoundary,
} from "./routes/app/developer/models/root";

const convert = (queryClient: QueryClient) => (m: any) => {
  const { clientLoader, clientAction, default: Component, ...rest } = m;
  return {
    ...rest,
    loader: clientLoader?.(queryClient),
    action: clientAction?.(queryClient),
    Component,
  };
};

export const createAppRouter = (queryClient: QueryClient) =>
  createBrowserRouter([
    {
      path: paths.landing.path,
      lazy: () => import("./routes/landing").then(convert(queryClient)),
    },
    {
      path: paths.auth.login.path,
      lazy: () => import("./routes/auth/login").then(convert(queryClient)),
    },
    {
      path: paths.app.root.path,
      element: <AppRoot />,
      ErrorBoundary: AppRootErrorBoundary,
      children: [
        {
          path: paths.app.home.path,
          lazy: () => import("./routes/app/home").then(convert(queryClient)),
        },
        {
          path: paths.app.admin.root.path,
          element: <AppAdminRoot />,
          ErrorBoundary: AppAdminRootErrorBoundary,
          children: [
            {
              path: paths.app.admin.workspaces.path,
              lazy: () =>
                import("./routes/app/admin/workspaces").then(
                  convert(queryClient),
                ),
            },
            {
              path: paths.app.admin.users.path,
              lazy: () =>
                import("./routes/app/admin/users").then(convert(queryClient)),
            },
          ],
        },
        {
          path: paths.app.developer.root.path,
          element: <AppDeveloperRoot />,
          ErrorBoundary: AppDeveloperRootErrorBoundary,
          children: [
            {
              path: paths.app.developer.specifications.root.path,
              element: <AppDeveloperSpecificationsRoot />,
              ErrorBoundary: AppDeveloperSpecificationsRootErrorBoundary,
              children: [
                {
                  path: paths.app.developer.specifications.list.path,
                  lazy: () =>
                    import("./routes/app/developer/specifications/list").then(
                      convert(queryClient),
                    ),
                },
                {
                  path: paths.app.developer.specifications.view.path,
                  lazy: () =>
                    import("./routes/app/developer/specifications/view").then(
                      convert(queryClient),
                    ),
                },
                {
                  path: paths.app.developer.specifications.create.path,
                  lazy: () =>
                    import("./routes/app/developer/specifications/create").then(
                      convert(queryClient),
                    ),
                },
              ],
            },
            {
              path: paths.app.developer.registry.path,
              lazy: () =>
                import("./routes/app/developer/registry").then(
                  convert(queryClient),
                ),
            },
            {
              path: paths.app.developer.models.root.path,
              element: <AppDeveloperModelsRoot />,
              ErrorBoundary: AppDeveloperModelsRootErrorBoundary,
              children: [
                {
                  path: paths.app.developer.models.list.path,
                  lazy: () =>
                    import("./routes/app/developer/models/list").then(
                      convert(queryClient),
                    ),
                },
                {
                  path: paths.app.developer.models.view.path,
                  lazy: () =>
                    import("./routes/app/developer/models/view").then(
                      convert(queryClient),
                    ),
                },
                {
                  path: paths.app.developer.models.create.path,
                  lazy: () =>
                    import("./routes/app/developer/models/create").then(
                      convert(queryClient),
                    ),
                },
              ],
            },
          ],
        },
        {
          path: paths.app.user.root.path,
          element: <AppUserRoot />,
          ErrorBoundary: AppUserRootErrorBoundary,
          children: [
            {
              path: paths.app.user.landingZone.path,
              lazy: () =>
                import("./routes/app/user/landing-zone").then(
                  convert(queryClient),
                ),
            },
            {
              path: paths.app.user.fileGroups.root.path,
              element: <AppUserFileGroupsRoot />,
              ErrorBoundary: AppUserFileGroupsRootErrorBoundary,
              children: [
                {
                  path: paths.app.user.fileGroups.list.path,
                  lazy: () =>
                    import("./routes/app/user/file-groups/list").then(
                      convert(queryClient),
                    ),
                },
                {
                  path: paths.app.user.fileGroups.view.path,
                  lazy: () =>
                    import("./routes/app/user/file-groups/view").then(
                      convert(queryClient),
                    ),
                },
                {
                  path: paths.app.user.fileGroups.create.path,
                  lazy: () =>
                    import("./routes/app/user/file-groups/create").then(
                      convert(queryClient),
                    ),
                },
              ],
            },
          ],
        },
      ],
    },
    {
      path: "*",
      lazy: () => import("./routes/not-found").then(convert(queryClient)),
    },
  ]);

export const AppRouter = () => {
  const queryClient = useQueryClient();

  const router = useMemo(() => createAppRouter(queryClient), [queryClient]);

  return <RouterProvider router={router} />;
};
