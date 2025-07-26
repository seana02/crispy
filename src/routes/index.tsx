import { createRoute } from "@tanstack/react-router"; import { invoke } from "@tauri-apps/api/core";
import { useState } from "react";
import reactLogo from "../assets/react.svg";
import "../App.css";
import { rootRoute } from "./__root";

export const indexRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: '/',
    component: App,
});

function App() {
    const [greetMsg, setGreetMsg] = useState("");
    const [name, setName] = useState("");

    async function greet() {
        // Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
        setGreetMsg(await invoke("greet", { name }));
    }

    return (
        <div className="flex-1 bg-black"></div>
    );
}
