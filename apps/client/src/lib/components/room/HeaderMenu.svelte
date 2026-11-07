<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { Icon, IconButton, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { setMuted } from '$lib/state/actions';
	import { ui, type Dialog } from '$lib/state/ui.svelte';
	import { isMuted } from '$lib/matrix/rooms';
	import { canRename } from '$lib/matrix/power';

	let { room, direct }: { room: Room; direct: boolean } = $props();
	let open = $state(false);

	const info = $derived.by(() => {
		void app.version;
		const client = app.client;
		return {
			muted: !!client && isMuted(client, room.roomId),
			admin: !!client && canRename(client, room)
		};
	});

	async function toggleMute() {
		open = false;
		const mute = !info.muted;
		try {
			await setMuted(room.roomId, mute);
			toast.success(mute ? t('menu.muted') : t('menu.unmuted'));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('menu.muteFailed'));
		}
	}

	function show(dialog: Dialog) {
		open = false;
		ui.dialog = dialog;
	}
</script>

<div class="relative">
	<IconButton variant="ghost" aria-label={t('menu.more')} aria-expanded={open} aria-haspopup="menu" onclick={() => (open = !open)}>
		<Icon icon={glyphs.dots} size={18} />
	</IconButton>
	{#if open}
		<button type="button" class="fixed inset-0 z-20 cursor-default" aria-label={t('menu.close')} onclick={() => (open = false)}></button>
		<div role="menu" class="absolute right-0 top-12 z-30 flex w-56 flex-col rounded-fc-lg bg-fc-component p-1 shadow-lg">
			<button type="button" role="menuitem" class="menu-item" onclick={toggleMute}>
				<Icon icon={info.muted ? glyphs.bell : glyphs.bellOff} size={16} />
				{t(info.muted ? (direct ? 'menu.unmuteDirect' : 'menu.unmuteChannel') : direct ? 'menu.muteDirect' : 'menu.muteChannel')}
			</button>
			{#if !direct}
				<button type="button" role="menuitem" class="menu-item" onclick={() => show('add-agent')}>
					<Icon icon={glyphs.mention} size={16} />{t('menu.addAgent')}
				</button>
			{/if}
			{#if info.admin && !direct}
				<button type="button" role="menuitem" class="menu-item" onclick={() => show('rename-room')}>
					<Icon icon={icons.edit} size={16} />{t('menu.renameChannel')}
				</button>
			{/if}
		</div>
	{/if}
</div>

<style>
	.menu-item {
		display: flex;
		min-height: 2.75rem;
		align-items: center;
		gap: 0.5rem;
		border-radius: var(--radius-fc-md);
		padding: 0 0.75rem;
		text-align: left;
		font-size: var(--text-fc-sm);
	}
	.menu-item:hover,
	.menu-item:focus-visible {
		background: var(--color-fc-surface);
	}
</style>
