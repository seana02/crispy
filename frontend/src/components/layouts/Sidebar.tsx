import { createSignal, For } from "solid-js";
import { A, useLocation } from "@solidjs/router";
import "src/styles/sidebar.css";
import { Banknote, House, Menu, PiggyBank, Settings, Tag, Wallet } from "lucide-solid";

export default function Sidebar() {
    const [collapsed, setCollapsed] = createSignal(false);
    const location = useLocation();

    const nav = [
        { href: "/", label: "Dashboard", icon: <House color="var(--color-text-primary)" /> },
        { href: "/transactions", label: "Transactions", icon: <Banknote color="var(--color-text-primary)" /> },
        { href: "/budgets", label: "Budgets", icon: <PiggyBank color="var(--color-text-primary)" /> },
        { href: "/accounts", label: "Accounts", icon: <Wallet color="var(--color-text-primary)" /> },
        { href: "/tags", label: "Tags", icon: <Tag color="var(--color-text-primary)" /> },
        { href: "/settings", label: "Settings", icon: <Settings color="var(--color-text-primary)" /> },
    ];

    return (
        <aside class={`sidebar ${collapsed() ? "collapsed" : ""}`}>
            <div class="sidebar-header">
                <button class="collapse-btn" onClick={() => setCollapsed(!collapsed())}>
                    <Menu color="var(--color-text-primary)" />
                </button>
            </div>
            <nav class="sidebar-nav">
                <For each={nav}>
                    {(item) => (
                        <A
                            href={item.href}
                            classList={{ active: location.pathname === item.href }}
                        >
                            {item.icon} <span>{item.label}</span>
                        </A>
                    )}
                </For>
            </nav>
        </aside>
    );
}
