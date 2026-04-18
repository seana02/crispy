import { createSignal, For } from "solid-js";
import "../styles/newItem.css";
import { createStore } from "solid-js/store";
import { domain } from "../../wailsjs/go/models";
import { useNavigate } from "@solidjs/router";
import { Currency } from "src/stores/currencyStore";
import { getAccountIdByName } from "src/stores/accountStore";

interface TransactionProps {
    id?: number
    description?: string
    date?: Date
    status?: domain.Status
    tags?: string
    postings?: domain.PostingDTO[]
    reference?: string
    dateCreated?: string
    dateUpdated?: string
    submit: (description: string, date: Date, status: domain.Status, tags: string[], postings: { id: number, accountID: number, account: string, amount: string, currency: string }[]) => Promise<void>
}

export default function TransactionForm(props: TransactionProps) {
    const navigate = useNavigate();
    const [description, setDescription] = createSignal(props.description || "");
    const [date, setDate] = createSignal(props.date || new Date())
    const [status, setStatus] = createSignal(props.status || domain.Status.Pending);
    const [tags, setTags] = createSignal(props.tags || "");
    const [postings, setPostings] = createStore(
        props.postings ? props.postings.map(p => ({ id: p.id || -1, accountID: p.accountID, account: p.accountName, amount: p.amount, currency: p.currency })) :
        [{ id: -1, accountID: -1, account: "", amount: "", currency: "USD" }]
    );
    const [error, setError] = createSignal("");

    async function handleSubmit(e: Event) {
        let failed = false;
        e.preventDefault();
        if (description().length <= 0) {
            setError("Description cannot be empty");
            failed = true;
            return;
        }
        if (postings.length <= 0) {
            setError("Postings cannot be empty");
            failed = true;
            return;
        }
        if (postings.reduce((x, p) => x + +p.amount, 0) != 0) {
            setError("Postings must be balanced");
            failed = true;
            return;
        }
        const accountIDPromises = postings.map(async (p, i) => {
            const newID = await getAccountIdByName(p.account);
            if (newID == -1) {
                setError("Account \"" + p.account + "\" does not exist");
            failed = true;
            }
            setPostings(i, "accountID", newID);
        })
        await Promise.all(accountIDPromises);
        if (!failed) {
            props.submit(description(), date(), status(), tags().split(" "), postings)
                .then(() => navigate("/transactions"));
        }
    }

    return (
        <form class="new-item-form" onSubmit={handleSubmit}>
            <h2 class="form-title">{props.id == -1 ? "New" : "Edit"} Transaction</h2>

            <div class="form-group-group">
                <div class="form-group">
                    <label for="transaction-date-input">Date</label>
                    <input
                        id="transaction-date-input"
                        type="date"
                        value={date().toLocaleDateString("en-CA")}
                        onInput={(e) => setDate(new Date(e.currentTarget.value+"T00:00"))}
                    />
                </div>
                <div class="form-group">
                    <label for="transaction-status-input">Status</label>
                    <div class="select-wrapper">
                        <select
                            id="transaction-status-input"
                            value={status()}
                            onChange={e => setStatus(e.currentTarget.value as domain.Status)}
                        >
                            <For each={Object.values(domain.Status)} >
                                {s => <option value={s}>{s}</option>}
                            </For>
                        </select>
                    </div>
                </div>
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
                <div class="form-group">
                    <label for="transaction-reference-input">Reference</label>
                    <input
                        id="transaction-reference-input"
                        type="text"
                        value={props.reference === "null" ? "" : props.reference || ""}
                        disabled
                    />
                </div>
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

            <div class="form-group">
                <label for="transaction-tags-input">Tags</label>
                <input
                    id="transaction-tags-input"
                    type="text"
                    value={tags()}
                    onInput={(e) => setTags(e.currentTarget.value)}
                />
            </div>

            <h3 style={{ margin: 0 }}>Postings</h3>

            <For each={postings}>
                {(p, i) =>
                    <div class="form-posting-row">
                        <div class="form-group posting-account-input">
                            <label for={"transaction-account-input-" + i}>Account</label>
                            <input
                                id={"transaction-account-input" + i}
                                type="text"
                                value={p.account}
                                onInput={(e) => setPostings(i(), "account", e.target.value)}
                            />
                        </div>
                        <div class="form-group posting-currency-input">
                            <label for={"transaction-currency-input" + i}>Currency</label>
                            <div class="select-wrapper">
                                <select
                                    id={"transaction-currency-input" + i}
                                    class="transaction-currency-input"
                                    value={p.currency}
                                    onChange={e => setPostings(i(), "currency", e.currentTarget.value)}
                                >
                                    <For each={Object.values(Currency)} >
                                        {s => <option value={s}>{s}</option>}
                                    </For>
                                </select>
                            </div>
                        </div>
                        <div class="form-group posting-amount-input">
                            <label for={"transaction-amount-input" + i}>Amount</label>
                            <input
                                id={"transaction-amount-input" + i}
                                class="transaction-amount-input"
                                type="text"
                                inputMode="decimal"
                                value={p.amount}
                                onInput={(e) => setPostings(i(), "amount", e.target.value)}
                            />
                        </div>
                    </div>
                }
            </For>

            {error() && <p class="form-error">{error()}</p>}

            <button class="add-posting-btn" onClick={e => {
                e.preventDefault();
                setPostings([...postings, { id: -1, accountID: -1, account: "", amount: "", currency: "USD" }])
            }}>Add Posting</button>
            <button type="submit" class="submit-btn">Submit</button>
        </form>
    );
}
