import './styles/Sidebar.css';

interface SidebarProps {
    activeTab: string;
    updateTab: (newTab: string) => void;
}

interface TabProps {
    text: string;
    status: TabStatus;
    onClick: () => void;
}

enum TabStatus {
    Active = "active",
    Inactive = "inactive",
    Hidden = "hidden"
}

export default function Sidebar(props: SidebarProps) {
    let tabs = [
        "Overview",
        "Transactions"
    ];
    return (
        <div className="sidebar">
            {tabs.map((t,i) => <Tab text={t} key={i} status={props.activeTab === t ? TabStatus.Active : TabStatus.Inactive} onClick={() => props.updateTab(t)} />)}
        </div>
    );
}

function Tab(props: TabProps) {
    return (
        <div className={`sidebar-tab ${props.status}`} onClick={props.onClick}>
            {props.text}
        </div>
    )
}
