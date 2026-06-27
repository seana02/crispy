import { createStore, reconcile } from "solid-js/store";
import { CreateTransaction, UpdateTransaction, GetTransactionList, DeleteTransaction, GetTransactionByID } from "../../wailsjs/go/main/App";
import { domain } from "wailsjs/go/models";
import { createEffect, createResource } from "solid-js";

const [txList, setTxList] = createStore<domain.TransactionDTO[]>([]);

const [dataResource, { refetch } ] = createResource(async () => {
    return await GetTransactionList();
});

createEffect(() => {
    if (dataResource()) setTxList(reconcile(dataResource()!));
})

const submitNewTransaction = async (
    desc: string,
    date: Date,
    status: domain.Status,
    tags: string[],
    postings: { accountID: number, account: string, amount: string, currency: string }[],
) => {
    try {
        let postingsDTO = await Promise.all(postings.map(async p => {
            return domain.PostingDTO.createFrom({
                accountID: p.accountID,
                accountName: p.account,
                amount: p.amount,
                currency: p.currency,
                dateCreated: new Date(),
                dateUpdated: new Date(),
                id: -1,
                transactionID: -1
            });
        }));
        let newDTO = domain.TransactionDTO.createFrom({
            id: -1,
            referenceID: null,
            description: desc,
            date,
            status,
            tags,
            postings: postingsDTO,
            dateCreated: new Date(),
            dateUpdated: new Date()
        });
        await CreateTransaction(newDTO);
        await refetch();
    } catch (err) {
        console.log(err);
    }
}

const submitEditTransaction = async (
    id: number,
    desc: string,
    date: Date,
    status: domain.Status,
    tags: string[],
    postings: domain.PostingDTO[],
) => {
    try {
        let newDTO = domain.TransactionDTO.createFrom({
            id: id,
            referenceID: null,
            description: desc,
            date,
            status,
            tags,
            postings: postings,
            dateCreated: new Date(),
            dateUpdated: new Date()
        });
        await UpdateTransaction(newDTO);
        await refetch();
    } catch (err) {
        console.log(err);
    }
}

const getTransactionById = async (id: number) => {
    return await GetTransactionByID(id, true);
};

const deleteTransactionById = async (id: number) => {
    try {
        await DeleteTransaction(id);
        await refetch();
    } catch (err) {
        console.log(err);
    }
}

export {
    txList,
    submitNewTransaction,
    submitEditTransaction,
    getTransactionById,
    deleteTransactionById,
};

