import { createMemo, createSignal } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getTransactions, deleteTransaction } from '../api/transaction';
import { Transaction } from '../api/types';
import { StyledList } from '../components/lists/StyledList';
import { ConfirmDialog } from '../components/layout/ConfirmDialog';
import { formatCurrency } from '../state/decimal';
import './Transactions.css';

export default function Transactions() {
    const navigate = useNavigate();

    const [search, setSearch] = createSignal('');
    const [dateFrom, setDateFrom] = createSignal('');
    const [dateTo, setDateTo] = createSignal('');
    const [statusFilter, setStatusFilter] = createSignal('');
    const [page, setPage] = createSignal(1);
    const [perPage] = createSignal(10);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = createSignal(false);
    const [transactionToDelete, setTransactionToDelete] = createSignal<string | null>(null);

    const [transactionsResource, { refetch }] = createResource(() => ({
        page: page(), per_page: perPage(),
        search: search() || undefined,
        date_from: dateFrom() || undefined,
        date_to: dateTo() || undefined,
    }), getTransactions);

    const transactions = createMemo(() => transactionsResource()?.data ?? []);
    const meta = createMemo(() => transactionsResource()?.meta ?? { page: 1, per_page: 10, total: 0, totalPages: 0 });

    const handleDelete = () => {
        const id = transactionToDelete();
        if (id) deleteTransaction(id).then(() => { refetch(); setIsDeleteDialogOpen(false); setTransactionToDelete(null); });
    };

    // Show first posting amount as the display amount
    const displayAmount = (tx: Transaction) => {
        const p = tx.postings?.[0];
        if (!p) return '';
        return formatCurrency(p.amount, p.currency);
    };

    const filterControls = (
        <div class="transactions-filters">
            <input type="date" value={dateFrom()} onInput={(e) => { setDateFrom(e.currentTarget.value); setPage(1); }} />
            <input type="date" value={dateTo()} onInput={(e) => { setDateTo(e.currentTarget.value); setPage(1); }} />
            <select value={statusFilter()} onChange={(e) => { setStatusFilter(e.currentTarget.value); setPage(1); }}>
                <option value="">All Status</option>
                <option value="cleared">Cleared</option>
                <option value="pending">Pending</option>
                <option value="posted">Posted</option>
            </select>
        </div>
    );

    const actionButtons = [
        { label: 'View', onClick: (item: Transaction) => navigate(`/transactions/${item.id}`), variant: 'secondary' as const },
        { label: 'Edit', onClick: (item: Transaction) => navigate(`/transactions/${item.id}/edit`), variant: 'secondary' as const },
        { label: 'Delete', onClick: (item: Transaction) => { setTransactionToDelete(item.id); setIsDeleteDialogOpen(true); }, variant: 'danger' as const },
    ];

    return (
        <div class="page">
            <div class="page-header">
                <h1>Transactions</h1>
                <button class="btn btn-primary" onClick={() => navigate('/transactions/new')}>+ New Transaction</button>
            </div>

            <StyledList
                items={transactions()}
                searchValue={search()}
                onSearchChange={(q) => { setSearch(q); setPage(1); }}
                searchPlaceholder="Search transactions..."
                filterControls={filterControls}
                pagination={{ page: page(), perPage: perPage(), total: meta().total, totalPages: meta().totalPages, onPageChange: setPage }}
                actions={actionButtons}
                renderRow={(tx) => (
                    <div class="transaction-row">
                        <div class="transaction-row-main">
                            <span class="transaction-date">{tx.date}</span>
                            <a class="transaction-description" href={`/transactions/${tx.id}`}>{tx.description}</a>
                        </div>
                        <span class="amount amount-positive">{displayAmount(tx)}</span>
                        {tx.status !== 'cleared' && (
                            <span class={`status-badge ${tx.status}`}>{tx.status}</span>
                        )}
                    </div>
                )}
            />

            <ConfirmDialog
                isOpen={isDeleteDialogOpen()}
                title="Delete Transaction"
                message="Are you sure you want to delete this transaction? This action cannot be undone."
                variant="danger"
                confirmText="Delete"
                cancelText="Cancel"
                onConfirm={handleDelete}
                onCancel={() => { setIsDeleteDialogOpen(false); setTransactionToDelete(null); }}
            />
        </div>
    );
}
