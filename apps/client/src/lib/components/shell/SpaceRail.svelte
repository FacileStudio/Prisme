<script lang="ts">
	import { IconButton, Icon, icons, toast } from '@facile/muse';
	import { goto } from '$app/navigation';
	import type { Room } from 'matrix-js-sdk';
	import type { Nav } from '$lib/state/nav';
	import { acceptSpace, roomHref } from '$lib/state/actions';
	import { ui } from '$lib/state/ui.svelte';
	import RailButton from './RailButton.svelte';

	let { nav }: { nav: Nav } = $props();

	async function accept(room: Room) {
		try {
			const space = await acceptSpace(room);
			if (!space) return;
			ui.spaceId = space.room_id;
			const general = space.channels.find((c) => c.name === 'general');
			if (general) await goto(roomHref(general.room_id));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not join the space');
		}
	}
</script>

<div class="flex h-full w-16 flex-col items-center gap-2 border-r border-fc-border bg-fc-page py-3" aria-label="Spaces">
	{#each nav.spaces as entry (entry.space.room_id)}
		<RailButton
			name={entry.name}
			active={ui.spaceId === entry.space.room_id}
			badge={entry.badge}
			onclick={() => (ui.spaceId = entry.space.room_id)}
		/>
	{/each}
	{#each nav.invites as room (room.roomId)}
		<RailButton name={room.name} invite badge={{ count: 1, strong: true, bold: true }} onclick={() => accept(room)} />
	{/each}
	<IconButton variant="ghost" aria-label="Create a space" title="Create a space" onclick={() => (ui.dialog = 'create-space')}>
		<Icon icon={icons.plus} size={20} />
	</IconButton>
</div>
