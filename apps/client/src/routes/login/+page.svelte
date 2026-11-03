<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Alert, Button, Icon, Spinner } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
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
	<title>Sign in · Prisme</title>
</svelte:head>

<main class="flex min-h-dvh items-center justify-center bg-fc-page px-4 text-fc-fg">
	<div class="flex w-full max-w-sm flex-col gap-6">
		<div class="flex flex-col items-center gap-3 text-center">
			<span class="flex size-14 items-center justify-center rounded-fc-lg bg-fc-accent text-fc-accent-fg">
				<Icon icon={glyphs.logo} size={28} />
			</span>
			<h1 class="text-fc-2xl font-semibold">Prisme</h1>
			<p class="text-fc-sm text-fc-fg-muted">The Facile team chat. Sign in with your Facile account.</p>
		</div>
		{#if error}
			<Alert tone="danger" title="Sign-in failed">{error}</Alert>
		{/if}
		{#if checking}
			<div class="flex justify-center"><Spinner label="Checking your session" /></div>
		{:else}
			<Button href={LOGIN_PATH} size="lg" icon={glyphs.login} class="w-full" data-sveltekit-reload>
				Sign in with Facile
			</Button>
			{#if config && !config.oidc_enabled}
				<p class="text-center text-fc-xs text-fc-fg-muted">
					Single sign-on is not configured on this server. Ask an administrator.
				</p>
			{/if}
		{/if}
	</div>
</main>
