<script lang="ts">
	import { Alert, Button, Icon, IconButton, Spinner, Tooltip, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { call } from '$lib/call/call.svelte';
	import ParticipantTile from './ParticipantTile.svelte';

	const roomName = $derived(app.client?.getRoom(call.roomId)?.name ?? t('call.defaultRoom'));
	const screens = $derived.by(() => {
		void call.version;
		return call.participants.filter((participant) => participant.isScreenShareEnabled);
	});
</script>

<section class="flex max-h-[60dvh] flex-col gap-3 bg-fc-component p-3" aria-label={t('call.panelTitle', { room: roomName })}>
	<div class="flex items-center gap-2">
		<Icon icon={glyphs.phone} size={18} class="text-fc-success" />
		<h2 class="min-w-0 flex-1 truncate text-fc-sm font-semibold">{t('call.panelTitle', { room: roomName })}</h2>
		{#if call.encrypted}
			<Tooltip text={t('call.e2eeNote')} class="contents">
				<span class="flex shrink-0 items-center gap-1 text-fc-xs text-fc-fg-muted">
					<Icon icon={glyphs.lock} size={14} />
					{t('call.e2ee')}
				</span>
			</Tooltip>
		{/if}
		{#if call.status === 'connecting'}<Spinner size="sm" label={t('call.connecting')} />{/if}
	</div>
	{#if call.status === 'disconnected'}
		<Alert tone="warning" title={t('call.left')}>
			<div class="flex flex-wrap items-center gap-2">
				<span>{call.error || t('call.dropped')}</span>
				<Button size="sm" icon={icons.refresh} onclick={() => call.rejoin().catch(() => undefined)}>{t('call.rejoin')}</Button>
			</div>
		</Alert>
	{:else if call.error}
		<p class="text-fc-xs text-fc-warning" role="status">{call.error}</p>
	{/if}
	{#if !call.canPlayAudio}
		<Button variant="secondary" icon={glyphs.volume} onclick={() => call.enableAudio()}>{t('call.enableSound')}</Button>
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
		<IconButton variant={call.micOn ? 'secondary' : 'ghost-danger'} aria-label={call.micOn ? t('call.micMute') : t('call.micUnmute')} aria-pressed={!call.micOn} onclick={() => call.toggle('mic')}>
			<Icon icon={glyphs.mic} size={20} />
		</IconButton>
		<IconButton variant={call.camOn ? 'secondary' : 'ghost'} aria-label={call.camOn ? t('call.camOff') : t('call.camOn')} aria-pressed={call.camOn} onclick={() => call.toggle('cam')}>
			<Icon icon={glyphs.camera} size={20} />
		</IconButton>
		{#if call.canShare}
			<IconButton variant={call.screenOn ? 'secondary' : 'ghost'} aria-label={call.screenOn ? t('call.shareStop') : t('call.share')} aria-pressed={call.screenOn} onclick={() => call.toggle('screen')}>
				<Icon icon={glyphs.screenShare} size={20} />
			</IconButton>
		{/if}
		<Button variant="danger" size="lg" icon={glyphs.endCall} onclick={() => call.leave()}>{t('call.leave')}</Button>
	</div>
</section>
