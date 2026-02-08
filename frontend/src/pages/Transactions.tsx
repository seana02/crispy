import { For } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { txList } from "../stores/transactionStore";
import { domain } from "../../wailsjs/go/models";
import "../styles/lists.css";
import { PencilLine, Trash2 } from "lucide-solid";

export default function Transactions() {
    const navigate = useNavigate();

    // let elements = [
    //     "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "l", "m",
    //     "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "l", "m",
    //     "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "l", "m",
    //     "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "l", "m",
    //     "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "l", "m",
    //     "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "l", "m",
    //     "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "l", "m",
    //     "a", "b", "c", "d", "e", "f", "g", "h", "j", "k", "l", "m",
    // ];
    return (
        <div id="transactions-page">
            <div class="filter-menu">
                <button class="add-btn" onClick={() => navigate("/transactions/add")}>New Transaction</button>
            </div>
            <div class="list-box">
                <For each={txList}>{e => <TransactionRow data={e} />}</For>
            </div>
        </div>
    );
}

function TransactionRow(props: { data: domain.TransactionDTO }) {
    const navigate = useNavigate();

    return (
        <div class="list-row">
            <div class="list-row-date">{new Date(props.data.date).toLocaleDateString("en-CA")}</div>
            <div class="list-row-desc">{props.data.description}</div>
            { props.data.status == domain.Status.Cleared ? <></> : <div class="list-row-status">{props.data.status}</div> }
            <button class="list-row-button list-row-edit" onClick={() => navigate("/transactions/"+props.data.id)}>
                <PencilLine color="var(--color-surface)" class="svg-small" />
            </button>
            <button class="list-row-button list-row-delete">
                <Trash2 color="var(--color-surface)" class="svg-small" />
            </button>
        </div>
    );
}
