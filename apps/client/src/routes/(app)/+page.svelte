<script lang="ts">
	import { goto } from '$app/navigation';
	import { EmptyState } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { roomHref } from '$lib/state/actions';
	import { lastRoom } from '$lib/state/last-room';

	const target = $derived.by(() => {
		void app.version;
		const client = app.client;
		if (!client) return null;
		const joined = (id: string) => client.getRoom(id)?.getMyMembership() === 'join';
		const last = lastRoom();
		if (last && joined(last)) return last;
		const channels = app.spaces.flatMap((s) => s.channels);
		return (channels.find((c) => c.name === 'general' && joined(c.room_id)) ?? channels.find((c) => joined(c.room_id)))?.room_id ?? null;
	});

	$effect(() => {
		if (target) void goto(roomHref(target), { replaceState: true });
	});
</script>

<div class="flex flex-1 items-center justify-center p-6">
	<EmptyState
		icon={glyphs.chat}
		title="Welcome to Prisme"
		description="Pick a channel or a conversation in the sidebar, or create a space to start."
	/>
</div>
