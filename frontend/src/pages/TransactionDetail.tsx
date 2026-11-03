import { createMemo, createSignal, Show, For, createEffect } from 'solid-js';
import { useParams, useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getTransaction, deleteTransaction, updateTransaction } from '../api/transaction';
import { formatCurrency, addDecimals } from '../state/decimal';
import { TagInput } from '../components/forms/TagInput';
import { ConfirmDialog } from '../components/layout/ConfirmDialog';
import './TransactionDetail.css';

export default function TransactionDetail() {
  const params = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = createSignal(false);
  const [isEditing, setIsEditing] = createSignal(false);
  const [editDescription, setEditDescription] = createSignal('');
  const [editDate, setEditDate] = createSignal('');
  const [editTags, setEditTags] = createSignal<{ key: string; value: string }[]>([]);

  const [tx, { refetch }] = createResource(() => getTransaction(params.id));

  createEffect(() => {
    const t = tx();
    if (t && !isEditing()) {
      setEditDescription(t.description);
      setEditDate(t.date);
      setEditTags(t.tags.map(tag => ({ key: tag.key, value: tag.value })));
    }
  });

  const handleSave = async () => {
    const t = tx();
    if (!t) return;
    await updateTransaction(t.id, { description: editDescription(), date: editDate(), tags: editTags() });
    setIsEditing(false);
    refetch();
  };

  const handleDelete = () => {
    deleteTransaction(params.id).then(() => navigate('/transactions'));
  };

  const totalAmount = createMemo(() => {
    const t = tx();
    if (!t) return '0';
    return t.postings.reduce((sum, p) => addDecimals(sum, p.amount), '0');
  });

  const firstCurrency = createMemo(() => tx()?.postings?.[0]?.currency ?? 'USD');

  return (
    <div class="tx-detail">
      <Show when={tx()} fallback={<p style="color: var(--color-text-secondary)">Loading...</p>}>
        {(t) => (
          <>
            <div class="tx-detail-header">
              <h1>Transaction</h1>
              <div class="tx-detail-actions">
                <button class="btn btn-primary" onClick={() => navigate(`/transactions/${t().id}/edit`)}>Edit</button>
                <button class="btn btn-danger" onClick={() => setIsDeleteDialogOpen(true)}>Delete</button>
                <button class="btn btn-secondary" onClick={() => navigate('/transactions')}>Back</button>
              </div>
            </div>

            <div class="tx-detail-grid">
              <div class="tx-detail-field">
                <div class="field-label">Date</div>
                <Show when={isEditing()} fallback={<div class="field-value">{t().date}</div>}>
                  <input type="date" value={editDate()} onInput={(e) => setEditDate(e.currentTarget.value)} />
                </Show>
              </div>
              <div class="tx-detail-field">
                <div class="field-label">Description</div>
                <Show when={isEditing()} fallback={<div class="field-value">{t().description}</div>}>
                  <input type="text" value={editDescription()} onInput={(e) => setEditDescription(e.currentTarget.value)} />
                </Show>
              </div>
              <div class="tx-detail-field">
                <div class="field-label">Total</div>
                <div class="field-value amount amount-positive">{formatCurrency(totalAmount(), firstCurrency())}</div>
              </div>
              <Show when={t().status !== 'cleared'}>
                <div class="tx-detail-field">
                  <div class="field-label">Status</div>
                  <div class="field-value"><span class={`status-badge ${t().status}`}>{t().status}</span></div>
                </div>
              </Show>
            </div>

            <div class="tx-detail-card">
              <div class="field-label" style="margin-bottom: 8px;">Tags</div>
              <Show when={isEditing()} fallback={
                <div class="tags-row">
                  <Show when={t().tags.length > 0} fallback={<span style="color: var(--color-text-secondary)">No tags</span>}>
                    <For each={t().tags}>
                      {(tag) => <a class="tag-chip" href={`/tags/${tag.key}/${tag.value}`}>{tag.key}:{tag.value}</a>}
                    </For>
                  </Show>
                </div>
              }>
                <TagInput value={editTags()} onChange={setEditTags} existingTags={[]} />
              </Show>
            </div>

            <div class="tx-detail-card">
              <h3>Postings</h3>
              <div class="postings-list">
                <For each={t().postings}>
                  {(posting) => (
                    <div class="posting-row">
                      <a class="posting-account" href={`/accounts/${posting.account_id}`}>
                        {posting.account_full_name || posting.account_id}
                      </a>
                      <span class="amount">{formatCurrency(posting.amount, posting.currency)}</span>
                    </div>
                  )}
                </For>
              </div>
            </div>

            <Show when={t().recurring_template_id}>
              <p style="color: var(--color-text-secondary); font-size: 13px;">
                Generated from <a class="tag-chip" href={`/recurring/${t().recurring_template_id}`}>recurring template</a>
              </p>
            </Show>
          </>
        )}
      </Show>

      <ConfirmDialog
        isOpen={isDeleteDialogOpen()}
        title="Delete Transaction"
        message="Are you sure you want to delete this transaction? This action cannot be undone."
        variant="danger"
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleDelete}
        onCancel={() => setIsDeleteDialogOpen(false)}
      />
    </div>
  );
}
