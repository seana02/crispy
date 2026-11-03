import { Show, For } from 'solid-js';
import './SelectInput.css';

interface SelectInputProps {
  value: string;
  onChange: (value: string) => void;
  options: Array<{ value: string; label: string }>;
  placeholder?: string;
  label?: string;
  error?: string;
  disabled?: boolean;
}

export function SelectInput(props: SelectInputProps) {
  return (
    <div class="select-root" classList={{ disabled: !!props.disabled, 'select-error': !!props.error }}>
      <Show when={props.label}><label class="select-label">{props.label}</label></Show>
      <select
        class="select-input"
        value={props.value}
        onChange={(e) => props.onChange(e.currentTarget.value)}
        disabled={props.disabled}
      >
        <Show when={props.placeholder}>
          <option value="">{props.placeholder}</option>
        </Show>
        <For each={props.options}>
          {(opt) => <option value={opt.value}>{opt.label}</option>}
        </For>
      </select>
      <Show when={props.error}><span class="select-error-msg">{props.error}</span></Show>
    </div>
  );
}
