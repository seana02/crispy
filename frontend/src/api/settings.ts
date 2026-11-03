import { fetchApi } from './client';
import { Setting } from './types';

export async function getSettings(category?: string): Promise<Setting[]> {
  const params = category ? `?category=${encodeURIComponent(category)}` : '';
  return fetchApi<Setting[]>(`/settings${params}`);
}

export async function getSetting(key: string): Promise<Setting> {
  return fetchApi<Setting>(`/settings/${encodeURIComponent(key)}`);
}

export async function updateSetting(key: string, data: {
  value: string;
}): Promise<Setting> {
  return fetchApi<Setting>(`/settings/${encodeURIComponent(key)}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

export async function updateSettings(data: Record<string, string>): Promise<void> {
  await fetchApi<void>('/settings/bulk', {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}
