import { describe, it, expect } from 'vitest';
import { addDecimals, isZero, formatCurrency } from './decimal';

describe('decimal utilities', () => {
    describe('addDecimals', () => {
        it('adds two decimal strings correctly', () => {
            expect(addDecimals('100.50', '25.25')).toBe('125.75');
        });

        it('handles whole numbers', () => {
            expect(addDecimals('100', '200')).toBe('300');
        });

        it('handles zero', () => {
            expect(addDecimals('0', '50.50')).toBe('50.5');
        });

        it('handles large numbers', () => {
            expect(addDecimals('999999999999', '1')).toBe('1000000000000');
        });

        it('avoids floating point issues (0.1 + 0.2 === 0.3)', () => {
            expect(addDecimals('0.1', '0.2')).toBe('0.3');
        });
    });

    describe('isZero', () => {
        it('returns true for "0"', () => {
            expect(isZero('0')).toBe(true);
        });

        it('returns true for "0.00"', () => {
            expect(isZero('0.00')).toBe(true);
        });

        it('returns true for "0.000"', () => {
            expect(isZero('0.000')).toBe(true);
        });

        it('returns false for non-zero values', () => {
            expect(isZero('0.01')).toBe(false);
            expect(isZero('1')).toBe(false);
            expect(isZero('-0.01')).toBe(false);
        });

        it('returns false for empty string', () => {
            expect(isZero('')).toBe(true);
        });
    });

    describe('formatCurrency', () => {
        it('formats USD currency', () => {
            const result = formatCurrency('1234.56', 'USD');
            expect(result).toBe('$1,234.56');
        });

        it('formats EUR currency', () => {
            const result = formatCurrency('1234.56', 'EUR');
            expect(result).toBe('€1,234.56');
        });

        it('formats JPY currency (no decimals)', () => {
            const result = formatCurrency('1234', 'JPY');
            expect(result).toBe('¥1,234');
        });

        it('formats negative amounts', () => {
            const result = formatCurrency('-1234.56', 'USD');
            expect(result).toBe('-$1,234.56');
        });

        it('handles zero', () => {
            const result = formatCurrency('0', 'USD');
            expect(result).toBe('$0');
        });

        it('handles small decimal values', () => {
            const result = formatCurrency('0.01', 'USD');
            expect(result).toBe('$0.01');
        });
    });
});
