<script lang="ts">
	import type { ChannelRecord, SpaceRecord, UserRecord } from '$lib/types';
	import CreateSpaceModal from './CreateSpaceModal.svelte';
	import CreateChannelModal from './CreateChannelModal.svelte';
	import AddAgentModal from './AddAgentModal.svelte';
	import NewDMModal from './NewDMModal.svelte';
	import CallModal from './CallModal.svelte';

	let {
		showCreateSpace,
		showCreateChannel,
		showAddAgent,
		showNewDM,
		activeCallUrl,
		activeSpaceId,
		activeRoomId,
		oncreatespace,
		oncreatechannel,
		oncreateagent,
		onstartdm,
		onclose
	}: {
		showCreateSpace: boolean;
		showCreateChannel: boolean;
		showAddAgent: boolean;
		showNewDM: boolean;
		activeCallUrl: string | null;
		activeSpaceId: number | null;
		activeRoomId: string | null;
		oncreatespace: (space: SpaceRecord) => void;
		oncreatechannel: (channel: ChannelRecord) => void;
		oncreateagent: () => void;
		onstartdm: (user: UserRecord) => void;
		onclose: () => void;
	} = $props();
</script>

{#if showCreateSpace}
	<CreateSpaceModal
		onclose={onclose}
		oncreated={oncreatespace}
	/>
{/if}

{#if showCreateChannel && activeSpaceId}
	<CreateChannelModal
		spaceId={activeSpaceId}
		onclose={onclose}
		oncreated={oncreatechannel}
	/>
{/if}

{#if showAddAgent}
	<AddAgentModal
		currentRoomId={activeRoomId}
		onclose={onclose}
		oncreated={oncreateagent}
	/>
{/if}

{#if showNewDM}
	<NewDMModal
		onclose={onclose}
		onstartdm={onstartdm}
	/>
{/if}

{#if activeCallUrl}
	<CallModal
		joinUrl={activeCallUrl}
		onclose={onclose}
	/>
{/if}
