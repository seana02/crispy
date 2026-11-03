import { fetchApi } from './client';
import { Tag, PaginatedResult } from './types';

export async function getTags(params?: {
  page?: number;
  per_page?: number;
  search?: string;
}): Promise<PaginatedResult<Tag>> {
  const paramsObj = new URLSearchParams();
  if (params?.page) paramsObj.append('page', params.page.toString());
  if (params?.per_page) paramsObj.append('per_page', params.per_page.toString());
  if (params?.search) paramsObj.append('search', params.search);

  return fetchApi<PaginatedResult<Tag>>(`/tags${paramsObj.toString() ? '?' + paramsObj.toString() : ''}`);
}

export async function createTag(data: {
  key: string;
  value: string;
}): Promise<Tag> {
  return fetchApi<Tag>('/tags', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function deleteTag(key: string, value: string): Promise<void> {
  await fetchApi<void>(`/tags/${encodeURIComponent(key)}/${encodeURIComponent(value)}`, {
    method: 'DELETE',
  });
}
