<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button, Icon, IconButton, Section, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import type { SpaceNav } from '$lib/state/nav';
	import { joinChannel, roomHref } from '$lib/state/actions';
	import { canRename } from '$lib/matrix/power';
	import { ui } from '$lib/state/ui.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import RoomRow from './RoomRow.svelte';

	let { nav }: { nav: SpaceNav } = $props();
	let joining = $state('');

	const info = $derived.by(() => {
		void app.version;
		const room = app.client?.getRoom(nav.space.room_id);
		return {
			admin: !!room && !!app.client && canRename(app.client, room),
			members: room?.getJoinedMemberCount() ?? 0
		};
	});

	const total = $derived(nav.joined.length + nav.joinable.length);

	async function join(roomId: string) {
		joining = roomId;
		try {
			await joinChannel(roomId);
			await goto(roomHref(roomId));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('space.joinFailed'));
		} finally {
			joining = '';
		}
	}
</script>

<header class="flex flex-col gap-4 sm:flex-row sm:items-center">
	<div class="flex min-w-0 flex-1 items-center gap-4">
		<span
			class="flex size-14 shrink-0 items-center justify-center rounded-fc-md bg-fc-accent text-fc-xl font-semibold text-fc-accent-fg"
			aria-hidden="true">{nav.name.charAt(0).toUpperCase()}</span
		>
		<div class="flex min-w-0 flex-col gap-1">
			<div class="flex min-w-0 items-center gap-1">
				<h1 class="truncate text-fc-2xl font-semibold text-fc-fg">{nav.name}</h1>
				{#if info.admin}
					<IconButton variant="ghost" aria-label={t('space.renameAria', { name: nav.name })} tooltip={t('space.renameTitle')} onclick={() => (ui.dialog = 'rename-space')}>
						<Icon icon={icons.edit} size={16} />
					</IconButton>
				{/if}
			</div>
			<p class="text-fc-sm text-fc-fg-muted">{t('space.channelCount', { count: total })} · {t('space.memberCount', { count: info.members })}</p>
		</div>
	</div>
	<div class="flex flex-wrap gap-2">
		<Button size="lg" variant="secondary" icon={glyphs.userPlus} onclick={() => (ui.dialog = 'invite')}>{t('space.invite')}</Button>
		<Button size="lg" icon={icons.plus} onclick={() => (ui.dialog = 'create-channel')}>{t('space.newChannel')}</Button>
	</div>
</header>

<Section title={t('space.yours')} card gap="tight" bodyClass="gap-0 p-2">
	{#if nav.joined.length > 0}
		<ul class="flex flex-col">
			{#each nav.joined as row (row.roomId)}
				<li><RoomRow href={roomHref(row.roomId)} name={row.name} badge={row.badge} muted={row.muted} icon={glyphs.hashtag} roomId={row.roomId} /></li>
			{/each}
		</ul>
	{:else}
		<p class="px-3 py-3 text-fc-sm text-fc-fg-muted">{t('space.noneJoined', { name: nav.name })}</p>
	{/if}
</Section>

{#if nav.joinable.length > 0}
	<Section title={t('space.browse')} description={t('space.browseBody', { name: nav.name })} card gap="tight" bodyClass="gap-0 p-2">
		<ul class="flex flex-col">
			{#each nav.joinable as row (row.roomId)}
				<li class="flex min-h-11 items-center gap-2 rounded-fc-md px-2 text-fc-sm text-fc-fg-muted hover:bg-fc-surface">
					<Icon icon={glyphs.hashtag} size={16} class="shrink-0" />
					<span class="min-w-0 flex-1 truncate">{row.name}</span>
					<Button size="sm" variant="ghost" icon={glyphs.login} disabled={joining === row.roomId} onclick={() => join(row.roomId)}>
						{joining === row.roomId ? t('common.joining') : t('common.join')}
					</Button>
				</li>
			{/each}
		</ul>
	</Section>
{/if}
