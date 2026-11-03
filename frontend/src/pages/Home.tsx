import { createMemo, createSignal, Show } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import { createResource } from 'solid-js';
import { getTransactions } from '../api/transaction';
import { formatCurrency } from '../state/decimal';
import { getBalanceSheet } from '../api/report';
import './Home.css';

export default function Home() {
    const navigate = useNavigate();

    const [balanceSheet] = createResource(() => getBalanceSheet());

    const [transactionsResource] = createResource(() => {
        return getTransactions({ page: 1, per_page: 5, search: '' });
    });

    const transactions = createMemo(() => transactionsResource()?.data || []);

    const netWorth = createMemo(() => {
        const bs = balanceSheet();
        if (!bs) return '0';
        return bs.net_worth ?? '0';
    });

    const recentTransactions = createMemo(() => transactions());

    const formatCurrencyValue = (value: string, currency: string = 'USD') => {
        const num = Number(value);
        if (isNaN(num)) return value;
        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: currency,
            minimumFractionDigits: 2,
            maximumFractionDigits: 2,
        }).format(num);
    };

    return (
        <div id="home">
            <h1 id="heading">Dashboard</h1>

            <div id="values" class="row">
                <div class="card">
                    <div class="label">Net Worth</div>
                    <div class="value net">{formatCurrencyValue(netWorth() || '0')}</div>
                    <div class="footnote">
                        As of today
                    </div>
                </div>

                <div class="card">
                    <div class="label">Assets</div>
                    <div class="value assets">
                        {formatCurrencyValue(balanceSheet()?.total_assets || '0')}
                    </div>
                    <div class="footnote">Total asset accounts</div>
                </div>

                <div class="card">
                    <div class="label">Liabilities</div>
                    <div class="value liabilities">
                        {formatCurrencyValue(balanceSheet()?.total_liab || '0')}
                    </div>
                    <div class="footnote">Total liability accounts</div>
                </div>
            </div>

            <div id="actions" class="row">
                <div class="card">
                    <div class="header">
                        <h2 class="title">Recent Transactions</h2>
                        <button id="viewall" class="color-primary" onClick={() => navigate('/transactions')}>View All</button>
                    </div>
                    <Show
                        when={recentTransactions().length > 0}
                        fallback={
                            <div style="padding: 16px; text-align: center; color: var(--color-text-secondary);">No recent transactions</div>
                        }
                    >
                        <div id="recent">
                            {recentTransactions().map((tx) => (
                                <div class="tx">
                                    <div style="flex: 1;">
                                        <div class="desc">{tx.description}</div>
                                        <div class="date">{tx.date}</div>
                                    </div>
                                    <div class="amount">{tx.postings[0] ? formatCurrencyValue(tx.postings[0].amount, tx.postings[0].currency) : ''}</div>
                                </div>
                            ))}
                        </div>
                    </Show>
                </div>

                <div class="card">
                    <h2 class="title">Quick Actions</h2>
                    <div id="action-grid">
                        <button onClick={() => navigate('/transactions/new')} class="quick-action-btn color-primary">+ New Transaction</button>
                        <button onClick={() => navigate('/accounts/new')} class="quick-action-btn color-surface">+ New Account</button>
                        <button onClick={() => navigate('/budgets/new')} class="quick-action-btn color-surface">+ New Budget</button>
                        <button onClick={() => navigate('/recurring/new')} class="quick-action-btn color-surface">+ Recurring</button>
                    </div>
                </div>
            </div>

            <div class="row">
                <div class="card">
                    <h2 class="title">Widget Grid Area</h2>
                    <div id="widget-area">
                        <p style="font-size: 16px; margin-bottom: 8px;">Widget Grid Placeholder</p>
                        <p style="font-size: 14px;">
                            Custom widgets can be added to create a personalized dashboard layout
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
}
