<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Room } from 'matrix-js-sdk';
	import { Button, EmptyState, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { acceptSpace, roomHref } from '$lib/state/actions';
	import { buildNav } from '$lib/state/nav';
	import { MANAGE_SPACES_HREF, MANAGE_SPACES_LABEL, ui } from '$lib/state/ui.svelte';
	import SpaceChannels from './SpaceChannels.svelte';
	import DmSection from './DmSection.svelte';
	import SectionTitle from './SectionTitle.svelte';

	let accepting = $state('');

	const nav = $derived.by(() => {
		void app.version;
		return app.client ? buildNav(app.client, app.spaces) : null;
	});

	const current = $derived(nav?.spaces.find((s) => s.space.room_id === ui.spaceId) ?? null);

	async function accept(room: Room) {
		accepting = room.roomId;
		try {
			const space = await acceptSpace(room);
			if (!space) return;
			ui.spaceId = space.room_id;
			const general = space.channels.find((c) => c.name === 'general');
			if (general) await goto(roomHref(general.room_id));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not join the space');
		} finally {
			accepting = '';
		}
	}
</script>

<svelte:head>
	<title>Chats · Prisme</title>
</svelte:head>

{#if nav}
	<div class="mx-auto flex w-full max-w-2xl flex-col gap-2 py-4 md:px-6 md:py-8">
		{#if nav.invites.length > 0}
			<SectionTitle title="Space invites" />
			<ul class="flex flex-col gap-1 px-2">
				{#each nav.invites as room (room.roomId)}
					<li class="flex min-h-11 items-center gap-3 rounded-fc-md bg-fc-component px-3 py-2">
						<span class="min-w-0 flex-1 truncate text-fc-sm font-medium">{room.name}</span>
						<Button size="lg" icon={glyphs.login} disabled={accepting === room.roomId} onclick={() => accept(room)}>
							{accepting === room.roomId ? 'Joining…' : 'Join'}
						</Button>
					</li>
				{/each}
			</ul>
		{/if}

		{#if current}
			<SpaceChannels nav={current} />
		{:else if nav.spaces.length === 0}
			<EmptyState
				icon={glyphs.chat}
				title="No space yet"
				description="A space groups the channels of a project or a team."
			/>
		{/if}

		<DmSection rows={nav.dms} />

		<SectionTitle title="Spaces" />
		<div class="flex flex-wrap gap-2 px-4">
			<Button size="lg" variant="outline" icon={icons.plus} onclick={() => (ui.dialog = 'create-space')}>Create a space</Button>
			<Button size="lg" variant="ghost" icon={icons.usersGroup} href={MANAGE_SPACES_HREF}>{MANAGE_SPACES_LABEL}</Button>
		</div>
	</div>
{/if}
