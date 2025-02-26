import { QueryClient, useQueryClient } from "@tanstack/react-query";
import { useMemo } from "react";
import { createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router/dom";

import {
  default as AppRoot,
  ErrorBoundary as AppRootErrorBoundary,
} from "./routes/app/root";

import { paths } from "@/config/paths";
import { ProtectedRoute } from "@/lib/authentication";
import { POLICIES } from "@/lib/authorization";

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
      path: paths.home.path,
      lazy: () => import("./routes/landing").then(convert(queryClient)),
    },
    {
      path: paths.auth.login.path,
      lazy: () => import("./routes/auth/login").then(convert(queryClient)),
    },
    {
      path: paths.app.root.path,
      element: (
        <ProtectedRoute>
          <AppRoot />
        </ProtectedRoute>
      ),
      ErrorBoundary: AppRootErrorBoundary,
      children: [
        {
          path: paths.app.admin.root.path,
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
          children: [
            {
              path: paths.app.developer.registry.path,
              lazy: () =>
                import("./routes/app/developer/registry").then(
                  convert(queryClient),
                ),
            },
            {
              path: paths.app.developer.models.path,
              lazy: () =>
                import("./routes/app/developer/models").then(
                  convert(queryClient),
                ),
            },
          ],
        },
        {
          path: paths.app.dashboard.path,
          lazy: () =>
            import("./routes/app/dashboard").then(convert(queryClient)),
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
