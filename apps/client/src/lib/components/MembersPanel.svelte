<script lang="ts">
	import type { RoomMember } from '$lib/types';

	let {
		members,
		onclose
	}: {
		members: RoomMember[];
		onclose: () => void;
	} = $props();

	const agents = $derived(members.filter((m) => m.isAgent));
	const humans = $derived(members.filter((m) => !m.isAgent));
</script>

<aside class="flex h-full w-64 flex-col border-l border-fc-border bg-fc-surface select-none shrink-0">
	<div class="flex h-14 items-center justify-between border-b border-fc-border px-4 font-bold text-fc-fg">
		<span class="text-fc-sm">Members ({members.length})</span>
		<button
			type="button"
			aria-label="Close"
			class="flex size-7 items-center justify-center rounded-fc-md text-fc-fg-muted hover:bg-fc-surface-hover hover:text-fc-fg"
			onclick={onclose}
		>
			<svg width="16" height="16" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
				<line x1="18" y1="6" x2="6" y2="18" />
				<line x1="6" y1="6" x2="18" y2="18" />
			</svg>
		</button>
	</div>

	<div class="flex-1 overflow-y-auto p-4 space-y-4">
		{#if agents.length > 0}
			<div>
				<h3 class="mb-2 text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
					Agents ({agents.length})
				</h3>
				<div class="space-y-1">
					{#each agents as agent (agent.userId)}
						<div class="flex items-center gap-2 rounded-fc-md p-1.5 hover:bg-fc-surface-hover">
							<div class="flex size-7 items-center justify-center rounded-full bg-fc-accent/20 text-fc-xs font-bold text-fc-accent">
								🤖
							</div>
							<div class="flex-1 min-w-0">
								<p class="truncate text-fc-sm font-semibold text-fc-fg">{agent.displayName}</p>
								<p class="truncate text-[10px] text-fc-fg-muted font-mono">{agent.userId}</p>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}

		<div>
			<h3 class="mb-2 text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
				Humans ({humans.length})
			</h3>
			<div class="space-y-1">
				{#each humans as human (human.userId)}
					<div class="flex items-center gap-2 rounded-fc-md p-1.5 hover:bg-fc-surface-hover">
						<div class="flex size-7 items-center justify-center rounded-full bg-fc-surface text-[11px] font-bold border border-fc-border">
							{human.displayName.slice(0, 2).toUpperCase()}
						</div>
						<div class="flex-1 min-w-0">
							<p class="truncate text-fc-sm font-semibold text-fc-fg">{human.displayName}</p>
							<p class="truncate text-[10px] text-fc-fg-muted font-mono">{human.userId}</p>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>
</aside>
