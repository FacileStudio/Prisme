<script lang="ts">
	import { goto } from '$app/navigation';
	import type { Room } from 'matrix-js-sdk';
	import { Button, Card, EmptyState, Icon, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { acceptSpace, roomHref } from '$lib/state/actions';
	import { buildNav } from '$lib/state/nav';
	import { MANAGE_SPACES_HREF, ui } from '$lib/state/ui.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import SpaceChannels from './SpaceChannels.svelte';

	let accepting = $state('');

	const nav = $derived.by(() => {
		void app.version;
		return app.client ? buildNav(app.client, app.spaces, app.agents) : null;
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
			toast.danger(err instanceof Error ? err.message : t('overview.joinFailed'));
		} finally {
			accepting = '';
		}
	}
</script>

<svelte:head>
	<title>{current ? t('overview.titleSpace', { name: current.name }) : t('overview.title')}</title>
</svelte:head>

{#if nav}
	<div class="mx-auto flex w-full max-w-3xl flex-col gap-8 px-4 py-6 md:px-6 md:py-10">
		{#each nav.invites as room (room.roomId)}
			<Card class="flex flex-col gap-3 sm:flex-row sm:items-center">
				<span class="flex size-10 shrink-0 items-center justify-center rounded-fc-md bg-fc-accent text-fc-accent-fg">
					<Icon icon={icons.mail} size={20} />
				</span>
				<p class="min-w-0 flex-1 text-fc-sm text-fc-fg">{t('overview.invited', { name: room.name })}</p>
				<Button size="lg" icon={glyphs.login} disabled={accepting === room.roomId} onclick={() => accept(room)}>
					{accepting === room.roomId ? t('common.joining') : t('overview.joinSpace')}
				</Button>
			</Card>
		{/each}

		{#if current}
			<SpaceChannels nav={current} />
		{:else if nav.spaces.length === 0}
			<EmptyState
				icon={icons.space}
				title={t('overview.none')}
				description={t('overview.noneBody')}
			/>
		{/if}

		<div class="flex flex-col gap-3 pt-6 sm:flex-row sm:items-center">
			<p class="min-w-0 flex-1 text-fc-sm text-fc-fg-muted">{t('overview.footer')}</p>
			<div class="flex flex-wrap gap-2">
				<Button size="lg" variant="ghost" icon={icons.space} href={MANAGE_SPACES_HREF}>{t('spaces.manage')}</Button>
				<Button size="lg" variant="secondary" icon={icons.plus} onclick={() => (ui.dialog = 'create-space')}>{t('overview.createSpace')}</Button>
			</div>
		</div>
	</div>
{/if}
