import { invoke, InvokeArgs, InvokeOptions } from "@tauri-apps/api/core";
import { TransactionData } from "./types";

function generate_queryFn<T>(invoke_cmd: string) {
    return async function(params: InvokeArgs = {}) {
        console.log('executing', invoke_cmd, 'with params', params);
        const result = await invoke(invoke_cmd, params);
        if (typeof result === 'string') {
            throw new Error(result);
        }
        return result as T;
    }
}

export let getTransactions = generate_queryFn<TransactionData[]>('get_transactions');
export let getTransactionDetails = generate_queryFn<TransactionData>('get_transaction_details');

