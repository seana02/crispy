import { createSignal, createMemo, Show, For } from 'solid-js';
import './AutocompleteInput.css';

interface AutocompleteInputProps {
  value: string;
  onChange: (value: string) => void;
  options: string[];
  placeholder?: string;
  label?: string;
  error?: string;
  disabled?: boolean;
  onSelect?: (value: string) => void;
  highlightMatch?: boolean;
}

export function AutocompleteInput(props: AutocompleteInputProps) {
  const [isFocused, setIsFocused] = createSignal(false);
  const [isOpen, setIsOpen] = createSignal(false);
  const [highlightedIndex, setHighlightedIndex] = createSignal(-1);

  const filteredOptions = createMemo(() => {
    const query = props.value.toLowerCase();
    if (!query) return [];
    return props.options.filter((opt) => opt.toLowerCase().includes(query));
  });

  const handleInputChange = (e: Event) => {
    const target = e.target as HTMLInputElement;
    props.onChange(target.value);
    setIsOpen(true);
    setHighlightedIndex(-1);
  };

  const handleSelect = (value: string) => {
    props.onChange(value);
    setIsOpen(false);
    props.onSelect?.(value);
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setHighlightedIndex(prev => (prev + 1) >= filteredOptions().length ? 0 : prev + 1);
      setIsOpen(true);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setHighlightedIndex(prev => (prev - 1) < 0 ? filteredOptions().length - 1 : prev - 1);
      setIsOpen(true);
    } else if (e.key === 'Enter') {
      if (highlightedIndex() >= 0 && filteredOptions().length > 0) {
        e.preventDefault();
        handleSelect(filteredOptions()[highlightedIndex()]);
        setIsOpen(false);
      }
    } else if (e.key === 'Escape') {
      setIsOpen(false);
    }
  };

  const handleBlur = () => {
    setIsFocused(false);
    setTimeout(() => setIsOpen(false), 150);
  };

  return (
    <div
      class="autocomplete-root"
      classList={{ 'autocomplete-focused': isFocused(), 'autocomplete-open': isOpen(), 'autocomplete-error': !!props.error, 'autocomplete-disabled': !!props.disabled }}
    >
      <Show when={props.label}>
        <label class="autocomplete-label">{props.label}</label>
      </Show>
      <div class="autocomplete-input-wrapper">
        <input
          class="autocomplete-input"
          type="text"
          value={props.value}
          onInput={handleInputChange}
          onKeyDown={handleKeyDown}
          onFocus={() => setIsFocused(true)}
          onBlur={handleBlur}
          placeholder={props.placeholder}
          disabled={props.disabled}
        />
      </div>
      <Show when={isOpen() && filteredOptions().length > 0}>
        <div class="autocomplete-dropdown">
          <For each={filteredOptions()}>
            {(option, index) => (
              <button
                class="autocomplete-option"
                classList={{ 'autocomplete-highlighted': index() === highlightedIndex() }}
                onClick={() => handleSelect(option)}
                type="button"
                onMouseEnter={() => setHighlightedIndex(index())}
              >
                {option}
              </button>
            )}
          </For>
        </div>
      </Show>
      <Show when={props.error}>
        <span class="autocomplete-error-msg">{props.error}</span>
      </Show>
    </div>
  );
}
