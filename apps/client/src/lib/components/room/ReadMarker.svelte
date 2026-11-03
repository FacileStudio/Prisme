<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';

	let { room }: { room: Room } = $props();
	let visible = $state(typeof document === 'undefined' || !document.hidden);

	const last = $derived.by(() => {
		void app.version;
		const events = room.getLiveTimeline().getEvents();
		return [...events].reverse().find((e) => !e.isSending() && e.getId()?.startsWith('$')) ?? null;
	});

	$effect(() => {
		const client = app.client;
		const event = last;
		if (!client || !event || !visible) return;
		const me = client.getUserId() ?? '';
		if (room.getEventReadUpTo(me) === event.getId()) return;
		void client.sendReadReceipt(event).catch(() => undefined);
	});
</script>

<svelte:document onvisibilitychange={() => (visible = !document.hidden)} />
