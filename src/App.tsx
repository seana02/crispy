import { useState } from "react";
import "./styles/App.css";
import Sidebar from "./Sidebar";
import Transaction from "./Transaction";
import TransactionEdit from "./TransactionEdit";

function App() {
    let [tab, setTab] = useState("TransactionEdit");

    // async function greet() {
    //     // Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
    // }

    let content;
    switch(tab) {
        case "Transactions": content = <Transaction updateTab={setTab}/>; break;
        case "TransactionEdit": content = <TransactionEdit updateTab={setTab}/>; break;
        case "Overview":
        default: content = <div className="overview">Overview</div>; break;
    }

    return (
        <div id="root" className="lovelace">
            <Sidebar activeTab={tab} updateTab={setTab} />
            <div id="main-content">
                {content}
            </div>
        </div>
    );

}

export default App;
