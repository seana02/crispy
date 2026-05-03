import { useParams } from "@solidjs/router";
import { createMemo, createResource, ErrorBoundary, Suspense } from "solid-js";
import AccountForm from "src/components/AccountForm";
import { getAccountById, submitEditAccount } from "src/stores/accountStore";
import { domain } from "wailsjs/go/models";

export default function EditAccount() {
    const params = useParams();

    const id = +params.id!;

    const [acct] = createResource(() => id, () => getAccountById(id));
    
    const data = createMemo(() => {
        const item = acct();
        if (!item) return null;
        return {
            parentID: item.parentID,
            name: item.name,
            type: item.type,
            currency: item.currency,
            description: item.description,
            active: item.active,
            dateCreated: item.dateCreated,
            dateUpdated: item.dateUpdated,
        };
    });

    return (
        <ErrorBoundary fallback={err => <p>Loading error: {err.message}</p>}>
            <Suspense fallback={<p>Loading</p>}>
                <AccountForm
                    id={id}
                    name={data()!.name}
                    type={data()!.type}
                    currency={data()!.currency}
                    description={data()!.description}
                    active={data()!.active}
                    dateCreated={data()!.dateCreated}
                    dateUpdated={data()!.dateUpdated}
                    submit={(name: string, description: string, type: domain.Type, currency: string, active: boolean) => {
                        return submitEditAccount(id, data()!.parentID, name, description, type, currency, active);
                    }}
                />
            </Suspense>
        </ErrorBoundary>
    );
}
