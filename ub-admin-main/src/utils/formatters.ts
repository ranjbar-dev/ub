/**
 * Builds a URL query string from an object of key-value pairs.
 * All keys and values are URL-encoded to prevent injection.
 *
 * @param params - Object with string key-value pairs
 * @returns Query string starting with '?' or empty string if no params
 */
export function queryStringer(params: Record<string, unknown>): string {
	const entries = Object.entries(params).filter(
		([, value]) => value !== undefined && value !== null,
	);
	if (entries.length === 0) {
		return '';
	}
	const qs = entries
		.map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
		.join('&');
	return `?${qs}`;
}
export const CurrencyFormater = (val: string) => {
	let value = val;
	if (typeof value == 'number') {
		value = value + ''
	}
	// @ts-expect-error — checking string against number 0 for edge case safety
	if (value === '0' || value === 0) {
		return '0.00'
	}
	if (!value) {
		return '-';
	}
	if (value.includes('%') || value.includes('0.00000')) {
		return value;
	}
	let trail = '';
	if (value.includes(' ')) {
		value = +value.split(' ')[0] + ' ' + value.split(' ')[1];
		trail = ' ' + value.split(' ')[1];
	} else {
		value = +value + '';
	}
	if (!value.split('.')[1]) {
		value = Number(value.split(' ')[0]).toFixed(2) + trail;
	} else if (value.split(' ')[0].split('.')[1].length < 3) {
		value = Number(value.split(' ')[0]).toFixed(2) + trail;
	}

	let separated = value.split('.');
	let comaSeparated = separated[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',');
	if (separated[1]) {
		return comaSeparated + '.' + separated[1];
	}
	return comaSeparated;
};

export function Format(value: number | string) {
	const v = value ? Number(value) : 0;
	if (v > 0) {
		return Intl.NumberFormat('de-DE').format(v).split('.').join(',');
	} else if (v === 0) {
		return '';
	} else {
		return (
			'(' +
			Intl.NumberFormat('de-DE')
				.format(-1 * v)
				.split('.')
				.join(',') +
			')'
		);
	}
}

export function FormatDate(value: number | string) {
	const v = value ? value + '' : '';
	if (v.length > 0) {
		return (
			v.substring(0, 4) + '-' + v.substring(4, 6) + '-' + v.substring(6, 8)
		);
	} else {
		return '';
	}
}
export const DatePrefixer = (number: number) => {
	return number < 10 ? '0' + number : '' + number;
};

export const vw = (percent: number, viewportWidth: number) => {
	const per = (percent / 100) * viewportWidth + '';
	return Number(per.split('.')[0]);
};
export const censor = (item: { value: string; from: number; to: number }) => {
	let transmitted = item.value.split('');
	for (var i = item.from; i <= item.to; i++) {
		transmitted[i] = '*';
	}
	return transmitted.join('');
};
export const CopyToClipboard = (text: string) => {
	let dummy = document.createElement('textarea');
	document.body.appendChild(dummy);
	dummy.value = text;
	dummy.select();
	document.execCommand('copy');
	document.body.removeChild(dummy);
};
export const Translator = (translateData: {
	intl: { formatMessage: (descriptor: { id: string; defaultMessage: string }) => string };
	containerPrefix: string;
	message: string;
}) => {
	return translateData.intl.formatMessage({
		id: translateData.containerPrefix + '.' + translateData.message,
		defaultMessage: 'ET.' + translateData.message,
	});
};
export const PairFormat = (pair: string) => {
	return pair.replace('-', '');
};
export const under = (key: string) => {
	return key
		.replace(/\.?([A-Z]+)/g, function (x, y) {
			return '_' + y.toLowerCase();
		})
		.replace(/^_/, '');
};

/**
 * Safely adds two financial string values and returns a formatted string.
 * Avoids floating-point arithmetic issues by using fixed-point arithmetic.
 *
 * @param a - First value as string (e.g., "0.001234")
 * @param b - Second value as string (e.g., "0.005678")
 * @param decimals - Number of decimal places (default: 8 for crypto)
 * @returns Formatted sum as string
 */
export function safeFinancialAdd(a: string | number, b: string | number, decimals: number = 8): string {
	const numA = parseFloat(String(a)) || 0;
	const numB = parseFloat(String(b)) || 0;
	const multiplier = Math.pow(10, decimals);
	const result = (Math.round(numA * multiplier) + Math.round(numB * multiplier)) / multiplier;
	return result.toFixed(decimals);
}
