import type { DownloadEntry } from '@facile/muse';

const RELEASES = 'https://github.com/FacileStudio/Prisme/releases';
const LATEST = `${RELEASES}/latest/download`;

export const RELEASES_URL = RELEASES;
export const CHECKSUMS_URL = `${LATEST}/desktop-checksums.txt`;

/** Every installer the release workflows attach, under names that carry no version. */
export const DOWNLOADS: DownloadEntry[] = [
	{ platform: 'macos', href: `${LATEST}/Prisme-macos-universal.dmg`, arch: 'universal' },
	{ platform: 'windows', href: `${LATEST}/Prisme-windows-x64-setup.exe`, arch: 'x64' },
	{ platform: 'linux', href: `${LATEST}/Prisme-linux-x64.AppImage`, label: 'AppImage', arch: 'x64' },
	{ platform: 'linux', href: `${LATEST}/Prisme-linux-x64.deb`, label: '.deb', arch: 'x64' },
	{ platform: 'android', href: `${LATEST}/Prisme-android.apk`, label: 'APK' }
];
