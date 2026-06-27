import type { JSX } from "solid-js";
import Sidebar from "src/components/layouts/Sidebar";
import Topbar from "src/components/layouts/Topbar";
import { RouteSectionProps } from "@solidjs/router";

export default function PageLayout(props: RouteSectionProps) {
    return (
        <div class="app-layout">
            <Sidebar />
            <div class="main-content">
                <Topbar />
                <main class="page-content">
                    {props.children}
                </main>
            </div>
        </div>
    );
}
