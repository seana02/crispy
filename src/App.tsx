import { useState } from "react";
import "./styles/App.css";
import Sidebar from "./Sidebar";
import Transaction from "./Transaction";
import TransactionEdit from "./TransactionEdit";
import { TransactionData } from "./types";

function App() {
    let [[tab, data], setTab] = useState<[string, TransactionData | null]>(["Transactions", null]);

    // async function greet() {
    //     // Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
    // }
    
    let content;
    switch(tab) {
        case "Transactions": content = <Transaction updateTab={setTab}/>; break;
        case "TransactionEdit": {
            if (data) {
                content = <TransactionEdit updateTab={setTab} data={data}/>;
            } else {
                content = <TransactionEdit updateTab={setTab}/>;
            }
            break;
        }
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
