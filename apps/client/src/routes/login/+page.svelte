<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Alert, AuthScreen, Icon } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { rememberLocale, t } from '$lib/i18n/index.svelte';
	import { LOGIN_PATH, prisme, type AuthConfig } from '$lib/api/prisme';

	let config = $state<AuthConfig | null>(null);
	let checking = $state(true);
	const error = $derived(page.url.searchParams.get('error') ?? '');

	onMount(async () => {
		const signedIn = await prisme.session().then(() => true, () => false);
		if (signedIn) {
			await goto('/', { replaceState: true });
			return;
		}
		config = await prisme.authConfig().catch(() => null);
		checking = false;
	});
</script>

<svelte:head>
	<title>{t('login.title')}</title>
</svelte:head>

<AuthScreen appName="Prisme" href={LOGIN_PATH} loading={checking} onLocaleChange={rememberLocale}>
	{#snippet logo()}
		<Icon icon={glyphs.logo} size={40} class="text-fc-fg" />
	{/snippet}
	{#snippet tagline()}
		{t('login.tagline')}
	{/snippet}
	{#snippet notice()}
		{#if error}
			<Alert tone="danger" title={t('login.failedTitle')}>{error}</Alert>
		{/if}
		{#if config && !config.oidc_enabled}
			<Alert tone="warning">{t('login.noSso')}</Alert>
		{/if}
	{/snippet}
</AuthScreen>
