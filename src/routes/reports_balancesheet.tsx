import { createRoute } from "@tanstack/react-router";
import { reportsRoute } from "./reports";
import { useQuery } from "@tanstack/react-query";
import { getBalanceSheet } from "@/api";
import { AccountTree, Wallet } from "@/types";
import { createColumnHelper, createTable } from "@tanstack/table-core";
import { useState } from "react";

export const reportsBalanceSheetRoute = createRoute({
    getParentRoute: () => reportsRoute,
    path: '/bs',
    component: BalanceSheet,
});

function BalanceSheet() {
    const { isPending, isError, data, error } = useQuery<AccountTree>({
        queryKey: ['balance_sheet'],
        queryFn: async () => {
            let data = await getBalanceSheet();

            return data;
        },
        staleTime: 0,
    });

    const [folded, setFolded] = useState<Set<string>>(new Set());

    const toggleFold = (path: string) => {
        setFolded(prev => {
            const next = new Set(prev);
            if (next.has(path)) next.delete(path);
            else next.add(path);
            return next;
        });
    }

    if (isPending) return <div>Pending...</div>
    if (isError) return <div>Error: {error.message}</div>

    let converted = convertDataToRows(data, folded);

    return (
        <div className="p-8 overflow-y-scroll h-full">
            <table className="w-full border text-xl font-mono">
                <tbody className="divide-y">
                    {converted.map((d, i) => {
                        return Row(d, i, toggleFold);
                    })}
                </tbody>
            </table>
        </div>
    );
}

enum RowStatus {
    Open,
    Closed,
    NoFold,
    Hidden,
}

type AccountRow = {
    fullpath: string,
    label: string,
    totalString: string,
    indents: number,
    status: RowStatus,
}

function Row(data: AccountRow, key: number, toggleFold: (path: string) => void) {
    let trClass = "divide-x";
    if (data.indents == 0) {
        trClass += " bg-blue-900";
    } else if (key % 2 == 0) {
        trClass += " bg-black";
    } else {
        trClass += " bg-gray-900";
    }
    if (data.status == RowStatus.Hidden) {
        trClass += " hidden";
    }
    return (
        <tr className={trClass} key={key}>
            <td className="w-[30px] text-center select-none" onClick={() => (data.status == RowStatus.Open || data.status == RowStatus.Closed) ? toggleFold(data.fullpath) : {}}>
                {data.status == RowStatus.Open ? "˅" : data.status == RowStatus.Closed ? "˃" : ""}
            </td>
            <td className="whitespace-pre p-2 w-1/2" style={{ paddingLeft: `${data.indents * 2 + 1}em` }}>
                {data.label + (data.status == RowStatus.Closed ? "..." : "")}
            </td>
            <td className="text-right p-2">{data.totalString}</td>
        </tr>
    );
}

function convertDataToRows(data: AccountTree, folded: Set<string>) {
    let rows: AccountRow[] = [];
    data.sub_accounts?.forEach(r => addRows(r, '', 0, false));
    return rows;

    function addRows(data: AccountTree, prefix: string, indents: number, hidden: boolean) {
        if (!data) return;
        if (data.sub_accounts?.length == 1) {
            addRows({
                label: data.label + ":" + data.sub_accounts[0].label,
                total_currency: data.total_currency,
                sub_accounts: data.sub_accounts[0].sub_accounts
            }, prefix, indents, hidden);
        } else {
            rows.push({
                fullpath: prefix + data.label,
                label: data.label,
                totalString: formatCurrency(data.total_currency),
                indents,
                status: hidden ? RowStatus.Hidden :
                    data.sub_accounts?.length == 0 ? RowStatus.NoFold :
                        folded.has(prefix + data.label) ? RowStatus.Closed :
                            RowStatus.Open
            });
            data.sub_accounts?.forEach(sub => addRows(sub, prefix + data.label + ':', indents + 1, hidden || folded.has(prefix + data.label)));
        }
    }
}

function formatCurrency(w: Wallet[]) {
    let output = "";
    w.forEach((c, i) => {
        if (i > 0) output += ", ";
        switch (c.ccy) {
            case "USD": output += `$${c.total_value}`; break;
            case "JPY": output += `¥${c.total_value}`; break;
            default: output += `${c.total_value} ${c.ccy}`; break;
        }
    });
    return output;
}
