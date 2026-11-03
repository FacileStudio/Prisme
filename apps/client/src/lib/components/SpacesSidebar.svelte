<script lang="ts">
	import type { SpaceRecord, UserProfile } from '$lib/types';

	let {
		spaces,
		activeSpaceId,
		activeTab,
		user,
		onselectspace,
		onselecttab,
		oncreatespace,
		onlogout
	}: {
		spaces: SpaceRecord[];
		activeSpaceId: number | null;
		activeTab: 'channels' | 'dms';
		user: UserProfile;
		onselectspace: (spaceId: number) => void;
		onselecttab: (tab: 'channels' | 'dms') => void;
		oncreatespace: () => void;
		onlogout: () => void;
	} = $props();

	function initials(name: string): string {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.slice(0, 2)
			.toUpperCase();
	}
</script>

<aside class="flex h-full w-18 flex-col items-center border-r border-fc-border bg-fc-panel py-3 select-none shrink-0">
	<div class="mb-4 flex flex-col items-center">
		<div class="flex size-11 items-center justify-center rounded-fc-lg bg-fc-accent text-fc-accent-fg font-bold shadow-sm" title="Prisme">
			<svg width="22" height="22" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
				<polygon points="12 2 2 22 22 22" />
			</svg>
		</div>
	</div>

	<div class="flex w-full flex-1 flex-col items-center gap-2 overflow-y-auto px-2">
		<button
			type="button"
			class={`flex size-11 items-center justify-center rounded-fc-lg text-fc-xs font-semibold transition-all ${
				activeTab === 'dms'
					? 'bg-fc-accent text-fc-accent-fg shadow-sm'
					: 'bg-fc-surface text-fc-fg hover:bg-fc-surface-hover'
			}`}
			title="Direct Messages"
			onclick={() => onselecttab('dms')}
		>
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
				<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
			</svg>
		</button>

		<div class="my-1 h-px w-8 bg-fc-border"></div>

		{#each spaces as space (space.id)}
			<button
				type="button"
				class={`group relative flex size-11 items-center justify-center rounded-fc-lg text-fc-xs font-bold transition-all ${
					activeTab === 'channels' && activeSpaceId === space.id
						? 'bg-fc-accent text-fc-accent-fg shadow-sm'
						: 'bg-fc-surface text-fc-fg hover:bg-fc-surface-hover'
				}`}
				title={space.name}
				onclick={() => {
					onselecttab('channels');
					onselectspace(space.id);
				}}
			>
				{initials(space.name)}
				{#if activeTab === 'channels' && activeSpaceId === space.id}
					<span class="absolute -left-2 h-6 w-1 rounded-r-full bg-fc-accent"></span>
				{/if}
			</button>
		{/each}

		<button
			type="button"
			class="flex size-11 items-center justify-center rounded-fc-lg border border-dashed border-fc-border text-fc-fg-muted transition-colors hover:border-fc-accent hover:text-fc-accent"
			title="Create Space"
			onclick={oncreatespace}
		>
			<svg width="20" height="20" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
				<line x1="12" y1="5" x2="12" y2="19" />
				<line x1="5" y1="12" x2="19" y2="12" />
			</svg>
		</button>
	</div>

	<div class="mt-auto flex flex-col items-center gap-2 pt-2">
		<div
			class="flex size-10 items-center justify-center overflow-hidden rounded-full bg-fc-surface text-fc-fg font-semibold border border-fc-border"
			title={`${user.name} (${user.email})`}
		>
			{#if user.avatar_url}
				<img src={user.avatar_url} alt={user.name} class="size-full object-cover" />
			{:else}
				{initials(user.name || user.email)}
			{/if}
		</div>

		<button
			type="button"
			class="flex size-8 items-center justify-center rounded-fc-md text-fc-fg-muted transition-colors hover:bg-fc-surface-hover hover:text-fc-danger"
			title="Log out"
			onclick={onlogout}
		>
			<svg width="16" height="16" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
				<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
				<polyline points="16 17 21 12 16 7" />
				<line x1="21" y1="12" x2="9" y2="12" />
			</svg>
		</button>
	</div>
</aside>
