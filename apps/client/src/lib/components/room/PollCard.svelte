<script lang="ts">
	import { PollEvent, type MatrixEvent, type Poll } from 'matrix-js-sdk';
	import { Button, ConfirmModal, Icon, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { tally, type PollInfo, type PollVote } from '$lib/chat/poll';
	import { percent, toggleAnswer } from '$lib/chat/poll-draft';
	import { endedAt, endPoll, pollModel, pollVotes, vote } from '$lib/matrix/polls';
	import { t } from '$lib/i18n/index.svelte';

	let { event, poll }: { event: MatrixEvent; poll: PollInfo } = $props();

	let model = $state.raw<Poll | undefined>(undefined);
	let votes = $state.raw<PollVote[]>([]);
	let endTs = $state<number | undefined>(undefined);
	let pending = $state<string[] | null>(null);
	let busy = $state(false);
	let confirming = $state(false);

	const me = $derived(app.client?.getUserId() ?? '');
	const room = $derived(app.client?.getRoom(event.getRoomId() ?? '') ?? null);
	const ended = $derived(endTs !== undefined);
	const result = $derived(tally(poll, votes, me, endTs));
	const picked = $derived(pending ?? result.mine);
	const showResults = $derived(ended || (poll.disclosed && result.mine.length > 0));
	const canEnd = $derived(!ended && event.getSender() === me);

	$effect(() => {
		const r = room;
		if (!r) return;
		let live = true;
		let current: Poll | undefined;
		const refresh = async () => {
			if (!current || !live) return;
			const next = await pollVotes(current).catch(() => [] as PollVote[]);
			if (!live) return;
			votes = next;
			endTs = endedAt(r, current);
			pending = null;
		};
		void pollModel(r, event).then((p) => {
			if (!live || !p) return;
			current = p;
			model = p;
			p.on(PollEvent.Responses, refresh);
			p.on(PollEvent.Update, refresh);
			p.on(PollEvent.End, refresh);
			void refresh();
		});
		return () => {
			live = false;
			current?.off(PollEvent.Responses, refresh);
			current?.off(PollEvent.Update, refresh);
			current?.off(PollEvent.End, refresh);
		};
	});

	async function choose(id: string) {
		if (!app.client || ended || busy) return;
		const next = toggleAnswer(picked, id, poll.maxSelections);
		pending = next;
		busy = true;
		try {
			await vote(app.client, event, next);
		} catch (err) {
			pending = null;
			toast.danger(err instanceof Error ? err.message : t('poll.voteFailed'));
		} finally {
			busy = false;
		}
	}

	async function end() {
		if (!app.client) return;
		confirming = false;
		try {
			await endPoll(app.client, event, t('poll.endText'));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('poll.sendFailed'));
		}
	}

	function share(id: string): number {
		return percent(result.counts.get(id) ?? 0, result.voters);
	}
</script>

<section class="mt-1 flex w-full max-w-md flex-col gap-3 rounded-fc-lg bg-fc-component p-3" aria-label={t('poll.label')}>
	<header class="flex items-start gap-3">
		<span class="flex size-9 shrink-0 items-center justify-center rounded-fc-pill bg-fc-accent/10 text-fc-accent" aria-hidden="true">
			<Icon icon={glyphs.poll} size={18} />
		</span>
		<div class="flex min-w-0 flex-1 flex-col">
			<h3 class="text-fc-sm font-semibold text-fc-fg [overflow-wrap:anywhere]">{poll.question}</h3>
			<p class="text-fc-xs text-fc-fg-muted">
				{ended ? t('poll.ended') : poll.maxSelections > 1 ? t('poll.several', { count: poll.maxSelections }) : t('poll.single')}
			</p>
		</div>
	</header>

	<ul class="flex flex-col gap-1.5" role="group" aria-label={poll.question}>
		{#each poll.answers as answer (answer.id)}
			{@const chosen = picked.includes(answer.id)}
			{@const count = result.counts.get(answer.id) ?? 0}
			{@const top = ended && count > 0 && count === result.leader}
			<li>
				<button
					type="button"
					aria-pressed={chosen}
					disabled={ended || !model}
					onclick={() => choose(answer.id)}
					class="relative isolate flex min-h-11 w-full items-center gap-3 overflow-hidden rounded-fc-md px-3 py-2 text-left text-fc-sm transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring disabled:cursor-default {chosen
						? 'bg-fc-accent/10 text-fc-fg'
						: 'bg-fc-surface text-fc-fg hover:bg-fc-field'}"
				>
					<span
						class="absolute inset-0 -z-10 origin-left rounded-fc-md transition-transform duration-500 ease-fc will-change-transform {top || chosen ? 'bg-fc-accent/20' : 'bg-fc-fg/5'}"
						style:transform="scaleX({showResults ? share(answer.id) / 100 : 0})"
						aria-hidden="true"
					></span>
					<span
						class="flex size-5 shrink-0 items-center justify-center transition-colors duration-200 ease-fc {poll.maxSelections > 1 ? 'rounded-fc-sm' : 'rounded-fc-pill'} {chosen
							? 'bg-fc-accent text-fc-accent-fg'
							: 'bg-fc-field'}"
						aria-hidden="true"
					>
						{#if chosen}<Icon icon={glyphs.voted} size={14} />{/if}
					</span>
					<span class="min-w-0 flex-1 [overflow-wrap:anywhere] {top ? 'font-semibold' : ''}">{answer.text}</span>
					<span
						class="w-9 shrink-0 text-right text-fc-xs tabular-nums text-fc-fg-muted transition-opacity duration-300 ease-fc {showResults ? 'opacity-100' : 'opacity-0'}"
						aria-hidden={!showResults}
					>{share(answer.id)}%</span>
				</button>
			</li>
		{/each}
	</ul>

	<footer class="flex items-center gap-2">
		<p class="min-w-0 flex-1 text-fc-xs text-fc-fg-muted" aria-live="polite">
			{#if !poll.disclosed && !ended}
				{t('poll.hiddenUntilEnd')}
			{:else}
				{t('poll.votes', { count: result.voters })}
			{/if}
		</p>
		{#if canEnd}
			<Button size="sm" variant="ghost" onclick={() => (confirming = true)}>{t('poll.end')}</Button>
		{/if}
	</footer>
</section>

<ConfirmModal open={confirming} title={t('poll.end')} description={t('poll.endConfirm')} confirmLabel={t('poll.end')} tone="danger" onConfirm={end} onCancel={() => (confirming = false)} />
