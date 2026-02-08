import { For } from "solid-js";
import { useNavigate } from "@solidjs/router";
import "../styles/lists.css";

export default function Accounts() {
    const navigate = useNavigate();

    return (
        <div id="accounts-page">
            <div class="filter-menu">
                <button class="add-btn" onClick={() => navigate("/accounts/add")}>New Account</button>
            </div>
            <div class="list-box">
                <For each={["a", "b"]}>{e => <TransactionRow data={e} />}</For>
            </div>
        </div>
    );
}

function TransactionRow(props: { data: string }) {
    return (
        <div class="list-row">
            {props.data}
        </div>
    );
}
