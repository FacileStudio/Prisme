const RELEASES = 'https://github.com/FacileStudio/Prisme/releases';
const LATEST = `${RELEASES}/latest/download`;

export type Platform = 'android' | 'macos' | 'windows' | 'linux';

export type Build = { platform: Platform; url: string; extra?: { url: string } };

export const RELEASES_URL = RELEASES;

/** Every installer the release workflows attach, under names that carry no version. */
export const BUILDS: readonly Build[] = [
	{ platform: 'android', url: `${LATEST}/Prisme-android.apk` },
	{ platform: 'macos', url: `${LATEST}/Prisme-macos-universal.dmg` },
	{ platform: 'windows', url: `${LATEST}/Prisme-windows-x64-setup.exe` },
	{ platform: 'linux', url: `${LATEST}/Prisme-linux-x64.AppImage`, extra: { url: `${LATEST}/Prisme-linux-x64.deb` } }
];

/** Guesses the visitor's platform from the user agent, or null when it cannot tell. */
export function detectPlatform(userAgent: string): Platform | null {
	if (/android/i.test(userAgent)) return 'android';
	if (/iphone|ipad/i.test(userAgent)) return null;
	if (/mac os x|macintosh/i.test(userAgent)) return 'macos';
	if (/windows/i.test(userAgent)) return 'windows';
	if (/linux|x11/i.test(userAgent)) return 'linux';
	return null;
}

/** Puts the visitor's own platform first, keeping the rest in order. */
export function ordered(platform: Platform | null): Build[] {
	return [...BUILDS].sort((a, b) => Number(b.platform === platform) - Number(a.platform === platform));
}
