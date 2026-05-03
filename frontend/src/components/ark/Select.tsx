import { Select as ArkSelect, SelectValueChangeDetails, createListCollection } from '@ark-ui/solid'
import { ChevronsUpDownIcon, XIcon } from 'lucide-solid'
import { createSignal } from 'solid-js'
import { Index, Portal } from 'solid-js/web'
import '../../styles/select.css';

interface Item {
    label: string
    value: string
    disabled?: boolean
}

interface SelectProps {
    onChange: (details: SelectValueChangeDetails) => void;
    label: string
    list: string[]
    value: string;
}

export default function Select(props: SelectProps) {
    const collection = createListCollection<Item>({ items: props.list.map(item => ({ label: item, value: item })) });

    return (
        <ArkSelect.Root collection={collection} value={props.value ? [props.value] : []} onValueChange={props.onChange}>
            <ArkSelect.Label>{props.label}</ArkSelect.Label>
            <ArkSelect.Control>
                <ArkSelect.Trigger>
                    <ArkSelect.ValueText />
                    <ArkSelect.Indicator>
                        <ChevronsUpDownIcon />
                    </ArkSelect.Indicator>
                </ArkSelect.Trigger>
            </ArkSelect.Control>
            <Portal>
                <ArkSelect.Positioner>
                    <ArkSelect.Content>
                        <Index each={collection.items}>
                            {(item) => (
                                <ArkSelect.Item item={item()}>
                                    <ArkSelect.ItemText>{item().label}</ArkSelect.ItemText>
                                    <ArkSelect.ItemIndicator>✓</ArkSelect.ItemIndicator>
                                </ArkSelect.Item>
                            )}
                        </Index>
                    </ArkSelect.Content>
                </ArkSelect.Positioner>
            </Portal>
            <ArkSelect.HiddenSelect />
        </ArkSelect.Root>
    )
}
