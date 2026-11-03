<script lang="ts">
	import { Track, type Participant } from 'livekit-client';
	import { Avatar } from '@facile/muse';

	let { participant, version }: { participant: Participant; version: number } = $props();
	let video: HTMLVideoElement | undefined = $state();

	const view = $derived.by(() => {
		void version;
		const camera = participant.getTrackPublication(Track.Source.Camera);
		return {
			name: participant.name || participant.identity,
			track: camera && !camera.isMuted ? (camera.videoTrack ?? null) : null,
			mic: participant.isMicrophoneEnabled,
			speaking: participant.isSpeaking,
			local: participant.isLocal
		};
	});

	$effect(() => {
		const track = view.track;
		const element = video;
		if (!track || !element) return;
		track.attach(element);
		return () => void track.detach(element);
	});
</script>

<div class="relative aspect-video overflow-hidden rounded-fc-md bg-fc-surface {view.speaking ? 'ring-2 ring-fc-success' : ''}">
	{#if view.track}
		<video bind:this={video} autoplay playsinline muted class="size-full object-cover"></video>
	{:else}
		<div class="flex size-full items-center justify-center">
			<Avatar name={view.name} size="lg" />
		</div>
	{/if}
	<div class="absolute bottom-1.5 left-1.5 flex max-w-[calc(100%-0.75rem)] items-center gap-1.5 rounded-fc-sm bg-fc-scrim px-2 py-0.5">
		<span class="size-2 shrink-0 rounded-fc-pill {view.mic ? 'bg-fc-success' : 'bg-fc-danger'}" aria-label={view.mic ? 'Microphone on' : 'Muted'} role="img"></span>
		<span class="truncate text-fc-xs text-white">{view.name}{view.local ? ' (you)' : ''}</span>
	</div>
</div>
