<script lang="ts">
	import { Input, Spinner } from '@facile/muse';
	import { prisme, type DirectoryUser } from '$lib/api/prisme';
	import { isBot } from '$lib/chat/users';
	import { app } from '$lib/state/app.svelte';
	import UserName from '$lib/components/common/UserName.svelte';
	import MxcAvatar from '$lib/components/common/MxcAvatar.svelte';

	type Props = { onpick: (user: DirectoryUser) => void; busyId?: string; exclude?: string[] };
	let { onpick, busyId = '', exclude = [] }: Props = $props();
	let query = $state('');
	let loading = $state(false);
	let users = $state.raw<DirectoryUser[]>([]);

	$effect(() => {
		const q = query.trim();
		const me = app.session?.matrix.user_id;
		loading = true;
		const timer = setTimeout(async () => {
			const found = await prisme.users(q).catch(() => []);
			users = found.filter((u) => u.user_id !== me && !isBot(u.user_id));
			loading = false;
		}, 200);
		return () => clearTimeout(timer);
	});
</script>

<Input bind:value={query} type="search" placeholder="Search people by name or email" aria-label="Search people" />
<ul class="flex max-h-72 flex-col overflow-y-auto" aria-busy={loading}>
	{#each users as user (user.user_id)}
		{@const taken = exclude.includes(user.user_id)}
		<li>
			<button
				type="button"
				disabled={taken || busyId === user.user_id}
				onclick={() => onpick(user)}
				class="flex min-h-11 w-full items-center gap-3 rounded-fc-md px-2 text-left text-fc-sm hover:bg-fc-surface disabled:opacity-50"
			>
				<MxcAvatar name={user.name} class="size-7 text-fc-xs" />
				<span class="flex min-w-0 flex-1 flex-col">
					<UserName userId={user.user_id} name={user.name} class="font-medium" />
					<span class="truncate text-fc-xs text-fc-fg-muted">{user.user_id}</span>
				</span>
				{#if busyId === user.user_id}<Spinner size="sm" />{:else if taken}<span class="text-fc-xs text-fc-fg-muted">Member</span>{/if}
			</button>
		</li>
	{:else}
		<li class="px-2 py-3 text-fc-sm text-fc-fg-muted">{loading ? 'Searching…' : 'Nobody matches.'}</li>
	{/each}
</ul>
