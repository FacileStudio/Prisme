<script lang="ts">
	import { Track, type Participant } from 'livekit-client';
	import { Avatar } from '@facile/muse';

	let { participant, version, screen = false }: { participant: Participant; version: number; screen?: boolean } =
		$props();
	let video: HTMLVideoElement | undefined = $state();

	const view = $derived.by(() => {
		void version;
		const publication = participant.getTrackPublication(screen ? Track.Source.ScreenShare : Track.Source.Camera);
		const name = participant.name || participant.identity;
		return {
			name,
			label: screen ? (participant.isLocal ? 'Your screen' : `${name}'s screen`) : name + (participant.isLocal ? ' (you)' : ''),
			track: publication && !publication.isMuted ? (publication.videoTrack ?? null) : null,
			mic: participant.isMicrophoneEnabled,
			speaking: participant.isSpeaking,
			local: participant.isLocal,
			wide: screen && !participant.isLocal
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

<div
	class="relative aspect-video overflow-hidden rounded-fc-md {screen ? 'bg-black' : 'bg-fc-surface'} {view.wide ? 'col-span-full max-h-[45dvh] w-full' : ''} {!screen && view.speaking ? 'ring-2 ring-fc-success' : ''}"
	role={screen ? 'img' : undefined}
	aria-label={screen ? view.label : undefined}
>
	{#if view.track}
		<video bind:this={video} autoplay playsinline muted class="size-full {screen ? 'object-contain' : 'object-cover'}"></video>
	{:else}
		<div class="flex size-full items-center justify-center">
			<Avatar name={view.name} size="lg" />
		</div>
	{/if}
	<div class="absolute bottom-1.5 left-1.5 flex max-w-[calc(100%-0.75rem)] items-center gap-1.5 rounded-fc-sm bg-fc-scrim px-2 py-0.5">
		{#if !screen}
			<span class="size-2 shrink-0 rounded-fc-pill {view.mic ? 'bg-fc-success' : 'bg-fc-danger'}" aria-label={view.mic ? 'Microphone on' : 'Muted'} role="img"></span>
		{/if}
		<span class="truncate text-fc-xs text-white">{view.label}</span>
	</div>
</div>
