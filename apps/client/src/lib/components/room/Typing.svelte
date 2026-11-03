<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';

	let { room }: { room: Room } = $props();

	const text = $derived.by(() => {
		void app.version;
		const me = app.client?.getUserId();
		const names = room.getMembers().filter((m) => m.typing && m.userId !== me).map((m) => m.name);
		if (names.length === 0) return '';
		if (names.length === 1) return `${names[0]} is typing…`;
		if (names.length === 2) return `${names[0]} and ${names[1]} are typing…`;
		return 'Several people are typing…';
	});
</script>

<div class="h-5 px-4 text-fc-xs text-fc-fg-muted" aria-live="polite">{text}</div>
