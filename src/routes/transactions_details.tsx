import { useQuery, useQueryClient } from "@tanstack/react-query";
import { createRoute, useRouter } from "@tanstack/react-router";
import { transactionRoute } from "./transactions";
import { PostingData, TransactionData } from "../types";
import { getTransactionDetails, updateTransaction } from "../api";
import { TransactionForm } from "@/components/transaction_form";

export const transactionDetailsRoute = createRoute({
    getParentRoute: () => transactionRoute,
    path: '/$id',
    component: TransactionDetails,
});

function TransactionDetails() {
    const queryClient = useQueryClient();

    const { id } = transactionDetailsRoute.useParams();
    const { isPending, isError, data, error } = useQuery<TransactionData>({
        queryKey: ['transaction_details', id],
        queryFn: () => getTransactionDetails({ id: +id }),
        staleTime: Infinity,
    });

    if (isPending) return <div>Pending...</div>
    if (isError) return <div>Error: {error.message}</div>

    const submit = (year: number, month: number, day: number, postings: PostingData[], desc: string, deleteList: number[]) => {
        updateTransaction({ id: +id, year, month, day, postings, desc, deleteList });
        queryClient.invalidateQueries({ queryKey: ['transaction_details', `${id}`] })
        queryClient.invalidateQueries({ queryKey: ['transactions'] })
    };

    return (
        <TransactionForm
            key={data.id}
            date={data.transaction_date}
            description={data.description}
            postings={data.postings}
            onSubmit={submit}
        />
    );
}
