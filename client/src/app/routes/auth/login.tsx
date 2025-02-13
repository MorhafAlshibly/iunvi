import { useNavigate, useSearchParams } from "react-router";
import { paths } from "@/config/paths";
import { useEffect } from "react";
import { useLogin } from "@/lib/authentication";

const LoginRoute = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const redirectTo = searchParams.get("redirectTo");
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
    login.mutate({});
  }, []);

  return <>Login window opened</>;
};

export default LoginRoute;
