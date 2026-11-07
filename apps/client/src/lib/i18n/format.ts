import { locale } from './index.svelte';

/** Formats a date or instant with `Intl.DateTimeFormat` in the current locale. */
export function formatDate(value: Date | string | number, options: Intl.DateTimeFormatOptions = {}): string {
	return new Intl.DateTimeFormat(locale(), options).format(new Date(value));
}

/** Formats the clock time of an instant in the current locale. */
export function formatTime(value: Date | string | number): string {
	return formatDate(value, { hour: '2-digit', minute: '2-digit' });
}

/** Formats a byte count as kilobytes or megabytes in the current locale. */
export function formatBytes(bytes: number): string {
	if (bytes <= 0) return '';
	const megabytes = bytes > 1_000_000;
	const value = megabytes ? bytes / 1_000_000 : Math.max(1, Math.round(bytes / 1000));
	const options = { style: 'unit', unit: megabytes ? 'megabyte' : 'kilobyte', maximumFractionDigits: megabytes ? 1 : 0 } as const;
	return new Intl.NumberFormat(locale(), options).format(value);
}

/** Joins names the way the current locale lists them: "Ann, Bob and Cy". */
export function formatList(names: string[]): string {
	return new Intl.ListFormat(locale(), { style: 'long', type: 'conjunction' }).format(names);
}
