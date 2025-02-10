import { useState } from "react";
import { invoke } from "@tauri-apps/api/core";
import "./App.css";

function App() {
    let [tab, setTab] = useState("overview");

    async function greet() {
        // Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
        setGreetMsg(await invoke("greet", { name }));
    }

    switch(tab) {
        case "overview": <Overview />
    }

    return (
        <div id="root">
            <div className="sidebar">

            </div>
        </div>
    );

}

export default App;
