import { createSignal, createEffect, createMemo, on, For } from "solid-js"
import { Combobox as ArkCombobox, createListCollection } from "@ark-ui/solid"
import 'src/styles/combobox.css';

function debounce<T extends (...args: any[]) => void>(fn: T, ms: number): T {
    let timer: ReturnType<typeof setTimeout>
    return ((...args: Parameters<T>) => {
        clearTimeout(timer)
        timer = setTimeout(() => fn(...args), ms)
    }) as T
}

interface ComboboxProps {
    fetchOptions: (input: string) => Promise<string[]>
    value: string
    onChange: (value: string) => void
    placeholder?: string
    debounceMs?: number
}

export default function Combobox(props: ComboboxProps) {
    const [inputValue, setInputValue] = createSignal(props.value ?? "")
    const [items, setItems] = createSignal<string[]>([])

    const collection = createMemo(() =>
        createListCollection({ items: items() })
    )

    // Sync parent value → input (e.g. parent resets the field)
    createEffect(() => {
        setInputValue(props.value ?? "")
    })

    const debouncedFetch = debounce(async (query: string) => {
        if (query === "") return;
        try {
            const results = await props.fetchOptions(query)
            setItems([query, ...results])
        } catch (err) {
            console.error("Combobox: fetchOptions failed", err)
            setItems([])
        }
    }, props.debounceMs ?? 300)

    // Fetch whenever the input changes, skip first run
    createEffect(
        on(inputValue, (val) => debouncedFetch(val), { defer: true })
    )

    function commit(val: string) {
        props.onChange(val)
    }

    return (
        <ArkCombobox.Root
            collection={collection()}
            inputValue={inputValue()}
            value={[props.value]}
            allowCustomValue
            onInputValueChange={(details) => {
                setInputValue(details.inputValue)
            }}
            // Fired when an item is selected from the dropdown
            onValueChange={(details) => {
                const picked = details.value[0] ?? ""
                setInputValue(picked)
                commit(picked)
            }}
        >
            <ArkCombobox.Label>Account</ArkCombobox.Label>
            <ArkCombobox.Control>
                <ArkCombobox.Input
                    placeholder={props.placeholder ?? "Type to search…"}
                    onKeyDown={(e) => {
                        // Commit on Enter whether or not an item is highlighted
                        if (e.key === "Enter") {
                            commit(inputValue())
                        }
                    }}
                    onBlur={() => {
                        // Commit on blur; also resets input to committed value
                        // in case the user typed something but didn't press Enter
                        commit(inputValue())
                    }}
                />
            </ArkCombobox.Control>
            <ArkCombobox.Positioner>
                <ArkCombobox.Content>
                    <For each={items()}>
                        {(item) => (
                            <ArkCombobox.Item item={item}>
                                <ArkCombobox.ItemText>{item}</ArkCombobox.ItemText>
                                <ArkCombobox.ItemIndicator>✓</ArkCombobox.ItemIndicator>
                            </ArkCombobox.Item>
                        )}
                    </For>
                </ArkCombobox.Content>
            </ArkCombobox.Positioner>
        </ArkCombobox.Root>
    )
}
