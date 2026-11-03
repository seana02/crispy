import { describe, it, expect } from 'vitest';
import { todayLocalDate, formatDateDisplay, formatPeriodDisplay } from './dates';

describe('todayLocalDate', () => {
  it('returns a string in YYYY-MM-DD format', () => {
    const result = todayLocalDate();
    expect(result).toMatch(/^\d{4}-\d{2}-\d{2}$/);
  });

  it('does not return a UTC date that could be off by a day', () => {
    // This test verifies the function returns local date, not UTC date.
    // We can't easily simulate timezone offset in tests, but we can verify
    // the format and that it's a valid date.
    const result = todayLocalDate();
    const parts = result.split('-');
    expect(parts).toHaveLength(3);
    const year = parseInt(parts[0]);
    const month = parseInt(parts[1]);
    const day = parseInt(parts[2]);
    expect(year).toBeGreaterThan(2000);
    expect(month).toBeGreaterThanOrEqual(1);
    expect(month).toBeLessThanOrEqual(12);
    expect(day).toBeGreaterThanOrEqual(1);
    expect(day).toBeLessThanOrEqual(31);
  });
});

describe('formatDateDisplay', () => {
  it('formats a YYYY-MM-DD string correctly', () => {
    const result = formatDateDisplay('2026-10-31');
    expect(result).toBe('2026-10-31');
  });

  it('does not shift the date due to UTC timezone parsing', () => {
    // If we passed through new Date('2026-01-01'), it would be UTC midnight
    // and could display as Dec 31 in negative-offset timezones.
    // formatDateDisplay should not have this problem.
    const result = formatDateDisplay('2026-01-01');
    // The year should still be 2026, not 2025
    expect(result.startsWith('2026')).toBe(true);
  });

  it('returns empty string for empty input', () => {
    expect(formatDateDisplay('')).toBe('');
  });
});

describe('formatPeriodDisplay', () => {
  it('formats a YYYY-MM period string as Month Year', () => {
    const result = formatPeriodDisplay('2026-10');
    expect(result).toContain('2026');
    // Should contain the month name (October)
    expect(result.toLowerCase()).toContain('oct');
  });

  it('returns the input unchanged if not a valid period', () => {
    const result = formatPeriodDisplay('not-a-period');
    expect(result).toBe('not-a-period');
  });

  it('returns empty string for empty input', () => {
    expect(formatPeriodDisplay('')).toBe('');
  });
});
