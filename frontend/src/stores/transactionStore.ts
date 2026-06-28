import { createStore, reconcile } from "solid-js/store";
import { CreateTransaction, UpdateTransaction, GetTransactionList, DeleteTransaction, GetTransactionByID, GetTransactionsByTag, RemoveTagFromTransaction, GetTransactionsByAccount } from "../../wailsjs/go/main/App";
import { domain } from "wailsjs/go/models";
import { createEffect, createResource } from "solid-js";
import { accountRefetch } from "src/stores/accountStore";
import { tagRefetch } from "src/stores/tagStore";

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
        accountRefetch();
        tagRefetch();
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

const fetchTransactionsByTag = async (tagId: number): Promise<domain.TransactionDTO[]> => {
    try {
        return await GetTransactionsByTag(tagId);
    } catch (err) {
        console.log(err);
        return [];
    }
}

const fetchTransactionsByAccount = async (acctId: number): Promise<domain.TransactionDTO[]> => {
    try {
        return await GetTransactionsByAccount(acctId);
    } catch (err) {
        console.log(err);
        return [];
    }
}

const removeTagFromTransaction = async (transactionId: number, tagId: number) => {
    try {
        await RemoveTagFromTransaction(transactionId, tagId);
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
    fetchTransactionsByTag,
    fetchTransactionsByAccount,
    removeTagFromTransaction,
    refetch as transactionRefetch,
};

