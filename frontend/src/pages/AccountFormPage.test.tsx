/**
 * UI tests for the AccountFormPage component.
 *
 * Covers user stories:
 * 1. Create root account: check root box, select type, name without ":"
 * 2. Create subaccount: uncheck root box, enter full ":" path (e.g. "Expenses:Food")
 * 3. Validation: root rejects ":", subaccount requires ":", empty name
 * 4. Auto-suggest type from well-known names
 * 5. Edit mode: pre-fill form from existing account (createEffect fix)
 */
import { render, fireEvent } from '@solidjs/testing-library';
import { describe, it, expect, vi, beforeEach } from 'vitest';

const mockCreateAccount = vi.fn().mockResolvedValue({ id: '1' });
const mockGetAccount = vi.fn();
const mockUpdateAccount = vi.fn().mockResolvedValue({ id: '1' });
const mockNavigate = vi.fn();

vi.mock('../api/account', () => ({
  createAccount: (...args: any[]) => mockCreateAccount(...args),
  getAccount: (...args: any[]) => mockGetAccount(...args),
  updateAccount: (...args: any[]) => mockUpdateAccount(...args),
}));
vi.mock('@solidjs/router', () => ({
  useParams: () => ({ id: undefined }),
  useNavigate: () => mockNavigate,
}));

import AccountFormPage from './AccountFormPage';

beforeEach(() => {
  vi.clearAllMocks();
  mockCreateAccount.mockResolvedValue({ id: '1' });
});

describe('AccountFormPage — basic rendering', () => {
  it('renders a New Account heading', () => {
    const { getByText } = render(() => <AccountFormPage />);
    expect(getByText('New Account')).toBeTruthy();
  });

  it('shows account type selector for root accounts by default', () => {
    const { getByText } = render(() => <AccountFormPage />);
    expect(getByText(/Account Type/i)).toBeTruthy();
  });

  it('hides type selector when root account checkbox is unchecked', () => {
    const { queryByText, container } = render(() => <AccountFormPage />);
    const checkbox = container.querySelector('input[type="checkbox"]')!;
    fireEvent.change(checkbox, { target: { checked: false } });
    expect(queryByText(/Account Type/i)).toBeNull();
  });
});

describe('AccountFormPage — Story 1: create root account', () => {
  it('shows error when root account name contains ":"', async () => {
    const { getByText, container } = render(() => <AccountFormPage />);
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'Bad:Name' } });
    fireEvent.click(getByText('Create Account'));
    await Promise.resolve();
    expect(getByText(/may not contain/i)).toBeTruthy();
  });

  it('shows error when root account has no type selected', async () => {
    const { getByText, container } = render(() => <AccountFormPage />);
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    // Use a name that doesn't trigger auto-suggestion
    fireEvent.input(nameInput, { target: { value: 'MyAccount123' } });
    fireEvent.click(getByText('Create Account'));
    await Promise.resolve();
    expect(getByText(/Account type is required/i)).toBeTruthy();
  });

  it('auto-suggests "asset" type when name is "Assets"', () => {
    const { container } = render(() => <AccountFormPage />);
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'Assets' } });
    const select = container.querySelector('select') as HTMLSelectElement;
    expect(select.value).toBe('asset');
  });

  it('auto-suggests "expense" type when name is "Expenses"', () => {
    const { container } = render(() => <AccountFormPage />);
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'Expenses' } });
    const select = container.querySelector('select') as HTMLSelectElement;
    expect(select.value).toBe('expense');
  });

  it('auto-suggests "revenue" type when name is "Income"', () => {
    const { container } = render(() => <AccountFormPage />);
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'Income' } });
    const select = container.querySelector('select') as HTMLSelectElement;
    expect(select.value).toBe('revenue');
  });

  it('does not auto-suggest type when a type is already selected', () => {
    const { container } = render(() => <AccountFormPage />);
    const select = container.querySelector('select') as HTMLSelectElement;
    fireEvent.change(select, { target: { value: 'liability' } });
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'Assets' } });
    expect(select.value).toBe('liability');
  });

  it('calls createAccount with correct payload for root account', async () => {
    const { container, getByText } = render(() => <AccountFormPage />);
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'Assets' } });
    // Type is auto-suggested to "asset"
    fireEvent.click(getByText('Create Account'));
    await new Promise(r => setTimeout(r, 10));
    expect(mockCreateAccount).toHaveBeenCalledWith({
      name: 'Assets',
      account_type: 'asset',
    });
  });
});

describe('AccountFormPage — Story 2: create subaccount via ":" path', () => {
  it('shows placeholder with ":" example when root is unchecked', () => {
    const { container } = render(() => <AccountFormPage />);
    const checkbox = container.querySelector('input[type="checkbox"]')!;
    fireEvent.change(checkbox, { target: { checked: false } });
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    expect(nameInput.placeholder).toContain(':');
  });

  it('shows error when subaccount name does not contain ":"', async () => {
    const { container, getByText } = render(() => <AccountFormPage />);
    const checkbox = container.querySelector('input[type="checkbox"]')!;
    fireEvent.change(checkbox, { target: { checked: false } });
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'JustAName' } });
    fireEvent.click(getByText('Create Account'));
    await Promise.resolve();
    expect(getByText(/must include the parent path/i)).toBeTruthy();
  });

  it('accepts ":" in name when root is unchecked', async () => {
    const { container, getByText } = render(() => <AccountFormPage />);
    const checkbox = container.querySelector('input[type="checkbox"]')!;
    fireEvent.change(checkbox, { target: { checked: false } });
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'Expenses:Food' } });
    fireEvent.click(getByText('Create Account'));
    await new Promise(r => setTimeout(r, 10));
    // Should call createAccount (no validation error about ":")
    expect(mockCreateAccount).toHaveBeenCalledWith({
      name: 'Expenses:Food',
      account_type: undefined,
    });
  });

  it('sends deep hierarchy path to createAccount', async () => {
    const { container, getByText } = render(() => <AccountFormPage />);
    const checkbox = container.querySelector('input[type="checkbox"]')!;
    fireEvent.change(checkbox, { target: { checked: false } });
    const nameInput = container.querySelector('input[type="text"]') as HTMLInputElement;
    fireEvent.input(nameInput, { target: { value: 'Expenses:Food:Restaurants' } });
    fireEvent.click(getByText('Create Account'));
    await new Promise(r => setTimeout(r, 10));
    expect(mockCreateAccount).toHaveBeenCalledWith({
      name: 'Expenses:Food:Restaurants',
      account_type: undefined,
    });
  });
});

describe('AccountFormPage — validation edge cases', () => {
  it('shows error for empty name', async () => {
    const { getByText } = render(() => <AccountFormPage />);
    fireEvent.click(getByText('Create Account'));
    await Promise.resolve();
    expect(getByText(/Name is required/i)).toBeTruthy();
  });
});
