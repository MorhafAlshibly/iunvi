import { configureAuth } from "react-query-auth";
import { ActiveUser, ROLES } from "@/types/user";
import { instance } from "@/lib/msal";
import { paths } from "@/config/paths";
import { loginRequest } from "@/config/authentication";
import { Navigate, useLocation } from "react-router-dom";

// api call definitions for auth (types, schemas, requests):
// these are not part of features as this is a module shared across features

const getUser = async (): Promise<ActiveUser> => {
  const user = instance.getActiveAccount()?.idTokenClaims;
  if (!user) {
    throw new Error("User not found");
  }

  return {
    objectId: user.oid as string,
    tenantId: user.tid as string,
    displayName: user.name as string,
    username: user.preferred_username as string,
    role: user.roles?.includes("Admin") ? ROLES.ADMIN : ROLES.USER,
  };
};

const login = async (): Promise<ActiveUser> => {
  await instance.loginPopup(loginRequest);
  return getUser();
};

const logout = async (): Promise<void> => {
  await instance.logoutPopup({
    account: instance.getActiveAccount(),
    postLogoutRedirectUri: paths.landing.getHref(),
  });
};

const authConfig = {
  userFn: getUser,
  loginFn: login,
  registerFn: login,
  logoutFn: logout,
};

export const { useUser, useLogin, useLogout, useRegister, AuthLoader } =
  configureAuth(authConfig);

export const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const user = useUser();
  const location = useLocation();

  if (!user.data) {
    return (
      <Navigate to={paths.auth.login.getHref(location.pathname)} replace />
    );
  }

  return children;
};
