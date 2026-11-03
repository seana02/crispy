import { createMemo, createSignal } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getTags } from '../api/tags';
import { Tag } from '../api/types';
import { StyledList } from '../components/lists/StyledList';

export default function Tags() {
  const navigate = useNavigate();
  const [search, setSearch] = createSignal('');
  const [page, setPage] = createSignal(1);
  const [perPage] = createSignal(10);

  const [tagsResource] = createResource(() =>
    getTags({ page: page(), per_page: perPage(), search: search() || undefined })
  );

  const tags = createMemo(() => tagsResource()?.data ?? []);
  const meta = createMemo(() => tagsResource()?.meta ?? { page: 1, per_page: 10, total: 0, totalPages: 0 });

  const actions = [
    { label: 'View', onClick: (item: Tag) => navigate(`/tags/${encodeURIComponent(item.key)}/${encodeURIComponent(item.value)}`), variant: 'secondary' as const },
  ];

  return (
    <div class="page">
      <div class="page-header"><h1>Tags</h1></div>
      <StyledList
        items={tags()}
        searchValue={search()}
        onSearchChange={(q) => { setSearch(q); setPage(1); }}
        searchPlaceholder="Search tags..."
        pagination={{ page: page(), perPage: perPage(), total: meta().total, totalPages: meta().totalPages, onPageChange: setPage }}
        actions={actions}
        renderRow={(tag) => (
          <a class="tag-chip" href={`/tags/${encodeURIComponent(tag.key)}/${encodeURIComponent(tag.value)}`}>
            {tag.key}:{tag.value}
          </a>
        )}
      />
    </div>
  );
}
