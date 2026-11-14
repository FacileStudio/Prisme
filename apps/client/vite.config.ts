import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	optimizeDeps: { exclude: ['@facile/muse', '@facile/mise-en-scene'] },
	server: {
		proxy: {
			'/api': {
				target: 'http://127.0.0.1:14090',
				changeOrigin: true
			}
		}
	},
	test: {
		include: ['src/**/*.test.ts'],
		environment: 'node'
	}
});
