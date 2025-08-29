import { createRoute } from "@tanstack/react-router";
import { reportsRoute } from "./reports";
import { useQuery } from "@tanstack/react-query";
import { getBalanceSheet } from "@/api";
import { AccountTree, Wallet } from "@/types";

export const reportsBalanceSheetRoute = createRoute({
    getParentRoute: () => reportsRoute,
    path: '/bs',
    component: BalanceSheet,
});

function BalanceSheet() {
    const { isPending, isError, data, error } = useQuery<AccountTree>({
        queryKey: ['balance_sheet'],
        queryFn: () => getBalanceSheet(),
        staleTime: Infinity,
    });

    if (isPending) return <div>Pending...</div>
    if (isError) return <div>Error: {error.message}</div>

    console.log(convertDataToRows(data));

    return (
        <div>
            {data.sub_accounts?.map(a => (
                <div className="flex">
                    <div>{a.label}</div>
                    <div>{formatCurrency(a.total_currency)}</div>
                </div>
            ))}
            <div className="flex whitespace-pre">
                <div>{data.label}</div>
                <div>{formatCurrency(data.total_currency)}</div>
            </div>
        </div>
    );
}

function BSTable() {

}

function convertDataToRows(data: AccountTree) {
    let rows: { label: string, totalString: string }[] = [];
    addRows(data);
    return rows;

    function addRows(data: AccountTree) {
        if (!data) return;
        rows.push({
            label: data.label,
            totalString: formatCurrency(data.total_currency),
        });
        data.sub_accounts?.forEach(sub => addRows(sub));
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
