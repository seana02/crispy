import { createSignal, onMount } from "solid-js";

type Theme = "professional-dark" | "sleek-neutral-dark" | "sleek-neutral-light"
const [theme, setTheme] = createSignal<Theme>("sleek-neutral-dark");

onMount(() => {
  const saved = localStorage.getItem("theme") as Theme | null;
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
