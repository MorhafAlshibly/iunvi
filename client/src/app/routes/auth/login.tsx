import { useNavigate, useSearchParams } from "react-router-dom";
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
      navigate(`${redirectTo ? `${redirectTo}` : paths.app.home.getHref()}`, {
        replace: true,
      });
    },
  });

  useEffect(() => {
    if (user.data) {
      navigate(paths.app.home.getHref(), { replace: true });
      return;
    }
    login.mutate({});
  }, [user.data]);

  return <>Login popup opened</>;
};

export default LoginRoute;
