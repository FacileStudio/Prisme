<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { SideBar, ConfirmModal } from '@facile/muse';
	import { backend } from '$lib/backend';
	import { matrix } from '$lib/matrix.svelte';
	import { buildNavPages, findChannelRoom } from '$lib/nav';
	import type { AgentRecord, ChatMessage, SessionResponse, SpaceRecord, UserRecord } from '$lib/types';
	import ChatHeader from '$lib/components/ChatHeader.svelte';
	import MessageList from '$lib/components/MessageList.svelte';
	import MessageComposer from '$lib/components/MessageComposer.svelte';
	import ThreadPanel from '$lib/components/ThreadPanel.svelte';
	import MembersPanel from '$lib/components/MembersPanel.svelte';
	import ChatModals from '$lib/components/ChatModals.svelte';

	let collapsed = $state(false);
	let session = $state<SessionResponse | null>(null);
	let spaces = $state<SpaceRecord[]>([]);
	let agents = $state<AgentRecord[]>([]);
	let dmRooms = $state<Array<{ roomId: string; name: string; isAgent: boolean }>>([]);

	let activeSpaceId = $state<number | null>(null);
	let activeRoomId = $state<string | null>(null);
	let activeRoomName = $state('');
	let activeThread = $state<ChatMessage | null>(null);
	let showMembers = $state(false);
	let activeCallUrl = $state<string | null>(null);

	let showCreateSpace = $state(false);
	let showCreateChannel = $state(false);
	let showAddAgent = $state(false);
	let showNewDM = $state(false);
	let showLogoutModal = $state(false);

	const activeSpace = $derived(spaces.find((s) => s.id === activeSpaceId) || null);
	const currentMessages = $derived(activeRoomId ? matrix.messagesByRoom[activeRoomId] || [] : []);
	const threadReplies = $derived(
		activeThread && activeRoomId
			? (matrix.messagesByRoom[activeRoomId] || []).filter((m) => m.threadRootId === activeThread?.id)
			: []
	);
	const currentMembers = $derived(activeRoomId ? matrix.membersByRoom[activeRoomId] || [] : []);
	const typingUsers = $derived(activeRoomId ? matrix.typingByRoom[activeRoomId] || [] : []);
	const mappedSpaces = $derived(spaces.map((s) => ({ id: String(s.id), name: s.name })));
	const pages = $derived(buildNavPages(activeSpace, dmRooms, agents, activeRoomId));

	onMount(() => { void initApp(); });
	onDestroy(() => { matrix.stop(); });

	async function initApp() {
		try {
			const sess = await backend.getSession();
			session = sess;
			const [loadedSpaces, loadedAgents] = await Promise.all([backend.listSpaces(), backend.listAgents()]);
			spaces = loadedSpaces;
			agents = loadedAgents;
			matrix.setKnownAgents(loadedAgents.map((a) => a.user_id));
			await matrix.start(sess.matrix);

			const target = new URLSearchParams(window.location.search).get('room');
			const match = target ? findChannelRoom(loadedSpaces, target) : null;
			if (match) {
				activeSpaceId = match.spaceId;
				selectRoom(match.channel.room_id, '#' + match.channel.name);
				return;
			}
			if (loadedSpaces[0]?.channels[0]) {
				activeSpaceId = loadedSpaces[0].id;
				selectRoom(loadedSpaces[0].channels[0].room_id, '#' + loadedSpaces[0].channels[0].name);
			}
		} catch {
			await goto('/login');
		}
	}

	function selectRoom(roomId: string, name: string) {
		activeRoomId = roomId;
		activeRoomName = name;
		activeThread = null;
		window.history.replaceState(null, '', `?room=${roomId}`);
		void matrix.loadRoomMembers(roomId);
	}

	function handleSpaceSelect(id: string | null) {
		if (!id) return;
		activeSpaceId = Number(id);
		const sp = spaces.find((s) => s.id === Number(id));
		if (sp?.channels[0]) selectRoom(sp.channels[0].room_id, '#' + sp.channels[0].name);
	}

	function handleNavClick(e: MouseEvent) {
		const a = (e.target as HTMLElement)?.closest('a');
		if (!a) return;
		const href = a.getAttribute('href');
		if (!href) return;

		if (href === '#create-channel') { e.preventDefault(); showCreateChannel = true; return; }
		if (href === '#new-dm') { e.preventDefault(); showNewDM = true; return; }
		if (href === '#add-agent') { e.preventDefault(); showAddAgent = true; return; }
		if (href === '#create-space') { e.preventDefault(); showCreateSpace = true; return; }
		if (href === '#logout') { e.preventDefault(); showLogoutModal = true; return; }
		if (href.startsWith('?room=')) {
			e.preventDefault();
			const roomId = href.slice(6);
			const ch = activeSpace?.channels.find((c) => c.room_id === roomId);
			if (ch) { selectRoom(ch.room_id, '#' + ch.name); return; }
			const dm = dmRooms.find((d) => d.roomId === roomId);
			if (dm) selectRoom(dm.roomId, dm.name);
			return;
		}
		if (href.startsWith('?agent=')) {
			e.preventDefault();
			const ag = agents.find((a) => a.user_id === href.slice(7));
			if (ag) void handleStartDM({ user_id: ag.user_id, name: ag.name });
		}
	}

	async function handleStartCall() {
		if (!activeRoomId) return;
		try {
			const res = await backend.startCall(activeRoomId, false);
			activeCallUrl = res.join_url;
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Could not start call');
		}
	}

	async function handleStartDM(user: UserRecord) {
		try {
			const roomId = await matrix.createDirectMessage(user.user_id);
			dmRooms = [...dmRooms.filter((d) => d.roomId !== roomId), { roomId, name: user.name, isAgent: false }];
			selectRoom(roomId, user.name);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Could not start direct message');
		}
	}

	async function handleLogout() {
		try {
			await backend.logout();
		} finally {
			matrix.stop();
			await goto('/login');
		}
	}
