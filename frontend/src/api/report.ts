import { fetchApi } from './client';
import { BalanceSheetReport, IncomeStatementReport, SpendingCategory } from './types';

export async function getBalanceSheet(params?: {
  date_from?: string;
  date_to?: string;
  depth?: number;
  tag_filters?: string;
}): Promise<BalanceSheetReport> {
  const url = new URL('/reports/balance-sheet', window.location.origin);
  if (params?.date_from) url.searchParams.append('date_from', params.date_from);
  if (params?.date_to) url.searchParams.append('date_to', params.date_to);
  if (params?.depth !== undefined) url.searchParams.append('depth', params.depth.toString());
  if (params?.tag_filters) url.searchParams.append('tag_filters', params.tag_filters);
  return fetchApi<BalanceSheetReport>(url.pathname + url.search);
}

export async function getIncomeStatement(params?: {
  period_start?: string;
  period_end?: string;
  depth?: number;
  tag_filters?: string;
}): Promise<IncomeStatementReport> {
  const url = new URL('/reports/income-statement', window.location.origin);
  if (params?.period_start) url.searchParams.append('period_start', params.period_start);
  if (params?.period_end) url.searchParams.append('period_end', params.period_end);
  if (params?.depth !== undefined) url.searchParams.append('depth', params.depth.toString());
  if (params?.tag_filters) url.searchParams.append('tag_filters', params.tag_filters);
  return fetchApi<IncomeStatementReport>(url.pathname + url.search);
}

export async function getSpendingByCategory(
  date_from: string,
  date_to: string
): Promise<SpendingCategory[]> {
  return fetchApi<SpendingCategory[]>(
    `/reports/spending-by-category?date_from=${encodeURIComponent(date_from)}&date_to=${encodeURIComponent(date_to)}`
  );
}
