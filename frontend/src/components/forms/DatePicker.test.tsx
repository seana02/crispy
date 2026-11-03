/**
 * UI tests for DatePicker component.
 * Verifies it outputs YYYY-MM-DD strings and never shifts dates via UTC.
 */
import { render, fireEvent } from '@solidjs/testing-library';
import { createSignal } from 'solid-js';
import { describe, it, expect } from 'vitest';
import { DatePicker } from './DatePicker';

function DatePickerWrapper() {
  const [date, setDate] = createSignal('2026-01-15');
  return (
    <div>
      <DatePicker value={date()} onChange={setDate} label="Date" />
      <span data-testid="value">{date()}</span>
    </div>
  );
}

describe('DatePicker', () => {
  it('renders with the provided date value', () => {
    const { getByDisplayValue } = render(() => <DatePickerWrapper />);
    expect(getByDisplayValue('2026-01-15')).toBeTruthy();
  });

  it('outputs a YYYY-MM-DD string when changed', () => {
    const { getByTestId, container } = render(() => <DatePickerWrapper />);
    const input = container.querySelector('input[type="date"]')!;
    fireEvent.input(input, { target: { value: '2026-06-01' } });
    expect(getByTestId('value').textContent).toBe('2026-06-01');
  });

  it('does not output a Date object — always a plain string', () => {
    const { getByTestId, container } = render(() => <DatePickerWrapper />);
    const input = container.querySelector('input[type="date"]')!;
    fireEvent.input(input, { target: { value: '2026-12-31' } });
    const value = getByTestId('value').textContent;
    // Must be a string in YYYY-MM-DD format, not a timestamp or Date object representation
    expect(value).toMatch(/^\d{4}-\d{2}-\d{2}$/);
    expect(value).toBe('2026-12-31');
  });

  it('January 1st is not shifted to December 31st (UTC timezone trap)', () => {
    const { getByTestId, container } = render(() => <DatePickerWrapper />);
    const input = container.querySelector('input[type="date"]')!;
    fireEvent.input(input, { target: { value: '2026-01-01' } });
    // If we were using new Date('2026-01-01') this could become Dec 31 in UTC-offset zones
    expect(getByTestId('value').textContent).toBe('2026-01-01');
  });
});
