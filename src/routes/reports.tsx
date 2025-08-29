import { createRoute, Link, Outlet } from "@tanstack/react-router";
import { rootRoute } from "./__root";
import { reportsBalanceSheetRoute } from "./reports_balancesheet";

export const reportsRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: 'reports',
    component: Reports,
})

function Reports() {
    return (
        <>
            <div className="min-w-1/8 w-1/8 border-r-2 border-gray-600 box-border flex flex-col p-4 gap-2 overflow-y-scroll">
                <ReportLink link={"bs"} text={"Balance Sheet"} />
            </div>
            <div className="h-full flex-1"><Outlet /></div>
        </>
    );
}

function ReportLink({ text, link, className }: { text: string, link: string, className?: string}) {
    return (
        <Link to={link}
            className={`px-4 py-2 flex items-center transition border-2 border-gray-700 hover:bg-gray-700 select-none overflow-hidden whitespace-nowrap ${className}`}
            activeOptions={{ exact: true, includeSearch: false }}
            activeProps={{ className: "bg-gray-800" }}
        >
            {text}
        </Link>
    );
}

reportsRoute.addChildren([
    reportsBalanceSheetRoute,
]);
