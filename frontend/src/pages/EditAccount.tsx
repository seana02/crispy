import { useParams } from "@solidjs/router";
import { createMemo, createResource, createSignal, ErrorBoundary, onMount, Show } from "solid-js";
import AccountForm from "src/components/AccountForm";
import TransactionList from "src/components/TransactionList";
import { getAccountById, submitEditAccount } from "src/stores/accountStore";
import { fetchTransactionsByAccount } from "src/stores/transactionStore";
import { domain } from "wailsjs/go/models";

export default function EditAccount() {
    const params = useParams();

    const id = +params.id!;

    const [acct] = createResource(() => id, () => getAccountById(id));
    const [acctTransactions, setAcctTransactions] = createSignal<domain.TransactionDTO[]>([]);

    const data = createMemo(() => {
        const item = acct();
        if (!item) return null;
        return {
            parentID: item.parentID,
            name: item.name,
            type: item.type,
            currency: item.currency,
            description: item.description,
            active: item.active,
            dateCreated: item.dateCreated,
            dateUpdated: item.dateUpdated,
        };
    });

    const loadTransactions = async () => {
        const txs = await fetchTransactionsByAccount(id);
        setAcctTransactions(txs || []);
    }

    onMount(() => {
        loadTransactions();
    });

    return (
        <ErrorBoundary fallback={err => <p>Loading error: {err.message}</p>}>
            <Show when={acct()} fallback={<p>Loading</p>}>
                <AccountForm
                    id={id}
                    name={data()!.name}
                    type={data()!.type}
                    currency={data()!.currency}
                    description={data()!.description}
                    active={data()!.active}
                    dateCreated={data()!.dateCreated}
                    dateUpdated={data()!.dateUpdated}
                    submit={(name: string, description: string, type: domain.Type, currency: string, active: boolean) => {
                        return submitEditAccount(id, data()!.parentID, name, description, type, currency, active);
                    }}
                />
                <div style={{ "border-bottom": "2px solid var(--color-text-secondary)", margin: "8px 0" }} />
                <h2 class="form-transaction-list">Transactions</h2>
                <TransactionList
                    transactions={acctTransactions()}
                />
            </Show>
        </ErrorBoundary>
    );
}
