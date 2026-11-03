import { createSignal, createResource, For, Show } from 'solid-js';
import { useParams, useNavigate } from '@solidjs/router';
import { createBudget, getBudget, updateBudget } from '../api/budget';
import { searchAccounts } from '../api/account';
import { AutocompleteInput } from '../components/forms/AutocompleteInput';
import { TagInput } from '../components/forms/TagInput';
import './Form.css';

export default function BudgetFormPage() {
  const params = useParams<{ id?: string }>();
  const navigate = useNavigate();
  const isEdit = !!params.id && params.id !== 'new';

  const [existingBudget] = createResource(() => isEdit ? getBudget(params.id!) : null);

  const [name, setName] = createSignal('');
  const [direction, setDirection] = createSignal<'spending' | 'saving'>('spending');
  const [period, setPeriod] = createSignal<'weekly' | 'monthly' | 'yearly'>('monthly');
  const [currency, setCurrency] = createSignal('USD');
  const [targetAmount, setTargetAmount] = createSignal('');
  const [tagFilters, setTagFilters] = createSignal<{ key: string; value: string }[]>([]);
  const [selectedAccounts, setSelectedAccounts] = createSignal<{ id: string; name: string }[]>([]);
  const [accountSearch, setAccountSearch] = createSignal('');
  const [accountOptions, setAccountOptions] = createSignal<{ id: string; name: string }[]>([]);
  const [error, setError] = createSignal('');
  const [saving, setSaving] = createSignal(false);

  // Populate from existing budget when editing
  createSignal(() => {
    const b = existingBudget();
    if (b) {
      setName(b.name);
      setDirection(b.direction as any);
      setPeriod(b.period as any);
      setCurrency(b.currency);
      setTargetAmount(b.target_amount);
    }
  });

  const handleAccountSearch = async (query: string) => {
    if (!query) return;
    const results = await searchAccounts(query);
    setAccountOptions(results.map(a => ({ id: String(a.id), name: a.full_name })));
  };

  const addAccount = (name: string) => {
    const found = accountOptions().find(a => a.name === name);
    if (found && !selectedAccounts().find(a => a.id === found.id)) {
      setSelectedAccounts(prev => [...prev, found]);
    }
    setAccountSearch('');
  };

  const removeAccount = (id: string) => {
    setSelectedAccounts(prev => prev.filter(a => a.id !== id));
  };

  const handleSubmit = async () => {
    setError('');
    if (!name()) { setError('Name is required.'); return; }
    if (!targetAmount() || isNaN(Number(targetAmount()))) { setError('Target amount must be a valid number.'); return; }
    if (selectedAccounts().length === 0) { setError('At least one account must be selected.'); return; }

    setSaving(true);
    try {
      const accountIds = selectedAccounts().map(a => Number(a.id));
      const tagIds: number[] = []; // tags are referenced by key:value — the backend handles lookup
      if (isEdit) {
        await updateBudget(params.id!, { name: name(), direction: direction(), period: period(), currency: currency(), target_amount: targetAmount(), account_ids: accountIds, tag_ids: tagIds });
        navigate(`/budgets/${params.id}`);
      } else {
        const budget = await createBudget({ name: name(), direction: direction(), period: period(), currency: currency(), target_amount: targetAmount(), account_ids: accountIds, tag_ids: tagIds });
        navigate(`/budgets/${budget.id}`);
      }
    } catch (e: any) {
      setError(e?.message ?? 'Failed to save budget.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div class="form-page">
      <h1>{isEdit ? 'Edit Budget' : 'New Budget'}</h1>

      <Show when={error()}>
        <div class="form-error">{error()}</div>
      </Show>

      <div class="form-section">
        <div class="form-field">
          <label>Budget Name</label>
          <input type="text" value={name()} onInput={e => setName(e.currentTarget.value)} placeholder="e.g. Monthly Food" />
        </div>

        <div class="form-row">
          <div class="form-field">
            <label>Direction</label>
            <select value={direction()} onChange={e => setDirection(e.currentTarget.value as any)}>
              <option value="spending">Spending (ceiling)</option>
              <option value="saving">Saving (goal)</option>
            </select>
          </div>
          <div class="form-field">
            <label>Period</label>
            <select value={period()} onChange={e => setPeriod(e.currentTarget.value as any)}>
              <option value="weekly">Weekly</option>
              <option value="monthly">Monthly</option>
              <option value="yearly">Yearly</option>
            </select>
          </div>
        </div>

        <div class="form-row">
          <div class="form-field">
            <label>Currency</label>
            <input type="text" value={currency()} onInput={e => setCurrency(e.currentTarget.value)} placeholder="USD" />
          </div>
          <div class="form-field">
            <label>Target Amount</label>
            <input type="text" value={targetAmount()} onInput={e => setTargetAmount(e.currentTarget.value)} placeholder="500.00" />
          </div>
        </div>

        <div class="form-field">
          <label>Tracked Accounts (required — includes all children)</label>
          <AutocompleteInput
            value={accountSearch()}
            options={accountOptions().map(a => a.name)}
            onChange={val => { setAccountSearch(val); handleAccountSearch(val); }}
            onSelect={addAccount}
            placeholder="Search and select accounts…"
            highlightMatch
          />
          <Show when={selectedAccounts().length > 0}>
            <div style="display: flex; flex-wrap: wrap; gap: 6px; margin-top: 8px">
              <For each={selectedAccounts()}>
                {acct => (
                  <span class="tag-chip" style="cursor: default">
                    {acct.name}
                    <button
                      type="button"
                      style="background: none; border: none; cursor: pointer; color: inherit; padding: 0 0 0 4px"
                      onClick={() => removeAccount(acct.id)}
                    >✕</button>
                  </span>
                )}
              </For>
            </div>
          </Show>
        </div>

        <div class="form-field">
          <label>Tag Filters (optional — narrow which transactions count)</label>
          <TagInput value={tagFilters()} onChange={setTagFilters} placeholder="key:value" />
        </div>
      </div>

      <div class="form-actions">
        <button class="btn btn-primary" onClick={handleSubmit} disabled={saving()}>
          {saving() ? 'Saving…' : isEdit ? 'Save Changes' : 'Create Budget'}
        </button>
        <button class="btn btn-secondary" onClick={() => navigate(isEdit ? `/budgets/${params.id}` : '/budgets')}>
          Cancel
        </button>
      </div>
    </div>
  );
}
