<script lang="ts">
	import '../app.css';
	import { Toaster, provideLocale } from '@facile/muse';
	import { bindLocale, localeSource } from '$lib/i18n/index.svelte';
	import { registerGlyphs } from '$lib/glyphs';
	import TitleBar from '$lib/components/native/TitleBar.svelte';
	import { startDesktop } from '$lib/native/desktop';
	import { isApp, isDesktop } from '$lib/native/platform';
	import { registerDeepLinks } from '$lib/native/signin';
	import { theme } from '$lib/state/theme.svelte';
	import { notify } from '$lib/state/notify.svelte';

	let { children } = $props();

	const i18n = provideLocale(localeSource);
	bindLocale(() => i18n.locale);

	const desktop = isDesktop();

	registerGlyphs();
	theme.restore();
	notify.restore();
	if (isApp()) registerDeepLinks();
	if (desktop) startDesktop();
</script>

{#if desktop}
	<TitleBar />
{/if}

{@render children()}

<Toaster />
