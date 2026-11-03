<script lang="ts">
	import type { DmRow } from '$lib/state/nav';
	import { roomHref } from '$lib/state/actions';
	import { ui } from '$lib/state/ui.svelte';
	import RoomRow from './RoomRow.svelte';
	import SectionTitle from './SectionTitle.svelte';
	import PresenceAvatar from '$lib/components/common/PresenceAvatar.svelte';
	import UserName from '$lib/components/common/UserName.svelte';

	let { rows }: { rows: DmRow[] } = $props();
</script>

<SectionTitle title="Direct messages" action="New message" onaction={() => (ui.dialog = 'new-dm')} />
<ul class="flex flex-col px-2">
	{#each rows as row (row.roomId)}
		<li>
			<RoomRow href={roomHref(row.roomId)} roomId={row.roomId} name={row.name} badge={row.badge}>
				{#snippet lead()}
					<PresenceAvatar userId={row.userId ?? ''} name={row.name} />
				{/snippet}
				{#snippet label()}
					<UserName userId={row.userId ?? ''} name={row.invite ? `${row.name} (invite)` : row.name} />
				{/snippet}
			</RoomRow>
		</li>
	{:else}
		<li class="px-2 py-2 text-fc-xs text-fc-fg-muted">No conversations yet.</li>
	{/each}
</ul>
