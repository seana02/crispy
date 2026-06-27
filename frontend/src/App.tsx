import { Route, Router } from '@solidjs/router';

import PageLayout from "src/components/layouts/PageLayout";
import Dashboard from "src/pages/Dashboard";
import Transactions from "src/pages/Transactions";
import NewTransaction from "src/pages/NewTransaction";
import EditTransaction from "src/pages/EditTransaction";
import Budgets from "src/pages/Budgets";
import Accounts from "src/pages/Accounts";
import NewAccount from "src/pages/NewAccount";
import Settings from "src/pages/Settings";
import EditAccount from 'src/pages/EditAccount';
import Tags from 'src/pages/Tags';
import EditTag from 'src/pages/EditTag';

export default function App() {
    return (
        <Router root={PageLayout}>
            <Route path="/" component={Dashboard} />
            <Route path="/transactions">
                <Route path="/" component={Transactions} />
                <Route path="/add" component={NewTransaction} />
                <Route path="/:id" component={EditTransaction} />
            </Route>
            <Route path="/budgets" component={Budgets} />
            <Route path="/accounts">
                <Route path="/" component={Accounts} />
                <Route path="/add" component={NewAccount} />
                <Route path="/:id" component={EditAccount} />
            </Route>
            <Route path="/tags">
                <Route path="/" component={Tags} />
                <Route path="/:id" component={EditTag} />
            </Route>
            <Route path="/settings" component={Settings} />
        </Router>
    );
};
