<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Alert, Button, Spinner, icons } from '@facile/muse';
	import Shell from '$lib/components/shell/Shell.svelte';
	import { app } from '$lib/state/app.svelte';

	let { children } = $props();

	onMount(() => {
		void app.boot();
		return () => app.stop();
	});

	$effect(() => {
		if (app.status === 'signed-out') void goto('/login', { replaceState: true });
	});
</script>

{#if app.status === 'ready' && app.client}
	<Shell>{@render children()}</Shell>
{:else if app.status === 'error'}
	<main class="flex min-h-dvh items-center justify-center bg-fc-page p-4 text-fc-fg">
		<div class="flex w-full max-w-md flex-col gap-4">
			<Alert tone="danger" title="Prisme could not start">{app.error}</Alert>
			<Button icon={icons.refresh} onclick={() => location.reload()}>Try again</Button>
		</div>
	</main>
{:else}
	<main class="flex min-h-dvh items-center justify-center bg-fc-page text-fc-fg">
		<Spinner label="Connecting to the chat server" />
	</main>
{/if}
