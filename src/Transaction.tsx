import './styles/Transaction.css';
import './styles/App.css';
// import { invoke } from "@tauri-apps/api/core";

interface TransactionProps {
    updateTab: (newTab: string) => void;
}

export default function Transaction(props: TransactionProps) {
    console.log(props);
    
    return (
        <div id="transaction">
            <h1>Transactions</h1>
        </div>
    );
}

