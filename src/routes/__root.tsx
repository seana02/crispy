import { createRootRoute, Link, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import '../App.css';
import { indexRoute } from ".";
import { transactionRoute } from "./transactions";
import Heading from "@/components/Heading";
import { reportsRoute } from "./reports";

export const rootRoute = createRootRoute({
    component: () => (
        <div className="flex flex-col h-screen">
            <Heading />
            <div className="flex flex-1 min-h-0">
                <div className="text-white text-lg px-4 py-6 flex flex-col gap-2 w-60 bg-black border-r-2 border-gray-600 box-border">
                    <LinkWrapper text="Home" link="/" />
                    <LinkWrapper text="Reports" link="/reports" />
                    <LinkWrapper text="Transactions" link="/transactions" />
                </div>
                <div className="h-full flex flex-1 bg-black text-white">
                    <Outlet />
                </div>
            </div>
            <TanStackRouterDevtools />
            <ReactQueryDevtools />
        </div>
    ),
});

function LinkWrapper({ text, link }: { text: string, link: string }) {
    return (
        <Link to={link} className="[&.active]:bg-gray-500 hover:bg-gray-600 transition duration-75 p-4 rounded-2xl select-none">
            {text}
        </Link>
    );
}

export const routeTree = rootRoute.addChildren([
    indexRoute,
    reportsRoute,
    transactionRoute,
]);

