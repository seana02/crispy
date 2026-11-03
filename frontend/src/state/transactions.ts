import { createSignal, createResource } from 'solid-js';
import { getTransactions, createTransaction as createTransactionApi, updateTransaction as updateTransactionApi, deleteTransaction as deleteTransactionApi } from '../api/transaction';
import { TransactionFilters } from '../api/types';

export function useTransactions() {
    const [page, setPage] = createSignal(1);
    const [pageSize, setPageSize] = createSignal(10);
    const [filters, setFilters] = createSignal<TransactionFilters>({});

    const [transactionsResource, { refetch }] = createResource(() => {
        const params: TransactionFilters = { page: page(), per_page: pageSize(), ...filters() };
        return getTransactions(params);
    });

    const createTransaction = async (data: {
        date: string;
        description: string;
        postings: Array<{ account_id: number; amount: string; currency: string }>;
        tags?: Array<{ key: string; value: string }>;
    }) => {
        await createTransactionApi(data);
        refetch();
    };

    const updateTransaction = async (id: string, data: {
        date?: string;
        description?: string;
        postings?: Array<{ account_id: number; amount: string; currency: string }>;
        tags?: Array<{ key: string; value: string }>;
    }) => {
        await updateTransactionApi(id, data);
        refetch();
    };

    const deleteTransaction = async (id: string) => {
        await deleteTransactionApi(id);
        refetch();
    };

    const setFilter = (key: keyof TransactionFilters, value: string | number | undefined) => {
        setFilters((prev) => {
            const next = { ...prev };
            if (value === undefined || value === '') {
                delete next[key];
            } else {
                (next as Record<string, string | number>)[key] = value;
            }
            return next;
        });
        setPage(1);
        refetch();
    };

    const setPageNumber = (newPage: number) => {
        setPage(newPage);
        refetch();
    };

    return {
        transactions: transactionsResource,
        refetch,
        createTransaction,
        updateTransaction,
        deleteTransaction,
        filters,
        setFilter,
        page,
        setPage: setPageNumber,
        pageSize,
        setPageSize,
        isLoading: () => transactionsResource.loading,
        total: () => transactionsResource()?.meta?.total ?? 0,
        totalPages: () => transactionsResource()?.meta?.totalPages ?? 0,
    };
}

