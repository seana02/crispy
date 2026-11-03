import { createMemo, JSX, For, Show } from 'solid-js';
import { Pagination } from './Pagination';
import './StyledList.css';

interface Action<T> {
  label: string;
  onClick: (item: T, event: MouseEvent) => void;
  icon?: string;
  variant?: 'primary' | 'secondary' | 'danger';
}

interface StyledListProps<T> {
  items: T[];
  renderRow: (item: T) => JSX.Element;
  searchValue?: string;
  onSearchChange?: (value: string) => void;
  searchPlaceholder?: string;
  filterControls?: JSX.Element;
  emptyMessage?: string;
  loading?: boolean;
  actions?: Action<T>[];
  pagination?: {
    page: number;
    perPage: number;
    total: number;
    totalPages: number;
    onPageChange: (page: number) => void;
  };
}

export function StyledList<T>(props: StyledListProps<T>) {
  const handleSearch = (e: Event) => {
    const target = e.target as HTMLInputElement;
    props.onSearchChange?.(target.value);
  };

  return (
    <div class="styled-list">
      <div class="styled-list-header">
        <Show when={props.onSearchChange}>
          <div class="styled-list-search">
            <input
              type="text"
              placeholder={props.searchPlaceholder ?? 'Search...'}
              value={props.searchValue ?? ''}
              onInput={handleSearch}
            />
          </div>
        </Show>
        <Show when={props.filterControls}>
          <div class="styled-list-filters">{props.filterControls}</div>
        </Show>
      </div>

      <div class="styled-list-items">
        <Show
          when={!props.loading}
          fallback={<div class="styled-list-loading">Loading...</div>}
        >
          <Show
            when={props.items.length > 0}
            fallback={<div class="styled-list-empty">{props.emptyMessage ?? 'No items found'}</div>}
          >
            <For each={props.items}>
              {(item) => (
                <div class="styled-list-row">
                  <div class="styled-list-row-content">{props.renderRow(item)}</div>
                  <Show when={props.actions && props.actions.length > 0}>
                    <div class="styled-list-row-actions">
                      <For each={props.actions}>
                        {(action) => (
                          <button
                            class={`btn-action btn-${action.variant ?? 'secondary'}`}
                            onClick={(e) => { e.stopPropagation(); action.onClick(item, e); }}
                            type="button"
                            aria-label={action.label}
                          >
                            {action.icon ?? action.label}
                          </button>
                        )}
                      </For>
                    </div>
                  </Show>
                </div>
              )}
            </For>
          </Show>
        </Show>
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
