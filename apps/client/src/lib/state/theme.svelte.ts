export type ThemePreference = 'light' | 'dark' | 'system';

const STORAGE_KEY = 'prisme.theme';

function isPreference(value: string | null): value is ThemePreference {
	return value === 'light' || value === 'dark' || value === 'system';
}

function readStored(): string | null {
	try {
		return localStorage.getItem(STORAGE_KEY);
	} catch {
		return null;
	}
}

class ThemeState {
	preference = $state<ThemePreference>('system');
	resolved = $state<'light' | 'dark'>('light');

	/** Reads the stored preference and follows the OS setting while on `system`. */
	restore() {
		const stored = readStored();
		this.preference = isPreference(stored) ? stored : 'system';
		this.apply();
		window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
			if (this.preference === 'system') this.apply();
		});
	}

	/** Stores and applies a preference. */
	set(next: ThemePreference) {
		this.preference = next;
		try {
			localStorage.setItem(STORAGE_KEY, next);
		} catch {
			return this.apply();
		}
		this.apply();
	}

	/**
	 * Writes `dark` or `light` on the root, never neither: muse paints dark from the
	 * media query on `:root:not(.light)`, so only an explicit `light` escapes it.
	 */
	private apply() {
		const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		const dark = this.preference === 'dark' || (this.preference === 'system' && prefersDark);
		this.resolved = dark ? 'dark' : 'light';
		document.documentElement.classList.toggle('dark', dark);
		document.documentElement.classList.toggle('light', !dark);
	}
}

export const theme = new ThemeState();
