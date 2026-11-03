import { createSignal, createResource, For, Show, createEffect } from 'solid-js';
import { useParams, useNavigate } from '@solidjs/router';
import { createTransaction, getTransaction, updateTransaction } from '../api/transaction';
import { searchAccounts } from '../api/account';
import { getSettings } from '../api/settings';
import { todayLocalDate } from '../state/dates';
import { isZero } from '../state/decimal';
import Big from 'big.js';
import { DatePicker } from '../components/forms/DatePicker';
import { TagInput } from '../components/forms/TagInput';
import './Form.css';

interface PostingDraft {
  account_id: number | null; // numeric ID from backend, null if new/unresolved
  account_name: string;      // display name (full_name)
  currency: string;
  amount: string;
}

const emptyPosting = (currency: string): PostingDraft => ({
  account_id: null,
  account_name: '',
  currency,
  amount: '',
});

export default function TransactionFormPage() {
  const params = useParams<{ id?: string }>();
  const navigate = useNavigate();
  const isEdit = !!params.id && params.id !== 'new';

  const [settings] = createResource(() => getSettings());
  const currencies = () => {
    const s = settings();
    if (!s) return ['USD'];
    const avail = s.find(x => x.key === 'available_currencies');
    if (avail) { try { return JSON.parse(avail.value) as string[]; } catch {} }
    const def = s.find(x => x.key === 'default_currency');
    return def ? [def.value] : ['USD'];
  };
  const defaultCurrency = () => currencies()[0] ?? 'USD';

  const [existingTx] = createResource(() => isEdit ? getTransaction(params.id!) : null);

  const [date, setDate] = createSignal(todayLocalDate());
  const [description, setDescription] = createSignal('');
  const [status, setStatus] = createSignal('pending');
  const [tags, setTags] = createSignal<{ key: string; value: string }[]>([]);

  // postings: stable signal, only updated on account select, currency change, or row add/remove.
  // NOT updated on every keystroke to avoid remounting focused inputs.
  const [postings, setPostings] = createSignal<PostingDraft[]>([
    emptyPosting('USD'),
    emptyPosting('USD'),
  ]);

  // Per-row account search text — separate from postings so typing doesn't remount rows
  const [accountSearchTexts, setAccountSearchTexts] = createSignal<string[]>(['', '']);

  // Shared autocomplete state — only one row is active at a time
  const [accountOptions, setAccountOptions] = createSignal<{ full_name: string; id: number }[]>([]);
  const [activeAutocompleteRow, setActiveAutocompleteRow] = createSignal<number | null>(null);
  const [highlightedOptionIndex, setHighlightedOptionIndex] = createSignal(-1);

  const [error, setError] = createSignal('');
  const [saving, setSaving] = createSignal(false);

  // Populate from existing transaction when editing
  createEffect(() => {
    const tx = existingTx();
    if (tx) {
      setDate(tx.date);
      setDescription(tx.description);
      setStatus(tx.status);
      setTags(tx.tags.map(t => ({ key: t.key, value: t.value })));
      const ps = tx.postings.map(p => ({
        account_id: Number(p.account_id),
        account_name: p.account_full_name || String(p.account_id),
        currency: p.currency,
        amount: p.amount,
      }));
      setPostings(ps);
      setAccountSearchTexts(ps.map(p => p.account_name));
      setInitialAmounts(ps.map(p => p.amount));
      setDisplayAmounts(ps.map(p => p.amount));
    }
  });

  const handleAccountSearch = async (rowIndex: number, query: string) => {
    setActiveAutocompleteRow(rowIndex);
    setHighlightedOptionIndex(-1);
    if (!query || query.length < 1) { setAccountOptions([]); return; }
    const results = await searchAccounts(query);
    setAccountOptions(results.map(a => ({ full_name: a.full_name, id: Number(a.id) })));
  };

  // Called when user selects an account from the dropdown
  const selectAccount = (rowIndex: number, fullName: string) => {
    const found = accountOptions().find(a => a.full_name === fullName);
    setPostings(prev => prev.map((p, i) =>
      i === rowIndex ? { ...p, account_id: found?.id ?? null, account_name: fullName } : p
    ));
    setAccountSearchTexts(prev => prev.map((t, i) => i === rowIndex ? fullName : t));
    setAccountOptions([]);
  };

  const updatePostingCurrency = (index: number, currency: string) => {
    setPostings(prev => prev.map((p, i) => i === index ? { ...p, currency } : p));
  };

  // Amount is read from the DOM on submit — uncontrolled to avoid remounting on keypress.
  // initialAmounts holds the value to set when a ref first mounts (used when editing).
  let amountRefs: (HTMLInputElement | undefined)[] = [];
  let currencyRefs: (HTMLSelectElement | undefined)[] = [];
  const [initialAmounts, setInitialAmounts] = createSignal<string[]>(['', '']);

  const addPosting = () => {
    setPostings(prev => [...prev, emptyPosting(defaultCurrency())]);
    setAccountSearchTexts(prev => [...prev, '']);
    amountRefs.push(undefined);
  };

  const removePosting = (index: number) => {
    if (postings().length <= 2) return;
    setPostings(prev => prev.filter((_, i) => i !== index));
    setAccountSearchTexts(prev => prev.filter((_, i) => i !== index));
    amountRefs.splice(index, 1);
  };

  // Read current amounts from DOM refs
  const getCurrentAmounts = () => amountRefs.map(ref => ref?.value ?? '');

  const computeBalance = () => {
    try {
      let sum = new Big(0);
      for (const amt of getCurrentAmounts()) {
        if (amt) sum = sum.plus(new Big(amt));
      }
      return sum;
    } catch { return null; }
  };

  // Reactive balance for display — based on postings signal amounts
  // (updated on blur of amount fields)
  const [displayAmounts, setDisplayAmounts] = createSignal<string[]>(['', '']);

  const isBalancedDisplay = () => {
    try {
      let sum = new Big(0);
      for (const amt of displayAmounts()) {
        if (amt) sum = sum.plus(new Big(amt));
      }
      return isZero(sum.toString());
    } catch { return false; }
  };

  const totalUnbalancedDisplay = () => {
    try {
      let sum = new Big(0);
      for (const amt of displayAmounts()) {
        if (amt) sum = sum.plus(new Big(amt));
      }
      return sum.abs().toString();
    } catch { return '?'; }
  };

  const handleAmountBlur = () => {
    setDisplayAmounts(getCurrentAmounts());
  };

  const handleSubmit = async () => {
    setError('');
    if (!description()) { setError('Description is required.'); return; }
    if (!date()) { setError('Date is required.'); return; }

    const amounts = getCurrentAmounts();
    const ps = postings();

    if (ps.some(p => !p.account_name)) {
      setError('All postings must have an account.');
      return;
    }

    const hasBlank = amounts.some(a => a === '');
    if (hasBlank && amounts[amounts.length - 1] !== '') {
      setError('Only the last posting amount may be left blank for auto-inference.');
      return;
    }

    const balance = computeBalance();
    if (balance && !balance.eq(0) && !hasBlank) {
      setError(`Postings do not balance. Difference: ${balance.abs().toString()}`);
      return;
    }

    setSaving(true);
    try {
      const postingInputs = ps.map((p, i) => ({
        account_id: p.account_id ?? 0,
        account_name: p.account_id ? '' : accountSearchTexts()[i], // send name only when no ID
        currency: p.currency,
        amount: amounts[i],
      }));

      if (isEdit) {
        await updateTransaction(params.id!, {
          date: date(), description: description(), status: status(),
          postings: postingInputs as any, tags: tags(),
        });
        navigate(`/transactions/${params.id}`);
      } else {
        const tx = await createTransaction({
          date: date(), description: description(), status: status(),
          postings: postingInputs as any, tags: tags(),
        });
        navigate(`/transactions/${tx.id}`);
      }
    } catch (e: any) {
      setError(e?.message ?? 'Failed to save transaction.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div class="form-page">
      <h1>{isEdit ? 'Edit Transaction' : 'New Transaction'}</h1>

      <Show when={error()}>
        <div class="form-error">{error()}</div>
      </Show>

      <div class="form-section">
        <div class="form-row">
          <div class="form-field">
            <label>Date</label>
            <DatePicker value={date()} onChange={setDate} />
          </div>
          <div class="form-field">
            <label>Status</label>
            <select value={status()} onChange={e => setStatus(e.currentTarget.value)}>
              <option value="pending">Pending</option>
              <option value="posted">Posted</option>
              <option value="cleared">Cleared</option>
            </select>
          </div>
        </div>

        <div class="form-field">
          <label>Description</label>
          <input
            type="text"
            value={description()}
            onInput={e => setDescription(e.currentTarget.value)}
            placeholder="e.g. Grocery run"
          />
        </div>

        <div class="form-field">
          <label>Tags</label>
          <TagInput value={tags()} onChange={setTags} placeholder="key:value" />
        </div>
      </div>

      <div class="postings-section">
        <div class="postings-header">
          <h3>Postings</h3>
          <span class={`balance-indicator ${isBalancedDisplay() ? 'balanced' : 'unbalanced'}`}>
            {isBalancedDisplay() ? '✓ Balanced' : `Unbalanced: ${totalUnbalancedDisplay()}`}
          </span>
        </div>

        <For each={postings()}>
          {(posting, i) => (
            <div class="posting-row">
              {/* Account autocomplete — uses its own local state, not postings signal */}
              <div class="form-field" style="margin: 0; position: relative">
                <input
                  class="autocomplete-input"
                  type="text"
                  value={accountSearchTexts()[i()]}
                  onInput={e => {
                    const val = e.currentTarget.value;
                    setAccountSearchTexts(prev => prev.map((t, idx) => idx === i() ? val : t));
                    handleAccountSearch(i(), val);
                  }}
                  onFocus={() => setActiveAutocompleteRow(i())}
                  onKeyDown={e => {
                    const opts = accountOptions();
                    const hi = highlightedOptionIndex();
                    const isOpen = activeAutocompleteRow() === i() && opts.length > 0;
                    if (e.key === 'ArrowDown') {
                      e.preventDefault();
                      if (isOpen) setHighlightedOptionIndex(Math.min(hi + 1, opts.length - 1));
                    } else if (e.key === 'ArrowUp') {
                      e.preventDefault();
                      if (isOpen) setHighlightedOptionIndex(Math.max(hi - 1, 0));
                    } else if ((e.key === 'Enter' || e.key === 'Tab') && isOpen && hi >= 0) {
                      e.preventDefault();
                      selectAccount(i(), opts[hi].full_name);
                      // After selecting, move focus to the currency select for this row
                      if (e.key === 'Tab') currencyRefs[i()]?.focus();
                    } else if (e.key === 'Escape') {
                      setAccountOptions([]);
                      setActiveAutocompleteRow(null);
                      setHighlightedOptionIndex(-1);
                    }
                  }}
                  onBlur={() => {
                    const text = accountSearchTexts()[i()];
                    if (text && !postings()[i()].account_name) {
                      setPostings(prev => prev.map((p, idx) =>
                        idx === i() ? { ...p, account_name: text } : p
                      ));
                    }
                    setTimeout(() => {
                      setAccountOptions([]);
                      setActiveAutocompleteRow(null);
                      setHighlightedOptionIndex(-1);
                    }, 150);
                  }}
                  placeholder="Account name…"
                />
                <Show when={activeAutocompleteRow() === i() && accountOptions().length > 0}>
                  <div class="autocomplete-dropdown">
                    <For each={accountOptions()}>
                      {(opt, optIdx) => (
                        <button
                          class={`autocomplete-option${highlightedOptionIndex() === optIdx() ? ' autocomplete-highlighted' : ''}`}
                          type="button"
                          onMouseDown={e => { e.preventDefault(); selectAccount(i(), opt.full_name); }}
                          onMouseEnter={() => setHighlightedOptionIndex(optIdx())}
                        >
                          {opt.full_name}
                        </button>
                      )}
                    </For>
                  </div>
                </Show>
              </div>

              <div class="form-field" style="margin: 0">
                <select
                  ref={el => { currencyRefs[i()] = el; }}
                  value={posting.currency}
                  onChange={e => updatePostingCurrency(i(), e.currentTarget.value)}
                >
                  <For each={currencies()}>
                    {c => <option value={c}>{c}</option>}
                  </For>
                </select>
              </div>

              {/* Amount is uncontrolled — no value= binding, so typing never causes remount */}
              <div class="form-field" style="margin: 0">
                <input
                  type="text"
                  ref={el => {
                    amountRefs[i()] = el;
                    // Set initial value when editing — uncontrolled inputs don't get value= prop
                    if (el && initialAmounts()[i()]) {
                      el.value = initialAmounts()[i()];
                    }
                  }}
                  placeholder={i() === postings().length - 1 ? 'auto' : '0.00'}
                  onBlur={handleAmountBlur}
                />
              </div>

              <button
                class="btn btn-secondary"
                style="padding: 8px; min-width: 36px"
                type="button"
                onClick={() => removePosting(i())}
                disabled={postings().length <= 2}
              >
                ✕
              </button>
            </div>
          )}
        </For>

        <button class="btn btn-secondary" style="width: 100%; margin-top: 8px" type="button" onClick={addPosting}>
          + Add Posting
        </button>
      </div>

      <div class="form-actions" style="margin-top: 24px">
        <button class="btn btn-primary" onClick={handleSubmit} disabled={saving()}>
          {saving() ? 'Saving…' : isEdit ? 'Save Changes' : 'Create Transaction'}
        </button>
        <button class="btn btn-secondary" onClick={() => navigate(isEdit ? `/transactions/${params.id}` : '/transactions')}>
          Cancel
        </button>
      </div>
    </div>
  );
}
