import { createSignal, createMemo, createEffect } from 'solid-js';
import { TextInput } from './TextInput';
import { SelectInput } from './SelectInput';
import { DatePicker } from './DatePicker';
import { AutocompleteInput } from './AutocompleteInput';
import { TagInput } from './TagInput';
import { formatCurrency, isZero } from '../../state/decimal';
import Big from 'big.js';
import { todayLocalDate } from '../../state/dates';

interface Posting {
  id?: string;
  account_id: string;
  amount: string;
  currency: string;
  is_dr: boolean;
}

interface TransactionFormProps {
  value?: {
    id?: string;
    date: string;
    description: string;
    postings: Posting[];
    tags?: { key: string; value: string }[];
  };
  onChange?: (data: {
    date: string;
    description: string;
    postings: Posting[];
    tags?: { key: string; value: string }[];
  }) => void;
  accounts: Array<{ id: string; name: string; type: string }>;
  currencies: string[];
  isEdit?: boolean;
}

export function TransactionForm(props: TransactionFormProps) {
  const [date, setDate] = createSignal(props.value?.date || todayLocalDate());
  const [description, setDescription] = createSignal(props.value?.description || '');
  const [tags, setTags] = createSignal<{ key: string; value: string }[]>(props.value?.tags || []);
  const [postings, setPostings] = createSignal<Posting[]>(
    props.value?.postings.length
      ? props.value.postings.map((p) => ({ ...p }))
      : [
          { account_id: '', amount: '', currency: props.currencies[0] || 'USD', is_dr: true },
          { account_id: '', amount: '', currency: props.currencies[0] || 'USD', is_dr: false },
        ],
  );

  const accountNames = createMemo(() =>
    props.accounts.map((a) => a.name),
  );

  const postingsList = () => postings();

  const postingsTotal = createMemo(() => {
    let sum = new Big(0);
    for (const p of postings()) {
      if (p.amount) {
        const amt = new Big(p.amount);
        sum = p.is_dr ? sum.plus(amt) : sum.minus(amt);
      }
    }
    return sum.abs().toString();
  });

  const balanced = createMemo(() => {
    return isZero(postingsTotal());
  });

  const handleAmountChange = (index: number, amount: string) => {
    const newPostings = [...postings()];
    newPostings[index].amount = amount;
    setPostings(newPostings);
  };

  const handleAccountChange = (index: number, accountId: string) => {
    const newPostings = [...postings()];
    newPostings[index].account_id = accountId;
    setPostings(newPostings);
  };

  const handleCurrencyChange = (index: number, currency: string) => {
    const newPostings = [...postings()];
    newPostings[index].currency = currency;
    setPostings(newPostings);
  };

  const toggleDR_CR = (index: number) => {
    const newPostings = [...postings()];
    newPostings[index].is_dr = !newPostings[index].is_dr;
    setPostings(newPostings);
  };

  const addPosting = () => {
    const newPosting: Posting = {
      account_id: '',
      amount: '',
      currency: props.currencies[0] || 'USD',
      is_dr: true,
    };
    setPostings([...postings(), newPosting]);
  };

  const removePosting = (index: number) => {
    if (postings().length <= 2) return;
    const newPostings = postings().filter((_, i) => i !== index);
    setPostings(newPostings);
  };

  const handleTagChange = (newTags: { key: string; value: string }[]) => {
    setTags(newTags);
  };

  const handleDataChange = () => {
    props.onChange?.({
      date: date(),
      description: description(),
      postings: postings(),
      tags: tags(),
    });
  };

  const rootStyle = {
    display: 'flex',
    flexDirection: 'column',
    gap: '24px',
  };

  const headerStyle = {
    display: 'grid',
    gridTemplateColumns: 'repeat(3, 1fr)',
    gap: '16px',
  };

  const fieldStyle = {
    display: 'flex',
    flexDirection: 'column',
    gap: '4px',
  };

  const postingsSectionStyle = {
    '.postingsHeader': {
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      marginBottom: '16px',
    },
    '.postingsTitle': {
      fontSize: '16px',
      fontWeight: '600',
      color: 'var(--color-text-primary)',
    },
    '.balanceInfo': {
      fontSize: '13px',
      fontWeight: '500',
    },
    '.balanceOK': {
      color: 'var(--color-success)',
    },
    '.balanceError': {
      color: 'var(--color-error)',
    },
  };

  const postingsListStyle = {
    display: 'flex',
    flexDirection: 'column',
    gap: '8px',
  };

  const postingRowStyle = {
    display: 'grid',
    gridTemplateColumns: '2fr 1fr 1fr 80px 40px',
    gap: '12px',
    alignItems: 'center',
  };

  const removeButtonStyle = {
    width: '32px',
    height: '32px',
    padding: '0',
    backgroundColor: 'var(--color-bg-elevated)',
    border: '1px solid var(--color-border)',
    borderRadius: 'var(--radius-md)',
    color: 'var(--color-text-secondary)',
    fontSize: '14px',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    transition: 'all 0.2s ease',
  };

  const drButtonStyle = {
    padding: '8px 12px',
    fontSize: '12px',
    fontWeight: '600',
    fontFamily: 'var(--font-family)',
    backgroundColor: 'var(--color-bg-elevated)',
    border: '1px solid var(--color-border)',
    borderRadius: 'var(--radius-md)',
    color: 'var(--color-text-secondary)',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    transition: 'all 0.2s ease',
    minWidth: '80px',
  };

  const drActiveStyle = {
    backgroundColor: 'var(--color-neon-green)',
    color: '#000',
    border: 'var(--color-neon-green)',
  };

  const drInactiveStyle = {
    backgroundColor: 'var(--color-neon-red)',
    color: '#000',
    border: 'var(--color-neon-red)',
  };

  const addPostingStyle = {
    marginTop: '16px',
  };

  const addPostingButtonStyle = {
    width: '100%',
    padding: '12px',
    fontSize: '14px',
    fontWeight: '500',
    fontFamily: 'var(--font-family)',
    backgroundColor: 'var(--color-bg-elevated)',
    color: 'var(--color-text-primary)',
    border: '1px dashed var(--color-border)',
    borderRadius: 'var(--radius-md)',
    cursor: 'pointer',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    gap: '8px',
  };

  createEffect(() => {
    handleDataChange();
  });

  return (
    <div style={rootStyle}>
      <div style={headerStyle}>
        <div style={fieldStyle}>
          <DatePicker value={date()} onChange={setDate} label="Date" />
        </div>
        <div style={fieldStyle}>
          <TextInput value={description()} onChange={setDescription} label="Description" placeholder="Enter description" />
        </div>
        <div style={fieldStyle}>
          <TagInput value={tags()} onChange={handleTagChange} label="Tags" existingTags={[]} />
        </div>
      </div>

      <div style={postingsSectionStyle}>
        <div style={postingsSectionStyle['.postingsHeader']}>
          <h3 style={postingsSectionStyle['.postingsTitle']}>Postings</h3>
          <span style={postingsSectionStyle['.balanceInfo']}>
            {balanced() ? (
              <span style={postingsSectionStyle['.balanceOK']}>Balanced</span>
            ) : (
              <span style={postingsSectionStyle['.balanceError']}>
                Unbalanced: {formatCurrency(postingsTotal(), 'USD')}
              </span>
            )}
          </span>
        </div>

        <div style={postingsListStyle}>
          {postingsList().map((posting, index) => (
            <div style={postingRowStyle}>
              <div style={fieldStyle}>
                <AutocompleteInput
                  value={posting.account_id}
                  options={accountNames()}
                  onChange={(val) => {
                    const account = props.accounts.find((a) => a.name === val);
                    if (account) handleAccountChange(index, account.id);
                  }}
                  placeholder="Select account..."
                  label="Account"
                  onSelect={(val) => {
                    const account = props.accounts.find((a) => a.name === val);
                    if (account) handleAccountChange(index, account.id);
                  }}
                />
              </div>

              <div style={fieldStyle}>
                <TextInput
                  value={posting.amount}
                  onChange={(val) => handleAmountChange(index, val)}
                  placeholder="0.00"
                  label="Amount"
                />
              </div>

              <div style={fieldStyle}>
                <SelectInput
                  value={posting.currency}
                  options={props.currencies.map((c) => ({ value: c, label: c }))}
                  onChange={(val) => handleCurrencyChange(index, val)}
                  label="Currency"
                />
              </div>

              <div style={fieldStyle}>
                <button
                  style={{
                    ...drButtonStyle,
                    ...(posting.is_dr ? drActiveStyle : drInactiveStyle),
                  }}
                  onClick={() => toggleDR_CR(index)}
                  type="button"
                >
                  {posting.is_dr ? 'DR' : 'CR'}
                </button>
              </div>

              <div style={fieldStyle}>
                {index > 1 && (
                  <button
                    style={removeButtonStyle}
                    onClick={() => removePosting(index)}
                    type="button"
                  >
                    ✕
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>

        <div style={addPostingStyle}>
          <button style={addPostingButtonStyle} onClick={addPosting} type="button">
            + Add Posting
          </button>
        </div>
      </div>
    </div>
  );
}
