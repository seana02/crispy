import { createStore, reconcile } from "solid-js/store";
import { CreateTransaction, UpdateTransaction, GetTransactionList } from "../../wailsjs/go/main/App";
import { domain } from "../../wailsjs/go/models";
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
    refetch();
}

const submitEditTransaction = async (
    id: number,
    desc: string,
    date: Date,
    status: domain.Status,
    tags: string[],
    postings: domain.PostingDTO[],
) => {
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
    refetch();
}

const getTransactionById = (id: number) => txList.find(i => i.id === id);

export { txList, submitNewTransaction, submitEditTransaction, getTransactionById, refetch };

