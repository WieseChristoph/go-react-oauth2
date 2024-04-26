import { type FC } from "react";

import { Outlet } from "@tanstack/react-router";

const Root: FC = () => {
  return (
    <>
      <main className="p-3">
        <Outlet />
      </main>
    </>
  );
};

export default Root;
