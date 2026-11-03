<script lang="ts">
	import type { AgentRecord, ChannelRecord, SpaceRecord } from '$lib/types';

	let {
		tab,
		activeSpace,
		activeRoomId,
		agents,
		dmRooms,
		onselectroom,
		oncreatechannel,
		onaddagent,
		onnewdm
	}: {
		tab: 'channels' | 'dms';
		activeSpace: SpaceRecord | null;
		activeRoomId: string | null;
		agents: AgentRecord[];
		dmRooms: Array<{ roomId: string; name: string; isAgent: boolean }>;
		onselectroom: (roomId: string, name: string) => void;
		oncreatechannel: () => void;
		onaddagent: () => void;
		onnewdm: () => void;
	} = $props();
</script>

<nav class="flex h-full w-60 flex-col border-r border-fc-border bg-fc-surface select-none shrink-0">
	{#if tab === 'channels' && activeSpace}
		<div class="flex h-14 items-center justify-between border-b border-fc-border px-4 font-semibold text-fc-fg">
			<span class="truncate text-fc-sm font-bold tracking-tight">{activeSpace.name}</span>
			<span class="rounded bg-fc-surface-hover px-1.5 py-0.5 text-[10px] text-fc-fg-muted">
				{activeSpace.slug}
			</span>
		</div>

		<div class="flex-1 overflow-y-auto p-2">
			<div class="mb-2 flex items-center justify-between px-2 pt-2 text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
				<span>Channels</span>
				<button
					type="button"
					class="flex size-5 items-center justify-center rounded hover:bg-fc-surface-hover hover:text-fc-fg"
					title="Create channel"
					onclick={oncreatechannel}
				>
					<svg width="14" height="14" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
						<line x1="12" y1="5" x2="12" y2="19" />
						<line x1="5" y1="12" x2="19" y2="12" />
					</svg>
				</button>
			</div>

			<div class="space-y-0.5">
				{#each activeSpace.channels as channel (channel.room_id)}
					<button
						type="button"
						class={`flex w-full items-center gap-2 rounded-fc-md px-2.5 py-1.5 text-left text-fc-sm transition-colors ${
							activeRoomId === channel.room_id
								? 'bg-fc-accent text-fc-accent-fg font-medium'
								: 'text-fc-fg hover:bg-fc-surface-hover'
						}`}
						onclick={() => onselectroom(channel.room_id, '#' + channel.name)}
					>
						<span class="text-fc-xs opacity-70">#</span>
						<span class="truncate">{channel.name}</span>
					</button>
				{/each}
			</div>

			<div class="mt-6 mb-2 flex items-center justify-between px-2 pt-2 text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
				<span>Agents</span>
				<button
					type="button"
					class="flex size-5 items-center justify-center rounded hover:bg-fc-surface-hover hover:text-fc-fg"
					title="Add an agent"
					onclick={onaddagent}
				>
					<svg width="14" height="14" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
						<line x1="12" y1="5" x2="12" y2="19" />
						<line x1="5" y1="12" x2="19" y2="12" />
					</svg>
				</button>
			</div>

			<div class="space-y-0.5">
				{#each agents as agent (agent.user_id)}
					<div class="flex items-center justify-between rounded-fc-md px-2.5 py-1.5 text-fc-sm text-fc-fg hover:bg-fc-surface-hover">
						<div class="flex items-center gap-2 truncate">
							<span class="flex size-4 items-center justify-center rounded-full bg-fc-accent/20 text-[10px] text-fc-accent">
								🤖
							</span>
							<span class="truncate font-medium">{agent.name}</span>
						</div>
						<span class="rounded bg-fc-surface px-1.5 py-0.5 text-[10px] text-fc-fg-muted font-mono">
							{agent.user_id.split(':')[0]}
						</span>
					</div>
				{/each}

				{#if agents.length === 0}
					<p class="px-2.5 py-1 text-fc-xs text-fc-fg-muted italic">
						No agents yet. Click + to add one.
					</p>
				{/if}
			</div>
		</div>
	{:else}
		<div class="flex h-14 items-center justify-between border-b border-fc-border px-4 font-semibold text-fc-fg">
			<span class="text-fc-sm font-bold tracking-tight">Direct Messages</span>
			<button
				type="button"
				class="flex size-6 items-center justify-center rounded-fc-md bg-fc-accent text-fc-accent-fg hover:opacity-90"
				title="New Message"
				onclick={onnewdm}
			>
				<svg width="14" height="14" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
					<line x1="12" y1="5" x2="12" y2="19" />
					<line x1="5" y1="12" x2="19" y2="12" />
				</svg>
			</button>
		</div>

		<div class="flex-1 overflow-y-auto p-2">
			<div class="space-y-0.5">
				{#each dmRooms as dm (dm.roomId)}
					<button
						type="button"
						class={`flex w-full items-center gap-2 rounded-fc-md px-2.5 py-1.5 text-left text-fc-sm transition-colors ${
							activeRoomId === dm.roomId
								? 'bg-fc-accent text-fc-accent-fg font-medium'
								: 'text-fc-fg hover:bg-fc-surface-hover'
						}`}
						onclick={() => onselectroom(dm.roomId, dm.name)}
					>
						<span class="flex size-5 items-center justify-center rounded-full bg-fc-border text-fc-xs font-bold">
							{dm.name[0] || '@'}
						</span>
						<span class="truncate">{dm.name}</span>
					</button>
				{/each}

				{#if dmRooms.length === 0}
					<p class="p-3 text-center text-fc-xs text-fc-fg-muted">
						No direct conversations yet. Click + to start one.
					</p>
				{/if}
			</div>
		</div>
	{/if}
</nav>
