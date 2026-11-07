<script lang="ts">
	import '../app.css';
	import { Toaster, provideLocale } from '@facile/muse';
	import { bindLocale, localeSource } from '$lib/i18n/index.svelte';
	import { registerGlyphs } from '$lib/glyphs';
	import { theme } from '$lib/state/theme.svelte';
	import { notify } from '$lib/state/notify.svelte';

	let { children } = $props();

	const i18n = provideLocale(localeSource);
	bindLocale(() => i18n.locale);

	registerGlyphs();
	theme.restore();
	notify.restore();
</script>

{@render children()}

<Toaster />
