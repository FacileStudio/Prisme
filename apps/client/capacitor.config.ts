import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
	appId: 'studio.facile.prisme',
	appName: 'Prisme',
	webDir: 'build',
	android: { allowMixedContent: false },
	server: { androidScheme: 'https' },
	plugins: {
		SystemBars: { insetsHandling: 'native', style: 'DEFAULT' },
		PushNotifications: { presentationOptions: [] }
	}
};

export default config;
