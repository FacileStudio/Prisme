<script lang="ts">
	import { goto } from '$app/navigation';
	import { Badge, EmptyState, Icon, Input, Page, PageHeader, Section, Skeleton, Spinner, Tabs, icons, toast } from '@facile/muse';
	import { prisme, type DirectoryUser } from '$lib/api/prisme';
	import { handle, isBot } from '$lib/chat/users';
	import { glyphs } from '$lib/glyphs';
	import { openDm } from '$lib/matrix/directs';
	import { shortTime } from '$lib/matrix/preview';
	import { app } from '$lib/state/app.svelte';
	import { roomHref } from '$lib/state/actions';
	import { buildNav, type DmRow } from '$lib/state/nav';
	import Avatar from '$lib/components/common/Avatar.svelte';
	import CountBadge from '$lib/components/common/CountBadge.svelte';

	type Filter = 'all' | 'people' | 'agents';

	let query = $state('');
	let filter = $state<Filter>('all');
	let directory = $state.raw<DirectoryUser[]>([]);
	let loading = $state(true);
	let busyId = $state('');

	const me = $derived(app.session?.matrix.user_id ?? '');

	const rows = $derived.by(() => {
		void app.version;
		return app.client ? buildNav(app.client, app.spaces, app.agents).dms : [];
	});

	const needle = $derived(query.trim().toLowerCase());
	const keep = (agent: boolean) => filter === 'all' || (filter === 'agents') === agent;
	const matches = (...fields: (string | null | undefined)[]) =>
		!needle || fields.some((field) => field?.toLowerCase().includes(needle));

	const conversations = $derived(rows.filter((row) => keep(row.agent) && matches(row.name, row.userId)));
	const talking = $derived(new Map(rows.filter((row) => row.userId && !row.invite).map((row) => [row.userId!, row.roomId])));
	const contacts = $derived(
		directory.filter((user) => user.user_id !== me && !isBot(user.user_id) && keep(!!user.agent || app.isAgent(user.user_id)))
	);
	const counts = $derived({
		people: directory.filter((u) => u.user_id !== me && !isBot(u.user_id) && !u.agent).length,
		agents: directory.filter((u) => u.agent).length
	});
	const tabs = $derived([
		{ id: 'all', label: 'Everyone' },
		{ id: 'people', label: 'People', badge: counts.people || undefined },
		{ id: 'agents', label: 'Agents', badge: counts.agents || undefined }
	]);

	$effect(() => {
		const q = query.trim();
		const timer = setTimeout(async () => {
			directory = await prisme.users(q).catch(() => []);
			loading = false;
		}, q ? 200 : 0);
		return () => clearTimeout(timer);
	});

	function owner(userId: string): string | null {
		const ownerId = app.agents.get(userId)?.owner_user_id;
		if (!ownerId) return null;
		if (ownerId === me) return 'Your agent';
		const name = app.client?.getUser(ownerId)?.displayName || handle(ownerId);
		return `${name}'s agent`;
	}

	function subtitle(row: DmRow): string {
		if (row.invite) return 'Wants to message you';
		return row.preview?.text ?? 'No messages yet';
	}

	async function message(user: DirectoryUser) {
		const existing = talking.get(user.user_id);
		if (existing) return goto(roomHref(existing));
		if (!app.client) return;
		busyId = user.user_id;
		try {
			await goto(roomHref(await openDm(app.client, user.user_id)));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not open the conversation');
		} finally {
			busyId = '';
		}
	}
</script>

<svelte:head>
	<title>Messages · Prisme</title>
</svelte:head>

<Page>
	<PageHeader title="Messages" description="Talk one to one with your teammates and with any agent." />

	<div class="flex flex-col gap-4">
		<Input bind:value={query} type="search" placeholder="Search people and agents" aria-label="Search people, agents and conversations" />
		<Tabs items={tabs} bind:value={filter} label="Who to show" />
	</div>

	<Section title="Conversations" card gap="tight" bodyClass="gap-0 p-2">
		{#if conversations.length > 0}
			<ul class="flex flex-col">
				{#each conversations as row (row.roomId)}
					<li>
						<a
							href={roomHref(row.roomId)}
							class="flex min-h-16 items-center gap-3 rounded-fc-md px-3 py-2 transition-colors hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-fc-ring"
						>
							<Avatar userId={row.userId ?? ''} name={row.name} size={40} presence={!row.agent} />
							<span class="flex min-w-0 flex-1 flex-col gap-0.5">
								<span class="flex min-w-0 items-center gap-2">
									<span class="truncate text-fc-sm {row.badge.bold ? 'font-semibold text-fc-fg' : 'font-medium text-fc-fg'}">{row.name}</span>
									{#if row.agent}<Badge tone="admin" class="shrink-0">Agent</Badge>{/if}
									{#if row.preview}
										<span class="ml-auto shrink-0 text-fc-xs text-fc-fg-muted">{shortTime(row.preview.at)}</span>
									{/if}
								</span>
								<span class="flex min-w-0 items-center gap-2">
									<span class="truncate text-fc-xs {row.badge.bold || row.invite ? 'text-fc-fg' : 'text-fc-fg-muted'}">{subtitle(row)}</span>
									<CountBadge badge={row.badge} label="unread messages" />
								</span>
							</span>
						</a>
					</li>
				{/each}
			</ul>
		{:else}
			<EmptyState
				bare
				icon={glyphs.chat}
				title={needle ? 'No conversation matches' : 'No conversations yet'}
				description={needle ? 'Try another name, or start one below.' : 'Pick someone below to send your first message.'}
			/>
		{/if}
	</Section>

	<Section title="Start a conversation" description="Everyone on Prisme you can message, people and agents alike.">
		{#if loading}
			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each [0, 1, 2] as slot (slot)}<Skeleton class="h-20 rounded-fc-lg" />{/each}
			</div>
		{:else if contacts.length > 0}
			<ul class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each contacts as user (user.user_id)}
					{@const agent = !!user.agent || app.isAgent(user.user_id)}
					{@const open = talking.has(user.user_id)}
					<li>
						<button
							type="button"
							disabled={busyId === user.user_id}
							onclick={() => message(user)}
							class="group flex min-h-20 w-full items-center gap-3 rounded-fc-lg bg-fc-component p-4 text-left transition-colors hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-fc-ring disabled:opacity-60"
						>
							<Avatar userId={user.user_id} name={user.name} mxc={user.avatar_url || null} size={40} presence={!agent} />
							<span class="flex min-w-0 flex-1 flex-col gap-0.5">
								<span class="flex min-w-0 items-center gap-2">
									<span class="truncate text-fc-sm font-medium text-fc-fg">{user.name}</span>
									{#if agent}<Badge tone="admin" class="shrink-0">Agent</Badge>{/if}
								</span>
								<span class="truncate text-fc-xs text-fc-fg-muted">{(agent && owner(user.user_id)) || handle(user.user_id)}</span>
							</span>
							{#if busyId === user.user_id}
								<Spinner size="sm" />
							{:else}
								<span class="flex shrink-0 items-center gap-1 text-fc-xs font-medium text-fc-fg-muted group-hover:text-fc-fg">
									{open ? 'Open' : 'Message'}
									<Icon icon={icons.arrow} size={16} />
								</span>
							{/if}
						</button>
					</li>
				{/each}
			</ul>
		{:else}
			<EmptyState
				icon={icons.usersGroup}
				title={needle ? 'Nobody matches' : 'Nobody here yet'}
				description={needle ? 'Search by name, email or @handle.' : 'People show up once they sign in to Prisme.'}
			/>
		{/if}
	</Section>
</Page>
