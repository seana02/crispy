import { fetchApi } from './client';
import { SearchResult } from './types';

export async function search(query: string): Promise<SearchResult> {
  return fetchApi<SearchResult>(`/search?q=${encodeURIComponent(query)}`);
}
