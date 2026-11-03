<script lang="ts">
	import type { MatrixEvent, Room } from 'matrix-js-sdk';
	import { toast } from '@facile/muse';
	import { app } from '$lib/state/app.svelte';
	import { toggleReaction } from '$lib/matrix/messages';
	import { reactionsFor, type Reaction } from '$lib/matrix/timeline';

	let { room, event }: { room: Room; event: MatrixEvent } = $props();

	const reactions = $derived.by(() => {
		void app.version;
		return reactionsFor(room, event, app.client?.getUserId() ?? '');
	});

	function who(reaction: Reaction): string {
		return reaction.senders.map((id) => room.getMember(id)?.name ?? id).join(', ');
	}

	async function toggle(reaction: Reaction) {
		if (!app.client) return;
		await toggleReaction(app.client, event, reaction.key, reaction.mine).catch(() => toast.danger('Could not react'));
	}
</script>

{#if reactions.length > 0}
	<div class="mt-1 flex flex-wrap gap-1">
		{#each reactions as reaction (reaction.key)}
			<button
				type="button"
				onclick={() => toggle(reaction)}
				title={who(reaction)}
				aria-pressed={!!reaction.mine}
				aria-label="{reaction.key} {reaction.count}, by {who(reaction)}"
				class="inline-flex h-7 items-center gap-1 rounded-fc-pill border px-2 text-fc-xs transition-colors {reaction.mine
					? 'border-fc-fg bg-fc-surface font-semibold'
					: 'border-fc-border hover:bg-fc-surface'}"
			>
				<span>{reaction.key}</span><span>{reaction.count}</span>
			</button>
		{/each}
	</div>
{/if}
