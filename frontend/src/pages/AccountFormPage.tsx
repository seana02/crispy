import { createSignal, createEffect, createResource, For, Show } from 'solid-js';
import { useParams, useNavigate } from '@solidjs/router';
import { createAccount, getAccount, updateAccount } from '../api/account';
import './Form.css';

const ACCOUNT_TYPES = ['asset', 'liability', 'revenue', 'expense', 'equity'] as const;

function suggestType(name: string): string {
    const lower = name.toLowerCase();
    if (['assets', 'asset', 'checking', 'savings', 'cash', 'bank', 'wallet'].includes(lower)) return 'asset';
    if (['liabilities', 'liability', 'credit', 'loan', 'mortgage', 'debt'].includes(lower)) return 'liability';
    if (['revenue', 'revenues', 'income', 'salary'].includes(lower)) return 'revenue';
    if (['expenses', 'expense', 'spending'].includes(lower)) return 'expense';
    if (['equity', 'opening', 'retained'].includes(lower)) return 'equity';
    return '';
}

export default function AccountFormPage() {
    const params = useParams<{ id?: string }>();
    const navigate = useNavigate();
    const isEdit = !!params.id && params.id !== 'new';

    const [existingAccount] = createResource(() => isEdit ? getAccount(params.id!) : null);

    const [name, setName] = createSignal('');
    const [accountType, setAccountType] = createSignal('');
    const [isRoot, setIsRoot] = createSignal(true);
    const [error, setError] = createSignal('');
    const [saving, setSaving] = createSignal(false);

    // Populate from existing account when editing
    createEffect(() => {
        const acct = existingAccount();
        if (acct) {
            setName(acct.full_name);
            setAccountType(acct.account_type ?? '');
            setIsRoot(!acct.parent_id);
        }
    });

    const handleNameChange = (val: string) => {
        setName(val);
        // Auto-suggest type for root accounts based on the name
        if (isRoot() && !accountType()) {
            setAccountType(suggestType(val));
        }
    };

    const handleSubmit = async () => {
        setError('');
        const n = name().trim();
        if (!n) { setError('Name is required.'); return; }

        if (isRoot()) {
            // Root accounts: no ":" allowed, type required
            if (n.includes(':')) { setError('Root account name may not contain ":".'); return; }
            if (!accountType()) { setError('Account type is required for root accounts.'); return; }
        } else {
            // Subaccounts: name must contain ":" to indicate hierarchy
            if (!n.includes(':')) { setError('Subaccount name must include the parent path separated by ":". For example: "Expenses:Food".'); return; }
        }

        setSaving(true);
        try {
            if (isEdit) {
                // For edit, send the leaf name and keep existing parent relationship
                const parts = n.split(':');
                const leafName = parts[parts.length - 1];
                await updateAccount(params.id!, { name: leafName, account_type: isRoot() ? accountType() : undefined });
                navigate(`/accounts/${params.id}`);
            } else {
                // For create: if subaccount, send full path with ":"; backend handles hierarchy
                const acct = await createAccount({
                    name: n,
                    account_type: isRoot() ? accountType() : undefined,
                });
                navigate(`/accounts/${acct.id}`);
            }
        } catch (e: any) {
            setError(e?.message ?? 'Failed to save account.');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div class="form-page">
            <h1>{isEdit ? 'Edit Account' : 'New Account'}</h1>

            <Show when={error()}>
                <div class="form-error">{error()}</div>
            </Show>

            <div class="form-section">
                <div class="form-field">
                    <label>
                        <input
                            type="checkbox"
                            checked={isRoot()}
                            onChange={e => setIsRoot(e.currentTarget.checked)}
                            style="margin-right: 8px"
                        />
                        Root account (top-level)
                    </label>
                </div>

                <div class="form-field">
                    <label>Account Name</label>
                    <input
                        type="text"
                        value={name()}
                        onInput={e => handleNameChange(e.currentTarget.value)}
                        placeholder={isRoot() ? 'e.g. Checking' : 'e.g. Expenses:Food:Restaurants'}
                    />
                </div>

                <Show when={isRoot()}>
                    <div class="form-field">
                        <label>Account Type</label>
                        <select value={accountType()} onChange={e => setAccountType(e.currentTarget.value)}>
                            <option value="">— select type —</option>
                            <For each={ACCOUNT_TYPES}>
                                {t => <option value={t}>{t.charAt(0).toUpperCase() + t.slice(1)}</option>}
                            </For>
                        </select>
                    </div>
                </Show>
            </div>

            <div class="form-actions">
                <button class="btn btn-primary" onClick={handleSubmit} disabled={saving()}>
                    {saving() ? 'Saving…' : isEdit ? 'Save Changes' : 'Create Account'}
                </button>
                <button class="btn btn-secondary" onClick={() => navigate(isEdit ? `/accounts/${params.id}` : '/accounts')}>
                    Cancel
                </button>
            </div>
        </div>
    );
}
