import { createMemo, Show } from 'solid-js';
import './BudgetProgress.css';

interface BudgetProgressProps {
  allocated: string;
  spent: string;
  remaining: string;
  percentage?: number;
  showLabel?: boolean;
}

export function BudgetProgress(props: BudgetProgressProps) {
  const rawPct = createMemo(() => {
    if (props.percentage !== undefined) return Math.max(0, props.percentage);
    const allocated = Number(props.allocated) || 0;
    const spent = Number(props.spent) || 0;
    if (allocated === 0) return 0;
    return (spent / allocated) * 100;
  });

  // Bar is capped at 100 for display, but rawPct is used for over-budget detection
  const displayPct = createMemo(() => Math.min(100, rawPct()));
  const over = createMemo(() => rawPct() > 100);

  const barColor = createMemo(() => {
    if (rawPct() > 100) return 'var(--color-error)';
    if (rawPct() > 80) return 'var(--color-warning)';
    return 'var(--color-success)';
  });

  return (
    <div class="budget-progress">
      <Show when={props.showLabel}>
        <div class="budget-progress-labels">
          <div class="budget-progress-label">
            <span class="label-name">Spent</span>
            <span class="label-value">{rawPct().toFixed(1)}%</span>
          </div>
          <div class="budget-progress-label">
            <span class="label-name">{over() ? 'Over' : 'Remaining'}</span>
            <span class="label-value" style={{ color: over() ? 'var(--color-error)' : undefined }}>
              {over() ? `+${((rawPct() - 100) * 0.01 * Number(props.allocated)).toFixed(2)}` : props.remaining}
            </span>
          </div>
        </div>
      </Show>
      <div class="budget-progress-bar-track">
        <div class="budget-progress-bar-fill" style={{ width: `${displayPct()}%`, 'background-color': barColor() }} />
      </div>
      <div class="budget-progress-stats">
        <div class="budget-stat"><span class="stat-name">Allocated</span><span class="stat-val">{props.allocated}</span></div>
        <div class="budget-stat"><span class="stat-name">Spent</span><span class="stat-val">{props.spent}</span></div>
        <div class="budget-stat">
          <span class="stat-name">{over() ? 'Overage' : 'Remaining'}</span>
          <span class="stat-val" style={{ color: over() ? 'var(--color-error)' : undefined }}>
            {over() ? `+${((rawPct() - 100) * 0.01 * Number(props.allocated)).toFixed(2)}` : props.remaining}
          </span>
        </div>
      </div>
    </div>
  );
}
