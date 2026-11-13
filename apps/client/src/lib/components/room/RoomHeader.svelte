<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { Icon, IconButton, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { ui } from '$lib/state/ui.svelte';
	import { dmPartner } from '$lib/matrix/rooms';
	import UserName from '$lib/components/common/UserName.svelte';
	import CallButton from './CallButton.svelte';
	import HeaderMenu from './HeaderMenu.svelte';

	let { room, direct }: { room: Room; direct: boolean } = $props();

	const info = $derived.by(() => {
		void app.version;
		const partner = direct && app.client ? dmPartner(app.client, room) : null;
		const topic = room.currentState.getStateEvents('m.room.topic', '')?.getContent().topic as string | undefined;
		return { name: room.name, partner, topic: topic ?? '', members: room.getJoinedMemberCount() };
	});
</script>

<header class="flex min-h-14 items-center gap-1 px-1 md:px-3">
	<IconButton variant="ghost" class="md:hidden" aria-label={t('header.openMenu')} onclick={() => (ui.drawer = true)}>
		<Icon icon={glyphs.menu} size={20} />
	</IconButton>
	<div class="flex min-w-0 flex-1 flex-col">
		<h1 class="flex min-w-0 items-center gap-1 text-fc-md font-semibold">
			{#if direct && info.partner}
				<UserName userId={info.partner} name={info.name} />
			{:else}
				<Icon icon={glyphs.hashtag} size={16} class="shrink-0 text-fc-fg-muted" />
				<span class="truncate">{info.name}</span>
			{/if}
		</h1>
		{#if info.topic}
			<p class="truncate text-fc-xs text-fc-fg-muted">{info.topic}</p>
		{/if}
	</div>
	<CallButton roomId={room.roomId} />
	<IconButton variant="ghost" aria-label={t('header.searchAria')} tooltip={t('common.search')} aria-pressed={ui.panel === 'search'} onclick={() => ui.toggle('search')}>
		<Icon icon={icons.search} size={18} />
	</IconButton>
	<IconButton variant="ghost" aria-label={t('header.membersAria', { count: info.members })} tooltip={t('common.members')} aria-pressed={ui.panel === 'members'} onclick={() => ui.toggle('members')}>
		<Icon icon={glyphs.users} size={18} />
	</IconButton>
	<HeaderMenu {room} {direct} />
</header>
