// import Palette from "lucide-solid/icons/palette";
import { Palette } from "lucide-solid";
import { theme, toggleTheme } from "../../stores/themeStore";
import "../../styles/topbar.css";

export default function Topbar() {
    return (
        <header class="topbar">
            <div class="topbar-left">
                <div class="title">Crispy</div>
                <input type="text" placeholder="Search..." class="topbar-search" />
            </div>
            <div class="topbar-right">
                <button class="theme-btn" onClick={toggleTheme}>
                    <Palette color="var(--color-text-primary)" size={16} />
                </button>
                <button class="add-btn">＋ Add</button>
            </div>
        </header>
    );
}
