import type { Catalogs, Message } from '@facile/muse';
import type { Row } from '../row';
import { calls } from './calls';
import { common } from './common';
import { dialogs } from './dialogs';
import { panels } from './panels';
import { room } from './room';
import { settings } from './settings';
import { shell } from './shell';

const rows = { ...common, ...shell, ...dialogs, ...room, ...panels, ...calls, ...settings };

export type AppKey = keyof typeof rows;
type AppCatalog = Record<AppKey, Message>;

function column(index: number): AppCatalog {
	const entries = Object.entries(rows as Record<string, Row>).map(([key, row]) => [key, row[index]]);
	return Object.fromEntries(entries) as AppCatalog;
}

/** The app's catalogs: English is the typed base, and every row carries all four locales. */
export const catalogs: Catalogs<AppCatalog> = {
	en: column(0),
	fr: column(1),
	es: column(2),
	de: column(3)
};
