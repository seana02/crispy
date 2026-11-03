import Big from 'big.js';

export function addDecimals(a: string, b: string): string {
    const result = new Big(a).plus(b);
    return result.toFixed();
}

export function isZero(a: string): boolean {
    if (!a || a.trim() === '') return true;
    return new Big(a).eq(0);
}

export function formatCurrency(amount: string, currency: string): string {
    const amountBig = new Big(amount);
    const amountNum = amountBig.toNumber();

    const formatter = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: currency,
        minimumFractionDigits: 0,
        maximumFractionDigits: 2,
    });

    return formatter.format(amountNum);
}

