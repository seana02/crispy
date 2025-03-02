
export interface PostingData {
    id: number;
    account: string;
    value: string;
    currency: string;
    comment: string;
}

export interface TransactionData {
    id: number | null;
    date: Date;
    postings: PostingData[];
    desc: string;
}

