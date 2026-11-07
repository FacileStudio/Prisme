<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { EventTimeline, MatrixEvent, type Room, type Thread } from 'matrix-js-sdk';
	import { Spinner } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { DraftState } from '$lib/state/draft.svelte';
	import { timelineItems } from '$lib/matrix/timeline';
	import MessageRow from '$lib/components/room/MessageRow.svelte';
	import Composer from '$lib/components/room/Composer.svelte';

	let { room, rootId }: { room: Room; rootId: string } = $props();
	const draft = new DraftState();
	let root = $state.raw<MatrixEvent | null>(null);
	let thread = $state.raw<Thread | null>(null);
	let scroller: HTMLElement | undefined = $state();

	const items = $derived.by(() => {
		void app.version;
		return thread ? timelineItems(thread.events.filter((e) => e.getId() !== rootId)) : [];
	});

	async function open() {
		const client = app.client;
		if (!client) return;
		const event = room.findEventById(rootId) ?? new MatrixEvent(await client.fetchRoomEvent(room.roomId, rootId));
		const found = room.getThread(rootId) ?? room.createThread(rootId, event, [], true);
		root = event;
		thread = found;
		for (let page = 0; page < 3 && found.liveTimeline.getPaginationToken(EventTimeline.BACKWARDS); page++) {
			await client.paginateEventTimeline(found.liveTimeline, { backwards: true, limit: 50 });
		}
		app.version++;
	}

	onMount(() => void open());

	$effect(() => {
		void items.length;
		void tick().then(() => scroller && (scroller.scrollTop = scroller.scrollHeight));
	});
</script>

<div bind:this={scroller} class="flex-1 overflow-y-auto overscroll-contain py-2">
	{#if root}
		<MessageRow {room} event={root} grouped={false} {draft} inThread />
		<div class="flex items-center gap-3 px-4 py-2 text-fc-xs text-fc-fg-muted">
			<span>{t('thread.replies', { count: items.length })}</span><span class="h-px flex-1 bg-fc-surface"></span>
		</div>
		{#each items as item (item.event.getId() ?? item.event.getTxnId())}
			<MessageRow {room} event={item.event} grouped={item.grouped} {draft} inThread />
		{/each}
	{:else}
		<div class="flex justify-center py-6"><Spinner label={t('thread.loading')} /></div>
	{/if}
</div>
<Composer {room} {draft} threadId={rootId} placeholder={t('thread.placeholder')} />
