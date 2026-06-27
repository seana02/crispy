import { createSignal, For } from "solid-js";
import { useNavigate } from "@solidjs/router";
import "src/styles/lists.css";
import { acctList, deleteAccountById } from "src/stores/accountStore";
import { PencilLine, Trash2 } from "lucide-solid";
import ConfirmationPopup from "src/components/ConfirmationPopup";

export default function Accounts() {
    const navigate = useNavigate();
    const [deleteConfirmation, setDeleteConfirmID] = createSignal(-1);

    return (
        <div id="accounts-page">
            <div class="filter-menu">
                <button class="add-btn" onClick={() => navigate("/accounts/add")}>New Account</button>
            </div>
            <div class="list-box">
                <For each={acctList}>{(data,i) => {
                    return (
                        <div class="list-row">
                            <div class="list-row-name">{data.name}</div>
                            <div class="list-row-desc">{data.description}</div>
                            <div class="list-middle-gap"></div>
                            <button class="list-row-button list-row-edit" onClick={() => navigate("/accounts/"+data.id)}>
                                <PencilLine class="svg-small" />
                            </button>
                            <button class="list-row-button list-row-delete" onClick={() => setDeleteConfirmID(i())}>
                                <Trash2 class="svg-small" />
                            </button>
                            <ConfirmationPopup isOpen={deleteConfirmation() == i()} onClose={() => setDeleteConfirmID(-1)}>
                                <div class="popup-text">Are you sure you want to delete the following Account?</div>
                                <div style={{ height: "8px" }} />
                                <div class="popup-name">{data.name}</div>
                                <div class="popup-description">{data.description}</div>
                                <div style={{ height: "12px" }} />
                                <div class="popup-button-wrapper">
                                    <button class="popup-button popup-delete" onClick={() => {
                                        deleteAccountById(data.id);
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
