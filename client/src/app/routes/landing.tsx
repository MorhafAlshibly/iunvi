import { useNavigate } from "react-router-dom";

import { Button } from "@/components/ui/button";
import { paths } from "@/config/paths";
import { useIsAuthenticated } from "@azure/msal-react";

const LandingRoute = () => {
  const navigate = useNavigate();
  const isAuthenticated = useIsAuthenticated();

  const handleStart = () => {
    if (isAuthenticated) {
      navigate(paths.app.home.getHref());
    } else {
      navigate(paths.auth.login.getHref());
    }
  };

  return (
    <>
      <div className="flex h-screen items-center">
        <div className="mx-auto max-w-7xl px-4 py-12 text-center sm:px-6 lg:px-8 lg:py-16">
          <h2 className="text-3xl font-extrabold tracking-tight sm:text-4xl">
            <span className="block">iunvi</span>
          </h2>
          <p>Data Ingestion And Modelling Self-Service</p>
          <div className="mt-8 flex justify-center">
            <div className="inline-flex rounded-md shadow">
              <Button onClick={handleStart}>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="size-6"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                  />
                </svg>
                Get started
              </Button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default LandingRoute;
