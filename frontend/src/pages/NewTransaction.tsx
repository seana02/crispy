import "../styles/newItem.css";
import { submitNewTransaction } from "../stores/transactionStore";
import TransactionForm from "src/components/TransactionForm";

export default function NewTransaction() {
    return (
        <TransactionForm
            id={-1}
            submit={submitNewTransaction}
        />
    );
}
