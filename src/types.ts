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

