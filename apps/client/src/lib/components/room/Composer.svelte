<script lang="ts">
	import { tick, untrack } from 'svelte';
	import type { Room } from 'matrix-js-sdk';
	import { Icon, IconButton, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { ComposerState, mentionables } from '$lib/state/composer.svelte';
	import type { DraftState } from '$lib/state/draft.svelte';
	import MentionMenu from './MentionMenu.svelte';
	import PendingFiles from './PendingFiles.svelte';
	import DraftBanner from './DraftBanner.svelte';

	type Props = { room: Room; draft: DraftState; threadId: string | null; placeholder?: string };
	let { room, draft, threadId, placeholder }: Props = $props();

	const composer = untrack(() => new ComposerState(app.client!, room, threadId));
	let input: HTMLTextAreaElement | undefined = $state();
	let picker: HTMLInputElement | undefined = $state();

	const members = $derived.by(() => {
		void app.version;
		return mentionables(room, app.client?.getUserId() ?? '');
	});
	const options = $derived(composer.completions(members));
	const hint = $derived(placeholder ?? `Message ${room.name}`);

	$effect(() => {
		void room.loadMembersIfNeeded();
	});

	$effect(() => {
		if (draft.editing) composer.loadEdit(draft);
		if (draft.editing || draft.replyTo) void tick().then(() => input?.focus());
	});

	async function pick(index: number) {
		const caret = composer.choose(options[index]);
		await tick();
		input?.setSelectionRange(caret, caret);
		input?.focus();
	}

	async function submit() {
		if (composer.sending || (!composer.text.trim() && composer.files.length === 0)) return;
		composer.sending = true;
		try {
			await composer.submit(draft);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not send the message');
		} finally {
			composer.sending = false;
		}
	}

	function onKeydown(event: KeyboardEvent) {
		if (options.length > 0 && ['ArrowDown', 'ArrowUp', 'Enter', 'Tab'].includes(event.key)) {
			event.preventDefault();
			const step = event.key === 'ArrowDown' ? 1 : event.key === 'ArrowUp' ? -1 : 0;
			if (step) composer.selected = (composer.selected + step + options.length) % options.length;
			else void pick(composer.selected);
		} else if (event.key === 'Enter' && !event.shiftKey && !event.isComposing) {
			event.preventDefault();
			void submit();
		} else if (event.key === 'Escape') draft.clear();
	}

	function onInput() {
		composer.caret = input?.selectionStart ?? composer.text.length;
		composer.selected = 0;
		composer.typed();
	}
</script>

<div class="relative px-3 pb-3 md:px-4">
	{#if options.length > 0}
		<MentionMenu {options} selected={composer.selected} onpick={pick} />
	{/if}
	<div class="flex flex-col rounded-fc-lg border border-fc-border bg-fc-bg focus-within:border-fc-ring">
		<DraftBanner {room} {draft} />
		<PendingFiles files={composer.files} onremove={(f) => (composer.files = composer.files.filter((x) => x !== f))} />
		<div class="flex items-end gap-1 p-1">
			<IconButton variant="ghost" aria-label="Attach a file" title="Attach a file" disabled={!!draft.editing} onclick={() => picker?.click()}>
				<Icon icon={glyphs.paperclip} size={18} />
			</IconButton>
			<textarea
				bind:this={input}
				bind:value={composer.text}
				oninput={onInput}
				onkeydown={onKeydown}
				onclick={onInput}
				rows="1"
				placeholder={hint}
				aria-label={hint}
				class="field-sizing-content max-h-48 min-h-11 flex-1 resize-none bg-transparent px-1 py-2.5 text-fc-md outline-none placeholder:text-fc-fg-muted md:text-fc-sm"
			></textarea>
			<IconButton variant="ghost" aria-label={draft.editing ? 'Save the edit' : 'Send'} title="Send" disabled={composer.sending} onclick={submit}>
				<Icon icon={glyphs.send} size={18} />
			</IconButton>
		</div>
	</div>
	<input bind:this={picker} type="file" multiple hidden onchange={(e) => {
		composer.files = [...composer.files, ...(e.currentTarget.files ?? [])];
		e.currentTarget.value = '';
	}} />
</div>
