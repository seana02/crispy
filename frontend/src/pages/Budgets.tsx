import { createMemo, createSignal, Show, For } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getBudgets, deleteBudget } from '../api/budget';
import { Budget } from '../api/types';
import { Pagination } from '../components/lists/Pagination';
import { ConfirmDialog } from '../components/layout/ConfirmDialog';
import { formatCurrency } from '../state/decimal';
import './Budgets.css';

export default function Budgets() {
  const navigate = useNavigate();
  const [page, setPage] = createSignal(1);
  const [perPage] = createSignal(10);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = createSignal(false);
  const [budgetToDelete, setBudgetToDelete] = createSignal<string | null>(null);

  const [budgetsResource, { refetch }] = createResource(() =>
    getBudgets({ page: page(), per_page: perPage() })
  );

  const budgets = createMemo(() => budgetsResource()?.data ?? []);
  const meta = createMemo(() => budgetsResource()?.meta ?? { page: 1, per_page: 10, total: 0, totalPages: 0 });

  const handleDelete = () => {
    const id = budgetToDelete();
    if (id) deleteBudget(id).then(() => { refetch(); setIsDeleteDialogOpen(false); setBudgetToDelete(null); });
  };

  return (
    <div class="page">
      <div class="page-header">
        <h1>Budgets</h1>
        <button class="btn btn-primary" onClick={() => navigate('/budgets/new')}>+ New Budget</button>
      </div>

      <Show when={!budgetsResource.loading} fallback={<p style="color: var(--color-text-secondary)">Loading...</p>}>
        <Show when={budgets().length > 0} fallback={<p style="color: var(--color-text-secondary); text-align: center; padding: 40px">No budgets found</p>}>
          <div class="budgets-grid">
            <For each={budgets()}>
              {(budget) => (
                <div class="budget-card">
                  <div class="budget-card-header">
                    <h3>{budget.name}</h3>
                  </div>
                  <div class="budget-card-meta">
                    <span class={`budget-direction-badge ${budget.direction}`}>{budget.direction}</span>
                    <span style="font-size: 12px; color: var(--color-text-secondary)">{budget.period}</span>
                  </div>
                  <div class="budget-target">
                    <span class="target-label">Target</span>
                    <span class="target-value">{formatCurrency(budget.target_amount, budget.currency)}</span>
                  </div>
                  <div class="budget-card-actions">
                    <button class="btn btn-secondary" onClick={() => navigate(`/budgets/${budget.id}`)}>View</button>
                    <button class="btn btn-danger" onClick={() => { setBudgetToDelete(budget.id); setIsDeleteDialogOpen(true); }}>Delete</button>
                  </div>
                </div>
              )}
            </For>
          </div>

          <Pagination
            page={page()}
            perPage={perPage()}
            total={meta().total}
            totalPages={meta().totalPages}
            onPageChange={setPage}
          />
        </Show>
      </Show>

      <ConfirmDialog
        isOpen={isDeleteDialogOpen()}
        title="Delete Budget"
        message="Are you sure you want to delete this budget? This action cannot be undone."
        variant="danger"
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleDelete}
        onCancel={() => { setIsDeleteDialogOpen(false); setBudgetToDelete(null); }}
      />
    </div>
  );
}
