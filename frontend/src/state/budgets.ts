import { createSignal, createResource } from 'solid-js';
import { getBudgets as getBudgetsApi, createBudget as createBudgetApi, updateBudget as updateBudgetApi, deleteBudget as deleteBudgetApi } from '../api/budget';
import { Budget, PaginatedResult } from '../api/types';

export function useBudgets() {
  const [page, setPage] = createSignal(1);
  const [pageSize, setPageSize] = createSignal(10);
  const [period, setPeriod] = createSignal<string | undefined>();

  const [budgetsResource, { refetch }] = createResource<PaginatedResult<Budget>>(() => {
    const params: { page?: number; per_page?: number; period?: string } = { page: page(), per_page: pageSize() };
    if (period()) params.period = period();
    return getBudgetsApi(params);
  });

  const createBudget = async (data: {
    name: string;
    direction: string;
    period: string;
    currency: string;
    target_amount: string;
    account_ids: number[];
    tag_ids: number[];
  }) => {
    await createBudgetApi(data);
    refetch();
  };

  const updateBudget = async (id: string, data: {
    name?: string;
    direction?: string;
    period?: string;
    currency?: string;
    target_amount?: string;
    account_ids?: number[];
    tag_ids?: number[];
  }) => {
    await updateBudgetApi(id, data);
    refetch();
  };

  const deleteBudget = async (id: string) => {
    await deleteBudgetApi(id);
    refetch();
  };

  const setPeriodFilter = (newPeriod: string | undefined) => {
    setPeriod(newPeriod);
    setPage(1);
    refetch();
  };

  const setPageNumber = (newPage: number) => {
    setPage(newPage);
    refetch();
  };

  return {
    budgets: budgetsResource,
    refetch,
    createBudget,
    updateBudget,
    deleteBudget,
    period,
    setPeriod: setPeriodFilter,
    page,
    setPage: setPageNumber,
    pageSize,
    setPageSize,
    isLoading: () => budgetsResource.loading,
    total: () => budgetsResource()?.meta?.total ?? 0,
    totalPages: () => budgetsResource()?.meta?.totalPages ?? 0,
  };
}

