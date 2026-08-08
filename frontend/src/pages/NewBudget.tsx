import BudgetForm from "src/components/BudgetForm";
import { submitNewBudget } from "src/stores/budgetStore";

export default function NewBudget() {
    return (
        <BudgetForm submit={submitNewBudget} />
    );
}
