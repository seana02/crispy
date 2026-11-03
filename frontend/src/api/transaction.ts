import { fetchApi } from './client';
import { Transaction, PostingInput, TagInput, TransactionFilters, PaginatedResult } from './types';

export async function getTransactions(filters?: TransactionFilters): Promise<PaginatedResult<Transaction>> {
  const params = new URLSearchParams();
  if (filters?.date_from) params.append('date_from', filters.date_from);
  if (filters?.date_to) params.append('date_to', filters.date_to);
  if (filters?.account_id) params.append('account_id', filters.account_id);
  if (filters?.tag_key) params.append('tag_key', filters.tag_key);
  if (filters?.tag_value) params.append('tag_value', filters.tag_value);
  if (filters?.search) params.append('search', filters.search);
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.per_page) params.append('per_page', filters.per_page.toString());

  return fetchApi<PaginatedResult<Transaction>>(`/transactions${params.toString() ? '?' + params.toString() : ''}`);
}

export async function getTransaction(id: string): Promise<Transaction> {
  return fetchApi<Transaction>(`/transactions/${id}`);
}

export async function createTransaction(data: {
  date: string;
  description: string;
  status?: string;
  postings: PostingInput[];
  tags?: TagInput[];
}): Promise<Transaction> {
  return fetchApi<Transaction>('/transactions', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateTransaction(id: string, data: {
  date?: string;
  description?: string;
  status?: string;
  postings?: PostingInput[];
  tags?: TagInput[];
}): Promise<Transaction> {
  return fetchApi<Transaction>(`/transactions/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function deleteTransaction(id: string): Promise<void> {
  await fetchApi<void>(`/transactions/${id}`, {
    method: 'DELETE',
  });
}

export async function addTagToTransaction(
  transactionId: string,
  tag: TagInput
): Promise<Transaction> {
  return fetchApi<Transaction>(`/transactions/${transactionId}/tags`, {
    method: 'POST',
    body: JSON.stringify(tag),
  });
}

export async function removeTagFromTransaction(
  transactionId: string,
  tagKey: string,
  tagValue: string
): Promise<Transaction> {
  return fetchApi<Transaction>(`/transactions/${transactionId}/tags/${encodeURIComponent(tagKey)}/${encodeURIComponent(tagValue)}`, {
    method: 'DELETE',
  });
}
