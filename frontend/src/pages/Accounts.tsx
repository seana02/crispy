import { batch, createEffect, createMemo, createSignal, onCleanup, Show } from 'solid-js';
import { A, useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getAccounts, deleteAccount } from '../api/account';
import { Account } from '../api/types';
import { StyledList } from '../components/lists/StyledList';
import { ConfirmDialog } from '../components/layout/ConfirmDialog';
import './Accounts.css';

export default function Accounts() {
    const navigate = useNavigate();
    const [search, setSearch] = createSignal('');
    const [debouncedSearch, setDebouncedSearch] = createSignal('');
    const [typeFilter, setTypeFilter] = createSignal('');
    const [page, setPage] = createSignal(1);
    const [perPage] = createSignal(10);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = createSignal(false);
    const [isDeleting, setIsDeleting] = createSignal(false);
    const [deleteError, setDeleteError] = createSignal<string | null>(null);
    const [accountToDelete, setAccountToDelete] = createSignal<string | null>(null);

    const [accountsResource, { refetch }] = createResource(
        () => ({
            page: page(),
            per_page: perPage(),
            search: debouncedSearch() || undefined,
            type: typeFilter() || undefined,
        }),
        params => getAccounts(params)
    );

    const accounts = createMemo(() => accountsResource.error ? [] : accountsResource()?.data ?? []);
    const meta = createMemo(() => accountsResource.error ? { page: 1, per_page: 10, total: 0, totalPages: 0 } : accountsResource()?.meta ?? { page: 1, per_page: 10, total: 0, totalPages: 0 });

    const closeDialog = () => {
        setIsDeleteDialogOpen(false);
        setAccountToDelete(null);
        setDeleteError(null);
    }

    const handleDelete = async () => {
        const id = accountToDelete();
        if (!id || isDeleting()) return;
        setIsDeleting(true);
        setDeleteError(null);
        try {
            await deleteAccount(id);
            const wasLastOnPage = accounts().length === 1 && page() > 1;
            closeDialog();
            if (wasLastOnPage) setPage(p => p - 1);
            else refetch();
        } catch (e) {
            setDeleteError(e instanceof Error ? e.message : 'Delete failed');
        } finally {
            setIsDeleting(false);
        }
    };

    createEffect(() => {
        const q = search();
        const t = setTimeout(() => {
            if (q !== debouncedSearch()) batch(() => { setDebouncedSearch(q); setPage(1); });
        }, 250);
        onCleanup(() => clearTimeout(t));
    });

    const filterControls = (
        <select class="accounts-filter-select" value={typeFilter()} onChange={(e) => { setTypeFilter(e.currentTarget.value); setPage(1); }} aria-label="Filter by account type">
            <option value="">All Types</option>
            <option value="asset">Asset</option>
            <option value="liability">Liability</option>
            <option value="revenue">Revenue</option>
            <option value="expense">Expense</option>
            <option value="equity">Equity</option>
        </select>
    );

    const actions = [
        { label: 'View', onClick: (item: Account) => navigate(`/accounts/${item.id}`), variant: 'secondary' as const },
        { label: 'Edit', onClick: (item: Account) => navigate(`/accounts/${item.id}/edit`), variant: 'secondary' as const },
        { label: 'Delete', onClick: (item: Account) => { setAccountToDelete(item.id); setIsDeleteDialogOpen(true); }, variant: 'danger' as const },
    ];

    return (
        <div class="page">
            <div class="page-header">
                <h1>Accounts</h1>
                <button class="btn btn-primary" onClick={() => navigate('/accounts/new')}>+ New Account</button>
            </div>

            <Show when={accountsResource.error}>
                <div class="error-banner" role="alert">
                    Couldn't load accounts.
                </div>
            </Show>

            <StyledList
                items={accounts()}
                searchValue={search()}
                onSearchChange={(q) => setSearch(q)}
                searchPlaceholder="Search accounts..."
                filterControls={filterControls}
                pagination={{ page: page(), perPage: perPage(), total: meta().total, totalPages: meta().totalPages, onPageChange: setPage }}
                actions={actions}
                loading={accountsResource.loading && !accountsResource()}
                renderRow={(account) => {
                    return (
                        <div class="account-row">
                            <Show when={account.account_type}>
                                <span class={`type-badge ${account.account_type}`}>{account.account_type}</span>
                            </Show>
                            <A class="account-short-name" href={`/accounts/${account.id}`}>{account.full_name}</A>
                        </div>
                    );
                }}
            />

            <ConfirmDialog
                isOpen={isDeleteDialogOpen()}
                title="Delete Account"
                message="Are you sure you want to delete this account? This action cannot be undone."
                variant="danger"
                confirmText="Delete"
                cancelText="Cancel"
                onConfirm={handleDelete}
                onCancel={closeDialog}
                error={deleteError()}
                isLoading={isDeleting()}
            />
        </div>
    );
}
