<script lang="ts">
	import { Avatar, Button, Icon, IconButton, Input, icons } from '@facile/muse';
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
	let timer: ReturnType<typeof setTimeout> | null = null;

	function handleSearch() {
		if (timer) clearTimeout(timer);
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

<svelte:window onkeydown={(e) => { if (e.key === 'Escape') onclose(); }} />

<div
	role="presentation"
	class="fixed inset-0 z-50 flex items-center justify-center bg-fc-scrim p-4"
	onclick={(e) => { if (e.target === e.currentTarget) onclose(); }}
>
	<div class="w-full max-w-md rounded-fc-xl bg-fc-component p-6 shadow-xl text-fc-fg">
		<div class="flex items-center justify-between pb-4">
			<h2 class="text-fc-lg font-bold text-fc-fg">New Direct Message</h2>
			<IconButton variant="ghost" class="size-8" aria-label="Close" onclick={onclose}>
				<Icon icon={icons.close} size={18} />
			</IconButton>
		</div>

		<div>
			<Input
				bind:value={query}
				oninput={handleSearch}
				placeholder="Search by name, email or @handle…"
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
						<div class="flex items-center gap-2.5 min-w-0">
							<Avatar name={u.name} size="sm" />
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
			<Button variant="ghost" onclick={onclose}>
				Cancel
			</Button>
		</div>
	</div>
</div>
