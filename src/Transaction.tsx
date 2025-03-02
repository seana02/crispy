import './styles/Transaction.css';
import './styles/App.css';
import { useEffect, useState } from 'react';
import { invoke } from "@tauri-apps/api/core";
import Collapsible from 'react-collapsible';
import { TransactionData, PostingData } from './types';

interface TransactionProps {
    updateTab: ([newTab, data]: [string, any]) => void;
}

interface Received {
    description: string,
    id: number,
    postings: PostingData[],
    transaction_date: number[],
}

export default function Transaction(props: TransactionProps) {

    let [limit, setLimit] = useState(10);
    let [offset, setOffset] = useState(0);
    let [transactionList, setTransactionList] = useState<TransactionData[]>([]);
    let [refresher, refresh] = useState(1);

    useEffect(() => {
        invoke('get_latest_transactions', { limit, offset }).then(rec => {
            let l = rec as Received[];
            let t_list: TransactionData[] = [];
            l.forEach(p => t_list.push({
                id: p.id,
                date: new Date(p.transaction_date[0], p.transaction_date[1] - 1, p.transaction_date[2]),
                postings: p.postings,
                desc: p.description
            }));
            setTransactionList(t_list);
        }, e => {
            console.log(e);
        });
    }, [refresher]);

    return (
        <div id="transaction">
            <h1>Transactions</h1>
            <div id="transaction-button-row">
                <button onClick={() => props.updateTab(['TransactionEdit', null])}>Add Transaction</button>
            </div>
            <div id="transaction-list">
                <div className="transaction-list-header">
                    <div className="transaction-list-cell transaction-list-id">ID</div>
                    <div style={{ textAlign: "center" }} className="transaction-list-cell transaction-list-date">Date</div>
                    <div className="transaction-list-cell transaction-list-desc">Description</div>
                </div>
                {transactionList.map((t,i) => <TransactionRow index={i} transaction={t} />)}
            </div>
        </div>
    );

    function TransactionRow(rowProps: TransactionRowProps) {
        let [deleteConfirmation, setDeleteConfirmation] = useState(false);
        let deleteButton = deleteConfirmation ? 
            (<button className="transaction-list-cell transaction-list-button" onClick={e => {
                e.stopPropagation();
                invoke('delete_transaction', { id: rowProps.transaction.id }).then(_ => {
                    console.log('Delete success', rowProps.transaction.id);
                    refresh(c => c+1);
                }, r => {
                    console.log('Delete error:', r);
                })
            }}>Confirm</button>) :
            (<button className="transaction-list-cell transaction-list-button" onClick={e => {
                e.stopPropagation();
                setDeleteConfirmation(true);
            }}>Delete</button>);
        let header = (
            <div className="transaction-list-row">
                <div className="transaction-list-cell transaction-list-id">{rowProps.transaction.id}.</div>
                <div className="transaction-list-cell transaction-list-date">{rowProps.transaction.date.toISOString().split("T")[0]}</div>
                <div className="transaction-list-cell transaction-list-desc">{rowProps.transaction.desc}</div>
                <button className="transaction-list-cell transaction-list-button" onClick={e => {
                    e.stopPropagation();
                    props.updateTab(["TransactionEdit", rowProps.transaction])
                }}>Edit</button>
                {deleteButton}
            </div>);
        
        return (
            <div className="transaction-list-item">
                <Collapsible trigger={header} transitionTime={300}>
                    {rowProps.transaction.postings.map(p => <PostingRow id={p.id} account={p.account} value={p.value} currency={p.currency} comment={p.comment} />)}
                </Collapsible>
                {/* Togglable postings */}
            </div>
        );
    }

}

interface TransactionRowProps {
    index: number;
    transaction: TransactionData;
}

function PostingRow(data: PostingData) {
    return (
        <div className="posting-list-row">
            <div className="posting-list-cell posting-list-account">{data.account}</div>
            <div className="posting-list-cell posting-list-value">{data.value}</div>
            <div className="posting-list-cell posting-list-currency">{data.currency}</div>
            <div className="posting-list-cell posting-list-comment">{data.comment}</div>
        </div>
    );
}

