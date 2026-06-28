import { useParams } from "@solidjs/router";
import { createMemo, createResource, createSignal, ErrorBoundary, Show, onMount } from "solid-js";
import TagForm from "src/components/TagForm";
import TransactionList from "src/components/TransactionList";
import ConfirmationPopup from "src/components/ConfirmationPopup";
import { getTagById, submitEditTag } from "src/stores/tagStore";
import { fetchTransactionsByTag, removeTagFromTransaction } from "src/stores/transactionStore";
import { domain } from "wailsjs/go/models";
import { Unlink } from "lucide-solid";
import 'src/styles/form.css';


export default function EditTag() {
    const params = useParams();

    const id = +params.id!;

    const [tag] = createResource(() => id, () => getTagById(id));
    const [tagTransactions, setTagTransactions] = createSignal<domain.TransactionDTO[]>([]);
    const [removeTagConfirmation, setRemoveTagConfirmation] = createSignal<number>(-1);

    const data = createMemo(() => {
        const item = tag();
        if (!item) return null;
        return {
            id: item.id,
            name: item.name
        };
    });

    const loadTransactions = async () => {
        const txs = await fetchTransactionsByTag(id);
        setTagTransactions(txs || []);
    };

    onMount(() => {
        loadTransactions();
    });

    return (
        <ErrorBoundary fallback={err => <p>Loading error: {err.message}</p>}>
            <Show when={tag()} fallback={<p>Loading</p>}>
                <TagForm
                    id={id}
                    name={data()!.name}
                    submit={(name) => submitEditTag(id, name)}
                />
                <div style={{ "border-bottom": "2px solid var(--color-text-secondary)", margin: "8px 0" }} />
                <h2 class="form-transaction-list">Transactions</h2>
                <TransactionList
                    transactions={tagTransactions()}
                    customButton={{
                        color: "var(--color-accent-info)",
                        icon: <Unlink class="svg-small" />,
                        onClick: (txId) => setRemoveTagConfirmation(txId)
                    }}
                />
                <Show when={tagTransactions().length > 0} fallback={<></>}>

                    {tagTransactions().map((tx) => (
                        <ConfirmationPopup
                            isOpen={removeTagConfirmation() === tx.id}
                            onClose={() => setRemoveTagConfirmation(-1)}
                        >
                            <div class="popup-text">Are you sure you want to remove this tag from the following transaction?</div>
                            <div style={{ height: "8px" }} />
                            <div class="popup-date">Date: {new Date(tx.date).toLocaleDateString("en-CA")}</div>
                            <div class="popup-description">{tx.description}</div>
                            <div style={{ height: "12px" }} />
                            <div class="popup-button-wrapper">
                                <button class="popup-button popup-delete" onClick={() => {
                                    removeTagFromTransaction(tx.id, id);
                                    setRemoveTagConfirmation(-1);
                                    loadTransactions();
                                }}>Remove Tag</button>
                                <button class="popup-button popup-cancel" onClick={() => setRemoveTagConfirmation(-1)}>Cancel</button>
                            </div>
                        </ConfirmationPopup>
                    ))}
                </Show>
            </Show>
        </ErrorBoundary>
    );
}
