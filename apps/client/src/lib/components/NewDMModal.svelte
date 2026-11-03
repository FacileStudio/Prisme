<script lang="ts">
	import { backend } from '$lib/backend';
	import type { UserRecord } from '$lib/types';

	let {
		onclose,
		onstartdm
	}: {
		onclose: () => void;
		onstartdm: (user: UserRecord) => void;
	} = $props();

	let query = $state('');
	let users = $state<UserRecord[]>([]);
	let searching = $state(false);
	let timer: any = null;

	function handleSearch() {
		clearTimeout(timer);
		const trimmed = query.trim();
		if (!trimmed) {
			users = [];
			searching = false;
			return;
		}
		searching = true;
		timer = setTimeout(async () => {
			try {
				users = await backend.searchUsers(trimmed);
			} catch {
				users = [];
			} finally {
				searching = false;
			}
		}, 250);
	}
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
	<div class="w-full max-w-md rounded-fc-lg border border-fc-border bg-fc-page p-6 shadow-xl">
		<div class="flex items-center justify-between pb-4">
			<h2 class="text-fc-lg font-bold text-fc-fg">New Direct Message</h2>
			<button
				type="button"
				aria-label="Close"
				class="rounded-fc-md p-1 text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg"
				onclick={onclose}
			>
				<svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
					<line x1="18" y1="6" x2="6" y2="18" />
					<line x1="6" y1="6" x2="18" y2="18" />
				</svg>
			</button>
		</div>

		<div>
			<input
				type="text"
				bind:value={query}
				oninput={handleSearch}
				placeholder="Search by name, email or @handle…"
				class="w-full rounded-fc-md border border-fc-border bg-fc-surface px-3 py-2 text-fc-sm text-fc-fg placeholder:text-fc-fg-muted focus:border-fc-accent focus:outline-none focus:ring-1 focus:ring-fc-accent"
			/>
		</div>

		<div class="mt-4 max-h-60 overflow-y-auto space-y-1">
			{#if searching}
				<p class="p-3 text-center text-fc-xs text-fc-fg-muted">Searching users…</p>
			{:else if users.length > 0}
				{#each users as u (u.user_id)}
					<button
						type="button"
						class="flex w-full items-center justify-between rounded-fc-md p-2 text-left hover:bg-fc-surface transition-colors"
						onclick={() => {
							onstartdm(u);
							onclose();
						}}
					>
						<div class="flex items-center gap-2 min-w-0">
							<div class="flex size-7 shrink-0 items-center justify-center rounded-full bg-fc-surface border border-fc-border text-fc-xs font-bold">
								{u.name.slice(0, 2).toUpperCase()}
							</div>
							<span class="truncate text-fc-sm font-semibold text-fc-fg">{u.name}</span>
						</div>
						<span class="text-fc-xs text-fc-fg-muted font-mono">{u.user_id}</span>
					</button>
				{/each}
			{:else if query.trim()}
				<p class="p-3 text-center text-fc-xs text-fc-fg-muted">No users found.</p>
			{:else}
				<p class="p-3 text-center text-fc-xs text-fc-fg-muted">Type to search for people to message.</p>
			{/if}
		</div>

		<div class="mt-4 flex justify-end pt-2">
			<button
				type="button"
				class="rounded-fc-md border border-fc-border px-4 py-2 text-fc-sm font-medium text-fc-fg hover:bg-fc-surface"
				onclick={onclose}
			>
				Cancel
			</button>
		</div>
	</div>
</div>
