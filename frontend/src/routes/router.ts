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

const routeTree = rootRoute.addChildren([
  createRoute({
    getParentRoute: () => rootRoute,
    path: "/",
    component: Home,
  }),
]);

export default createRouter({ routeTree });
