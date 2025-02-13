import { configureAuth } from "react-query-auth";
import { Navigate, useLocation } from "react-router";
import { User, ROLES } from "@/types/user";
import { instance } from "@/lib/msal";
import { paths } from "@/config/paths";
import { loginRequest } from "@/config/authentication";

// api call definitions for auth (types, schemas, requests):
// these are not part of features as this is a module shared across features

const getUser = async (): Promise<User> => {
  const user = instance.getActiveAccount()?.idTokenClaims;
  console.log("Got user", user);
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

const login = async (): Promise<User> => {
  await instance.loginPopup(loginRequest);
  console.log("Logged in");
  return getUser();
};

const logout = async (): Promise<void> => {
  await instance.logoutRedirect({
    postLogoutRedirectUri: paths.home.getHref(),
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
