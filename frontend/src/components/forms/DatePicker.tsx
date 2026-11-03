import { Show } from 'solid-js';
import './DatePicker.css';

interface DatePickerProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
  error?: string;
  disabled?: boolean;
  min?: string;
  max?: string;
}

export function DatePicker(props: DatePickerProps) {
  return (
    <div class="date-picker-root" classList={{ disabled: !!props.disabled, 'date-picker-error': !!props.error }}>
      <Show when={props.label}><label class="date-picker-label">{props.label}</label></Show>
      <input
        class="date-picker-input"
        type="date"
        value={props.value}
        onInput={(e) => props.onChange(e.currentTarget.value)}
        disabled={props.disabled}
        min={props.min}
        max={props.max}
      />
      <Show when={props.error}><span class="date-picker-error-msg">{props.error}</span></Show>
    </div>
  );
}
