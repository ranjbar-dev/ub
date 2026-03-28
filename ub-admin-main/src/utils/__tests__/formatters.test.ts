import {
	queryStringer,
	CurrencyFormater,
	Format,
	FormatDate,
	DatePrefixer,
	censor,
	safeFinancialAdd,
} from '../formatters';

// ─── queryStringer ────────────────────────────────────────────────────────────

describe('queryStringer', () => {
	it('returns empty string for empty object', () => {
		expect(queryStringer({})).toBe('');
	});

	it('returns ?key=value for single param', () => {
		expect(queryStringer({ page: '1' })).toBe('?page=1');
	});

	it('joins multiple params with &', () => {
		expect(queryStringer({ page: '1', limit: '10' })).toBe('?page=1&limit=10');
	});

	it('encodes special characters in keys and values', () => {
		expect(queryStringer({ 'foo bar': 'hello world' })).toBe('?foo%20bar=hello%20world');
	});

	it('encodes & in values', () => {
		expect(queryStringer({ q: 'a&b' })).toBe('?q=a%26b');
	});

	it('skips null values', () => {
		expect(queryStringer({ a: 'x', b: null })).toBe('?a=x');
	});

	it('skips undefined values', () => {
		expect(queryStringer({ a: 'x', b: undefined })).toBe('?a=x');
	});

	it('returns empty string when all values are null/undefined', () => {
		expect(queryStringer({ a: null, b: undefined })).toBe('');
	});

	it('coerces non-string values via String()', () => {
		expect(queryStringer({ n: 42 })).toBe('?n=42');
	});
});

// ─── safeFinancialAdd ─────────────────────────────────────────────────────────

describe('safeFinancialAdd', () => {
	it('avoids floating-point error: 0.1 + 0.2', () => {
		expect(safeFinancialAdd('0.1', '0.2')).toBe('0.30000000');
	});

	it('adds integers correctly', () => {
		expect(safeFinancialAdd('1', '2')).toBe('3.00000000');
	});

	it('handles zeros', () => {
		expect(safeFinancialAdd('0', '0')).toBe('0.00000000');
	});

	it('handles large crypto precision sum', () => {
		expect(safeFinancialAdd('999999.99999999', '0.00000001')).toBe('1000000.00000000');
	});

	it('supports custom decimals', () => {
		expect(safeFinancialAdd('0.1', '0.2', 2)).toBe('0.30');
	});

	it('adds two decimal strings', () => {
		expect(safeFinancialAdd('100.50', '200.25')).toBe('300.75000000');
	});

	it('handles negative first operand', () => {
		expect(safeFinancialAdd('-1', '2')).toBe('1.00000000');
	});

	it('handles negative second operand', () => {
		expect(safeFinancialAdd('5', '-3')).toBe('2.00000000');
	});

	it('accepts number inputs', () => {
		expect(safeFinancialAdd(0.1, 0.2)).toBe('0.30000000');
	});

	it('treats invalid strings as 0', () => {
		expect(safeFinancialAdd('abc', '1')).toBe('1.00000000');
	});

	it('uses 8 decimal places by default', () => {
		const result = safeFinancialAdd('1', '1');
		expect(result.split('.')[1]).toHaveLength(8);
	});
});

// ─── CurrencyFormater ─────────────────────────────────────────────────────────

describe('CurrencyFormater', () => {
	it("returns '0.00' for string '0'", () => {
		expect(CurrencyFormater('0')).toBe('0.00');
	});

	it("returns '-' for empty string", () => {
		expect(CurrencyFormater('')).toBe('-');
	});

	it('returns value as-is when it contains %', () => {
		expect(CurrencyFormater('12.5%')).toBe('12.5%');
	});

	it('returns value as-is when it contains 0.00000', () => {
		expect(CurrencyFormater('0.000001')).toBe('0.000001');
	});

	it('formats integer with 2 decimal places', () => {
		expect(CurrencyFormater('1234')).toBe('1,234.00');
	});

	it('formats value with 1 decimal place to 2', () => {
		expect(CurrencyFormater('1234.5')).toBe('1,234.50');
	});

	it('preserves more than 2 decimal places', () => {
		expect(CurrencyFormater('1234.5678')).toBe('1,234.5678');
	});

	it('formats large number with commas', () => {
		expect(CurrencyFormater('1000000')).toBe('1,000,000.00');
	});

	it('formats value with exactly 2 decimal places', () => {
		expect(CurrencyFormater('99.99')).toBe('99.99');
	});
});

// ─── Format ───────────────────────────────────────────────────────────────────

describe('Format', () => {
	it('formats positive number with German-style comma thousands separator', () => {
		// de-DE uses '.' for thousands then we swap '.' → ',' so 1000 → '1,000'
		expect(Format(1000)).toBe('1,000');
	});

	it('returns empty string for zero', () => {
		expect(Format(0)).toBe('');
	});

	it('returns empty string for empty string input (treated as 0)', () => {
		expect(Format('')).toBe('');
	});

	it('wraps negative numbers in parentheses', () => {
		expect(Format(-1000)).toBe('(1,000)');
	});

	it('formats large positive number', () => {
		expect(Format(1000000)).toBe('1,000,000');
	});

	it('formats negative large number in parentheses', () => {
		expect(Format(-1000000)).toBe('(1,000,000)');
	});
});

// ─── FormatDate ───────────────────────────────────────────────────────────────

describe('FormatDate', () => {
	it("converts 'YYYYMMDD' string to 'YYYY-MM-DD'", () => {
		expect(FormatDate('20231225')).toBe('2023-12-25');
	});

	it('converts numeric YYYYMMDD to YYYY-MM-DD', () => {
		expect(FormatDate(20230101)).toBe('2023-01-01');
	});

	it('returns empty string for empty string input', () => {
		expect(FormatDate('')).toBe('');
	});

	it('returns empty string for 0', () => {
		expect(FormatDate(0)).toBe('');
	});

	it('handles minimal 8-char date string', () => {
		expect(FormatDate('20000229')).toBe('2000-02-29');
	});
});

// ─── DatePrefixer ─────────────────────────────────────────────────────────────

describe('DatePrefixer', () => {
	it('pads single-digit number with leading zero', () => {
		expect(DatePrefixer(5)).toBe('05');
	});

	it('does not pad double-digit number', () => {
		expect(DatePrefixer(12)).toBe('12');
	});

	it('pads 0 with leading zero', () => {
		expect(DatePrefixer(0)).toBe('00');
	});

	it('does not pad 10', () => {
		expect(DatePrefixer(10)).toBe('10');
	});

	it('handles single digit 9', () => {
		expect(DatePrefixer(9)).toBe('09');
	});
});

// ─── censor ───────────────────────────────────────────────────────────────────

describe('censor', () => {
	it('replaces characters in range with asterisks', () => {
		expect(censor({ value: 'hello', from: 1, to: 3 })).toBe('h***o');
	});

	it('censors from start of string', () => {
		expect(censor({ value: 'hello', from: 0, to: 2 })).toBe('***lo');
	});

	it('censors to end of string', () => {
		expect(censor({ value: 'hello', from: 3, to: 4 })).toBe('hel**');
	});

	it('censors entire string', () => {
		expect(censor({ value: 'abc', from: 0, to: 2 })).toBe('***');
	});

	it('censors single character', () => {
		expect(censor({ value: 'hello', from: 2, to: 2 })).toBe('he*lo');
	});

	it('works with email-like strings', () => {
		expect(censor({ value: 'user@example.com', from: 1, to: 3 })).toBe('u***@example.com');
	});
});
