import { useNavigate } from "@solidjs/router";
import { domain } from "../../wailsjs/go/models";
import { createSignal, For } from "solid-js";
import { Currency } from "src/stores/currencyStore";
import "../styles/form.css";
import Select from "src/components/ark/Select";
import Checkbox from "src/components/ark/Checkbox";

interface AccountProps {
    id?: number
    parentID?: number
    name?: string
    type?: domain.Type
    currency?: string
    description?: string
    active?: boolean
    dateCreated?: string
    dateUpdated?: string
    submit: (name: string, description: string, type: domain.Type, currency: string, active: boolean) => Promise<void>
}

export default function AccountForm(props: AccountProps) {
    const navigate = useNavigate();

    const [name, setName] = createSignal(props.name || "")
    const [description, setDescription] = createSignal(props.description || "");
    const [type, setType] = createSignal(props.type || domain.Type.Asset)
    const [currency, setCurrency] = createSignal(props.currency || "USD");
    const [active, setActive] = createSignal(props.active ?? true);
    const [error, setError] = createSignal("");

    async function handleSubmit(e: Event) {
        e.preventDefault();
        if (name().length <= 0) {
            setError("Description cannot be empty");
            return;
        }
        props.submit(name(), description(), type(), currency(), active())
            .then(() => navigate("/accounts"));
    }

    return (
        <form class="new-item-form" onSubmit={handleSubmit}>
            <h2 class="form-title">{props.id === -1 ? "New" : "Edit"} Account</h2>

            <div class="form-group-group">
                <div class="form-group grow">
                    <label for="account-name-input">Name</label>
                    <input
                        id="account-name-input"
                        type="string"
                        value={name()}
                        onInput={(e) => setName(e.currentTarget.value)}
                    />
                </div>
            </div>

            <div class="form-group-group">
                <div class="form-group grow">
                    <label for="transaction-description-input">Description</label>
                    <input
                        id="transaction-description-input"
                        type="text"
                        value={description()}
                        onInput={(e) => setDescription(e.currentTarget.value)}
                        placeholder="e.g., Grocery shopping"
                    />
                </div>
            </div>

            <div class="form-group-group">
                <Select
                    label="Type"
                    value={type()}
                    list={Object.values(domain.Type).map(t => t.toString())}
                    onChange={e => setType(e.value[0] as domain.Type)}
                />
                <Select
                    label="Currency"
                    value={currency()}
                    list={Object.values(Currency).map(c => c.toString())}
                    onChange={e => setCurrency(e.value[0])}
                />
                <Checkbox
                    value={active()}
                    onChange={b => setActive(b)}
                />
                <div class="form-group">
                    <label for="transaction-date-created">Date Created</label>
                    <input
                        id="transaction-date-created"
                        type="text"
                        value={props.dateCreated || ""}
                        disabled
                    />
                </div>
                <div class="form-group">
                    <label for="transaction-date-updated">Date Updated</label>
                    <input
                        id="transaction-date-updated"
                        type="text"
                        value={props.dateUpdated || ""}
                        disabled
                    />
                </div>
            </div>

            {error() && <p class="form-error">{error()}</p>}

            <button type="submit" class="submit-btn">Submit</button>
        </form>
    );
}
