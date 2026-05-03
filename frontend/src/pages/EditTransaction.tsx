import { createMemo, createResource, ErrorBoundary, Suspense } from "solid-js";
import { submitEditTransaction, getTransactionById } from "../stores/transactionStore";
import { domain } from "../../wailsjs/go/models";
import { useParams } from "@solidjs/router";
import TransactionForm from "src/components/TransactionForm";

export default function EditTransaction() {
    const params = useParams();

    const id = +params.id!;

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
                    id={id}
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

