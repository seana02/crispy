import { createRoute } from "@tanstack/react-router";
import { reportsRoute } from "./reports";
import { useQuery } from "@tanstack/react-query";
import { getBalanceSheet } from "@/api";
import { AccountTree, Wallet } from "@/types";
import { createColumnHelper, createTable } from "@tanstack/table-core";

export const reportsBalanceSheetRoute = createRoute({
    getParentRoute: () => reportsRoute,
    path: '/bs',
    component: BalanceSheet,
});

function BalanceSheet() {
    const { isPending, isError, data, error } = useQuery<AccountTree>({
        queryKey: ['balance_sheet'],
        queryFn: () => getBalanceSheet(),
        staleTime: 0,
    });

    if (isPending) return <div>Pending...</div>
    if (isError) return <div>Error: {error.message}</div>

    let converted = convertDataToRows(data);

    return (
        <div className="p-8 overflow-y-scroll h-full">
            <table className="w-full border text-xl font-mono">
                <tbody className="divide-y">
                    {converted.map((d, i) => {
                        if (i == 0) return <></>;
                        return Row(d, i);
                    })}
                </tbody>
            </table>
        </div>
    );
}

type AccountRow = {
    fullpath: string,
    label: string,
    totalString: string,
    indents: number,
}

function Row(data: AccountRow, key: number) {
    let trClass = "divide-x";
    if (data.indents == 0) {
        trClass += " bg-blue-900";
    } else if (key % 2 == 0) {
        trClass += " bg-black";
    } else {
        trClass += " bg-gray-900";
    }
    return (
        <tr className={trClass} key={key}>
            <td className="whitespace-pre p-2 w-1/2" style={{ paddingLeft: `${data.indents*2+1}em` }}>{data.label}</td>
            <td className="text-right p-2">{data.totalString}</td>
        </tr>
    );
}

function convertDataToRows(data: AccountTree) {
    let rows: AccountRow[] = [];
    addRows(data, '', -1);
    return rows;

    function addRows(data: AccountTree, prefix: string, indents: number) {
        if (!data) return;
        rows.push({
            fullpath: prefix + (data.label === 'root' ? '' : data.label),
            label: data.label === 'root' ? 'Total' : data.label,
            totalString: formatCurrency(data.total_currency),
            indents,
        });
        data.sub_accounts?.forEach(sub => addRows(sub, prefix+(data.label === 'root' ? '' : data.label+':'), indents+1));
    }
}

function formatCurrency(w: Wallet[]) {
    let output = "";
    w.forEach((c,i) => {
        if (i > 0) output += ", ";
        switch(c.ccy) {
            case "USD": output += `$${c.total_value}`; break;
            case "JPY": output += `¥${c.total_value}`; break;
            default: output += `${c.total_value} ${c.ccy}`; break;
        }
    });
    return output;
}
