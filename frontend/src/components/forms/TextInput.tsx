import { createSignal, createEffect, on, Show } from 'solid-js';
import './TextInput.css';

interface TextInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  label?: string;
  error?: string;
  disabled?: boolean;
}

export function TextInput(props: TextInputProps) {
  const [localValue, setLocalValue] = createSignal(props.value);

  createEffect(on(() => props.value, (val) => setLocalValue(val)));

  const handleInputChange = (e: Event) => {
    const val = (e.target as HTMLInputElement).value;
    setLocalValue(val);
    props.onChange(val);
  };

  return (
    <div class="text-input-root" classList={{ disabled: !!props.disabled, 'text-input-error': !!props.error }}>
      <Show when={props.label}><label class="text-input-label">{props.label}</label></Show>
      <input
        class="text-input"
        type="text"
        value={localValue()}
        onInput={handleInputChange}
        placeholder={props.placeholder}
        disabled={props.disabled}
      />
      <Show when={props.error}><span class="text-input-error-msg">{props.error}</span></Show>
    </div>
  );
}
