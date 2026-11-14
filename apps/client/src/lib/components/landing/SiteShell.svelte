<script lang="ts">
	import type { Snippet } from 'svelte';
	import { LandingFooter, LandingHeader, SignInWithFacile } from '@facile/muse';
	import { LOGIN_PATH } from '$lib/api/prisme';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';

	let { children }: { children: Snippet } = $props();

	const links = $derived([
		{ href: '/#features', label: t('landing.nav.features') },
		{ href: '/downloads', label: t('landing.nav.apps') }
	]);

	const columns = $derived([
		{
			title: t('landing.footer.product'),
			links: [
				{ href: '/#features', label: t('landing.nav.features') },
				{ href: '/downloads', label: t('downloads.menu') }
			]
		}
	]);
</script>

<div class="h-dvh overflow-y-auto bg-fc-page text-fc-fg">
	<LandingHeader appName="Prisme" icon={glyphs.logo} {links} skipHref="#main">
		{#snippet actions()}
			<SignInWithFacile href={LOGIN_PATH} size="sm" />
		{/snippet}
	</LandingHeader>

	<main id="main">
		{@render children()}
	</main>

	<LandingFooter appName="Prisme" icon={glyphs.logo} tagline={t('landing.tagline')} {columns} />
</div>
