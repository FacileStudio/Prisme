<script lang="ts">
	import type { MatrixEvent, Room } from 'matrix-js-sdk';
	import { goto } from '$app/navigation';
	import { Input, Spinner } from '@facile/muse';
	import { app } from '$lib/state/app.svelte';
	import { spaceOf } from '$lib/state/nav';
	import { roomHref } from '$lib/state/actions';
	import { stripReplyFallback } from '$lib/chat/markdown';

	let { room }: { room: Room } = $props();
	let term = $state('');
	let busy = $state(false);
	let error = $state('');
	let results = $state.raw<MatrixEvent[] | null>(null);

	const space = $derived(spaceOf(app.spaces, room.roomId));
	const scope = $derived(space ? `in ${space.name}` : 'in this conversation');

	async function search(event: SubmitEvent) {
		event.preventDefault();
		const client = app.client;
		if (!client || !term.trim()) return;
		const rooms = space ? space.channels.map((c) => c.room_id).filter((id) => client.getRoom(id)?.getMyMembership() === 'join') : [room.roomId];
		busy = true;
		error = '';
		try {
			const found = await client.searchRoomEvents({ term: term.trim(), filter: { rooms } });
			results = found.results.map((r) => r.context.getEvent());
		} catch (err) {
			error = err instanceof Error ? err.message : 'Search failed';
		} finally {
			busy = false;
		}
	}

	function label(event: MatrixEvent) {
		const target = app.client?.getRoom(event.getRoomId() ?? '');
		const sender = target?.getMember(event.getSender() ?? '')?.name ?? event.getSender();
		return { room: target?.name ?? '', sender, body: stripReplyFallback(String(event.getContent().body ?? '')) };
	}
</script>

<form class="border-b border-fc-border p-3" onsubmit={search} role="search">
	<Input bind:value={term} type="search" placeholder="Search {scope}" aria-label="Search {scope}" />
</form>
<div class="flex-1 overflow-y-auto overscroll-contain p-2">
	{#if busy}
		<div class="flex justify-center py-6"><Spinner label="Searching" /></div>
	{:else if error}
		<p class="p-3 text-fc-sm text-fc-danger" role="alert">{error}</p>
	{:else if results && results.length === 0}
		<p class="p-3 text-fc-sm text-fc-fg-muted">No messages match “{term}”.</p>
	{:else if results}
		<ul class="flex flex-col gap-1">
			{#each results as event (event.getId())}
				{@const info = label(event)}
				<li>
					<button type="button" class="flex w-full flex-col gap-0.5 rounded-fc-md p-2 text-left hover:bg-fc-surface" onclick={() => goto(roomHref(event.getRoomId() ?? ''))}>
						<span class="text-fc-xs text-fc-fg-muted"><strong class="font-semibold text-fc-fg">{info.sender}</strong> in #{info.room} · {new Date(event.getTs()).toLocaleString()}</span>
						<span class="line-clamp-3 text-fc-sm">{info.body}</span>
					</button>
				</li>
			{/each}
		</ul>
	{:else}
		<p class="p-3 text-fc-sm text-fc-fg-muted">Search the messages {scope}.</p>
	{/if}
</div>
