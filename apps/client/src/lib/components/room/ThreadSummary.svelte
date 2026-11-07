<script lang="ts">
	import type { MatrixEvent, Room } from 'matrix-js-sdk';
	import { Icon } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { ui } from '$lib/state/ui.svelte';
	import { t } from '$lib/i18n/index.svelte';

	let { room, event }: { room: Room; event: MatrixEvent } = $props();

	const summary = $derived.by(() => {
		void app.version;
		const id = event.getId();
		const thread = id ? (room.getThread(id) ?? event.getThread()) : undefined;
		if (!thread || thread.length === 0) return null;
		const last = thread.replyToEvent ?? thread.lastReply();
		const sender = last?.getSender() ?? '';
		return { count: thread.length, last: room.getMember(sender)?.rawDisplayName || sender, ts: last?.getTs() ?? 0 };
	});
</script>

{#if summary}
	<button
		type="button"
		onclick={() => ui.openThread(event.getId()!)}
		class="mt-1 inline-flex min-h-9 items-center gap-2 rounded-fc-md px-2 text-fc-xs font-semibold text-fc-info hover:bg-fc-info/10"
	>
		<Icon icon={glyphs.thread} size={14} />
		{t('thread.replies', { count: summary.count })}
		<span class="font-normal text-fc-fg-muted">· {t('thread.lastFrom', { name: summary.last })}</span>
	</button>
{/if}
