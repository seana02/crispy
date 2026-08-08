import { useNavigate } from "@solidjs/router";
import { createSignal, Show } from "solid-js";
import { tagList } from "src/stores/tagStore";
import "src/styles/form.css";
import Select from "src/components/ark/Select";
import Combobox from "src/components/ark/Combobox";
import DatePicker from "src/components/ark/DatePicker";

interface BudgetFormProps {
    id?: number;
    name?: string;
    description?: string;
    target?: string;
    tagId?: number;
    frequency?: string;
    interval?: number;
    dayOfWeek?: number;
    dayOfMonth?: number;
    startDate?: Date;
    submit: (
        name: string,
        description: string,
        target: string,
        tagId: number,
        frequency: string,
        interval: number,
        dayOfWeek: number,
        dayOfMonth: number,
        startDate: string,
    ) => Promise<void>;
}

export default function BudgetForm(props: BudgetFormProps) {
    const navigate = useNavigate();

    const [name, setName] = createSignal(props.name || "");
    const [description, setDescription] = createSignal(props.description || "");
    const [target, setTarget] = createSignal(props.target || "");
    const [tagName, setTagName] = createSignal(
        tagList.find((t) => t.id === props.tagId)?.name || ""
    );
    const [frequency, setFrequency] = createSignal(props.frequency || "Monthly");
    const [interval, setInterval] = createSignal(props.interval || 1);
    const [dayOfWeek, setDayOfWeek] = createSignal(props.dayOfWeek || 0);
    const [dayOfMonth, setDayOfMonth] = createSignal(props.dayOfMonth || 1);
    const [startDate, setStartDate] = createSignal(props.startDate || new Date());
    const [error, setError] = createSignal("");

    const daysOfWeek = [
        "Sunday",
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
    ];

    const frequencyDisplayMap: Record<string, string> = {
        Daily: "Day",
        Weekly: "Week",
        Monthly: "Month",
    };
    const displayToFrequency: Record<string, string> = {
        Day: "Daily",
        Week: "Weekly",
        Month: "Monthly",
    };
    const s = () => interval() > 1 ? "s" : "";
    const frequencyDisplayList = () => ["Day", "Week", "Month"].map(d => `${d}${s()}`);
    const frequencyDisplay = () => `${frequencyDisplayMap[frequency()] || "Month"}${s()}`;

    const handleSubmit = async (e: Event) => {
        e.preventDefault();
        setError("");

        if (!name() || name().trim() === "") {
            setError("Name cannot be empty");
            return;
        }

        const targetNum = parseFloat(target());
        if (isNaN(targetNum) || targetNum <= 0) {
            setError("Target must be a positive number");
            return;
        }

        try {
            const resolvedTagId = tagList.find((t) => t.name === tagName())?.id ?? -1;
            await props.submit(
                name(),
                description(),
                target(),
                resolvedTagId,
                frequency(),
                interval(),
                dayOfWeek(),
                dayOfMonth(),
                startDate().toISOString().split("T")[0],
            );
            navigate("/budgets");
        } catch (err) {
            console.error(err);
            setError("Failed to save budget");
        }
    };

    return (
        <form class="new-item-form" onSubmit={handleSubmit}>
            <h2 class="form-title">{props.id ? "Edit" : "New"} Budget</h2>

            <div class="form-group-group">
                <div class="form-group grow">
                    <label for="budget-name-input">Name</label>
                    <input
                        id="budget-name-input"
                        type="text"
                        value={name()}
                        onInput={(e) => setName(e.currentTarget.value)}
                    />
                </div>
            </div>

            <div class="form-group-group">
                <div class="form-group grow">
                    <label for="budget-description-input">Description</label>
                    <input
                        id="budget-description-input"
                        type="text"
                        value={description()}
                        onInput={(e) => setDescription(e.currentTarget.value)}
                        placeholder="e.g., Groceries monthly"
                    />
                </div>
            </div>

            <div class="form-group-group">
                <div class="form-group grow">
                    <Combobox
                        label="Tag"
                        value={tagName()}
                        fetchOptions={async (input) =>
                            tagList
                                .filter((t) => t.name.toLowerCase().includes(input.toLowerCase()))
                                .map((t) => t.name)
                        }
                        onChange={(val) => setTagName(val)}
                    />
                </div>
                <div class="form-group">
                    <label for="budget-target-input">Target Amount</label>
                    <input
                        id="budget-target-input"
                        type="text"
                        value={target()}
                        onInput={(e) => setTarget(e.currentTarget.value)}
                    />
                </div>
            </div>

            <div class="form-group-group">
                <div class="form-group">
                    <label>Reset Frequency</label>
                    <div style="display: flex; gap: 0.5rem; align-items: center;">
                        <span>Reset every</span>
                        <input
                            type="text"
                            value={interval()}
                            onInput={(e) => setInterval(Number(e.currentTarget.value))}
                            style="width: 60px;"
                        />
                        <Select
                            label=""
                            value={frequencyDisplay()}
                            list={frequencyDisplayList()}
                            onChange={(e) => setFrequency(displayToFrequency[e.value[0].replace(/s$/, "")] || "Monthly")}
                        />
                        <Show when={frequency() === "Weekly"}>
                            <span>on</span>
                            <Select
                                label=""
                                value={daysOfWeek[dayOfWeek()] || ""}
                                list={daysOfWeek}
                                onChange={(e) =>
                                    setDayOfWeek(daysOfWeek.indexOf(e.value[0]))
                                }
                            />
                        </Show>
                        <Show when={frequency() === "Monthly"}>
                            <span>on day</span>
                            <input
                                type="text"
                                value={dayOfMonth()}
                                onInput={(e) => setDayOfMonth(Number(e.currentTarget.value))}
                                style="width: 50px;"
                            />
                        </Show>
                        <span>starting from</span>
                        <DatePicker
                            label=" "
                            value={startDate()}
                            onChange={(e) => {
                                setStartDate(new Date(e.value[0].year, e.value[0].month - 1, e.value[0].day));
                            }}
                        />
                    </div>
                </div>
            </div>

            {error() && <p class="form-error">{error()}</p>}

            <button type="submit" class="submit-btn">Submit</button>
        </form>
    );
}
