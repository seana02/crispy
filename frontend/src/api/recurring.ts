import { fetchApi } from './client';
import { RecurringTemplate, TemplatePosting, PaginatedResult } from './types';

export async function getRecurringTemplates(params?: {
  page?: number;
  per_page?: number;
}): Promise<PaginatedResult<RecurringTemplate>> {
  const url = new URL('/recurring-templates', window.location.origin);
  if (params?.page) url.searchParams.append('page', params.page.toString());
  if (params?.per_page) url.searchParams.append('per_page', params.per_page.toString());
  return fetchApi<PaginatedResult<RecurringTemplate>>(url.pathname + url.search);
}

export async function getRecurringTemplate(id: string): Promise<RecurringTemplate> {
  return fetchApi<RecurringTemplate>(`/recurring-templates/${id}`);
}

export async function createRecurringTemplate(data: {
  name: string;
  frequency: string;
  start_date: string;
  end_date?: string;
  postings: TemplatePosting[];
}): Promise<RecurringTemplate> {
  return fetchApi<RecurringTemplate>('/recurring-templates', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function updateRecurringTemplate(
  id: string,
  data: {
    name?: string;
    frequency?: string;
    start_date?: string;
    end_date?: string;
    postings?: TemplatePosting[];
  }
): Promise<RecurringTemplate> {
  return fetchApi<RecurringTemplate>(`/recurring-templates/${id}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  });
}

export async function deleteRecurringTemplate(id: string): Promise<void> {
  await fetchApi<void>(`/recurring-templates/${id}`, {
    method: 'DELETE',
  });
}

export async function createTransactionFromTemplate(id: string): Promise<void> {
  await fetchApi<void>(`/recurring-templates/${id}/execute`, {
    method: 'POST',
  });
}

export async function previewSchedule(data: {
  schedule_type: string;
  interval_n: number;
  day_of_month?: number | null;
  start_date: string;
  end_date?: string | null;
  count?: number;
}): Promise<{ scheduled_dates: string[] }> {
  return fetchApi<{ scheduled_dates: string[] }>('/recurring-templates/preview', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}
