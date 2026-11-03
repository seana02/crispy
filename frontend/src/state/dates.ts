import { format } from 'date-fns';

/**
 * Returns today's date as a YYYY-MM-DD string in the local timezone.
 * Does NOT use new Date().toISOString() which returns UTC and can be
 * off by a day near midnight in negative-offset timezones.
 */
export function todayLocalDate(): string {
  return format(new Date(), 'yyyy-MM-dd');
}

/**
 * Formats a YYYY-MM-DD date string for display.
 * Does NOT pass through new Date() to avoid timezone shifting.
 * Uses simple string parsing to extract year, month, day.
 */
export function formatDateDisplay(dateStr: string, dateFormat: string = 'yyyy-MM-dd'): string {
  if (!dateStr) return '';
  // Parse the YYYY-MM-DD string manually to avoid timezone issues
  const [year, month, day] = dateStr.split('-').map(Number);
  // Create a date at noon local time to avoid any timezone edge cases
  const date = new Date(year, month - 1, day, 12, 0, 0);
  return format(date, dateFormat);
}

/**
 * Formats a YYYY-MM period string (e.g. "2026-10") as a human-readable
 * month and year (e.g. "October 2026").
 */
export function formatPeriodDisplay(period: string): string {
  if (!period) return '';
  const parts = period.split('-');
  if (parts.length === 2) {
    const [year, month] = parts.map(Number);
    // Create at noon to avoid timezone issues
    const date = new Date(year, month - 1, 1, 12, 0, 0);
    return date.toLocaleString('default', { month: 'long', year: 'numeric' });
  }
  return period;
}
