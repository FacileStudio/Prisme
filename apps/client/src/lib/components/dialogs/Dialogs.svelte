<script lang="ts">
	import { page } from '$app/state';
	import { app } from '$lib/state/app.svelte';
	import { ui } from '$lib/state/ui.svelte';
	import CreateSpaceDialog from './CreateSpaceDialog.svelte';
	import CreateChannelDialog from './CreateChannelDialog.svelte';
	import InviteDialog from './InviteDialog.svelte';
	import NewDmDialog from './NewDmDialog.svelte';
	import AddAgentDialog from './AddAgentDialog.svelte';
	import RenameDialog from './RenameDialog.svelte';
	import SettingsDialog from './SettingsDialog.svelte';

	const roomId = $derived(page.params.id ? decodeURIComponent(page.params.id) : '');
	const space = $derived(app.spaces.find((s) => s.room_id === ui.spaceId) ?? null);
	const close = () => (ui.dialog = null);
</script>

{#if ui.dialog === 'create-space'}
	<CreateSpaceDialog onclose={close} />
{:else if ui.dialog === 'create-channel' && space}
	<CreateChannelDialog {space} onclose={close} />
{:else if ui.dialog === 'invite' && space}
	<InviteDialog {space} onclose={close} />
{:else if ui.dialog === 'new-dm'}
	<NewDmDialog onclose={close} />
{:else if ui.dialog === 'add-agent' && roomId}
	<AddAgentDialog {roomId} onclose={close} />
{:else if ui.dialog === 'rename-room' && roomId}
	<RenameDialog {roomId} kind="channel" onclose={close} />
{:else if ui.dialog === 'rename-space' && space}
	<RenameDialog roomId={space.room_id} kind="space" onclose={close} />
{:else if ui.dialog === 'settings'}
	<SettingsDialog onclose={close} />
{/if}
