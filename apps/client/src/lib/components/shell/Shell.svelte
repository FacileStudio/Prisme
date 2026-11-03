<script lang="ts">
	import type { Snippet } from 'svelte';
	import { goto } from '$app/navigation';
	import { app } from '$lib/state/app.svelte';
	import { notify } from '$lib/state/notify.svelte';
	import { roomHref } from '$lib/state/actions';
	import { ui } from '$lib/state/ui.svelte';
	import { call } from '$lib/call/call.svelte';
	import Sidebar from './Sidebar.svelte';
	import Dialogs from '$lib/components/dialogs/Dialogs.svelte';
	import SidePanel from '$lib/components/panels/SidePanel.svelte';
	import CallPanel from '$lib/components/call/CallPanel.svelte';

	let { children }: { children: Snippet } = $props();

	const joinedSpaces = $derived.by(() => {
		void app.version;
		return (app.client?.getRooms() ?? [])
			.filter((r) => r.isSpaceRoom() && r.getMyMembership() === 'join')
			.map((r) => r.roomId)
			.sort()
			.join(',');
	});

	$effect(() => {
		void joinedSpaces;
		void app.refreshSpaces();
	});

	$effect(() => {
		const client = app.client;
		if (!client) return;
		void client.setPresence({ presence: 'online' }).catch(() => undefined);
		return notify.watch(client, (roomId) => void goto(roomHref(roomId)));
	});
</script>

<div class="flex h-dvh overflow-hidden bg-fc-page text-fc-fg">
	<div class="hidden md:flex">
		<Sidebar />
	</div>
	{#if ui.drawer}
		<div class="fixed inset-0 z-40 flex md:hidden">
			<button type="button" class="absolute inset-0 bg-fc-scrim" aria-label="Close the menu" onclick={() => (ui.drawer = false)}></button>
			<div class="relative flex h-full max-w-[90vw] shadow-lg">
				<Sidebar />
			</div>
		</div>
	{/if}
	<main class="flex min-w-0 flex-1 flex-col">
		{#if call.active}
			<CallPanel />
		{/if}
		{@render children()}
	</main>
	<SidePanel />
</div>

<Dialogs />
