export interface Account {
  id: string;
  name: string;
  full_name: string;
  parent_id: string | null;
  account_type: string;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: string;
  date: string;
  description: string;
  status: string;
  recurring_template_id: string | null;
  postings: Posting[];
  tags: Tag[];
  created_at: string;
  updated_at: string;
}

export interface Posting {
  id: string;
  transaction_id: string;
  account_id: string;
  account_full_name: string; // populated on reads
  currency: string;
  amount: string;
}

export interface Tag {
  id: string;
  key: string;
  value: string;
}

export interface Budget {
  id: string;
  name: string;
  direction: string;
  period: string;
  currency: string;
  target_amount: string;
  account_ids: string[];
  account_names: string[]; // full_names of tracked accounts
  tag_ids: string[];
  created_at: string;
  updated_at: string;
}

export interface ExchangeRate {
  id: string;
  base_currency: string;
  quote_currency: string;
  rate: string;
  effective_date: string;
  source: string;
  created_at: string;
}

export interface RecurringTemplate {
  id: string;
  description: string;
  status: string;
  schedule_type: string;
  interval_n: number;
  day_of_month: number | null;
  start_date: string;
  end_date: string | null;
  last_generated_date: string | null;
  is_active: boolean;
  postings: TemplatePosting[];
  tags: Tag[];
  created_at: string;
  updated_at: string;
}

export interface TemplatePosting {
  id: string;
  template_id: string;
  account_id: string;
  account_full_name: string;
  currency: string;
  amount: string;
}

export interface WidgetConfig {
  id: string;
  widget_type: string;
  position_x: number;
  position_y: number;
  width: number;
  height: number;
  config: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface Setting {
  key: string;
  value: string;
}

export interface PaginatedResult<T> {
  data: T[];
  meta: {
    page: number;
    per_page: number;
    total: number;
    totalPages: number;
  };
}

export interface SearchResult {
  accounts: Account[];
  transactions: Transaction[];
  tags: Tag[];
}

export interface SpendingCategory {
  account_id: string;
  account_name: string;
  amount: string;
  currency: string;
}

export interface PostingInput {
  account_id: number;
  account_name?: string; // optional — if provided and account_id is 0, backend creates the account
  currency: string;
  amount: string;
}

export interface TagInput {
  key: string;
  value: string;
}

export interface TransactionFilters {
  date_from?: string;
  date_to?: string;
  account_id?: string;
  tag_key?: string;
  tag_value?: string;
  search?: string;
  page?: number;
  per_page?: number;
}

export type Currency = string;
export type Amount = string;
export type DateString = string;

export interface BalanceSheetReport {
  asset_rows: BalanceSheetRow[];
  liability_rows: BalanceSheetRow[];
  total_assets: string;
  total_liab: string;
  net_worth?: string;
  comparison_date?: string;
  comparison_net_worth?: string;
  trend?: string;
  trend_percentage?: string;
}

export interface BalanceSheetRow {
  id: string;
  name: string;
  full_name: string;
  balance: string;
  type: 'asset' | 'liability';
  currency: string;
  depth?: number;
}

export interface IncomeStatementReport {
  revenue_rows: IncomeStatementRow[];
  expense_rows: IncomeStatementRow[];
  total_revenue: string;
  total_expense: string;
  net_income: string;
  period_start: string;
  period_end: string;
}

export interface IncomeStatementRow {
  id: string;
  name: string;
  full_name: string;
  balance: string;
  type: 'revenue' | 'expense';
  currency: string;
  depth?: number;
}
