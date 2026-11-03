import { createMemo, createSignal, Show } from 'solid-js';
import { createResource } from 'solid-js';
import { getBalanceSheet } from '../api/report';
import { formatCurrency } from '../state/decimal';
import { ReportTable, Column } from '../components/lists/ReportTable';
import { todayLocalDate } from '../state/dates';
import { BalanceSheetReport } from '../api/types';
import './Report.css';

export default function BalanceSheet() {
  const [date, setDate] = createSignal(todayLocalDate());
  const [periodGrouping, setPeriodGrouping] = createSignal<'weekly' | 'monthly' | 'yearly'>('monthly');
  const [depth, setDepth] = createSignal(3);
  const [tagFilter, setTagFilter] = createSignal('');
  const [view, setView] = createSignal<'table' | 'chart'>('table');

  const [report, { refetch }] = createResource<BalanceSheetReport>(() => {
    const params: Record<string, any> = { date_from: date(), date_to: date(), depth: depth() };
    if (tagFilter()) params.tag_filters = tagFilter();
    return getBalanceSheet(params);
  });

  const rows = createMemo(() => {
    const data = report();
    if (!data) return [];
    const assets = (data.asset_rows ?? []).map(r => ({ ...r, row_type: 'asset' }));
    const liabs = (data.liability_rows ?? []).map(r => ({ ...r, row_type: 'liability' }));
    return [...assets, ...liabs];
  });

  const columns: Column[] = [
    {
      key: 'name',
      label: 'Account',
      render: (item) => <a class="account-link" href={`/accounts/${item.id}`}>{item.name}</a>,
    },
    {
      key: 'row_type',
      label: 'Type',
      render: (item) => <span class={`account-type-badge ${item.row_type}`}>{item.row_type}</span>,
    },
    {
      key: 'balance',
      label: 'Balance',
      align: 'right' as const,
      render: (item) => <span class="mono">{formatCurrency(item.balance ?? '0', item.currency ?? 'USD')}</span>,
    },
  ];

  return (
    <div class="report-page">
      <h1>Balance Sheet</h1>

      <div class="report-controls">
        <div class="report-controls-left">
          <div class="report-field">
            <label>As of</label>
            <input type="date" value={date()} onInput={(e) => { setDate(e.currentTarget.value); refetch(); }} />
          </div>
          <div class="report-field">
            <label>Period Grouping</label>
            <select value={periodGrouping()} onChange={(e) => setPeriodGrouping(e.currentTarget.value as any)}>
              <option value="weekly">Weekly</option>
              <option value="monthly">Monthly</option>
              <option value="yearly">Yearly</option>
            </select>
          </div>
          <div class="report-field">
            <label>Depth</label>
            <input type="number" min="1" max="5" value={depth()} onInput={(e) => setDepth(Math.min(5, Math.max(1, parseInt(e.currentTarget.value) || 1)))} />
          </div>
          <div class="report-field">
            <label>Tag Filter</label>
            <input type="text" placeholder="key:value" value={tagFilter()} onInput={(e) => setTagFilter(e.currentTarget.value)} />
          </div>
        </div>
        <div class="report-controls-right">
          <button class={`btn ${view() === 'chart' ? 'btn-primary' : 'btn-secondary'}`} onClick={() => setView(v => v === 'table' ? 'chart' : 'table')}>
            {view() === 'table' ? 'Chart View' : 'Table View'}
          </button>
          <button class="btn btn-primary" onClick={() => refetch()} disabled={report.loading}>
            {report.loading ? 'Loading...' : 'Refresh'}
          </button>
        </div>
      </div>

      <Show when={!report.loading && report()} fallback={<div class="report-loading">Loading balance sheet...</div>}>
        <div class="report-summary">
          <div class="summary-card">
            <div class="summary-label">Total Assets</div>
            <div class="summary-value assets">{formatCurrency(report()?.total_assets ?? '0', 'USD')}</div>
          </div>
          <div class="summary-card">
            <div class="summary-label">Total Liabilities</div>
            <div class="summary-value liabilities">{formatCurrency(report()?.total_liab ?? '0', 'USD')}</div>
          </div>
          <div class="summary-card">
            <div class="summary-label">Net Worth</div>
            <div class="summary-value net-worth">{formatCurrency(report()?.net_worth ?? '0', 'USD')}</div>
          </div>
        </div>

        <Show when={view() === 'table'} fallback={<div class="chart-placeholder"><p>Chart view coming soon</p><p>This will display interactive charts for your balance sheet data</p></div>}>
          <ReportTable data={rows()} columns={columns} emptyMessage="No accounts found" />
        </Show>
      </Show>
    </div>
  );
}
