import { CreateAccount } from "../../wailsjs/go/main/App";
import { domain } from "../../wailsjs/go/models";
import { createResource } from "solid-js";

// const [acctList, { refetch, mutate }] = createResource(async () => {
//     return await GetAccountList();
// })

const submitNewAccount = async (
    name: string,
    description: string,
    type: domain.Type,
    currency: string,
    active: boolean,
) => {
    CreateAccount(domain.AccountDTO.createFrom({ id: -1, parentID: -1, name, type, currency, description, active, dateCreated: new Date(), dateUpdated: new Date() }))
}

export { submitNewAccount };

