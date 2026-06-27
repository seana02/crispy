import { useNavigate } from "@solidjs/router"
import { createSignal } from "solid-js";
import 'src/styles/form.css';

interface TagProps {
    id?: number
    name?: string
    submit: (name: string) => Promise<void>
}

export default function TagForm(props: TagProps) {
    const navigate = useNavigate();

    const [name, setName] = createSignal(props.name || "");
    const [error, setError] = createSignal("");

    async function handleSubmit(e: Event) {
        e.preventDefault();
        if (name().length <= 0) {
            setError("Name canno be empty");
            return;
        }
        props.submit(name()).then(() => navigate("/tags"));
    }

    return (
        <form class="new-item-form" onSubmit={handleSubmit}>
            <h2 class="form-title">{props.id === -1 ? "New" : "Edit"} Tag</h2>

            <div class="form-group-group">
                <div class="form-group grow">
                    <label for="tag-name-input">Name</label>
                    <input
                        id="tag-name-input"
                        type="string"
                        value={name()}
                        onInput={(e) => setName(e.currentTarget.value)}
                    />
                </div>
            </div>

            {error() && <p class="form-error">{error()}</p>}

            <button type="submit" class="submit-btn">Submit</button>
        </form>
    );
}
