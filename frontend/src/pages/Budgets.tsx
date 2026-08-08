import { useNavigate } from "@solidjs/router";
import { createMemo, createSignal, For } from "solid-js";
import "src/styles/lists.css";
import "src/styles/budgets.css";
import { budgetList, deleteBudgetById } from "src/stores/budgetStore";
import { PencilLine, Trash2 } from "lucide-solid";
import ConfirmationPopup from "src/components/ConfirmationPopup";

export default function Budgets() {
    const navigate = useNavigate();
    const [deleteConfirmation, setDeleteConfirmID] = createSignal(-1);

    const getDaysUntil = (nextReset: string) => {
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        const resetDate = new Date(nextReset);
        resetDate.setHours(0, 0, 0, 0);
        const diffTime = resetDate.getTime() - today.getTime();
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        if (diffDays < 0) return "resets today";
        if (diffDays === 0) return "resets today";
        if (diffDays < 7) return `resets in ${diffDays}d`;

        const weeks = Math.floor(diffDays / 7);
        if (weeks < 4) return `resets in ${weeks}w`;

        const months = Math.floor(diffDays / 30);
        return `resets in ${months}mo`;
    };

    const getProgressColor = (currentSpend: string, target: string) => {
        const current = parseFloat(currentSpend);
        const targetNum = parseFloat(target);
        const pct = (current / targetNum) * 100;

        if (pct > 100) return "var(--color-accent-error)";
        if (pct > 80) return "var(--color-accent-warn)";
        return "var(--color-primary)";
    };

    const getProgressWidth = (currentSpend: string, target: string) => {
        const current = parseFloat(currentSpend);
        const targetNum = parseFloat(target);
        const pct = (current / targetNum) * 100;
        return `${Math.min(pct, 100)}%`;
    };

    return (
        <div id="budgets-page">
            <div class="filter-menu">
                <button class="add-btn" onClick={() => navigate("/budgets/add")}>
                    New Budget
                </button>
            </div>
            <div class="list-box">
                <For each={budgetList}>
                    {(budget, i) => {
                        const progressColor = getProgressColor(budget.currentSpend, budget.target);
                        const progressWidth = getProgressWidth(budget.currentSpend, budget.target);

                        return (
                            <div class="budget-block">
                                <div class="budget-header">
                                    <div class="list-row-name">{budget.name}</div>
                                    <div class="list-row-desc">{budget.description}</div>
                                    <div class="list-middle-gap"></div>
                                    <button
                                        class="list-row-button list-row-edit"
                                        onClick={() => navigate(`/budgets/${budget.id}`)}
                                    >
                                        <PencilLine class="svg-small" />
                                    </button>
                                    <button
                                        class="list-row-button list-row-delete"
                                        onClick={() => setDeleteConfirmID(i())}
                                    >
                                        <Trash2 class="svg-small" />
                                    </button>
                                    <ConfirmationPopup
                                        isOpen={deleteConfirmation() == i()}
                                        onClose={() => setDeleteConfirmID(-1)}
                                    >
                                        <div class="popup-text">
                                            Are you sure you want to delete the following Budget?
                                        </div>
                                        <div style={{ height: "8px" }} />
                                        <div class="popup-name">{budget.name}</div>
                                        <div class="popup-description">{budget.description}</div>
                                        <div style={{ height: "12px" }} />
                                        <div class="popup-button-wrapper">
                                            <button
                                                class="popup-button popup-delete"
                                                onClick={() => {
                                                    deleteBudgetById(budget.id);
                                                    setDeleteConfirmID(-1);
                                                }}
                                            >
                                                Delete
                                            </button>
                                            <button
                                                class="popup-button popup-cancel"
                                                onClick={() => setDeleteConfirmID(-1)}
                                            >
                                                Cancel
                                            </button>
                                        </div>
                                    </ConfirmationPopup>
                                </div>
                                <div class="budget-progress-row">
                                    <div class="budget-progress-track">
                                        <div
                                            class="budget-progress-fill"
                                            style={{
                                                width: progressWidth,
                                                backgroundColor: progressColor,
                                            }}
                                        />
                                    </div>
                                    <span>
                                        {Math.min(
                                            (parseFloat(budget.currentSpend) /
                                                parseFloat(budget.target)) *
                                                100,
                                            100,
                                        ).toFixed(0)}
                                        %
                                    </span>
                                </div>
                                <div class="budget-meta">
                                    <span>
                                        ${parseFloat(budget.currentSpend).toFixed(2)} / $
                                        {parseFloat(budget.target).toFixed(2)}
                                    </span>
                                    <span>{getDaysUntil(budget.nextReset)}</span>
                                </div>
                            </div>
                        );
                    }}
                </For>
            </div>
        </div>
    );
}
