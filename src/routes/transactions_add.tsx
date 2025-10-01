import { useNavigate, createRoute, useSearch } from "@tanstack/react-router";
import { transactionRoute } from "./transactions";
import { TransactionForm } from "@/components/transaction_form";
import { PostingData } from "@/types";
import { useQueryClient } from "@tanstack/react-query";
import { addTransaction } from "@/api";
import { z } from 'zod';

const postingSchema = z.object({
    id: z.number(),
    account: z.string(),
    value: z.string(),
    currency: z.string(),
    comment: z.string(),
})

const def: {
    id: number
    account: string
    value: string
    currency: string
    comment: string
}[] = []

export function getTodayDatestring() {
    let today = new Date();
    return today.getFullYear() + "-" + String(today.getMonth()+1).padStart(2, '0') + "-" + String(today.getDate()).padStart(2, '0');
}

export const transactionAddRoute = createRoute({
    getParentRoute: () => transactionRoute,
    path: '/add',
    validateSearch: z.object({
        date: z.string().optional().default(getTodayDatestring()),
        description: z.string().optional().default(""),
        postings: z.preprocess(x => {
            if (typeof x !== 'string') return [];
            try {
                return JSON.parse(x);
            } catch {
                return [];
            }
        }, z.array(postingSchema).default([])),
        duplicated: z.boolean().optional().default(false)
    }),
    component: TransactionAdd,
})

function TransactionAdd() {
    const queryClient = useQueryClient();
    const navigate = useNavigate();

    const search = useSearch({ from: transactionRoute.id })

    const submit = async (year: number, month: number, day: number, postings: PostingData[], desc: string, _: number[]) => {
        console.log('submitted');
        const addedID = await addTransaction({ year, month, day, postings, desc });
        await queryClient.invalidateQueries({ queryKey: ['transactions'] });
        navigate({ to: `../${addedID}` });
    };

    return (
        <TransactionForm
            date={search.date}
            description={search.description}
            postings={search.postings}
            onSubmit={submit}
            duplicated={search.duplicated}
        />
    );
}