</script>

<svelte:head>
	<title>{activeRoomName ? `${activeRoomName} — Prisme` : 'Prisme'}</title>
</svelte:head>

{#if session}
	<div class="flex h-dvh w-screen overflow-hidden bg-fc-page p-3 gap-3 font-sans text-fc-fg">
		<SideBar
			bind:collapsed
			title="Prisme"
			icon="solar:chat-round-dots-bold-duotone"
			spaces={mappedSpaces}
			activeSpaceId={activeSpaceId ? String(activeSpaceId) : null}
			onSpaceSelect={handleSpaceSelect}
			manageSpacesHref="#create-space"
			manageSpacesLabel="Create Space"
			personalSpaceLabel={null}
			{pages}
			user={{ name: session.user.name, avatar: session.user.avatar_url || undefined }}
			userHref="#logout"
			onclick={handleNavClick}
		/>

		<main class="flex min-w-0 flex-1 flex-col overflow-hidden rounded-fc-lg bg-fc-component">
			{#if activeRoomId}
				<ChatHeader
					roomName={activeRoomName}
					callsEnabled={session.calls_enabled}
					hasActiveCall={!!activeCallUrl}
					onstartcall={handleStartCall}
					ontogglemembers={() => (showMembers = !showMembers)}
				/>
				<MessageList
					messages={currentMessages}
					currentUserId={session.matrix.user_id}
					onreact={(id, em) => { if (activeRoomId) void matrix.sendReaction(activeRoomId, id, em); }}
					onreplyinthread={(msg) => (activeThread = msg)}
					onjoincall={(_, url) => (activeCallUrl = url)}
				/>
				<MessageComposer
					roomName={activeRoomName}
					{typingUsers}
					onsend={async (text) => { if (activeRoomId) await matrix.sendMessage(activeRoomId, text); }}
					ontyping={(typing) => { if (activeRoomId) void matrix.sendTyping(activeRoomId, typing); }}
				/>
			{:else}
				<div class="m-auto flex flex-col items-center justify-center p-6 text-center text-fc-fg-muted">
					<p class="text-fc-base font-medium">Select a channel or conversation to start chatting</p>
				</div>
			{/if}
		</main>

		{#if activeThread && activeRoomId}
			<ThreadPanel
				rootMessage={activeThread}
				replies={threadReplies}
				currentUserId={session.matrix.user_id}
				onclose={() => (activeThread = null)}
				onsendreply={async (t) => { if (activeRoomId && activeThread) await matrix.sendMessage(activeRoomId, t, activeThread.id); }}
			/>
		{/if}

		{#if showMembers && activeRoomId}
			<MembersPanel members={currentMembers} onclose={() => (showMembers = false)} />
		{/if}
	</div>
{/if}

<ChatModals
	{showCreateSpace} {showCreateChannel} {showAddAgent} {showNewDM}
	{activeCallUrl} {activeSpaceId} {activeRoomId}
	oncreatespace={(sp) => {
		spaces = [...spaces, sp];
		activeSpaceId = sp.id;
		if (sp.channels[0]) selectRoom(sp.channels[0].room_id, '#' + sp.channels[0].name);
	}}
	oncreatechannel={(ch) => {
		if (activeSpace) {
			activeSpace.channels = [...activeSpace.channels, ch];
			selectRoom(ch.room_id, '#' + ch.name);
		}
	}}
	oncreateagent={async () => {
		agents = await backend.listAgents();
		matrix.setKnownAgents(agents.map((a) => a.user_id));
		if (activeRoomId) void matrix.loadRoomMembers(activeRoomId);
	}}
	onstartdm={handleStartDM}
	onclose={() => {
		showCreateSpace = false; showCreateChannel = false;
		showAddAgent = false; showNewDM = false; activeCallUrl = null;
	}}
/>

<ConfirmModal
	bind:open={showLogoutModal}
	title="Log out"
	description="Are you sure you want to log out of Prisme?"
	confirmLabel="Log out"
	tone="danger"
	onConfirm={handleLogout}
/>
