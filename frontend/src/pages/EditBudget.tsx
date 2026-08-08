import { useParams } from "@solidjs/router";
import { createMemo, createResource, createSignal, ErrorBoundary, onMount, Show } from "solid-js";
import BudgetForm from "src/components/BudgetForm";
import { getBudgetById, submitEditBudget } from "src/stores/budgetStore";

export default function EditBudget() {
    const params = useParams();

    const id = +params.id!;

    const [budget] = createResource(() => id, () => getBudgetById(id));

    const data = createMemo(() => {
        const item = budget();
        if (!item) return null;
        return {
            id: item.id,
            name: item.name,
            description: item.description,
            target: item.target,
            tagId: item.tagId,
            frequency: item.frequency,
            interval: item.interval,
            dayOfWeek: item.dayOfWeek,
            dayOfMonth: item.dayOfMonth,
            startDate: item.startDate,
        };
    });

    return (
        <ErrorBoundary fallback={err => <p>Loading error: {err.message}</p>}>
            <Show when={budget()} fallback={<p>Loading</p>}>
                <BudgetForm
                    id={id}
                    name={data()!.name}
                    description={data()!.description}
                    target={data()!.target}
                    tagId={data()!.tagId}
                    frequency={data()!.frequency}
                    interval={data()!.interval}
                    dayOfWeek={data()!.dayOfWeek}
                    dayOfMonth={data()!.dayOfMonth}
                    startDate={new Date(data()!.startDate)}
                    submit={(name: string, description: string, target: string, tagId: number, frequency: string, interval: number, dayOfWeek: number, dayOfMonth: number, startDate: string) => {
                        return submitEditBudget(id, name, description, target, tagId, frequency, interval, dayOfWeek, dayOfMonth, startDate);
                    }}
                />
            </Show>
        </ErrorBoundary>
    );
}
