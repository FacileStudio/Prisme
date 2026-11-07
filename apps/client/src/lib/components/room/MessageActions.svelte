<script lang="ts">
	import type { MatrixEvent, Room } from 'matrix-js-sdk';
	import { ConfirmModal, Icon, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import type { DraftState } from '$lib/state/draft.svelte';
	import { ui } from '$lib/state/ui.svelte';
	import { deleteMessage, toggleReaction } from '$lib/matrix/messages';
	import { callInfo, reactionsFor } from '$lib/matrix/timeline';
	import EmojiPicker from './EmojiPicker.svelte';

	type Props = { room: Room; event: MatrixEvent; draft: DraftState; inThread: boolean };
	let { room, event, draft, inThread }: Props = $props();
	let picking = $state(false);
	let confirming = $state(false);

	const mine = $derived(event.getSender() === app.client?.getUserId());
	const editable = $derived(mine && !callInfo(event) && event.getContent().msgtype === 'm.text');

	async function react(key: string) {
		picking = false;
		if (!app.client) return;
		const me = app.client.getUserId() ?? '';
		const existing = reactionsFor(room, event, me).find((r) => r.key === key)?.mine;
		await toggleReaction(app.client, event, key, existing).catch(() => toast.danger(t('actions.reactFailed')));
	}

	async function remove() {
		confirming = false;
		if (app.client) await deleteMessage(app.client, event).catch(() => toast.danger(t('actions.deleteFailed')));
	}
</script>

<div
	class="absolute -top-4 right-3 z-10 hidden items-center rounded-fc-md bg-fc-component shadow-sm group-hover:flex group-focus-within:flex {picking ? 'flex' : ''}"
	role="toolbar"
	aria-label={t('actions.toolbar')}
>
	<button type="button" class="action" aria-label={t('actions.reactAria')} title={t('actions.react')} onclick={() => (picking = !picking)}>
		<Icon icon={glyphs.smile} size={16} />
	</button>
	<button type="button" class="action" aria-label={t('actions.reply')} title={t('actions.reply')} onclick={() => draft.reply(event)}>
		<Icon icon={glyphs.reply} size={16} />
	</button>
	{#if !inThread}
		<button type="button" class="action" aria-label={t('actions.replyThread')} title={t('actions.replyThread')} onclick={() => ui.openThread(event.getId()!)}>
			<Icon icon={glyphs.thread} size={16} />
		</button>
	{/if}
	{#if editable}
		<button type="button" class="action" aria-label={t('actions.edit')} title={t('actions.edit')} onclick={() => draft.edit(event)}>
			<Icon icon={glyphs.pen} size={16} />
		</button>
	{/if}
	{#if mine}
		<button type="button" class="action hover:text-fc-danger" aria-label={t('common.delete')} title={t('common.delete')} onclick={() => (confirming = true)}>
			<Icon icon={icons.remove} size={16} />
		</button>
	{/if}
	{#if picking}
		<div class="absolute right-0 top-11">
			<EmojiPicker onpick={react} />
		</div>
	{/if}
</div>

<ConfirmModal
	open={confirming}
	title={t('actions.confirmTitle')}
	description={t('actions.confirmBody')}
	confirmLabel={t('common.delete')}
	tone="danger"
	onConfirm={remove}
	onCancel={() => (confirming = false)}
/>

<style>
	.action {
		display: inline-flex;
		width: 2.5rem;
		height: 2.5rem;
		align-items: center;
		justify-content: center;
		color: var(--color-fc-fg-muted);
	}
	.action:hover,
	.action:focus-visible {
		color: var(--color-fc-fg);
		background: var(--color-fc-surface);
	}
</style>
