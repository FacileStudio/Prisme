<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { cubicOut } from 'svelte/easing';
	import { Icon, Spinner, Tooltip, prefersReducedMotion, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { Recorder } from '$lib/state/recorder.svelte';
	import { clock } from '$lib/chat/voice';
	import { sendVoice } from '$lib/matrix/voice';
	import { t } from '$lib/i18n/index.svelte';

	type Props = { room: Room; threadId?: string | null; active?: boolean; disabled?: boolean };
	let { room, threadId = null, active = $bindable(false), disabled = false }: Props = $props();

	const recorder = new Recorder();
	const supported = Recorder.supported();
	let sending = $state(false);
	const ROUND =
		'inline-flex size-11 shrink-0 items-center justify-center rounded-fc-pill focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring disabled:cursor-not-allowed disabled:opacity-40';

	$effect(() => {
		active = recorder.phase !== 'idle' || sending;
	});

	function morph(node: Element, { from = 0.6, origin = 'right center' }: { from?: number; origin?: string } = {}) {
		const duration = prefersReducedMotion() ? 0 : 240;
		(node as HTMLElement).style.transformOrigin = origin;
		return {
			duration,
			easing: cubicOut,
			css: (u: number) => `opacity: ${u}; transform: scale(${from + (1 - from) * u}, 1)`
		};
	}

	function pop(node: Element) {
		const duration = prefersReducedMotion() ? 0 : 180;
		return { duration, easing: cubicOut, css: (u: number) => `opacity: ${u}; transform: scale(${0.6 + 0.4 * u})` };
	}

	$effect(() => () => recorder.cancel());

	async function start() {
		try {
			await recorder.start(() => void finish());
		} catch {
			toast.danger(t('voice.denied'));
		}
	}

	async function finish() {
		if (sending) return;
		const recording = await recorder.stop();
		if (!recording) {
			toast.info(t('voice.tooShort'));
			return;
		}
		if (!app.client) return;
		sending = true;
		try {
			const last = threadId ? room.getThread(threadId)?.lastReply()?.getId() : undefined;
			await sendVoice(app.client, { roomId: room.roomId, threadId }, recording, t('voice.body'), last);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('voice.sendFailed'));
		} finally {
			sending = false;
		}
	}

	function onKey(event: KeyboardEvent) {
		if (event.key === 'Escape' && recorder.phase === 'recording' && !sending) {
			event.preventDefault();
			recorder.cancel();
		}
	}
</script>

<svelte:window onkeydown={onKey} />

{#if active}
	<div class="flex min-w-0 flex-1 items-center gap-2 will-change-transform" role="group" aria-label={t('voice.recording')} in:morph={{ from: 0.4 }}>
		<Tooltip text={t('voice.cancel')} describe={false} class="contents">
			<button type="button" aria-label={t('voice.cancel')} disabled={sending} class="{ROUND} text-fc-fg-muted hover:bg-fc-surface hover:text-fc-danger" onclick={() => recorder.cancel()}>
				<Icon icon={glyphs.trash} size={20} />
			</button>
		</Tooltip>
		<div class="flex h-11 min-w-0 flex-1 items-center gap-3 rounded-fc-pill bg-fc-surface px-4">
			<span class="relative flex size-2.5 shrink-0" aria-hidden="true">
				<span class="absolute inline-flex size-full animate-ping rounded-fc-pill bg-fc-danger opacity-60 motion-reduce:animate-none"></span>
				<span class="relative inline-flex size-2.5 rounded-fc-pill bg-fc-danger"></span>
			</span>
			<span class="w-10 shrink-0 text-fc-sm font-medium tabular-nums text-fc-fg" aria-live="off">{clock(recorder.elapsed)}</span>
			<span class="flex h-6 min-w-0 flex-1 items-center justify-end gap-[2px] overflow-hidden" aria-hidden="true">
				{#each recorder.levels as level, i (i)}
					<span class="h-full w-[3px] shrink-0 origin-center rounded-fc-pill bg-fc-fg/50 will-change-transform" style:transform="scaleY({Math.max(0.12, level).toFixed(3)})"></span>
				{/each}
			</span>
		</div>
		<Tooltip text={t('voice.send')} describe={false} class="contents">
			<button type="button" aria-label={t('voice.send')} disabled={sending || recorder.phase !== 'recording'} class="{ROUND} bg-fc-accent text-fc-accent-fg hover:opacity-90" onclick={finish}>
				{#if sending}
					<Spinner size="sm" />
				{:else}
					<Icon icon={glyphs.send} size={18} />
				{/if}
			</button>
		</Tooltip>
	</div>
{:else if supported}
	<Tooltip text={t('voice.record')} describe={false} class="contents">
		<button type="button" aria-label={t('voice.record')} {disabled} class="{ROUND} text-fc-fg-muted transition-colors duration-200 ease-fc hover:bg-fc-surface hover:text-fc-fg" onclick={start} in:pop>
			<Icon icon={glyphs.mic} size={20} />
		</button>
	</Tooltip>
{/if}
