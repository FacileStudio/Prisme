<script lang="ts">
	import { page } from '$app/state';
	import { Icon, IconButton, icons } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { ui } from '$lib/state/ui.svelte';
	import ThreadPanel from './ThreadPanel.svelte';
	import MembersPanel from './MembersPanel.svelte';
	import SearchPanel from './SearchPanel.svelte';

	const room = $derived.by(() => {
		void app.version;
		const id = page.params.id ? decodeURIComponent(page.params.id) : '';
		return id ? (app.client?.getRoom(id) ?? null) : null;
	});
	const title = $derived(t(ui.panel === 'thread' ? 'panel.thread' : ui.panel === 'members' ? 'panel.members' : 'panel.search'));
	const closeLabel = $derived(
		t(ui.panel === 'thread' ? 'panel.closeThread' : ui.panel === 'members' ? 'panel.closeMembers' : 'panel.closeSearch')
	);
</script>

{#if ui.panel && room}
	<aside class="fixed inset-0 z-30 flex flex-col bg-fc-page md:static md:z-auto md:w-96" aria-label={title}>
		<div class="flex min-h-14 items-center pl-4 pr-1">
			<h2 class="flex-1 text-fc-md font-semibold">{title}</h2>
			<IconButton variant="ghost" aria-label={closeLabel} onclick={() => ui.closePanel()}>
				<Icon icon={icons.close} size={20} />
			</IconButton>
		</div>
		{#if ui.panel === 'thread' && ui.threadId}
			{#key ui.threadId}
				<ThreadPanel {room} rootId={ui.threadId} />
			{/key}
		{:else if ui.panel === 'members'}
			<MembersPanel {room} />
		{:else if ui.panel === 'search'}
			<SearchPanel {room} />
		{/if}
	</aside>
{/if}
