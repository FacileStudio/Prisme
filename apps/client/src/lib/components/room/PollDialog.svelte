<script lang="ts">
	import { tick } from 'svelte';
	import type { Room } from 'matrix-js-sdk';
	import { Field, Icon, IconButton, Input, Switch, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { MAX_ANSWERS } from '$lib/chat/poll';
	import { draftReady, MAX_ANSWER, MAX_QUESTION, MIN_ANSWERS } from '$lib/chat/poll-draft';
	import { sendPoll } from '$lib/matrix/polls';
	import { t } from '$lib/i18n/index.svelte';
	import FormDialog from '$lib/components/dialogs/FormDialog.svelte';

	let { room, threadId = null, onclose }: { room: Room; threadId?: string | null; onclose: () => void } = $props();

	let question = $state('');
	let answers = $state(['', '']);
	let multiple = $state(false);
	let hidden = $state(false);
	let busy = $state(false);
	let list: HTMLOListElement | undefined = $state();

	const draft = $derived({ question, answers, multiple, disclosed: !hidden });
	const filled = $derived(answers.map((a) => a.trim()).filter(Boolean));
	const duplicate = $derived(new Set(filled).size !== filled.length);

	async function add() {
		if (answers.length >= MAX_ANSWERS) return;
		answers = [...answers, ''];
		await tick();
		list?.querySelector<HTMLInputElement>('li:last-child input')?.focus();
	}

	function remove(index: number) {
		answers = answers.filter((_, i) => i !== index);
	}

	function onKey(event: KeyboardEvent, index: number) {
		if (event.key !== 'Enter' || index !== answers.length - 1 || !answers[index].trim()) return;
		event.preventDefault();
		void add();
	}

	async function submit() {
		if (!app.client || !draftReady(draft)) return;
		busy = true;
		try {
			const last = threadId ? room.getThread(threadId)?.lastReply()?.getId() : undefined;
			await sendPoll(app.client, { roomId: room.roomId, threadId }, draft, last);
			onclose();
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('poll.sendFailed'));
		} finally {
			busy = false;
		}
	}
</script>

<FormDialog title={t('poll.title')} submitLabel={t('poll.create')} icon={glyphs.poll} {busy} disabled={!draftReady(draft)} onsubmit={submit} {onclose}>
	<Field label={t('poll.question')}>
		{#snippet children({ id })}
			<Input {id} bind:value={question} maxlength={MAX_QUESTION} placeholder={t('poll.questionPlaceholder')} required />
		{/snippet}
	</Field>

	<fieldset class="flex flex-col gap-2">
		<legend class="mb-2 text-fc-sm font-medium text-fc-fg">{t('poll.options')}</legend>
		<ol bind:this={list} class="flex flex-col gap-2">
			{#each answers as _, i (i)}
				<li class="flex items-center gap-2">
					<span class="flex size-7 shrink-0 items-center justify-center rounded-fc-pill bg-fc-surface text-fc-xs font-semibold text-fc-fg-muted" aria-hidden="true">{i + 1}</span>
					<Input
						bind:value={answers[i]}
						maxlength={MAX_ANSWER}
						placeholder={t('poll.option', { n: i + 1 })}
						aria-label={t('poll.option', { n: i + 1 })}
						onkeydown={(e: KeyboardEvent) => onKey(e, i)}
						class="flex-1"
					/>
					{#if answers.length > MIN_ANSWERS}
						<IconButton variant="ghost" aria-label={t('poll.removeOption', { n: i + 1 })} tooltip={t('poll.removeOption', { n: i + 1 })} onclick={() => remove(i)}>
							<Icon icon={icons.close} size={16} />
						</IconButton>
					{/if}
				</li>
			{/each}
		</ol>
		{#if answers.length < MAX_ANSWERS}
			<button
				type="button"
				class="flex min-h-11 items-center gap-2 self-start rounded-fc-pill px-3 text-fc-sm font-medium text-fc-accent hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-fc-ring"
				onclick={add}
			>
				<Icon icon={icons.plus} size={16} />{t('poll.addOption')}
			</button>
		{/if}
		{#if duplicate}
			<p class="text-fc-xs text-fc-danger" role="alert">{t('poll.duplicate')}</p>
		{/if}
	</fieldset>

	<div class="flex flex-col gap-3 rounded-fc-lg bg-fc-surface p-3">
		<Switch bind:checked={multiple} label={t('poll.multiple')} />
		<Switch bind:checked={hidden} label={t('poll.hidden')} />
	</div>
</FormDialog>
