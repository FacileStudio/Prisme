<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { Icon, Tooltip, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import type { DraftState } from '$lib/state/draft.svelte';
	import { t } from '$lib/i18n/index.svelte';

	let { room, draft }: { room: Room; draft: DraftState } = $props();

	const target = $derived(draft.editing ?? draft.replyTo);
	const who = $derived(target ? (room.getMember(target.getSender() ?? '')?.name ?? target.getSender()) : '');
</script>

{#if target}
	<div class="flex items-center gap-2 rounded-fc-pill bg-fc-surface py-1 pl-3 pr-1 text-fc-xs text-fc-fg-muted">
		<Icon icon={draft.editing ? glyphs.pen : glyphs.reply} size={14} />
		<span class="min-w-0 flex-1 truncate">
			{#if draft.editing}{t('draft.editing')}{:else}{t('draft.replying', { name: String(who), text: String(target.getContent().body ?? '') })}{/if}
		</span>
		<Tooltip text={t('common.cancel')} describe={false} class="contents">
			<button type="button" class="flex size-9 items-center justify-center rounded-fc-pill hover:bg-fc-surface" aria-label={t('common.cancel')} onclick={() => draft.clear()}>
				<Icon icon={icons.close} size={14} />
			</button>
		</Tooltip>
	</div>
{/if}
