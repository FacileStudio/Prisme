<script lang="ts">
	import { app } from '$lib/state/app.svelte';
	import { buildNav } from '$lib/state/nav';
	import { HOME, ui } from '$lib/state/ui.svelte';
	import SpaceRail from './SpaceRail.svelte';
	import SpaceChannels from './SpaceChannels.svelte';
	import DmSection from './DmSection.svelte';
	import UserCard from './UserCard.svelte';

	const nav = $derived.by(() => {
		void app.version;
		return app.client ? buildNav(app.client, app.spaces) : null;
	});

	const current = $derived(nav?.spaces.find((s) => s.space.room_id === ui.spaceId) ?? null);

	$effect(() => {
		if (!nav || current) return;
		if (ui.spaceId === HOME && nav.spaces.length > 0) ui.spaceId = nav.spaces[0].space.room_id;
	});
</script>

{#if nav}
	<div class="flex h-full">
		<SpaceRail {nav} />
		<nav class="flex h-full w-64 flex-col border-r border-fc-border bg-fc-component" aria-label="Channels and messages">
			<div class="flex-1 overflow-y-auto overscroll-contain pb-4">
				{#if current}
					<SpaceChannels nav={current} />
				{/if}
				<DmSection rows={nav.dms} />
			</div>
			<UserCard />
		</nav>
	</div>
{/if}
