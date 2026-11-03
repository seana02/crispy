import { createMemo, createSignal, Show, For } from 'solid-js';
import { useParams, useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getRecurringTemplate, deleteRecurringTemplate, previewSchedule } from '../api/recurring';
import { RecurringTemplate, TemplatePosting } from '../api/types';
import { DatePicker } from '../components/forms/DatePicker';
import { ConfirmDialog } from '../components/layout/ConfirmDialog';
import { todayLocalDate } from '../state/dates';
import { formatCurrency } from '../state/decimal';

type ScheduleType = 'every_n_days' | 'every_n_weeks' | 'every_n_months' | 'monthly_on_day';

export default function RecurringTransactionDetail() {
  const params = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isNew = params.id === 'new';

  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = createSignal(false);

  const [templateResource] = createResource(() => isNew ? null : getRecurringTemplate(params.id));
  const template = createMemo(() => templateResource());

  // Schedule preview
  const [previewDates, setPreviewDates] = createSignal<string[]>([]);
  const [scheduleType, setScheduleType] = createSignal<ScheduleType>('every_n_months');
  const [intervalN, setIntervalN] = createSignal(1);
  const [dayOfMonth, setDayOfMonth] = createSignal<number | null>(null);
  const [startDate, setStartDate] = createSignal(todayLocalDate());
  const [endDate, setEndDate] = createSignal('');

  const handlePreview = async () => {
    const dates = await previewSchedule({
      schedule_type: scheduleType(),
      interval_n: intervalN(),
      day_of_month: dayOfMonth(),
      start_date: startDate(),
      end_date: endDate() || null,
      count: 8,
    });
    setPreviewDates(dates.scheduled_dates ?? []);
  };

  const handleDelete = () => {
    if (!isNew) deleteRecurringTemplate(params.id).then(() => navigate('/recurring'));
  };

  return (
    <div class="page">
      <div class="page-header">
        <h1>{isNew ? 'New Recurring Transaction' : 'Recurring Transaction'}</h1>
        <div style="display: flex; gap: 8px">
          <Show when={!isNew}>
            <button class="btn btn-danger" onClick={() => setIsDeleteDialogOpen(true)}>Delete</button>
          </Show>
          <button class="btn btn-secondary" onClick={() => navigate('/recurring')}>Back</button>
        </div>
      </div>

      <Show when={!isNew && template()} fallback={isNew ? null : <p style="color: var(--color-text-secondary)">Loading...</p>}>
        {(t) => (
          <div class="card" style="display: flex; flex-direction: column; gap: 12px">
            <div>
              <div style="font-size: 12px; color: var(--color-text-secondary); margin-bottom: 4px">Description</div>
              <div style="font-size: 16px; color: var(--color-text-primary)">{t().description}</div>
            </div>
            <div>
              <div style="font-size: 12px; color: var(--color-text-secondary); margin-bottom: 4px">Schedule</div>
              <div style="color: var(--color-text-primary)">{t().schedule_type.replace(/_/g, ' ')} · every {t().interval_n}{t().day_of_month ? ` on day ${t().day_of_month}` : ''}</div>
            </div>
            <div>
              <div style="font-size: 12px; color: var(--color-text-secondary); margin-bottom: 4px">Period</div>
              <div style="color: var(--color-text-primary)">{t().start_date} → {t().end_date ?? 'indefinite'}</div>
            </div>
            <Show when={t().postings.length > 0}>
              <div>
                <div style="font-size: 12px; color: var(--color-text-secondary); margin-bottom: 8px">Postings</div>
                <For each={t().postings}>
                  {(p) => (
                    <div style="display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid var(--color-border)">
                      <a class="account-link" href={`/accounts/${p.account_id}`}>{p.account_id}</a>
                      <span class="amount">{formatCurrency(p.amount, p.currency)}</span>
                    </div>
                  )}
                </For>
              </div>
            </Show>
          </div>
        )}
      </Show>

      <div class="card">
        <h3 style="margin: 0 0 16px; font-size: 16px; color: var(--color-text-primary)">Schedule Preview</h3>
        <div style="display: flex; gap: 12px; flex-wrap: wrap; align-items: flex-end; margin-bottom: 16px">
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Schedule Type</label>
            <select style="padding: 8px 12px; background: var(--color-bg-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={scheduleType()} onChange={(e) => setScheduleType(e.currentTarget.value as ScheduleType)}>
              <option value="every_n_days">Every N Days</option>
              <option value="every_n_weeks">Every N Weeks</option>
              <option value="every_n_months">Every N Months</option>
              <option value="monthly_on_day">Monthly on Day</option>
            </select>
          </div>
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Interval (N)</label>
            <input type="number" min="1" style="padding: 8px 12px; width: 80px; background: var(--color-bg-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={intervalN()} onInput={(e) => setIntervalN(parseInt(e.currentTarget.value) || 1)} />
          </div>
          <Show when={scheduleType() === 'monthly_on_day'}>
            <div style="display: flex; flex-direction: column; gap: 4px">
              <label style="font-size: 12px; color: var(--color-text-secondary)">Day of Month</label>
              <input type="number" min="1" max="31" style="padding: 8px 12px; width: 80px; background: var(--color-bg-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={dayOfMonth() ?? ''} onInput={(e) => setDayOfMonth(parseInt(e.currentTarget.value) || null)} />
            </div>
          </Show>
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Start Date</label>
            <DatePicker value={startDate()} onChange={setStartDate} />
          </div>
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">End Date (optional)</label>
            <DatePicker value={endDate()} onChange={setEndDate} />
          </div>
          <button class="btn btn-primary" onClick={handlePreview}>Preview Dates</button>
        </div>

        <Show when={previewDates().length > 0}>
          <div>
            <div style="font-size: 12px; color: var(--color-text-secondary); margin-bottom: 8px">Next {previewDates().length} occurrences:</div>
            <For each={previewDates()}>
              {(d) => <div style="padding: 4px 0; color: var(--color-text-primary); font-family: var(--font-family-mono); font-size: 13px">{d}</div>}
            </For>
          </div>
        </Show>
      </div>

      <ConfirmDialog
        isOpen={isDeleteDialogOpen()}
        title="Delete Template"
        message="Are you sure you want to delete this recurring template?"
        variant="danger"
        confirmText="Delete"
        cancelText="Cancel"
        onConfirm={handleDelete}
        onCancel={() => setIsDeleteDialogOpen(false)}
      />
    </div>
  );
}
