import { createSignal, onMount } from "solid-js";

type Themes = "professional-dark" | "sleek-neutral-dark" | "sleek-neutral-light"
const [theme, setTheme] = createSignal<Themes>("sleek-neutral-dark");

onMount(() => {
  const saved = localStorage.getItem("theme") as Themes | null;
  if (saved) setTheme(saved);
  document.documentElement.setAttribute("data-theme", theme());
});

export function toggleTheme() {
  const next = theme() === "sleek-neutral-dark" ? "professional-dark" : "sleek-neutral-dark";
  setTheme(next);
  document.documentElement.setAttribute("data-theme", next);
  localStorage.setItem("theme", next);
}

export { theme };
