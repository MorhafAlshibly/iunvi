import {
  NavLink,
  useLocation,
  useNavigate,
  useNavigation,
} from "react-router-dom";
import { AppSidebar } from "../app-sidebar";
import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbSeparator,
  BreadcrumbPage,
} from "../ui/breadcrumb";
import { SidebarProvider, SidebarInset, SidebarTrigger } from "../ui/sidebar";
import { Separator } from "../ui/separator";
import { paths } from "@/config/paths";

export function DashboardLayout({ children }: { children: React.ReactNode }) {
  const path = useLocation().pathname;
  const navigate = useNavigate();
  // Search path object to get the current path
  let curr: Record<string, { [key: string]: any }> = paths.app;
  const pathObjects = [] as Array<{
    text: string;
    href: string;
  }>;
  const splitPath = path.replace(/^\/app\//, "").split("/");
  const pathObjectsLength = pathObjects.length;
  for (let i = 0; i < splitPath.length; i++) {
    const pathPart = splitPath[i];
    for (const key in curr) {
      const item = curr[key as keyof typeof curr];
      if (item.root) {
        if (item.root.path !== pathPart) continue;
      } else if (item.path !== pathPart) continue;
      // Convert hyphen casing to spaced text
      const pathName = pathPart
        .replace(/-/g, " ")
        .replace(/([a-z])([A-Z])/g, "$1 $2")
        .replace(/\b\w/g, (char) => char.toUpperCase());
      let href = "";
      if (i > 0) {
        href = `/${splitPath.slice(1, i + 1).join("/")}`;
      }
      pathObjects.push({
        text: pathName,
        href: href,
      });
      curr = item;
      break;
    }
    if (pathObjectsLength === pathObjects.length) {
      // If no match is found, break the loop
      break;
    }
  }

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
          <div className="flex items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <Breadcrumb>
              <BreadcrumbList>
                {pathObjects.map((item, index) => {
                  if (index === pathObjects.length - 1) {
                    return (
                      <BreadcrumbItem key={index}>
                        <BreadcrumbPage>{item.text}</BreadcrumbPage>
                      </BreadcrumbItem>
                    );
                  }
                  return (
                    <>
                      <BreadcrumbItem key={index} className="hidden md:block">
                        <BreadcrumbLink onClick={() => navigate(item.href)}>
                          {item.text}
                        </BreadcrumbLink>
                      </BreadcrumbItem>
                      <BreadcrumbSeparator className="hidden md:block" />
                    </>
                  );
                })}
              </BreadcrumbList>
            </Breadcrumb>
          </div>
        </header>
        <div className="flex flex-1 flex-col gap-4 p-4 pt-0">{children}</div>
      </SidebarInset>
    </SidebarProvider>
  );
}
