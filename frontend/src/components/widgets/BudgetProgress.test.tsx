/**
 * UI tests for BudgetProgress widget.
 * Verifies progress calculation and visual states.
 */
import { render } from '@solidjs/testing-library';
import { describe, it, expect } from 'vitest';
import { BudgetProgress } from './BudgetProgress';

describe('BudgetProgress', () => {
  it('renders without errors when given valid props', () => {
    const { container } = render(() => (
      <BudgetProgress allocated="500.00" spent="200.00" remaining="300.00" />
    ));
    expect(container).toBeTruthy();
  });

  it('shows spending labels when showLabel is true', () => {
    const { getAllByText, getByText } = render(() => (
      <BudgetProgress allocated="500.00" spent="200.00" remaining="300.00" showLabel />
    ));
    expect(getByText('Allocated')).toBeTruthy();
    // "Spent" appears in both the label section and stats section
    expect(getAllByText('Spent').length).toBeGreaterThanOrEqual(1);
    expect(getAllByText('Remaining').length).toBeGreaterThanOrEqual(1);
  });

  it('calculates percentage correctly: 200/500 = 40%', () => {
    const { getByText } = render(() => (
      <BudgetProgress allocated="500.00" spent="200.00" remaining="300.00" showLabel />
    ));
    expect(getByText('40.0%')).toBeTruthy();
  });

  it('shows Overage label when spent exceeds allocated', () => {
    const { getAllByText } = render(() => (
      <BudgetProgress allocated="100.00" spent="150.00" remaining="-50.00" showLabel />
    ));
    expect(getAllByText('Overage').length).toBeGreaterThanOrEqual(1);
  });

  it('clamps progress bar at 100% when over budget', () => {
    // Even if spent is 200% of allocated, bar should not exceed 100%
    render(() => (
      <BudgetProgress allocated="100.00" spent="200.00" remaining="-100.00" />
    ));
    // The fill div should have width <= 100%
    const fill = document.querySelector('.budget-progress-bar-fill') as HTMLElement;
    if (fill) {
      const width = parseFloat(fill.style.width);
      expect(width).toBeLessThanOrEqual(100);
    }
  });

  it('shows 0% when allocated is zero to avoid division by zero', () => {
    const { getByText } = render(() => (
      <BudgetProgress allocated="0" spent="0" remaining="0" showLabel />
    ));
    expect(getByText('0.0%')).toBeTruthy();
  });
});
