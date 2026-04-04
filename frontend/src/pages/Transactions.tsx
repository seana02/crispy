import { createSignal, For } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { deleteTransactionById, txList } from "../stores/transactionStore";
import { domain } from "../../wailsjs/go/models";
import "../styles/lists.css";
import { PencilLine, Trash2 } from "lucide-solid";
import ConfirmationPopup from "src/components/ConfirmationPopup";

export default function Transactions() {
    const navigate = useNavigate();
    const [deleteConfirmation, setDeleteConfirmID] = createSignal(-1);

    return (
        <div id="transactions-page">
            <div class="filter-menu">
                <button class="add-btn" onClick={() => navigate("/transactions/add")}>New Transaction</button>
            </div>
            <div class="list-box">
                <For each={txList}>{(data,i) => {
                    let dateString = new Date(data.date).toLocaleDateString("en-CA");
                    return (
                        <div class="list-row">
                            <div class="list-row-date">{dateString}</div>
                            <div class="list-row-desc">{data.description}</div>
                            { data.status == domain.Status.Cleared ? <></> : <div class="list-row-status">{data.status}</div> }
                            <button class="list-row-button list-row-edit" onClick={() => navigate("/transactions/"+data.id)}>
                                <PencilLine class="svg-small" />
                            </button>
                                <button class="list-row-button list-row-delete" onClick={() => setDeleteConfirmID(i())}>
                                    <Trash2 class="svg-small" />
                                </button>
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
            </div>
        </div>
    );
}

