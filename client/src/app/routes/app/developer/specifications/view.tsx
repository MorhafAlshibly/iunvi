const SpecificationsViewRoute = () => {
  return (
    <>
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

export default SpecificationsViewRoute;
