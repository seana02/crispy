import { createMemo, createSignal, Show } from 'solid-js';
import { useParams, useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getBudget } from '../api/budget';
import { formatCurrency } from '../state/decimal';

export default function BudgetDetail() {
  const params = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [budgetResource] = createResource(() => getBudget(params.id));
  const budget = createMemo(() => budgetResource());

  return (
    <div class="page">
      <Show when={budget()} fallback={<p style="color: var(--color-text-secondary)">Loading...</p>}>
        {(b) => (
          <>
            <div class="page-header">
              <div>
                <h1>{b().name}</h1>
                <p style="color: var(--color-text-secondary); font-size: 13px; margin-top: 4px">
                  {b().direction} · {b().period} · {b().currency}
                </p>
              </div>
              <div style="display: flex; gap: 8px">
                <button class="btn btn-primary" onClick={() => navigate(`/budgets/${b().id}/edit`)}>Edit</button>
                <button class="btn btn-secondary" onClick={() => navigate('/budgets')}>Back</button>
              </div>
            </div>

            <div class="card">
              <div style="margin-bottom: 8px; font-size: 12px; color: var(--color-text-secondary)">Target Amount</div>
              <div class="amount" style="font-size: 24px">{formatCurrency(b().target_amount, b().currency)}</div>
            </div>

            <div class="card">
              <h3 style="margin: 0 0 12px; font-size: 16px; color: var(--color-text-primary)">Tracked Accounts</h3>
              <p style="color: var(--color-text-secondary); font-size: 13px">
                {b().account_ids.length} account(s) tracked
              </p>
            </div>
          </>
        )}
      </Show>
    </div>
  );
}
