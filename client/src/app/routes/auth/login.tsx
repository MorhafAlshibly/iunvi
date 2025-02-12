import { useNavigate, useSearchParams } from "react-router";
import { useMsal, useIsAuthenticated } from "@azure/msal-react";
import { paths } from "@/config/paths";
import { useEffect } from "react";

const LoginRoute = () => {
  const navigate = useNavigate();
  const { instance } = useMsal();
  const isAuthenticated = useIsAuthenticated();

  useEffect(() => {
    if (isAuthenticated) {
      navigate(paths.app.dashboard.getHref(), { replace: true });
    }
    const authenticate = async () => {
      await instance.initialize();
      console.log(import.meta.env.VITE_AZUREADSCOPE);
      await instance.loginRedirect({
        scopes: [import.meta.env.VITE_AZUREADSCOPE as string],
      });
    };
    authenticate();
  }, [isAuthenticated, instance, navigate]);

  return <></>;
};

export default LoginRoute;
