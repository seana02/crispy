import { createStore, reconcile } from "solid-js/store";
import { CreateAccount, DeleteAccount, GetAccountIDByName, GetAccountList, SearchAccount, UpdateAccount } from "wailsjs/go/main/App";
import { domain } from "wailsjs/go/models";
import { createEffect, createResource } from "solid-js";

const [acctList, setAcctList] = createStore<domain.AccountDTO[]>([]);

const [dataResource, { refetch, }] = createResource(async () => {
    return await GetAccountList();
})

createEffect(() => {
    if (dataResource()) setAcctList(reconcile(dataResource()!));
})

const submitNewAccount = async (
    name: string,
    description: string,
    type: domain.Type,
    currency: string,
    active: boolean,
) => {
    try {
        let newDTO = domain.AccountDTO.createFrom({
            id: -1,
            parentID: -1,
            name,
            type,
            currency,
            description,
            active,
            dateCreated: new Date(),
            dateUpdated: new Date()
        });
        await CreateAccount(newDTO);
        await refetch();
    } catch (err) {
        console.log(err);
    }
}

const submitEditAccount = async (
    id: number,
    parentID: number,
    name: string,
    description: string,
    type: domain.Type,
    currency: string,
    active: boolean,
) => {
    try {
        let newDTO = domain.AccountDTO.createFrom({
            id,
            parentID,
            name,
            type,
            currency,
            description,
            active,
            dateCreated: new Date(),
            dateUpdated: new Date(),
        });
        await UpdateAccount(newDTO);
        await refetch();
    } catch (err) {
        console.log(err);
    }
}

const getAccountById = (id: number) => acctList.find(i => i.id === id);
const getAccountIdByName = async (name: string) => GetAccountIDByName(name);

const deleteAccountById = async (id: number) => {
    try {
        await DeleteAccount(id);
        await refetch();
    } catch (err) {
        console.log(err);
    }
}

const searchAccount = async (searchStr: string) => SearchAccount(searchStr);

export {
    acctList,
    submitNewAccount,
    submitEditAccount,
    getAccountById,
    getAccountIdByName,
    deleteAccountById,
    searchAccount,
};

