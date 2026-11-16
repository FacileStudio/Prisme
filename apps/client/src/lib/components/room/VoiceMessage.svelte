<script lang="ts">
	import { cubicOut } from 'svelte/easing';
	import { Icon, Spinner, prefersReducedMotion } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { bars, clock, type VoiceInfo } from '$lib/chat/voice';
	import { loadMedia } from '$lib/matrix/media';
	import { t } from '$lib/i18n/index.svelte';

	let { voice, mine = false }: { voice: VoiceInfo; mine?: boolean } = $props();

	const SPEEDS = [1, 1.5, 2];

	let audio: HTMLAudioElement | undefined = $state();
	let src = $state('');
	let loading = $state(false);
	let failed = $state(false);
	let playing = $state(false);
	let position = $state(0);
	let length = $state(0);
	let speed = $state(1);

	const heights = $derived(bars(voice.waveform));
	const total = $derived(length > 0 && Number.isFinite(length) ? length * 1000 : voice.duration);
	const progress = $derived(total > 0 ? Math.min(1, (position * 1000) / total) : 0);

	function swap(_node: Element) {
		const duration = prefersReducedMotion() ? 0 : 160;
		return { duration, easing: cubicOut, css: (u: number) => `opacity: ${u}; transform: scale(${0.7 + 0.3 * u})` };
	}

	async function ensure(): Promise<boolean> {
		if (src) return true;
		loading = true;
		try {
			src = await loadMedia(voice.url);
			await new Promise((resolve) => requestAnimationFrame(resolve));
			return true;
		} catch {
			failed = true;
			return false;
		} finally {
			loading = false;
		}
	}

	async function toggle() {
		if (!(await ensure()) || !audio) return;
		if (audio.paused) {
			document.querySelectorAll<HTMLAudioElement>('audio[data-voice]').forEach((other) => other !== audio && other.pause());
			audio.playbackRate = speed;
			await audio.play().catch(() => (failed = true));
		} else {
			audio.pause();
		}
	}

	async function seek(event: Event) {
		const value = Number((event.currentTarget as HTMLInputElement).value);
		if (!(await ensure()) || !audio) return;
		audio.currentTime = (value / 1000) * (total / 1000);
		position = audio.currentTime;
	}

	$effect(() => {
		if (!playing || !audio) return;
		const el = audio;
		let frame = requestAnimationFrame(function step() {
			position = el.currentTime;
			frame = requestAnimationFrame(step);
		});
		return () => cancelAnimationFrame(frame);
	});

	function cycleSpeed() {
		speed = SPEEDS[(SPEEDS.indexOf(speed) + 1) % SPEEDS.length];
		if (audio) audio.playbackRate = speed;
	}
</script>

<div class="mt-1 flex w-full max-w-sm items-center gap-3 rounded-fc-pill bg-fc-component py-1.5 pr-3 pl-1.5">
	<button
		type="button"
		onclick={toggle}
		disabled={failed}
		aria-label={playing ? t('voice.pause') : t('voice.play')}
		class="flex size-11 shrink-0 items-center justify-center rounded-fc-pill transition-transform focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring active:scale-95 disabled:opacity-40 {mine
			? 'bg-fc-accent text-fc-accent-fg'
			: 'bg-fc-fg text-fc-bg'}"
	>
		{#if loading}
			<Spinner size="sm" />
		{:else}
			{#key playing}<span class="flex" in:swap><Icon icon={playing ? glyphs.pause : glyphs.play} size={16} /></span>{/key}
		{/if}
	</button>

	<div class="relative flex h-8 min-w-0 flex-1 items-center">
		<div class="pointer-events-none relative h-full w-full" aria-hidden="true">
			<div class="flex h-full w-full items-center gap-[2px]">
				{#each heights as h, i (i)}
					<span class="min-w-[2px] flex-1 rounded-fc-pill bg-fc-fg/20" style:height="{Math.round(h * 100)}%"></span>
				{/each}
			</div>
			<div class="absolute inset-0 flex items-center gap-[2px] will-change-[clip-path]" style:clip-path="inset(0 {(100 - progress * 100).toFixed(2)}% 0 0)">
				{#each heights as h, i (i)}
					<span class="min-w-[2px] flex-1 rounded-fc-pill bg-fc-accent" style:height="{Math.round(h * 100)}%"></span>
				{/each}
			</div>
		</div>
		<input
			type="range"
			min="0"
			max="1000"
			step="1"
			value={Math.round(progress * 1000)}
			oninput={seek}
			disabled={failed}
			aria-label={t('voice.seek')}
			aria-valuetext={`${clock(position * 1000)} / ${clock(total)}`}
			class="absolute inset-0 h-full w-full cursor-pointer appearance-none opacity-0 focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-fc-ring"
		/>
	</div>

	<span class="w-10 shrink-0 text-right text-fc-xs tabular-nums text-fc-fg-muted">
		{clock(playing || position > 0 ? position * 1000 : total)}
	</span>
	<button
		type="button"
		onclick={cycleSpeed}
		aria-label={t('voice.speed')}
		class="group -my-1.5 -mr-2 flex size-11 shrink-0 items-center justify-center rounded-fc-pill focus-visible:outline-2 focus-visible:outline-fc-ring"
	>
		<span class="flex h-7 min-w-9 items-center justify-center rounded-fc-pill bg-fc-surface px-2 text-fc-xs font-semibold tabular-nums text-fc-fg-muted group-hover:text-fc-fg">{speed}×</span>
	</button>

	{#if src}
		<audio
			bind:this={audio}
			data-voice
			{src}
			preload="metadata"
			ontimeupdate={() => (position = audio?.currentTime ?? 0)}
			onloadedmetadata={() => (length = audio?.duration ?? 0)}
			onplay={() => (playing = true)}
			onpause={() => (playing = false)}
			onended={() => ((playing = false), (position = 0))}
			onerror={() => (failed = true)}
		></audio>
	{/if}
</div>
{#if failed}
	<p class="mt-1 text-fc-xs text-fc-danger">{t('voice.failed')}</p>
{/if}
