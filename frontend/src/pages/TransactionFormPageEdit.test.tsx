/**
 * Regression tests for TransactionFormPage in EDIT mode.
 *
 * Bugs covered:
 * 1. Edit form did not show amounts — amount inputs are uncontrolled, initial values
 *    were never set from the existing transaction's postings.
 * 2. Account names showed numeric IDs ("5", "6") instead of full names.
 *    (covered separately in TransactionFormPage.test.tsx but verified here too)
 */
import { render } from '@solidjs/testing-library';
import { describe, it, expect, vi } from 'vitest';

// These mocks are hoisted — useParams returns an edit ID.
vi.mock('@solidjs/router', () => ({
  useParams: () => ({ id: '7' }),
  useNavigate: () => vi.fn(),
}));

vi.mock('../api/transaction', () => ({
  createTransaction: vi.fn(),
  getTransaction: vi.fn().mockResolvedValue({
    id: '7',
    date: '2026-03-10',
    description: 'Coffee run',
    status: 'cleared',
    recurring_template_id: null,
    tags: [],
    postings: [
      {
        id: '20',
        transaction_id: '7',
        account_id: '5',
        account_full_name: 'Assets:Checking',
        currency: 'USD',
        amount: '-4.50',
      },
      {
        id: '21',
        transaction_id: '7',
        account_id: '6',
        account_full_name: 'Expenses:Food',
        currency: 'USD',
        amount: '4.50',
      },
    ],
    created_at: '2026-03-10T10:00:00Z',
    updated_at: '2026-03-10T10:00:00Z',
  }),
  updateTransaction: vi.fn().mockResolvedValue({ id: '7' }),
}));

vi.mock('../api/account', () => ({
  searchAccounts: vi.fn().mockResolvedValue([]),
}));

vi.mock('../api/settings', () => ({
  getSettings: vi.fn().mockResolvedValue([
    { key: 'default_currency', value: 'USD' },
    { key: 'available_currencies', value: '["USD"]' },
  ]),
}));

import TransactionFormPage from './TransactionFormPage';

describe('TransactionFormPage — edit mode amount regression', () => {
  it('shows Edit Transaction heading when given an id param', () => {
    const { getByText } = render(() => <TransactionFormPage />);
    expect(getByText('Edit Transaction')).toBeTruthy();
  });

  it('pre-fills amount inputs from existing transaction postings', async () => {
    const { container } = render(() => <TransactionFormPage />);

    // Wait for the async resource (getTransaction) to resolve and createEffect to run
    await new Promise(r => setTimeout(r, 30));

    const amountInputs = Array.from(
      container.querySelectorAll('input[placeholder="0.00"], input[placeholder="auto"]')
    ) as HTMLInputElement[];

    // Regression: before fix, amounts were empty because uncontrolled inputs
    // were never populated from existingTx().postings.
    expect(amountInputs.length).toBeGreaterThanOrEqual(2);
    expect(amountInputs[0].value).toBe('-4.50');
    expect(amountInputs[1].value).toBe('4.50');
  });

  it('pre-fills account names from account_full_name, not account_id', async () => {
    const { container } = render(() => <TransactionFormPage />);

    await new Promise(r => setTimeout(r, 30));

    const accountInputs = Array.from(
      container.querySelectorAll('input[placeholder="Account name…"]')
    ) as HTMLInputElement[];

    expect(accountInputs.length).toBeGreaterThanOrEqual(2);

    // Regression: before fix, values were "5" and "6" (raw numeric IDs).
    expect(accountInputs[0].value).toBe('Assets:Checking');
    expect(accountInputs[1].value).toBe('Expenses:Food');
    expect(accountInputs[0].value).not.toBe('5');
    expect(accountInputs[1].value).not.toBe('6');
  });
});
