import { createMemo, createSignal, Show, For } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getRecurringTemplates, deleteRecurringTemplate } from '../api/recurring';
import { RecurringTemplate } from '../api/types';
import { ConfirmDialog } from '../components/layout/ConfirmDialog';
import { formatDateDisplay } from '../state/dates';
import { StyledList } from '../components/lists/StyledList';

export default function RecurringTransactions() {
  const navigate = useNavigate();
  const [page, setPage] = createSignal(1);
  const [perPage] = createSignal(10);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = createSignal(false);
  const [templateToDelete, setTemplateToDelete] = createSignal<string | null>(null);

  const [templatesResource, { refetch }] = createResource(() =>
    getRecurringTemplates({ page: page(), per_page: perPage() })
  );

  const templates = createMemo(() => templatesResource()?.data ?? []);
  const meta = createMemo(() => templatesResource()?.meta ?? { page: 1, per_page: 10, total: 0, totalPages: 0 });

  const handleDelete = () => {
    const id = templateToDelete();
    if (id) deleteRecurringTemplate(id).then(() => { refetch(); setIsDeleteDialogOpen(false); setTemplateToDelete(null); });
  };

  const scheduleSummary = (t: RecurringTemplate) => {
    const type = t.schedule_type.replace(/_/g, ' ');
    return `Every ${t.interval_n} ${type}${t.day_of_month ? ` on the ${t.day_of_month}th` : ''}`;
  };

  const actions = [
    { label: 'View', onClick: (item: RecurringTemplate) => navigate(`/recurring/${item.id}`), variant: 'secondary' as const },
    { label: 'Delete', onClick: (item: RecurringTemplate) => { setTemplateToDelete(item.id); setIsDeleteDialogOpen(true); }, variant: 'danger' as const },
  ];

  return (
    <div class="page">
      <div class="page-header">
        <h1>Recurring Transactions</h1>
        <button class="btn btn-primary" onClick={() => navigate('/recurring/new')}>+ New</button>
      </div>

      <StyledList
        items={templates()}
        pagination={{ page: page(), perPage: perPage(), total: meta().total, totalPages: meta().totalPages, onPageChange: setPage }}
        actions={actions}
        emptyMessage="No recurring templates found"
        renderRow={(t) => (
          <div style="display: flex; flex-direction: column; gap: 4px">
            <span style="font-weight: 500; color: var(--color-text-primary)">{t.description}</span>
            <span style="font-size: 13px; color: var(--color-text-secondary)">{scheduleSummary(t)}</span>
            <div style="display: flex; gap: 8px; font-size: 12px; color: var(--color-text-secondary)">
              <span>Start: {formatDateDisplay(t.start_date)}</span>
              <Show when={t.end_date}><span>End: {formatDateDisplay(t.end_date!)}</span></Show>
              <span class={`status-badge ${t.is_active ? 'posted' : 'pending'}`}>{t.is_active ? 'Active' : 'Paused'}</span>
            </div>
          </div>
        )}
      />

      <ConfirmDialog
        isOpen={isDeleteDialogOpen()}
        title="Delete Template"
        message="Are you sure you want to delete this recurring template?"
        variant="danger"
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleDelete}
        onCancel={() => { setIsDeleteDialogOpen(false); setTemplateToDelete(null); }}
      />
    </div>
  );
}
