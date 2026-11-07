<script lang="ts">
	import { MatrixEvent, type Room } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import { stripReplyFallback } from '$lib/chat/markdown';

	let { room, eventId }: { room: Room; eventId: string } = $props();
	let fetched = $state.raw<MatrixEvent | null>(null);

	const local = $derived.by(() => {
		void app.version;
		return room.findEventById(eventId) ?? null;
	});
	const target = $derived(local ?? fetched);

	$effect(() => {
		if (local || !app.client) return;
		void app.client.fetchRoomEvent(room.roomId, eventId).then(
			(raw) => (fetched = new MatrixEvent(raw)),
			() => undefined
		);
	});

	const preview = $derived.by(() => {
		if (!target) return { name: '', text: t('quote.original') };
		const sender = target.getSender() ?? '';
		const body = target.isRedacted() ? t('preview.deleted') : stripReplyFallback(String(target.getContent().body ?? ''));
		return { name: room.getMember(sender)?.name ?? sender, text: body };
	});

	function jump() {
		document.querySelector(`[data-event-id="${CSS.escape(eventId)}"]`)?.scrollIntoView({ behavior: 'smooth', block: 'center' });
	}
</script>

<button type="button" onclick={jump} class="mb-0.5 flex max-w-full items-center gap-2 bg-fc-surface rounded-fc-sm px-2 py-0.5 text-left text-fc-xs text-fc-fg-muted hover:text-fc-fg">
	{#if preview.name}<span class="shrink-0 font-semibold">{preview.name}</span>{/if}
	<span class="truncate">{preview.text}</span>
</button>
