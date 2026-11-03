import { createMemo } from 'solid-js';

interface BreadcrumbItem {
  label: string;
  path?: string;
  active?: boolean;
}

interface BreadcrumbsProps {
  items: BreadcrumbItem[];
}

export function Breadcrumbs(props: BreadcrumbsProps) {
  const filteredItems = createMemo(() => {
    return props.items.filter((item) => item.label && !item.label.trim().startsWith('/'));
  });

  if (filteredItems().length === 0) return null;

  return (
    <nav style="padding: 12px 0;">
      <ol style="display: flex; align-items: center; gap: 6px; font-size: 14px;">
        {filteredItems().map((item, index) => (
          <li style="display: flex; align-items: center; gap: 6px; color: var(--color-text-secondary);">
            {item.active ? (
              <span style="color: var(--color-text-primary); font-weight: 600;">{item.label}</span>
            ) : (
              <a href={item.path || '#'} style="color: var(--color-text-secondary); text-decoration: none; transition: color 0.2s ease;">
                {item.label}
              </a>
            )}
            {index < filteredItems().length - 1 && (
              <span style="font-size: 12px; color: var(--color-text-secondary);">/</span>
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
}
