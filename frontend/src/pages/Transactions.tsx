import { useNavigate } from "@solidjs/router";
import { txList } from "src/stores/transactionStore";
import "src/styles/lists.css";
import TransactionList from "src/components/TransactionList";

export default function Transactions() {
    const navigate = useNavigate();

    return (
        <div id="transactions-page">
            <div class="filter-menu">
                <button class="add-btn" onClick={() => navigate("/transactions/add")}>New Transaction</button>
            </div>
            <TransactionList transactions={txList} />
        </div>
    );
}

