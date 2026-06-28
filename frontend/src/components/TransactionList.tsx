import { For, Show, JSX, createSignal } from "solid-js";
import { domain } from "wailsjs/go/models";
import "src/styles/lists.css";
import { useNavigate } from "@solidjs/router";
import { PencilLine, Trash2 } from "lucide-solid";
import ConfirmationPopup from "src/components/ConfirmationPopup";
import { deleteTransactionById } from "src/stores/transactionStore";

interface CustomButton {
    color: string;
    icon: JSX.Element;
    onClick: (id: number) => void;
}

interface TransactionListProps {
    transactions: domain.TransactionDTO[];
    customButton?: CustomButton;
}

export default function TransactionList(props: TransactionListProps) {
    const navigate = useNavigate();
    const [deleteConfirmation, setDeleteConfirmID] = createSignal(-1);

    return (
        <div class="list-box">
            <Show
                when={props.transactions.length > 0}
                fallback={<div class="list-row">No Transactions</div>}
            >
                <For each={props.transactions}>{(data, i) => {
                    const dateString = new Date(data.date).toLocaleDateString("en-CA");
                    return (
                        <div class="list-row">
                            <div class="list-row-date">{dateString}</div>
                            <div class="list-row-desc">{data.description}</div>
                            <div class="list-middle-gap"></div>
                            <div class="list-row-status">
                                {data.status == domain.Status.Cleared ? <></> : data.status}
                            </div>
                            <button
                                class="list-row-button list-row-edit"
                                onClick={() => navigate("/transactions/"+data.id)}
                            >
                                <PencilLine class="svg-small" />
                            </button>
                            <button
                                class="list-row-button list-row-delete"
                                onClick={() => setDeleteConfirmID(i())}
                            >
                                <Trash2 class="svg-small" />
                            </button>
                            <Show when={props.customButton}>
                                <button
                                    class="list-row-button"
                                    style={{
                                        "background-color": props.customButton!.color,
                                        "color": "var(--color-surface)"
                                    }}
                                    onClick={() => props.customButton!.onClick(data.id)}
                                >
                                    {props.customButton!.icon}
                                </button>
                            </Show>
                            <ConfirmationPopup isOpen={deleteConfirmation() == i()} onClose={() => setDeleteConfirmID(-1)}>
                                <div class="popup-text">Are you sure you want to delete the following transaction?</div>
                                <div style={{ height: "8px" }} />
                                <div class="popup-date">Date: {dateString}</div>
                                <div class="popup-description">{data.description}</div>
                                <div style={{ height: "12px" }} />
                                <div class="popup-button-wrapper">
                                    <button class="popup-button popup-delete" onClick={() => {
                                        deleteTransactionById(data.id);
                                        setDeleteConfirmID(-1);
                                    }}>Delete</button>
                                    <button class="popup-button popup-cancel" onClick={() => setDeleteConfirmID(-1)}>Cancel</button>
                                </div>
                            </ConfirmationPopup>
                        </div>
                    );
                }}</For>
            </Show>
        </div>
    );
}
