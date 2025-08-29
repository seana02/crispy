export interface TransactionData {
    id: number,
    transaction_date: string,
    description: string,
    postings: PostingData[],
};

export interface PostingData {
    id: number,
    account: string,
    value: string,
    currency: string,
    comment: string,
};

export interface AccountTree {
    label: string,
    total_currency: Wallet[],
    sub_accounts: AccountTree[] | undefined,
}[]

export interface Wallet {
    ccy: string,
    total_value: string,
}

