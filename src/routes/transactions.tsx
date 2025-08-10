import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createRoute, Link, Outlet } from "@tanstack/react-router";
import { rootRoute } from "./__root";
import { TransactionData } from "../types";
import { delete_transaction, getTransactions } from "../api";
import { transactionDetailsRoute } from "./transactions_details";
import { transactionIndexRoute } from "./transactions_index";
import { transactionAddRoute } from "./transactions_add";
import { ConfirmationButton } from "@/components/ConfirmationButton";
import { MdDelete } from 'react-icons/md';
import { MdDeleteForever } from 'react-icons/md';

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

    return (
        <div className='h-full flex flex-1 bg-black text-white'>
            <div className='min-w-1/4 w-1/4 border-r-2 border-gray-600 box-border flex flex-col p-4 gap-2 overflow-y-scroll'>
                <div className="flex gap-2">
                    <TransactionLink link={'/transactions'} text={'Return'} />
                    <TransactionLink link={'add'} text={'Add'} />
                </div>
                {
                    isPending ? <div>Pending</div> : isError ? <div>{"Error: " + error.message}</div> :
                    Object.entries(data)?.map(([_, d]) => (
                            <div className="flex gap-2 h-[40px]">
                                <TransactionLink key={d.id} link={`${d.id}`} text={`(${d.transaction_date}) ${d.description}`} />
                                <ConfirmationButton
                                    className="h-[40px] w-[40px] box-border p-2 border-2 border-red-700 hover:bg-red-700 transition select-none flex items-center justify-center"
                                    onConfirm={() => {
                                        delete_transaction({ id: d.id });
                                        queryClient.invalidateQueries({ queryKey: ['transactions'] });
                                        queryClient.invalidateQueries({ queryKey: ['transaction_details', `${d.id}` ]})
                                    }}
                                    baseText={<MdDelete />}
                                    confirmationText={<MdDeleteForever />}
                                />
                            </div>
                        )).reverse()
                }
            </div>
            <div className='h-full flex-1'><Outlet /></div>
        </div>
    );
}

function TransactionLink({ text, link, className }: { text: string, link: string, className?: string }) {
    return (
        <Link to={link}
            className={`p-4 flex items-center transition border-2 border-gray-700 hover:bg-gray-700 flex-1 select-none overflow-hidden  whitespace-nowrap ${className}`}
            activeOptions={{ exact: true, includeSearch: false }}
            activeProps={{ className: 'bg-gray-800' }}
        >
            {text}
        </Link>
    );
}

transactionRoute.addChildren([
    transactionIndexRoute,
    transactionDetailsRoute,
    transactionAddRoute,
]);

