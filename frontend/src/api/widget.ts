import { fetchApi } from './client';
import { WidgetConfig, PaginatedResult } from './types';

export async function getWidgets(params?: {
  page?: number;
  per_page?: number;
}): Promise<PaginatedResult<WidgetConfig>> {
  const url = new URL('/widgets', window.location.origin);
  if (params?.page) url.searchParams.append('page', params.page.toString());
  if (params?.per_page) url.searchParams.append('per_page', params.per_page.toString());
  return fetchApi<PaginatedResult<WidgetConfig>>(url.pathname + url.search);
}

export async function getWidget(id: string): Promise<WidgetConfig> {
  return fetchApi<WidgetConfig>(`/widgets/${id}`);
}

export async function createWidget(data: {
  type: string;
  title: string;
  settings: Record<string, any>;
  layout: {
    x: number;
    y: number;
    w: number;
    h: number;
  };
}): Promise<WidgetConfig> {
  return fetchApi<WidgetConfig>('/widgets', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateWidget(
  id: string,
  data: {
    title?: string;
    settings?: Record<string, any>;
    layout?: {
      x?: number;
      y?: number;
      w?: number;
      h?: number;
    };
  }
): Promise<WidgetConfig> {
  return fetchApi<WidgetConfig>(`/widgets/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

export async function deleteWidget(id: string): Promise<void> {
  await fetchApi<void>(`/widgets/${id}`, {
    method: 'DELETE',
  });
}

export async function updateWidgetOrder(widgets: { id: string; x: number; y: number }[]): Promise<void> {
  await fetchApi<void>('/widgets/order', {
    method: 'PATCH',
    body: JSON.stringify(widgets),
  });
}
