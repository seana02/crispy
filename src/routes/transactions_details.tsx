import { useQuery } from "@tanstack/react-query";
import { createRoute } from "@tanstack/react-router";
import { transactionRoute } from "./transactions";
import { PostingData, TransactionData } from "../types";
import { getTransactionDetails } from "../api";
import { FieldApi, useForm } from "@tanstack/react-form";
import { useEffect } from "react";
import { TransactionForm } from "@/components/transaction_form";

export const transactionDetailsRoute = createRoute({
    getParentRoute: () => transactionRoute,
    path: '/$id',
    component: TransactionDetails,
});

function TransactionDetails() {
    const formDefault = {
        date: new Date(),
        description: "",
        postings: [] as PostingData[],
    }
    const transactionForm = useForm({
        defaultValues: formDefault,
        onSubmit: async ({ value }) => console.log(value)
    });

    const { id } = transactionDetailsRoute.useParams();
    const { isPending, isError, data, error } = useQuery<TransactionData>({
        queryKey: ['transaction_details', id],
        queryFn: () => getTransactionDetails({ id: +id }),
    });

    useEffect(() => {
        if (data) {
            transactionForm.reset({
                date: data.transaction_date,
                description: data.description,
                postings: data.postings,
            });
        }
    }, [data]);

    if (isPending) return <div>Pending...</div>
    if (isError) return <div>Error: {error.message}</div>

    return (
        <TransactionForm
            date={data.transaction_date}
            description={data.description}
            postings={data.postings}
            onSubmit={() => console.log("submitted")}
        />
    );
}
