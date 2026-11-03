import { createSignal, createMemo, Show, For } from 'solid-js';
import './TagInput.css';

interface TagInputProps {
  value: Array<{ key: string; value: string }>;
  onChange: (tags: Array<{ key: string; value: string }>) => void;
  existingTags?: Array<{ key: string; value: string }>;
  placeholder?: string;
  label?: string;
  error?: string;
  disabled?: boolean;
}

export function TagInput(props: TagInputProps) {
  const [inputValue, setInputValue] = createSignal('');
  const [showAutocomplete, setShowAutocomplete] = createSignal(false);

  const existingTags = () => props.existingTags ?? [];

  const parsedParts = createMemo(() => {
    const val = inputValue();
    const colonIndex = val.lastIndexOf(':');
    if (colonIndex === -1) return { key: val, value: '' };
    return { key: val.slice(0, colonIndex), value: val.slice(colonIndex + 1) };
  });

  const autocompleteOptions = createMemo(() => {
    if (!showAutocomplete()) return [];
    const key = parsedParts().key.toLowerCase();
    if (!key) return [];
    return existingTags().filter(t => t.key.toLowerCase().includes(key) || t.value.toLowerCase().includes(key)).slice(0, 5);
  });

  const addTag = (key: string, value: string) => {
    if (key.trim() && value.trim()) {
      props.onChange([...props.value, { key: key.trim(), value: value.trim() }]);
      setInputValue('');
      setShowAutocomplete(false);
    }
  };

  const removeTag = (index: number) => {
    props.onChange(props.value.filter((_, i) => i !== index));
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Enter' || e.key === ',') {
      e.preventDefault();
      const { key, value } = parsedParts();
      addTag(key, value);
    } else if (e.key === 'Backspace' && !inputValue() && props.value.length > 0) {
      removeTag(props.value.length - 1);
    }
  };

  return (
    <div class="tag-input-root" classList={{ disabled: !!props.disabled }}>
      <Show when={props.label}><label class="tag-input-label">{props.label}</label></Show>
      <div class="tag-input-chips">
        <For each={props.value}>
          {(tag, i) => (
            <span class="tag-chip-item">
              {tag.key}:{tag.value}
              <button class="tag-chip-remove" onClick={() => removeTag(i())} type="button">✕</button>
            </span>
          )}
        </For>
        <input
          class="tag-input-field"
          type="text"
          value={inputValue()}
          onInput={(e) => { setInputValue(e.currentTarget.value); setShowAutocomplete(true); }}
          onKeyDown={handleKeyDown}
          onBlur={() => setTimeout(() => setShowAutocomplete(false), 150)}
          placeholder={props.placeholder ?? 'key:value'}
          disabled={props.disabled}
        />
      </div>
      <Show when={showAutocomplete() && autocompleteOptions().length > 0}>
        <div class="tag-autocomplete">
          <For each={autocompleteOptions()}>
            {(tag) => (
              <button class="tag-autocomplete-item" onClick={() => addTag(tag.key, tag.value)} type="button">
                {tag.key}:{tag.value}
              </button>
            )}
          </For>
        </div>
      </Show>
      <span class="tag-input-hint">Type key:value and press Enter</span>
      <Show when={props.error}><span class="tag-input-error-msg">{props.error}</span></Show>
    </div>
  );
}
