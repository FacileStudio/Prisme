const INTRO_KEY = 'prisme.intro-seen';
const WELCOME_KEY = 'prisme.welcomed';

function read(key: string): boolean {
	try {
		return localStorage.getItem(key) === '1';
	} catch {
		return false;
	}
}

function mark(key: string) {
	try {
		localStorage.setItem(key, '1');
	} catch {
		return;
	}
}

/** True once this device has reached the sign-in slide of the app intro. */
export const introSeen = () => read(INTRO_KEY);

/** Remembers that the intro reached its sign-in slide, so the next start opens there. */
export const markIntroSeen = () => mark(INTRO_KEY);

/** True once this device has been through the welcome step after its first sign-in. */
export const welcomed = () => read(WELCOME_KEY);

/** Remembers that the welcome step is done on this device. */
export const markWelcomed = () => mark(WELCOME_KEY);
