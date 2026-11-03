import { createMemo, createSignal, For, Show } from 'solid-js';
import { createResource } from 'solid-js';
import { getSettings, updateSettings } from '../api/settings';

export default function Settings() {
  const [saving, setSaving] = createSignal(false);
  const [settingsResource, { refetch }] = createResource(() => getSettings());
  const settings = createMemo(() => settingsResource() ?? []);

  // Local editable state — record of key → value
  const [local, setLocal] = createSignal<Record<string, string>>({});

  // Sync local state when settings load
  createMemo(() => {
    const state: Record<string, string> = {};
    settings().forEach((s) => { state[s.key] = s.value; });
    setLocal(state);
  });

  const get = (key: string) => local()[key] ?? '';
  const set = (key: string, value: string) => setLocal(prev => ({ ...prev, [key]: value }));

  const handleSave = async () => {
    setSaving(true);
    try { await updateSettings(local()); refetch(); }
    finally { setSaving(false); }
  };

  const currencies = ['USD', 'EUR', 'GBP', 'JPY', 'CAD', 'AUD', 'CHF', 'CNY'];
  const dateFormats = ['YYYY-MM-DD', 'MM/DD/YYYY', 'DD/MM/YYYY', 'DD.MM.YYYY'];
  const months = ['January','February','March','April','May','June','July','August','September','October','November','December'];

  return (
    <div class="page">
      <div class="page-header">
        <h1>Settings</h1>
        <button class="btn btn-primary" onClick={handleSave} disabled={saving()}>
          {saving() ? 'Saving...' : 'Save'}
        </button>
      </div>

      <div class="card">
        <h2 style="margin: 0 0 16px; font-size: 18px; color: var(--color-text-primary)">Currency</h2>
        <div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 16px">
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Default Currency</label>
            <select style="padding: 8px 12px; background: var(--color-bg-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={get('default_currency')} onChange={(e) => set('default_currency', e.currentTarget.value)}>
              <For each={currencies}>{(c) => <option value={c}>{c}</option>}</For>
            </select>
          </div>
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Display Currency</label>
            <select style="padding: 8px 12px; background: var(--color-bg-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={get('display_currency')} onChange={(e) => set('display_currency', e.currentTarget.value)}>
              <For each={currencies}>{(c) => <option value={c}>{c}</option>}</For>
            </select>
          </div>
        </div>
      </div>

      <div class="card">
        <h2 style="margin: 0 0 16px; font-size: 18px; color: var(--color-text-primary)">Display</h2>
        <div style="display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 16px">
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Date Format</label>
            <select style="padding: 8px 12px; background: var(--color-bg-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={get('date_format')} onChange={(e) => set('date_format', e.currentTarget.value)}>
              <For each={dateFormats}>{(f) => <option value={f}>{f}</option>}</For>
            </select>
          </div>
          <div style="display: flex; flex-direction: column; gap: 4px">
            <label style="font-size: 12px; color: var(--color-text-secondary)">Fiscal Year Start Month</label>
            <select style="padding: 8px 12px; background: var(--color-bg-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={get('fiscal_year_start_month')} onChange={(e) => set('fiscal_year_start_month', e.currentTarget.value)}>
              <For each={months}>{(m, i) => <option value={String(i() + 1)}>{m}</option>}</For>
            </select>
          </div>
        </div>
      </div>

      <div class="card">
        <h2 style="margin: 0 0 16px; font-size: 18px; color: var(--color-text-primary)">Appearance</h2>
        <div style="display: flex; flex-direction: column; gap: 4px; max-width: 220px">
          <label style="font-size: 12px; color: var(--color-text-secondary)">Theme</label>
          <select style="padding: 8px 12px; background: var(--color-bg-surface); border: 1px solid var(--color-border); border-radius: var(--radius-md); color: var(--color-text-primary); font-family: var(--font-family)" value={get('theme')} onChange={(e) => set('theme', e.currentTarget.value)}>
            <option value="dark">Neon Night (Dark)</option>
          </select>
        </div>
      </div>
    </div>
  );
}
