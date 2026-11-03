<script lang="ts">
	import { Alert, Button, Icon, IconButton, Spinner, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { call } from '$lib/call/call.svelte';
	import ParticipantTile from './ParticipantTile.svelte';

	const roomName = $derived(app.client?.getRoom(call.roomId)?.name ?? 'Call');
	const screens = $derived.by(() => {
		void call.version;
		return call.participants.filter((participant) => participant.isScreenShareEnabled);
	});
</script>

<section class="flex max-h-[60dvh] flex-col gap-3 border-b border-fc-border bg-fc-component p-3" aria-label="Call in {roomName}">
	<div class="flex items-center gap-2">
		<Icon icon={glyphs.phone} size={18} class="text-fc-success" />
		<h2 class="min-w-0 flex-1 truncate text-fc-sm font-semibold">Call in {roomName}</h2>
		{#if call.encrypted}
			<span class="flex shrink-0 items-center gap-1 text-fc-xs text-fc-fg-muted" title="Only the people in this room hold the key. Echo and LiveKit cannot see or hear the call.">
				<Icon icon={glyphs.lock} size={14} />
				End-to-end encrypted
			</span>
		{/if}
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
	<div class="grid min-h-0 auto-rows-max grid-cols-2 gap-2 overflow-y-auto sm:grid-cols-3 lg:grid-cols-4">
		{#each screens as participant (participant.identity)}
			<ParticipantTile {participant} version={call.version} screen />
		{/each}
		{#each call.participants as participant (participant.identity)}
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
		{#if call.canShare}
			<IconButton variant={call.screenOn ? 'default' : 'ghost'} aria-label={call.screenOn ? 'Stop sharing the screen' : 'Share the screen'} title={call.screenOn ? 'Stop sharing the screen' : 'Share the screen'} aria-pressed={call.screenOn} onclick={() => call.toggle('screen')}>
				<Icon icon={glyphs.screenShare} size={20} />
			</IconButton>
		{/if}
		<Button variant="danger" size="lg" icon={glyphs.endCall} onclick={() => call.leave()}>Leave</Button>
	</div>
</section>
