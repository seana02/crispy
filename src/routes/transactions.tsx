import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createRoute, Link, Outlet } from "@tanstack/react-router";
import { rootRoute } from "./__root";
import { TransactionData } from "../types";
import { getTransactions } from "../api";
import { transactionDetailsRoute } from "./transactions_details";
import { transactionIndexRoute } from "./transactions_index";

export const transactionRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: 'transactions',
    component: Transactions,
});

function Transactions() {
    const queryClient = useQueryClient();

    const { isPending, isError, data, error } = useQuery<TransactionData[]>({
        queryKey: ['transactions'],
        queryFn: getTransactions
    });

    if (!isPending && !isError) Object.entries(data)?.forEach(([_, d]) => console.log(d));

    return (
        <div className='h-full flex flex-1 bg-black text-white'>
            <div className='w-1/4 border-r-2 border-gray-600 box-border flex flex-col p-4 gap-2 overflow-y-scroll'>
                <TransactionLink link={'/transactions'} text={'Return'} />
                {
                    isPending ? <div>Pending</div> : isError ? <div>{"Error: " + error.message}</div> :
                    Object.entries(data)?.map(([_, d]) => (<TransactionLink key={d.id} link={`${d.id}`} text={d.description} />)).reverse()
                }
            </div>
            <div className='h-full flex-1'><Outlet /></div>
        </div>
    );
}

function TransactionLink({ text, link }: { text: string, link: string }) {
    return (
        <Link to={link} className='p-4 transition border-2 border-gray-700 hover:bg-gray-700'>
            {text}
        </Link>
    );
}

transactionRoute.addChildren([
    transactionIndexRoute,
    transactionDetailsRoute,
]);

