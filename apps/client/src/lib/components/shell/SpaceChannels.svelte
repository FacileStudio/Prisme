<script lang="ts">
	import { goto } from '$app/navigation';
	import { Icon, IconButton, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import type { SpaceNav } from '$lib/state/nav';
	import { joinChannel, roomHref } from '$lib/state/actions';
	import { canRename } from '$lib/matrix/power';
	import { ui } from '$lib/state/ui.svelte';
	import RoomRow from './RoomRow.svelte';
	import SectionTitle from './SectionTitle.svelte';

	let { nav }: { nav: SpaceNav } = $props();
	let joining = $state('');

	const admin = $derived.by(() => {
		void app.version;
		const room = app.client?.getRoom(nav.space.room_id);
		return !!room && !!app.client && canRename(app.client, room);
	});

	async function join(roomId: string) {
		joining = roomId;
		try {
			await joinChannel(roomId);
			await goto(roomHref(roomId));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not join the channel');
		} finally {
			joining = '';
		}
	}
</script>

<div class="flex min-h-11 items-center gap-1 py-1 pl-4 pr-1">
	<h2 class="min-w-0 flex-1 truncate text-fc-xl font-semibold">{nav.name}</h2>
	<IconButton variant="ghost" aria-label="Invite people to {nav.name}" title="Invite people" onclick={() => (ui.dialog = 'invite')}>
		<Icon icon={glyphs.userPlus} size={18} />
	</IconButton>
	{#if admin}
		<IconButton variant="ghost" aria-label="Rename {nav.name}" title="Rename the space" onclick={() => (ui.dialog = 'rename-space')}>
			<Icon icon={icons.edit} size={18} />
		</IconButton>
	{/if}
</div>

<SectionTitle title="Channels" action="Create a channel" onaction={() => (ui.dialog = 'create-channel')} />
<ul class="flex flex-col px-2">
	{#each nav.joined as row (row.roomId)}
		<li><RoomRow href={roomHref(row.roomId)} name={row.name} badge={row.badge} muted={row.muted} icon={glyphs.hashtag} roomId={row.roomId} /></li>
	{/each}
</ul>

{#if nav.joinable.length > 0}
	<SectionTitle title="More channels" />
	<ul class="flex flex-col px-2">
		{#each nav.joinable as row (row.roomId)}
			<li>
				<button
					type="button"
					class="flex min-h-11 w-full items-center gap-2 rounded-fc-md px-2 text-left text-fc-sm text-fc-fg-muted hover:bg-fc-surface disabled:opacity-50"
					onclick={() => join(row.roomId)}
					disabled={joining === row.roomId}
				>
					<Icon icon={glyphs.hashtag} size={16} />
					<span class="min-w-0 flex-1 truncate">{row.name}</span>
					<span class="text-fc-xs font-medium">{joining === row.roomId ? 'Joining…' : 'Join'}</span>
				</button>
			</li>
		{/each}
	</ul>
{/if}
