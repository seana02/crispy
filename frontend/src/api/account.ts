import { fetchApi } from './client';
import { Account, PaginatedResult } from './types';

export async function getAccounts(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  type?: string;
}): Promise<PaginatedResult<Account>> {
  const p = new URLSearchParams();
  if (params?.page) p.append('page', params.page.toString());
  if (params?.per_page) p.append('per_page', params.per_page.toString());
  if (params?.search) p.append('search', params.search);
  if (params?.type) p.append('account_type', params.type);
  return fetchApi<PaginatedResult<Account>>(`/accounts${p.toString() ? '?' + p.toString() : ''}`);
}

export async function getAccount(id: string): Promise<Account> {
  return fetchApi<Account>(`/accounts/${id}`);
}

export async function searchAccounts(query: string): Promise<Account[]> {
  return fetchApi<Account[]>(`/accounts/search?q=${encodeURIComponent(query)}`);
}

export async function createAccount(data: {
  name: string;
  parent_id?: string | null;
  account_type?: string;
}): Promise<Account> {
  return fetchApi<Account>('/accounts', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateAccount(id: string, data: {
  name?: string;
  parent_id?: string | null;
  account_type?: string;
}): Promise<Account> {
  return fetchApi<Account>(`/accounts/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteAccount(id: string): Promise<void> {
  await fetchApi<void>(`/accounts/${id}`, {
    method: 'DELETE',
  });
}
