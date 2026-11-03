import { fetchApi } from './client';
import { Budget, PaginatedResult } from './types';

export async function getBudgets(params?: {
  period?: string;
  page?: number;
  per_page?: number;
}): Promise<PaginatedResult<Budget>> {
  const url = new URL('/budgets', window.location.origin);
  if (params?.period) url.searchParams.append('period', params.period);
  if (params?.page) url.searchParams.append('page', params.page.toString());
  if (params?.per_page) url.searchParams.append('per_page', params.per_page.toString());
  return fetchApi<PaginatedResult<Budget>>(url.pathname + url.search);
}

export async function getBudget(id: string): Promise<Budget> {
  return fetchApi<Budget>(`/budgets/${id}`);
}

export async function createBudget(data: {
  name: string;
  direction: string;
  period: string;
  currency: string;
  target_amount: string;
  account_ids: number[];
  tag_ids: number[];
}): Promise<Budget> {
  return fetchApi<Budget>('/budgets', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateBudget(id: string, data: {
  name?: string;
  direction?: string;
  period?: string;
  currency?: string;
  target_amount?: string;
  account_ids?: number[];
  tag_ids?: number[];
}): Promise<Budget> {
  return fetchApi<Budget>(`/budgets/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

export async function deleteBudget(id: string): Promise<void> {
  await fetchApi<void>(`/budgets/${id}`, {
    method: 'DELETE',
  });
}

export async function getBudgetSpending(budgetId: string): Promise<{
  allocated: string;
  spent: string;
  remaining: string;
  by_category: Record<string, string>;
}> {
  return fetchApi<{
    allocated: string;
    spent: string;
    remaining: string;
    by_category: Record<string, string>;
  }>(`/budgets/${budgetId}/spending`);
}
