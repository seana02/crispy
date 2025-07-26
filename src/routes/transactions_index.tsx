import { createRoute } from "@tanstack/react-router";
import { transactionRoute } from "./transactions";

export const transactionIndexRoute = createRoute({
    getParentRoute: () => transactionRoute,
    path: '/',
    component: () => <div>Click on a transaction</div>,
});

