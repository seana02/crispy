import { createMemo, createSignal, Show, For } from 'solid-js';
import { useParams, useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getAccount, deleteAccount } from '../api/account';
import { getTransactions } from '../api/transaction';
import { Transaction } from '../api/types';
import { ConfirmDialog } from '../components/layout/ConfirmDialog';
import { formatCurrency } from '../state/decimal';
import { StyledList } from '../components/lists/StyledList';
import './AccountDetail.css';

export default function AccountDetail() {
    const params = useParams<{ id: string }>();
    const navigate = useNavigate();

    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = createSignal(false);
    const [search, setSearch] = createSignal('');
    const [page, setPage] = createSignal(1);
    const [perPage] = createSignal(10);

    const [account] = createResource(() => getAccount(params.id));
    const [transactionsResource] = createResource(() =>
        getTransactions({ page: page(), per_page: perPage(), search: search() || undefined, account_id: params.id })
    );

    const transactions = createMemo(() => transactionsResource()?.data ?? []);
    const meta = createMemo(() => transactionsResource()?.meta ?? { page: 1, per_page: 10, total: 0, totalPages: 0 });

    const handleDelete = () => {
        deleteAccount(params.id).then(() => navigate('/accounts'));
        setIsDeleteDialogOpen(false);
    };

    const displayAmount = (tx: Transaction) => {
        const p = tx.postings?.[0];
        return p ? formatCurrency(p.amount, p.currency) : '';
    };

    const actions = [
        { label: 'View', onClick: (item: Transaction) => navigate(`/transactions/${item.id}`), variant: 'secondary' as const },
    ];

    return (
        <div class="account-detail">
            <Show when={account()} fallback={<p style="color: var(--color-text-secondary)">Loading...</p>}>
                {(acc) => (
                    <>
                        <div class="account-detail-header">
                            <div>
                                <h1 class="account-detail-title">
                                    {acc().name}
                                    <span class={`type-badge ${acc().account_type}`}>{acc().account_type}</span>
                                </h1>
                                <Show when={acc().parent_id}>
                                    <p class="account-full-name">{acc().full_name}</p>
                                </Show>
                            </div>
                            <div class="account-detail-actions">
                                <button class="btn btn-primary" onClick={() => navigate(`/accounts/${acc().id}/edit`)}>Edit</button>
                                <button class="btn btn-danger" onClick={() => setIsDeleteDialogOpen(true)}>Delete</button>
                                <button class="btn btn-secondary" onClick={() => navigate('/accounts')}>Back</button>
                            </div>
                        </div>

                        <div class="account-detail-stats">
                            <div class="stat-card">
                                <div class="stat-label">Full Path</div>
                                <div class="stat-value" style="font-size: 14px">{acc().full_name}</div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-label">Type</div>
                                <div class="stat-value"><span class={`type-badge ${acc().account_type}`}>{acc().account_type}</span></div>
                            </div>
                            <div class="stat-card">
                                <div class="stat-label">Created</div>
                                <div class="stat-value" style="font-size: 14px">{acc().created_at?.slice(0, 10)}</div>
                            </div>
                        </div>

                        <h2 style="margin: 0 0 16px; font-size: 20px; color: var(--color-text-primary)">Transactions</h2>

                        <StyledList
                            items={transactions()}
                            searchValue={search()}
                            onSearchChange={(q) => { setSearch(q); setPage(1); }}
                            searchPlaceholder="Search transactions..."
                            pagination={{ page: page(), perPage: perPage(), total: meta().total, totalPages: meta().totalPages, onPageChange: (p) => { setPage(p); } }}
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
                    </>
                )}
            </Show>

            <ConfirmDialog
                isOpen={isDeleteDialogOpen()}
                title="Delete Account"
                message="Are you sure you want to delete this account? This action cannot be undone."
                variant="danger"
                confirmText="Delete"
                cancelText="Cancel"
                onConfirm={handleDelete}
                onCancel={() => setIsDeleteDialogOpen(false)}
            />
        </div>
    );
}
