import { fetchApi } from './client';
import { ExchangeRate } from './types';

export async function getExchangeRates(params?: {
  base_currency?: string;
  quote_currency?: string;
  date_from?: string;
  date_to?: string;
}): Promise<ExchangeRate[]> {
  const p = new URLSearchParams();
  if (params?.base_currency) p.append('base_currency', params.base_currency);
  if (params?.quote_currency) p.append('quote_currency', params.quote_currency);
  if (params?.date_from) p.append('date_from', params.date_from);
  if (params?.date_to) p.append('date_to', params.date_to);
  return fetchApi<ExchangeRate[]>(`/exchange-rates${p.toString() ? '?' + p.toString() : ''}`);
}

export async function getExchangeRate(base: string, quote: string, date: string): Promise<ExchangeRate> {
  return fetchApi<ExchangeRate>(
    `/exchange-rates/${encodeURIComponent(base)}/${encodeURIComponent(quote)}/${encodeURIComponent(date)}`
  );
}

export async function createExchangeRate(data: {
  base_currency: string;
  quote_currency: string;
  rate: string;
  effective_date: string;
  source?: string;
}): Promise<ExchangeRate> {
  return fetchApi<ExchangeRate>('/exchange-rates', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function deleteExchangeRate(id: string): Promise<void> {
  await fetchApi<void>(`/exchange-rates/${id}`, {
    method: 'DELETE',
  });
}
