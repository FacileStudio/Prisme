<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Button } from '@facile/muse';
	import { backend, ssoLoginUrl } from '$lib/backend';

	let error = $state('');

	$effect(() => {
		error = page.url.searchParams.get('error') ?? '';
	});

	onMount(async () => {
		try {
			await backend.getSession();
			await goto('/');
		} catch {
			return;
		}
	});
</script>

<svelte:head>
	<title>Log in — Prisme</title>
</svelte:head>

<div class="flex min-h-dvh">
	<div class="hidden flex-col bg-black px-12 py-10 lg:flex lg:w-1/2">
		<a href="/" class="mb-auto flex items-center gap-3 text-white" aria-label="Prisme home">
			<svg width="24" height="24" viewBox="0 0 24 24" fill="none" class="stroke-white stroke-2">
				<polygon points="12 2 2 22 22 22" />
			</svg>
			<span class="text-fc-xl font-semibold tracking-tight">Prisme</span>
		</a>

		<div class="mb-auto">
			<h2 class="text-fc-3xl leading-tight font-semibold tracking-tight text-white">
				Your team.<br />Your chat.
			</h2>
			<p class="mt-4 max-w-xs text-fc-sm leading-relaxed text-white/50">
				Slack-shaped team chat built on Matrix with AI assistants as first-class teammates.
			</p>
		</div>

		<p class="text-fc-xs text-white/30">
			&copy; {new Date().getFullYear()} Prisme by Facile.
		</p>
	</div>

	<div class="flex w-full flex-col items-center justify-center bg-fc-page px-8 py-12 lg:w-1/2">
		<div class="w-full max-w-sm">
			<div class="mb-8">
				<h1 class="text-fc-2xl font-semibold tracking-tight text-fc-fg">
					Welcome to Prisme
				</h1>
				<p class="mt-1.5 text-fc-sm text-fc-fg-muted">
					Sign in with your Facile account to access your spaces and channels.
				</p>
			</div>

			{#if error}
				<div role="alert" class="mb-4 rounded-fc-md bg-fc-danger/10 px-3 py-2 text-fc-sm text-fc-danger">
					{error}
				</div>
			{/if}

			<Button href={ssoLoginUrl} data-sveltekit-reload size="lg" class="w-full">
				Sign in with Facile
			</Button>
		</div>
	</div>
</div>
