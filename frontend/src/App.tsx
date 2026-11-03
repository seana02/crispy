import { Router, Route } from '@solidjs/router';
import { Layout } from './components/layout/Layout';
import Home from './pages/Home';
import Accounts from './pages/Accounts';
import AccountDetail from './pages/AccountDetail';
import AccountFormPage from './pages/AccountFormPage';
import Transactions from './pages/Transactions';
import TransactionDetail from './pages/TransactionDetail';
import TransactionFormPage from './pages/TransactionFormPage';
import Budgets from './pages/Budgets';
import BudgetDetail from './pages/BudgetDetail';
import BudgetFormPage from './pages/BudgetFormPage';
import BalanceSheet from './pages/BalanceSheet';
import IncomeStatement from './pages/IncomeStatement';
import Tags from './pages/Tags';
import TagDetail from './pages/TagDetail';
import RecurringTransactions from './pages/RecurringTransactions';
import RecurringTransactionDetail from './pages/RecurringTransactionDetail';
import ExchangeRates from './pages/ExchangeRates';
import Settings from './pages/Settings';

const App = () => {
  return (
    <Router root={Layout}>
      <Route path="/" component={Home} />

      {/* Accounts */}
      <Route path="/accounts" component={Accounts} />
      <Route path="/accounts/new" component={AccountFormPage} />
      <Route path="/accounts/:id/edit" component={AccountFormPage} />
      <Route path="/accounts/:id" component={AccountDetail} />

      {/* Transactions */}
      <Route path="/transactions" component={Transactions} />
      <Route path="/transactions/new" component={TransactionFormPage} />
      <Route path="/transactions/:id/edit" component={TransactionFormPage} />
      <Route path="/transactions/:id" component={TransactionDetail} />

      {/* Budgets */}
      <Route path="/budgets" component={Budgets} />
      <Route path="/budgets/new" component={BudgetFormPage} />
      <Route path="/budgets/:id/edit" component={BudgetFormPage} />
      <Route path="/budgets/:id" component={BudgetDetail} />

      {/* Reports */}
      <Route path="/reports/balance-sheet" component={BalanceSheet} />
      <Route path="/reports/income-statement" component={IncomeStatement} />

      {/* Tags */}
      <Route path="/tags" component={Tags} />
      <Route path="/tags/:key/:value" component={TagDetail} />

      {/* Recurring */}
      <Route path="/recurring" component={RecurringTransactions} />
      <Route path="/recurring/new" component={RecurringTransactionDetail} />
      <Route path="/recurring/:id" component={RecurringTransactionDetail} />

      {/* Exchange Rates */}
      <Route path="/exchange-rates" component={ExchangeRates} />

      {/* Settings */}
      <Route path="/settings" component={Settings} />
    </Router>
  );
};

export default App;
