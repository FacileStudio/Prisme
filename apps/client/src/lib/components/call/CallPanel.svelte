<script lang="ts">
	import { Alert, Button, Icon, IconButton, Spinner, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { call } from '$lib/call/call.svelte';
	import ParticipantTile from './ParticipantTile.svelte';

	const roomName = $derived(app.client?.getRoom(call.roomId)?.name ?? 'Call');
	const tiles = $derived.by(() => {
		void call.version;
		return call.participants;
	});
</script>

<section class="flex max-h-[60dvh] flex-col gap-3 border-b border-fc-border bg-fc-component p-3" aria-label="Call in {roomName}">
	<div class="flex items-center gap-2">
		<Icon icon={glyphs.phone} size={18} class="text-fc-success" />
		<h2 class="min-w-0 flex-1 truncate text-fc-sm font-semibold">Call in {roomName}</h2>
		{#if call.status === 'connecting'}<Spinner size="sm" label="Connecting" />{/if}
	</div>
	{#if call.status === 'disconnected'}
		<Alert tone="warning" title="You left the call">
			<div class="flex flex-wrap items-center gap-2">
				<span>{call.error || 'The connection dropped.'}</span>
				<Button size="sm" icon={icons.refresh} onclick={() => call.rejoin().catch(() => undefined)}>Rejoin</Button>
			</div>
		</Alert>
	{:else if call.error}
		<p class="text-fc-xs text-fc-warning" role="status">{call.error}</p>
	{/if}
	{#if !call.canPlayAudio}
		<Button variant="outline" icon={glyphs.volume} onclick={() => call.enableAudio()}>Enable sound</Button>
	{/if}
	<div class="grid min-h-0 grid-cols-2 gap-2 overflow-y-auto sm:grid-cols-3 lg:grid-cols-4">
		{#each tiles as participant (participant.identity)}
			<ParticipantTile {participant} version={call.version} />
		{/each}
	</div>
	<div class="flex items-center justify-center gap-2">
		<IconButton variant={call.micOn ? 'default' : 'danger'} aria-label={call.micOn ? 'Mute the microphone' : 'Unmute the microphone'} aria-pressed={!call.micOn} onclick={() => call.toggle('mic')}>
			<Icon icon={glyphs.mic} size={20} />
		</IconButton>
		<IconButton variant={call.camOn ? 'default' : 'ghost'} aria-label={call.camOn ? 'Turn the camera off' : 'Turn the camera on'} aria-pressed={call.camOn} onclick={() => call.toggle('cam')}>
			<Icon icon={glyphs.camera} size={20} />
		</IconButton>
		<Button variant="danger" size="lg" icon={glyphs.endCall} onclick={() => call.leave()}>Leave</Button>
	</div>
</section>
