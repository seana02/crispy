import { useNavigate, createRoute } from "@tanstack/react-router";
import { transactionRoute } from "./transactions";
import { TransactionForm } from "@/components/transaction_form";
import { PostingData } from "@/types";
import { useQueryClient } from "@tanstack/react-query";
import { addTransaction } from "@/api";

export const transactionAddRoute = createRoute({
    getParentRoute: () => transactionRoute,
    path: '/add',
    component: TransactionAdd,
})

function TransactionAdd() {
    const queryClient = useQueryClient();
    const navigate = useNavigate();

    const submit = async (year: number, month: number, day: number, postings: PostingData[], desc: string, _: number[]) => {
        console.log('submitted');
        const addedID = await addTransaction({ year, month, day, postings, desc });
        await queryClient.invalidateQueries({ queryKey: ['transactions'] });
        navigate({ to: `../${addedID}` });
    };

    let todayDate = new Date();

    return (
        <TransactionForm
            date={todayDate.getFullYear() + "-" + String(todayDate.getMonth()+1).padStart(2, '0') + "-" + String(todayDate.getDate()).padStart(2, '0')}
            description={""}
            postings={[]}
            onSubmit={submit}
        />
    );
}
