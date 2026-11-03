import { createMemo, createSignal, Show, For } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getExchangeRates, deleteExchangeRate, createExchangeRate } from '../api/exchange_rate';
import { ExchangeRate } from '../api/types';
import { ConfirmDialog } from '../components/layout/ConfirmDialog';

export default function ExchangeRates() {
  const navigate = useNavigate();

  const [baseCurrency, setBaseCurrency] = createSignal('');
  const [quoteCurrency, setQuoteCurrency] = createSignal('');
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = createSignal(false);
  const [rateToDelete, setRateToDelete] = createSignal<string | null>(null);
  const [isAdding, setIsAdding] = createSignal(false);
  const [newBase, setNewBase] = createSignal('');
  const [newQuote, setNewQuote] = createSignal('');
  const [newRate, setNewRate] = createSignal('');
  const [newDate, setNewDate] = createSignal('');

  const [ratesResource, { refetch }] = createResource(() =>
    getExchangeRates({
      base_currency: baseCurrency() || undefined,
      quote_currency: quoteCurrency() || undefined,
    })
  );

  const rates = createMemo(() => ratesResource() ?? []);

  const handleDelete = () => {
    const id = rateToDelete();
    if (id) deleteExchangeRate(id).then(() => { refetch(); setIsDeleteDialogOpen(false); setRateToDelete(null); });
  };

  const handleAdd = async () => {
    await createExchangeRate({ base_currency: newBase(), quote_currency: newQuote(), rate: newRate(), effective_date: newDate(), source: 'manual' });
    setIsAdding(false);
    setNewBase(''); setNewQuote(''); setNewRate(''); setNewDate('');
    refetch();
  };

  return (
    <div class="page">
      <div class="page-header">
        <h1>Exchange Rates</h1>
        <button class="btn btn-primary" onClick={() => setIsAdding(true)}>+ Add Rate</button>
      </div>

      <Show when={isAdding()}>
        <div class="card" style="display: flex; gap: 12px; flex-wrap: wrap; align-items: flex-end">
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Base</label>
            <input style="padding: 8px 12px; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={newBase()} onInput={(e) => setNewBase(e.currentTarget.value)} placeholder="USD" />
          </div>
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Quote</label>
            <input style="padding: 8px 12px; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={newQuote()} onInput={(e) => setNewQuote(e.currentTarget.value)} placeholder="JPY" />
          </div>
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Rate</label>
            <input style="padding: 8px 12px; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family-mono)" value={newRate()} onInput={(e) => setNewRate(e.currentTarget.value)} placeholder="149.50" />
          </div>
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Effective Date</label>
            <input type="date" style="padding: 8px 12px; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={newDate()} onInput={(e) => setNewDate(e.currentTarget.value)} />
          </div>
          <button class="btn btn-primary" onClick={handleAdd}>Save</button>
          <button class="btn btn-secondary" onClick={() => setIsAdding(false)}>Cancel</button>
        </div>
      </Show>

      <div style="display: flex; gap: 8px; flex-wrap: wrap; align-items: center">
        <input style="padding: 8px 12px; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" placeholder="Base currency" value={baseCurrency()} onInput={(e) => setBaseCurrency(e.currentTarget.value)} />
        <input style="padding: 8px 12px; background: var(--color-bg-elevated); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" placeholder="Quote currency" value={quoteCurrency()} onInput={(e) => setQuoteCurrency(e.currentTarget.value)} />
        <button class="btn btn-secondary" onClick={() => refetch()}>Filter</button>
      </div>

      <div class="styled-list">
        <div class="styled-list-items">
          <Show when={!ratesResource.loading} fallback={<div class="styled-list-loading">Loading...</div>}>
            <Show when={rates().length > 0} fallback={<div class="styled-list-empty">No exchange rates found</div>}>
              <For each={rates()}>
                {(rate) => (
                  <div class="styled-list-row">
                    <div style="display: flex; justify-content: space-between; align-items: center; gap: 12px; flex-wrap: wrap">
                      <div style="display: flex; gap: 12px; align-items: center">
                        <span style="font-weight: 600; color: var(--color-text-primary)">{rate.base_currency}/{rate.quote_currency}</span>
                        <span class="amount">{rate.rate}</span>
                        <span style="color: var(--color-text-secondary); font-size: 12px">{rate.effective_date}</span>
                        <span class="type-badge equity">{rate.source}</span>
                      </div>
                      <button class="btn-action btn-danger" onClick={() => { setRateToDelete(rate.id); setIsDeleteDialogOpen(true); }}>Delete</button>
                    </div>
                  </div>
                )}
              </For>
            </Show>
          </Show>
        </div>
      </div>

      <ConfirmDialog
        isOpen={isDeleteDialogOpen()}
        title="Delete Rate"
        message="Are you sure you want to delete this exchange rate?"
        variant="danger"
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleDelete}
        onCancel={() => { setIsDeleteDialogOpen(false); setRateToDelete(null); }}
      />
    </div>
  );
}
