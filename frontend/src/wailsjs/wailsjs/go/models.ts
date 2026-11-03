export namespace model {
	
	export interface Account {
	    id: number;
	    name: string;
	    full_name: string;
	    parent_id?: number;
	    account_type: string;
	    created_at: string;
	    updated_at: string;
	}
	export interface StringDecimal {
	
	}
	export interface Budget {
	    id: number;
	    name: string;
	    direction: string;
	    period: string;
	    currency: string;
	    target_amount: StringDecimal;
	    created_at: string;
	    updated_at: string;
	    account_ids: number[];
	    tag_ids: number[];
	}
	export interface ExchangeRate {
	    id: number;
	    base_currency: string;
	    quote_currency: string;
	    rate: StringDecimal;
	    effective_date: string;
	    source: string;
	    created_at: string;
	}
	export interface PaginatedResult_github_com_crispy_2_backend_internal_model_Account_ {
	    items: Account[];
	    total: number;
	    page: number;
	    page_size: number;
	    total_pages: number;
	}
	export interface PaginatedResult_github_com_crispy_2_backend_internal_model_ExchangeRate_ {
	    items: ExchangeRate[];
	    total: number;
	    page: number;
	    page_size: number;
	    total_pages: number;
	}
	export interface Tag {
	    id: number;
	    key: string;
	    value: string;
	}
	export interface Posting {
	    id: number;
	    transaction_id: number;
	    account_id: number;
	    currency: string;
	    amount: StringDecimal;
	}
	export interface Transaction {
	    id: number;
	    date: string;
	    description: string;
	    status: string;
	    recurring_template_id?: number;
	    created_at: string;
	    updated_at: string;
	    postings: Posting[];
	    tags: Tag[];
	}
	export interface PaginatedResult_github_com_crispy_2_backend_internal_model_Transaction_ {
	    items: Transaction[];
	    total: number;
	    page: number;
	    page_size: number;
	    total_pages: number;
	}
	
	export interface TemplatePosting {
	    id: number;
	    template_id: number;
	    account_id: number;
	    currency: string;
	    amount: StringDecimal;
	}
	export interface RecurringTemplate {
	    id: number;
	    description: string;
	    status: string;
	    schedule_type: string;
	    interval_n: number;
	    day_of_month?: number;
	    start_date: string;
	    end_date?: string;
	    last_generated_date?: string;
	    is_active: boolean;
	    created_at: string;
	    updated_at: string;
	    postings: TemplatePosting[];
	    tags: Tag[];
	}
	export interface Setting {
	    key: string;
	    value: string;
	}
	
	
	
	
	export interface WidgetConfig {
	    id: number;
	    widget_type: string;
	    position_x: number;
	    position_y: number;
	    width: number;
	    height: number;
	    config: string;
	    created_at: string;
	    updated_at: string;
	}

}

export namespace service {
	
	export interface BalanceSheetRow {
	    account_id: number;
	    account_name: string;
	    balance: model.StringDecimal;
	    currency: string;
	    depth: number;
	}
	export interface BalanceSheetReport {
	    asset_rows: BalanceSheetRow[];
	    liability_rows: BalanceSheetRow[];
	    total_assets: model.StringDecimal;
	    total_liabilities: model.StringDecimal;
	}
	
	export interface IncomeStatementRow {
	    account_id: number;
	    account_name: string;
	    amount: model.StringDecimal;
	    currency: string;
	    depth: number;
	}
	export interface IncomeStatementReport {
	    revenue_rows: IncomeStatementRow[];
	    expense_rows: IncomeStatementRow[];
	    total_revenue: model.StringDecimal;
	    total_expenses: model.StringDecimal;
	    net_income: model.StringDecimal;
	}
	
	export interface NetWorthWithTrend {
	    as_of_date: string;
	    net_worth: model.StringDecimal;
	    comparison_date?: string;
	    comparison_net_worth?: model.StringDecimal;
	    trend?: model.StringDecimal;
	    trend_percentage?: string;
	}
	export interface PostingInput {
	    account_id: number;
	    currency: string;
	    amount: string;
	}
	export interface SearchResult {
	    accounts: model.Account[];
	    transactions: model.Transaction[];
	    tags: model.Tag[];
	}
	export interface SpendingCategory {
	    account_id: number;
	    account_name: string;
	    amount: model.StringDecimal;
	    currency: string;
	}
	export interface TagInput {
	    key: string;
	    value: string;
	}
	export interface TransactionFilters {
	    Description: string;
	    DateFrom: string;
	    DateTo: string;
	    AccountID?: number;
	    TagKey: string;
	    TagValue: string;
	    Status: string;
	    Currency: string;
	}

}

