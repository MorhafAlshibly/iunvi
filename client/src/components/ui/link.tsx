import { Link as RouterLink, LinkProps } from "react-router-dom";

import { cn } from "@/utils/cn";

export const Link = ({ className, children, ...props }: LinkProps) => {
  return (
    <RouterLink
      className={cn("text-neutral-600 hover:text-neutral-900", className)}
      {...props}
    >
      {children}
    </RouterLink>
  );
};
