<script lang="ts">
	import { Button, Section, icons } from '@facile/muse';
	import type { DmRow } from '$lib/state/nav';
	import { roomHref } from '$lib/state/actions';
	import { ui } from '$lib/state/ui.svelte';
	import RoomRow from './RoomRow.svelte';
	import PresenceAvatar from '$lib/components/common/PresenceAvatar.svelte';
	import UserName from '$lib/components/common/UserName.svelte';

	let { rows }: { rows: DmRow[] } = $props();
</script>

<Section title="Direct messages" card gap="tight" bodyClass="gap-0 p-2">
	{#snippet actions()}
		<Button size="sm" variant="ghost" icon={icons.plus} onclick={() => (ui.dialog = 'new-dm')}>New message</Button>
	{/snippet}
	{#if rows.length > 0}
		<ul class="flex flex-col">
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
			{/each}
		</ul>
	{:else}
		<p class="px-3 py-3 text-fc-sm text-fc-fg-muted">No conversations yet.</p>
	{/if}
</Section>
