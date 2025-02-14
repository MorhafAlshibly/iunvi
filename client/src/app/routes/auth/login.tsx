import { useNavigate, useSearchParams } from "react-router";
import { paths } from "@/config/paths";
import { useEffect } from "react";
import { useLogin, useUser } from "@/lib/authentication";

const LoginRoute = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const redirectTo = searchParams.get("redirectTo");
  const user = useUser();
  const login = useLogin({
    onSuccess: () => {
      navigate(
        `${redirectTo ? `${redirectTo}` : paths.app.dashboard.getHref()}`,
        {
          replace: true,
        },
      );
    },
  });

  useEffect(() => {
    if (user.data) {
      navigate(paths.app.dashboard.getHref(), { replace: true });
      return;
    }
    login.mutate({});
  }, [user.data]);

  return <>Login popup opened</>;
};

export default LoginRoute;
