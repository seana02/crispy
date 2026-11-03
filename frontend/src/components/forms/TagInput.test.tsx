/**
 * UI tests for TagInput component.
 * Verifies the key:value chip input workflow.
 */
import { render, fireEvent } from '@solidjs/testing-library';
import { createSignal } from 'solid-js';
import { describe, it, expect } from 'vitest';
import { TagInput } from './TagInput';

function TagInputWrapper() {
  const [tags, setTags] = createSignal<{ key: string; value: string }[]>([]);
  return (
    <div>
      <TagInput value={tags()} onChange={setTags} placeholder="key:value" />
      <span data-testid="count">{tags().length}</span>
    </div>
  );
}

describe('TagInput', () => {
  it('renders with empty state', () => {
    const { getByPlaceholderText } = render(() => <TagInputWrapper />);
    expect(getByPlaceholderText('key:value')).toBeTruthy();
  });

  it('adds a tag chip when Enter is pressed with key:value', () => {
    const { getByPlaceholderText, getByTestId } = render(() => <TagInputWrapper />);
    const input = getByPlaceholderText('key:value');
    fireEvent.input(input, { target: { value: 'trip:japan' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(getByTestId('count').textContent).toBe('1');
  });

  it('does not add a chip for input without colon separator', () => {
    const { getByPlaceholderText, getByTestId } = render(() => <TagInputWrapper />);
    const input = getByPlaceholderText('key:value');
    fireEvent.input(input, { target: { value: 'nocolon' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(getByTestId('count').textContent).toBe('0');
  });

  it('removes a chip when the remove button is clicked', () => {
    const { getByPlaceholderText, getByTestId, getByText } = render(() => <TagInputWrapper />);
    const input = getByPlaceholderText('key:value');
    fireEvent.input(input, { target: { value: 'trip:japan' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(getByTestId('count').textContent).toBe('1');
    // Click the remove button (✕)
    fireEvent.click(getByText('✕'));
    expect(getByTestId('count').textContent).toBe('0');
  });

  it('removes last chip when Backspace is pressed on empty input', () => {
    const { getByPlaceholderText, getByTestId } = render(() => <TagInputWrapper />);
    const input = getByPlaceholderText('key:value');
    fireEvent.input(input, { target: { value: 'trip:japan' } });
    fireEvent.keyDown(input, { key: 'Enter' });
    expect(getByTestId('count').textContent).toBe('1');
    // Clear the input then press Backspace
    fireEvent.input(input, { target: { value: '' } });
    fireEvent.keyDown(input, { key: 'Backspace' });
    expect(getByTestId('count').textContent).toBe('0');
  });
});
