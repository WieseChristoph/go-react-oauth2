import {
  createRootRoute,
  createRoute,
  createRouter,
} from "@tanstack/react-router";

import Root from "./Root";
import Home from "./Home";

const rootRoute = createRootRoute({
  component: Root,
});

const homeRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: Home,
});

const routeTree = rootRoute.addChildren([homeRoute]);

export default createRouter({ routeTree });
