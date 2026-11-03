import { createSignal, createMemo, JSX } from 'solid-js';
import { RouteSectionProps, useLocation } from '@solidjs/router';
import { Breadcrumbs } from './Breadcrumbs';
import { House, Wallet, Banknote, PiggyBank, ChartLine, Tag, Calendar, BanknoteCheck, BanknoteArrowUp, Settings } from "lucide-solid";
import './Layout.css';

const A = (props: { href: string; class?: string; children?: any }) => {
    return <a {...props as any} />;
};

const Sidebar = () => {
    const [isCollapsed, setIsCollapsed] = createSignal(false);
    const [expandedGroups, setExpandedGroups] = createSignal<Set<string>>(new Set(['reports']));

    const toggleCollapse = () => setIsCollapsed(!isCollapsed());
    const toggleGroup = (group: string) => {
        const expanded = new Set(expandedGroups());
        if (expanded.has(group)) {
            expanded.delete(group);
        } else {
            expanded.add(group);
        }
        setExpandedGroups(expanded);
    };

    const navItems: Array<{
        label: string;
        path?: string;
        icon?: JSX.Element;
        children?: { label: string; path: string }[];
    }> = [
            { label: 'Home', path: '/', icon: <House /> },
            { label: 'Accounts', path: '/accounts', icon: <Wallet /> },
            { label: 'Transactions', path: '/transactions', icon: <Banknote /> },
            { label: 'Budgets', path: '/budgets', icon: <PiggyBank /> },
            {
                label: 'Reports',
                icon: <ChartLine />,
                children: [
                    { label: 'Balance Sheet', path: '/reports/balance-sheet' },
                    { label: 'Income Statement', path: '/reports/income-statement' },
                ],
            },
            { label: 'Tags', path: '/tags', icon: <Tag /> },
            { label: 'Recurring Transactions', path: '/recurring', icon: <Calendar /> },
            { label: 'Exchange Rates', path: '/exchange-rates', icon: <BanknoteArrowUp /> },
            { label: 'Settings', path: '/settings', icon: <Settings /> },
        ];

    const renderNavItem = (item: typeof navItems[0], depth = 0) => {
        const padding = `${depth * 12 + 16}px`;

        if (item.children) {
            const isExpanded = expandedGroups().has(item.label);
            return (
                <div>
                    <div
                        onClick={() => toggleGroup(item.label)}
                        class="nav-item nav-group"
                        classList={{ 'nav-group-expanded': isExpanded }}
                        style={{ 'padding-left': padding }}
                    >
                        <span class="nav-icon">{item.icon}</span>
                        <span class="nav-text">{item.label}</span>
                        <span class="nav-arrow">{isExpanded ? '▼' : '◀'}</span>
                    </div>
                    <div
                        class="nav-group-children"
                        classList={{ 'nav-group-children-expanded': isExpanded }}
                    >
                        {item.children!.map(child => renderNavItem(child, depth + 1))}
                    </div>
                </div>
            );
        }

        return (
            <A
                class="nav-item"
                href={item.path!}
            >
                <span class="nav-icon">{item.icon}</span>
                <span class="nav-text">{item.label}</span>
            </A>
        );
    };

    return (
        <aside
            class="sidebar"
            classList={{ 'sidebar-collapsed': isCollapsed() }}
        >
            <div class="sidebar-header">
                <h1 class="brand">Crispy</h1>
                <button onClick={toggleCollapse} class="sidebar-toggle">
                    {isCollapsed() ? '◀' : '▶'}
                </button>
            </div>
            <nav class="nav-menu">{navItems.map(item => renderNavItem(item))}</nav>
        </aside>
    );
};

const TopBar = () => {
    const [isSearchExpanded, setIsSearchExpanded] = createSignal(false);
    const [searchQuery, setSearchQuery] = createSignal('');
    const location = useLocation();

    const breadcrumbs = createMemo(() => {
        const segments = location.pathname.split('/').filter(Boolean);
        const items: { label: string; path?: string; active?: boolean }[] = [{ label: 'Home', path: '/' }];

        for (let i = 0; i < segments.length; i++) {
            const segment = segments[i];
            const path = '/' + segments.slice(0, i + 1).join('/');
            const label = segment
                .split('-')
                .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
                .join(' ');

            items.push({
                label: i === segments.length - 1 ? label : `${label} #${segment}`,
                path: i === segments.length - 1 ? undefined : path,
                active: i === segments.length - 1,
            });
        }

        return items;
    });

    return (
        <header class="top-bar">
            <div class="top-bar-left">
                <div class="breadcrumbs">
                    <Breadcrumbs items={breadcrumbs()} />
                </div>
            </div>
            <div class="top-bar-center">
                <div
                    class="search-bar"
                    classList={{ 'search-bar-expanded': isSearchExpanded() }}
                >
                    <input
                        type="text"
                        placeholder="Search..."
                        value={searchQuery()}
                        onInput={(e) => setSearchQuery(e.currentTarget.value)}
                        onFocus={() => setIsSearchExpanded(true)}
                        onBlur={() => setIsSearchExpanded(false)}
                    />
                </div>
            </div>
            <div class="top-bar-right">
                <button class="btn btn-primary">+ New Transaction</button>
            </div>
        </header>
    );
};

const Layout = (props: RouteSectionProps) => {
    return (
        <div class="app-layout">
            <Sidebar />
            <div class="main-content">
                <TopBar />
                <main class="main-content-inner">{props.children}</main>
            </div>
        </div>
    );
};

export { Layout };
export default Layout;
