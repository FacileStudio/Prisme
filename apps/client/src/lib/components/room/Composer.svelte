<script lang="ts">
	import { tick, untrack } from 'svelte';
	import type { Room } from 'matrix-js-sdk';
	import { Icon, IconButton, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { segments } from '$lib/chat/mentions';
	import { activeMention } from '$lib/chat/complete';
	import { app } from '$lib/state/app.svelte';
	import { ComposerState, mentionables } from '$lib/state/composer.svelte';
	import type { DraftState } from '$lib/state/draft.svelte';
	import MentionMenu from './MentionMenu.svelte';
	import PendingFiles from './PendingFiles.svelte';
	import DraftBanner from './DraftBanner.svelte';
	import ComposerField from './ComposerField.svelte';

	type Props = { room: Room; draft: DraftState; threadId: string | null; placeholder?: string };
	let { room, draft, threadId, placeholder }: Props = $props();

	const composer = untrack(() => new ComposerState(app.client!, room, threadId));
	const STEPS: Record<string, number> = { ArrowDown: 1, ArrowUp: -1 };
	const menuId = `mentions-${Math.random().toString(36).slice(2)}`;
	let input: HTMLTextAreaElement | undefined = $state();
	let picker: HTMLInputElement | undefined = $state();

	const members = $derived.by(() => {
		void app.version;
		return mentionables(room, app.client?.getUserId() ?? '', app.agents);
	});
	const options = $derived(composer.completions(members));
	const parts = $derived(segments(composer.text, composer.known(members)));
	const hint = $derived(placeholder ?? `Message ${room.name}`);
	const open = $derived(options.length > 0);

	$effect(() => {
		void room.loadMembersIfNeeded();
		return () => composer.typing.stop();
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
			await composer.submit(draft, members);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not send the message');
		} finally {
			composer.sending = false;
		}
	}

	function menuKey(event: KeyboardEvent): boolean {
		const step = STEPS[event.key];
		if (step) composer.selected = (composer.selected + step + options.length) % options.length;
		else if (event.key === 'Enter' || event.key === 'Tab') void pick(composer.selected);
		else if (event.key === 'Escape') composer.dismissed = activeMention(composer.text, composer.caret)?.start ?? -1;
		else return false;
		return true;
	}

	function onKeydown(event: KeyboardEvent) {
		if (event.isComposing) return;
		if (open && !event.shiftKey && menuKey(event)) event.preventDefault();
		else if (event.key === 'Enter' && !event.shiftKey) {
			event.preventDefault();
			void submit();
		} else if (event.key === 'Escape') draft.clear();
	}

	function onCaret() {
		const caret = input?.selectionStart ?? composer.text.length;
		if (caret !== composer.caret) composer.selected = 0;
		composer.caret = caret;
	}

	function onInput() {
		onCaret();
		composer.typed();
	}
</script>

<div class="relative px-3 pb-3 md:px-4">
	{#if open}
		<MentionMenu id={menuId} {options} selected={composer.selected} onpick={pick} onhover={(i) => (composer.selected = i)} />
	{/if}
	<div class="flex flex-col rounded-fc-lg border border-fc-border bg-fc-bg focus-within:border-fc-ring">
		<DraftBanner {room} {draft} />
		<PendingFiles files={composer.files} onremove={(f) => (composer.files = composer.files.filter((x) => x !== f))} />
		<div class="flex items-end gap-1 p-1">
			<IconButton variant="ghost" aria-label="Attach a file" title="Attach a file" disabled={!!draft.editing} onclick={() => picker?.click()}>
				<Icon icon={glyphs.paperclip} size={18} />
			</IconButton>
			<ComposerField
				bind:input
				bind:value={composer.text}
				{parts}
				oninput={onInput}
				onkeydown={onKeydown}
				onkeyup={onCaret}
				onclick={onCaret}
				onblur={() => composer.typing.stop()}
				placeholder={hint}
				aria-label={hint}
				role="combobox"
				aria-autocomplete="list"
				aria-expanded={open}
				aria-controls={open ? menuId : undefined}
				aria-activedescendant={open ? `${menuId}-${composer.selected}` : undefined}
			/>
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
