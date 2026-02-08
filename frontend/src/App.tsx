import { Route, Router } from '@solidjs/router';

import PageLayout from "./components/layouts/PageLayout";
import Dashboard from "./pages/Dashboard";
import Transactions from "./pages/Transactions";
import NewTransaction from "./pages/NewTransaction";
import EditTransaction from "./pages/EditTransaction";
import Budgets from "./pages/Budgets";
import Accounts from "./pages/Accounts";
import NewAccount from "./pages/NewAccount";
import Settings from "./pages/Settings";

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
            </Route>
            <Route path="/settings" component={Settings} />
        </Router>
    );
};
