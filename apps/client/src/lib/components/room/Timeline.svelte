<script lang="ts">
	import { tick, untrack } from 'svelte';
	import { EventTimeline, type Room } from 'matrix-js-sdk';
	import { Button, Spinner } from '@facile/muse';
	import { app } from '$lib/state/app.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import type { DraftState } from '$lib/state/draft.svelte';
	import { timelineItems } from '$lib/matrix/timeline';
	import MessageRow from './MessageRow.svelte';
	import SeenBy from './SeenBy.svelte';

	let { room, draft }: { room: Room; draft: DraftState } = $props();
	let scroller: HTMLElement | undefined = $state();
	let loading = $state(false);
	let stalled = $state(false);
	let atBottom = true;

	const items = $derived.by(() => {
		void app.version;
		return timelineItems(room.getLiveTimeline().getEvents());
	});
	const canLoadMore = $derived.by(() => {
		void app.version;
		return !!room.getLiveTimeline().getPaginationToken(EventTimeline.BACKWARDS);
	});

	async function loadOlder() {
		if (loading || stalled || !canLoadMore || !app.client || !scroller) return;
		loading = true;
		const height = scroller.scrollHeight;
		try {
			await app.client.paginateEventTimeline(room.getLiveTimeline(), { backwards: true, limit: 30 });
			await tick();
			scroller.scrollTop += scroller.scrollHeight - height;
		} catch {
			stalled = true;
		} finally {
			loading = false;
		}
	}

	function onScroll() {
		if (!scroller) return;
		atBottom = scroller.scrollHeight - scroller.scrollTop - scroller.clientHeight < 80;
		if (scroller.scrollTop < 200) void loadOlder();
	}

	$effect(() => {
		void items.length;
		if (!scroller || !atBottom) return;
		void tick().then(() => scroller && (scroller.scrollTop = scroller.scrollHeight));
	});

	$effect(() => {
		if (scroller && items.length < 15 && canLoadMore) untrack(() => void loadOlder());
	});
</script>

<div bind:this={scroller} onscroll={onScroll} class="relative flex-1 overflow-y-auto overscroll-contain py-2" role="log" aria-label={t('timeline.label')} aria-live="polite">
	{#if loading}
		<div class="flex justify-center py-2"><Spinner size="sm" label={t('timeline.loadingOlder')} /></div>
	{:else if stalled}
		<div class="flex justify-center py-2">
			<Button size="sm" variant="ghost" onclick={() => ((stalled = false), loadOlder())}>{t('timeline.loadOlder')}</Button>
		</div>
	{:else if !canLoadMore}
		<p class="px-4 py-6 text-center text-fc-xs text-fc-fg-muted">{t('timeline.start', { name: room.name })}</p>
	{/if}
	{#each items as item, index (item.event.getId() ?? item.event.getTxnId())}
		{#if item.day}
			<div class="flex items-center gap-3 px-4 py-2 text-fc-xs font-medium text-fc-fg-muted" role="separator">
				<span class="h-px flex-1 bg-fc-surface"></span>{item.day}<span class="h-px flex-1 bg-fc-surface"></span>
			</div>
		{/if}
		<MessageRow {room} event={item.event} grouped={item.grouped} {draft} />
		{#if index === items.length - 1}
			<SeenBy {room} event={item.event} />
		{/if}
	{/each}
</div>
