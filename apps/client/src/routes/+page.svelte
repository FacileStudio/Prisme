<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { backend } from '$lib/backend';
	import { matrix } from '$lib/matrix.svelte';
	import type { AgentRecord, ChatMessage, SessionResponse, SpaceRecord, UserRecord } from '$lib/types';
	import SpacesSidebar from '$lib/components/SpacesSidebar.svelte';
	import ChannelsSidebar from '$lib/components/ChannelsSidebar.svelte';
	import ChatHeader from '$lib/components/ChatHeader.svelte';
	import MessageList from '$lib/components/MessageList.svelte';
	import MessageComposer from '$lib/components/MessageComposer.svelte';
	import ThreadPanel from '$lib/components/ThreadPanel.svelte';
	import MembersPanel from '$lib/components/MembersPanel.svelte';
	import ChatModals from '$lib/components/ChatModals.svelte';

	let session = $state<SessionResponse | null>(null);
	let spaces = $state<SpaceRecord[]>([]);
	let agents = $state<AgentRecord[]>([]);
	let dmRooms = $state<Array<{ roomId: string; name: string; isAgent: boolean }>>([]);

	let activeSpaceId = $state<number | null>(null);
	let activeRoomId = $state<string | null>(null);
	let activeRoomName = $state('');
	let activeTab = $state<'channels' | 'dms'>('channels');

	let activeThread = $state<ChatMessage | null>(null);
	let showMembers = $state(false);
	let activeCallUrl = $state<string | null>(null);

	let showCreateSpace = $state(false);
	let showCreateChannel = $state(false);
	let showAddAgent = $state(false);
	let showNewDM = $state(false);

	const activeSpace = $derived(spaces.find((s) => s.id === activeSpaceId) || null);
	const currentMessages = $derived(activeRoomId ? matrix.messagesByRoom[activeRoomId] || [] : []);
	const threadReplies = $derived(
		activeThread && activeRoomId
			? (matrix.messagesByRoom[activeRoomId] || []).filter((m) => m.threadRootId === activeThread?.id)
			: []
	);
	const currentMembers = $derived(activeRoomId ? matrix.membersByRoom[activeRoomId] || [] : []);
	const typingUsers = $derived(activeRoomId ? matrix.typingByRoom[activeRoomId] || [] : []);

	onMount(() => {
		void initApp();
	});

	onDestroy(() => {
		matrix.stop();
	});

	async function initApp() {
		try {
			const sess = await backend.getSession();
			session = sess;
			const [loadedSpaces, loadedAgents] = await Promise.all([
				backend.listSpaces(),
				backend.listAgents()
			]);
			spaces = loadedSpaces;
			agents = loadedAgents;
			matrix.setKnownAgents(loadedAgents.map((a) => a.user_id));
			await matrix.start(sess.matrix);

			if (loadedSpaces[0]) {
				activeSpaceId = loadedSpaces[0].id;
				if (loadedSpaces[0].channels[0]) {
					selectRoom(loadedSpaces[0].channels[0].room_id, '#' + loadedSpaces[0].channels[0].name);
				}
			}
		} catch {
			await goto('/login');
		}
	}

	function selectRoom(roomId: string, name: string) {
		activeRoomId = roomId;
		activeRoomName = name;
		activeThread = null;
		void matrix.loadRoomMembers(roomId);
	}

	async function handleSendMessage(text: string) {
		if (!activeRoomId) return;
		await matrix.sendMessage(activeRoomId, text);
	}

	async function handleReaction(messageId: string, emoji: string) {
		if (!activeRoomId) return;
		await matrix.sendReaction(activeRoomId, messageId, emoji);
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
			activeTab = 'dms';
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
	<div class="flex h-dvh w-screen overflow-hidden bg-fc-page font-sans text-fc-fg">
		<SpacesSidebar
			{spaces}
			{activeSpaceId}
			{activeTab}
			user={session.user}
			onselectspace={(id) => {
				activeSpaceId = id;
				const sp = spaces.find((s) => s.id === id);
				if (sp?.channels[0]) selectRoom(sp.channels[0].room_id, '#' + sp.channels[0].name);
			}}
			onselecttab={(tab) => {
				activeTab = tab;
				if (tab === 'channels' && activeSpace?.channels[0]) {
					selectRoom(activeSpace.channels[0].room_id, '#' + activeSpace.channels[0].name);
				}
			}}
			oncreatespace={() => (showCreateSpace = true)}
			onlogout={handleLogout}
		/>

		<ChannelsSidebar
			tab={activeTab}
			{activeSpace}
			{activeRoomId}
			{agents}
			{dmRooms}
			onselectroom={selectRoom}
			oncreatechannel={() => (showCreateChannel = true)}
			onaddagent={() => (showAddAgent = true)}
			onnewdm={() => (showNewDM = true)}
		/>

		<main class="flex min-w-0 flex-1 flex-col overflow-hidden bg-fc-page">
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
					onreact={handleReaction}
					onreplyinthread={(msg) => (activeThread = msg)}
					onjoincall={(_, url) => (activeCallUrl = url)}
				/>

				<MessageComposer
					roomName={activeRoomName}
					{typingUsers}
					onsend={handleSendMessage}
					ontyping={(isTyping) => {
						if (activeRoomId) void matrix.sendTyping(activeRoomId, isTyping);
					}}
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
				onsendreply={async (text) => {
					if (activeRoomId && activeThread) await matrix.sendMessage(activeRoomId, text, activeThread.id);
				}}
			/>
		{/if}

		{#if showMembers && activeRoomId}
			<MembersPanel
				members={currentMembers}
				onclose={() => (showMembers = false)}
			/>
		{/if}
	</div>
{/if}

<ChatModals
	{showCreateSpace}
	{showCreateChannel}
	{showAddAgent}
	{showNewDM}
	{activeCallUrl}
	{activeSpaceId}
	{activeRoomId}
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
		showCreateSpace = false;
		showCreateChannel = false;
		showAddAgent = false;
		showNewDM = false;
		activeCallUrl = null;
	}}
/>
