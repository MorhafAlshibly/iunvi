import { Button } from "@/components/ui/button";
import { Link } from "@/components/ui/link";
import { paths } from "@/config/paths";

const SpecificationsListRoute = () => {
  return (
    <>
      <div className="flex justify-end">
        <Link to={paths.app.developer.specifications.create.getHref()}>
          Create Specification
        </Link>
      </div>
      <div>
        <div className="p-4">
          {/* {workspaces.map((workspace) => (
            <>
              <div key={workspace.id} className="flex text-sm">
                <span className="flex-1 content-center">{workspace.name}</span>
                <span className="flex-1 text-right">
                  <EditWorkspace
                    onSubmit={(name) => editWorkspaceFn(workspace.id, name)}
                  />
                </span>
              </div>
              <Separator className="my-2" />
            </>
          ))} */}
        </div>
      </div>
    </>
  );
};

export default SpecificationsListRoute;
