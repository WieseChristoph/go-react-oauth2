import { Outlet } from "@tanstack/react-router";

const Root: React.FC = () => {
  return (
    <>
      <main className="p-3">
        <Outlet />
      </main>
    </>
  );
};

export default Root;
