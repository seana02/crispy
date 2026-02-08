import { createEffect, createMemo, createResource, createSignal, ErrorBoundary, For, Suspense } from "solid-js";
import "../styles/newItem.css";
import { createStore } from "solid-js/store";
import { submitEditTransaction, getTransactionById } from "../stores/transactionStore";
import { domain } from "../../wailsjs/go/models";
import { useNavigate, useParams } from "@solidjs/router";
import TransactionForm from "src/components/TransactionForm";
import { CurrencyIcon } from "lucide-solid";
import { UpdateTransaction } from "wailsjs/go/main/App";

interface EditTransactionProps {
    reference?: string
    dateCreated?: string
    dateUpdated?: string
}

export default function EditTransaction(props: EditTransactionProps) {
    const params = useParams();
    const navigate = useNavigate();

    const id = +params.id;

    const source = createMemo(() => getTransactionById(id));

    const [tx] = createResource(() => id, () => getTransactionById(id));

    const data = createMemo(() => {
        const item = tx();
        if (!item) return null;
        return {
            description: item.description,
            date: new Date(item.date),
            status: item.status,
            tags: item.tags?.join(' '),
            postings: item.postings,
            reference: item.referenceID,
            dateCreated: item.dateCreated,
            dateUpdated: item.dateUpdated
        };
    });

    return (
        <ErrorBoundary fallback={err => <p>Loading error: {err.message}</p>}>
            <Suspense fallback={<p>Loading...</p>}>
                <TransactionForm
                    id={+params.id}
                    description={data()!.description}
                    date={data()!.date}
                    status={data()!.status}
                    tags={data()!.tags}
                    postings={data()!.postings}
                    reference={"" + data()!.reference}
                    dateCreated={data()!.dateCreated}
                    dateUpdated={data()!.dateUpdated}
                    submit={(desc: string, date: Date, status: domain.Status, tags: string[], postings: { id: number, accountID: number, account: string, amount: string, currency: string }[]) => {
                        let postingsDTO = postings.map(p => domain.PostingDTO.createFrom({
                            id: p.id,
                            accountID: p.accountID,
                            accountName: p.account,
                            amount: p.amount,
                            currency: p.currency,
                            dateCreated: new Date(),
                            dateUpdated: new Date(),
                            transactionID: id,
                        }));
                        return submitEditTransaction(id, desc, date, status, tags, postingsDTO);
                    }}
                />
            </Suspense>
        </ErrorBoundary>
    );
}

