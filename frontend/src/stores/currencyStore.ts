
export const Currency = {
    USD: "USD",
    JPY: "JPY",
    EUR: "EUR",
    CAD: "CAD",
} as const;

export type Currency = typeof Currency[keyof typeof Currency];
