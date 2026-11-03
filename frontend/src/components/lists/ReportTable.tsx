import { createSignal, createMemo, For, Show, JSX } from 'solid-js';
import { Pagination } from './Pagination';
import './ReportTable.css';

export interface Column {
  key: string;
  label: string;
  render?: (item: any, index: number) => JSX.Element;
  sortable?: boolean;
  align?: 'left' | 'center' | 'right';
  width?: string;
}

interface ReportTableProps<T> {
  data: T[];
  columns: Column[];
  sorting?: { key: string; direction: 'asc' | 'desc' };
  onSort?: (key: string, direction: 'asc' | 'desc') => void;
  pagination?: {
    page: number;
    perPage: number;
    total: number;
    totalPages: number;
    onPageChange: (page: number) => void;
  };
  onSearch?: (query: string) => void;
  search?: string;
  loading?: boolean;
  emptyMessage?: string;
}

export function ReportTable<T>(props: ReportTableProps<T>) {
  const [searchValue, setSearchValue] = createSignal(props.search ?? '');

  const sortedData = createMemo(() => {
    if (!props.sorting?.key) return props.data;
    const { key, direction } = props.sorting;
    return [...props.data].sort((a, b) => {
      const aVal = (a as any)[key];
      const bVal = (b as any)[key];
      if (aVal < bVal) return direction === 'asc' ? -1 : 1;
      if (aVal > bVal) return direction === 'asc' ? 1 : -1;
      return 0;
    });
  });

  const handleSort = (key: string) => {
    if (!props.onSort) return;
    const newDir = props.sorting?.key === key && props.sorting.direction === 'asc' ? 'desc' : 'asc';
    props.onSort(key, newDir);
  };

  const colSpan = () => props.columns.length;

  return (
    <div class="report-table-container">
      <Show when={props.onSearch}>
        <div class="report-table-header">
          <div class="report-table-search">
            <input
              type="text"
              placeholder="Search..."
              value={searchValue()}
              onInput={(e) => { setSearchValue(e.currentTarget.value); props.onSearch!(e.currentTarget.value); }}
            />
          </div>
        </div>
      </Show>

      <div class="report-table-scroll">
        <table class="report-table">
          <thead>
            <tr>
              <For each={props.columns}>
                {(col) => (
                  <th
                    class={`${col.align ? `align-${col.align}` : ''} ${col.sortable ? 'sortable' : ''}`}
                    style={col.width ? { width: col.width } : undefined}
                    onClick={() => col.sortable && handleSort(col.key)}
                  >
                    {col.label}
                  </th>
                )}
              </For>
            </tr>
          </thead>
          <tbody>
            <Show
              when={!props.loading}
              fallback={
                <tr><td class="loading-cell" colspan={colSpan()}>Loading...</td></tr>
              }
            >
              <Show
                when={sortedData().length > 0}
                fallback={
                  <tr><td class="empty-cell" colspan={colSpan()}>{props.emptyMessage ?? 'No data found'}</td></tr>
                }
              >
                <For each={sortedData()}>
                  {(item, index) => (
                    <tr>
                      <For each={props.columns}>
                        {(col) => (
                          <td class={col.align ? `align-${col.align}` : ''}>
                            {col.render ? col.render(item, index()) : (item as any)[col.key]}
                          </td>
                        )}
                      </For>
                    </tr>
                  )}
                </For>
              </Show>
            </Show>
          </tbody>
        </table>
      </div>

      <Show when={props.pagination}>
        <Pagination
          page={props.pagination!.page}
          perPage={props.pagination!.perPage}
          total={props.pagination!.total}
          totalPages={props.pagination!.totalPages}
          onPageChange={props.pagination!.onPageChange}
        />
      </Show>
    </div>
  );
}
