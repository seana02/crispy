import AccountForm from "src/components/AccountForm";
import { submitNewAccount } from "src/stores/accountStore";

export default function NewAccount() {
    return (
        <AccountForm
            submit={submitNewAccount}
        />
    );
}
