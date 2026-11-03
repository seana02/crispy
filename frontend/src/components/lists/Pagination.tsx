import { createSignal, createMemo, on, For } from 'solid-js';
import './Pagination.css';

interface PaginationProps {
  page: number;
  perPage: number;
  total: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  onPerPageChange?: (perPage: number) => void;
}

export function Pagination(props: PaginationProps) {
  const [currentPage, setCurrentPage] = createSignal(props.page);

  createMemo(on(() => props.page, (val) => setCurrentPage(val)));

  const pages = createMemo<(number | '...')[]>(() => {
    const total = props.totalPages;
    if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);
    const cur = currentPage();
    if (cur <= 4) return [1, 2, 3, 4, 5, '...', total];
    if (cur >= total - 3) return [1, '...', total - 4, total - 3, total - 2, total - 1, total];
    return [1, '...', cur - 1, cur, cur + 1, '...', total];
  });

  const go = (page: number | '...') => {
    if (page === '...' || page === currentPage()) return;
    setCurrentPage(page);
    props.onPageChange(page);
  };

  return (
    <div class="pagination">
      <div class="pagination-pages">
        <button class="pagination-btn" onClick={() => go(currentPage() - 1)} disabled={currentPage() <= 1} type="button">«</button>
        <For each={pages()}>
          {(p) => (
            <button
              class={`pagination-btn${p === currentPage() ? ' active' : ''}`}
              onClick={() => go(p)}
              disabled={p === '...'}
              type="button"
            >
              {p}
            </button>
          )}
        </For>
        <button class="pagination-btn" onClick={() => go(currentPage() + 1)} disabled={currentPage() >= props.totalPages} type="button">»</button>
      </div>
      <div class="pagination-info">
        Page {currentPage()} of {props.totalPages} ({props.total} items)
      </div>
      <div class="pagination-per-page">
        <span>Per page:</span>
        <select value={props.perPage} onChange={(e) => props.onPerPageChange?.(Number(e.currentTarget.value))}>
          <option value={10}>10</option>
          <option value={25}>25</option>
          <option value={50}>50</option>
          <option value={100}>100</option>
        </select>
      </div>
    </div>
  );
}
