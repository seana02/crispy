import { createMemo, createSignal } from 'solid-js';
import { useParams, useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getTransactions } from '../api/transaction';
import { Transaction } from '../api/types';
import { formatCurrency } from '../state/decimal';
import { StyledList } from '../components/lists/StyledList';

export default function TagDetail() {
  const routeParams = useParams<{ key: string; value: string }>();
  const navigate = useNavigate();
  const [page, setPage] = createSignal(1);
  const [perPage] = createSignal(10);

  const [transactionsResource] = createResource(() =>
    getTransactions({ page: page(), per_page: perPage(), tag_key: routeParams.key, tag_value: routeParams.value })
  );

  const transactions = createMemo(() => transactionsResource()?.data ?? []);
  const meta = createMemo(() => transactionsResource()?.meta ?? { page: 1, per_page: 10, total: 0, totalPages: 0 });

  const displayAmount = (tx: Transaction) => {
    const p = tx.postings?.[0];
    return p ? formatCurrency(p.amount, p.currency) : '';
  };

  const actions = [
    { label: 'View', onClick: (item: Transaction) => navigate(`/transactions/${item.id}`), variant: 'secondary' as const },
  ];

  return (
    <div class="page">
      <div class="page-header">
        <div>
          <h1>Tag: <span class="tag-chip">{routeParams.key}:{routeParams.value}</span></h1>
          <p style="margin-top: 4px; font-size: 13px; color: var(--color-text-secondary)">Transactions with this tag</p>
        </div>
        <button class="btn btn-secondary" onClick={() => navigate('/tags')}>Back to Tags</button>
      </div>

      <StyledList
        items={transactions()}
        pagination={{ page: page(), perPage: perPage(), total: meta().total, totalPages: meta().totalPages, onPageChange: setPage }}
        actions={actions}
        renderRow={(tx) => (
          <div class="tx-row">
            <div class="tx-row-info">
              <span class="tx-date">{tx.date}</span>
              <a class="tx-desc" href={`/transactions/${tx.id}`}>{tx.description}</a>
            </div>
            <span class="amount amount-positive">{displayAmount(tx)}</span>
          </div>
        )}
      />
    </div>
  );
}
