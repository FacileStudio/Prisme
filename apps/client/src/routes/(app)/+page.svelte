<script lang="ts">
	import { goto } from '$app/navigation';
	import { app } from '$lib/state/app.svelte';
	import { roomHref } from '$lib/state/actions';
	import { lastRoom } from '$lib/state/last-room';
	import { BROWSE_PATH } from '$lib/state/nav';
	import ChatsHome from '$lib/components/shell/ChatsHome.svelte';

	const phone = window.matchMedia('(max-width: 767px)').matches;

	const target = $derived.by(() => {
		void app.version;
		const client = app.client;
		if (!client) return null;
		if (phone) return BROWSE_PATH;
		const joined = (id: string) => client.getRoom(id)?.getMyMembership() === 'join';
		const last = lastRoom();
		if (last && joined(last)) return roomHref(last);
		const channels = app.spaces.flatMap((s) => s.channels);
		const room = channels.find((c) => c.name === 'general' && joined(c.room_id)) ?? channels.find((c) => joined(c.room_id));
		return room ? roomHref(room.room_id) : null;
	});

	$effect(() => {
		if (target) void goto(target, { replaceState: true });
	});
</script>

<ChatsHome />
