<script lang="ts">
	import { tick, untrack } from 'svelte';
	import type { Room } from 'matrix-js-sdk';
	import { Icon, Select, Tooltip, icons, toast } from '@facile/muse';
	import { EXPIRIES } from '$lib/chat/capsule';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import { segments } from '$lib/chat/mentions';
	import { activeMention } from '$lib/chat/complete';
	import { app } from '$lib/state/app.svelte';
	import { ComposerState, mentionables } from '$lib/state/composer.svelte';
	import type { DraftState } from '$lib/state/draft.svelte';
	import MentionMenu from './MentionMenu.svelte';
	import PendingFiles from './PendingFiles.svelte';
	import DraftBanner from './DraftBanner.svelte';
	import ComposerField from './ComposerField.svelte';
	import ComposerMenu from './ComposerMenu.svelte';
	import EmojiPanel from './EmojiPanel.svelte';
	import { visibleEntries, type ComposerContext } from '$lib/state/composer-menu';
	import { ComposerPopups } from '$lib/state/composer-popups.svelte';

	type Props = { room: Room; draft: DraftState; threadId: string | null; placeholder?: string };
	let { room, draft, threadId, placeholder }: Props = $props();

	const composer = untrack(() => new ComposerState(app.client!, room, threadId));
	const STEPS: Record<string, number> = { ArrowDown: 1, ArrowUp: -1 };
	const menuId = `mentions-${Math.random().toString(36).slice(2)}`;
	let input: HTMLTextAreaElement | undefined = $state();
	let picker: HTMLInputElement | undefined = $state();
	let plus: HTMLButtonElement | undefined = $state();
	let smiley: HTMLButtonElement | undefined = $state();
	const popups = new ComposerPopups({ plus: () => plus, smiley: () => smiley, input: () => input });
	const addMenuId = `composer-menu-${Math.random().toString(36).slice(2)}`;

	const members = $derived.by(() => {
		void app.version;
		return mentionables(room, app.client?.getUserId() ?? '', app.agents);
	});
	const options = $derived(composer.completions(members));
	const parts = $derived(segments(composer.text, composer.known(members)));
	const hint = $derived(composer.sealed ? t('capsule.placeholder') : (placeholder ?? t('composer.placeholder', { name: room.name })));
	const context: ComposerContext = {
		get room() {
			return room;
		},
		get threadId() {
			return threadId;
		},
		get editing() {
			return !!draft.editing;
		},
		composer,
		pickFiles: () => picker?.click()
	};
	const entries = $derived(visibleEntries(context));
	const expiries = $derived(EXPIRIES.map((value) => ({ value, label: t(`capsule.${value}`) })));
	const open = $derived(options.length > 0);
	const canSend = $derived(!composer.sending && (!!composer.text.trim() || (!composer.sealed && composer.files.length > 0)));
	const ROUND =
		'inline-flex size-11 shrink-0 items-center justify-center rounded-fc-pill focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring disabled:cursor-not-allowed disabled:opacity-40';

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

	async function insertEmoji(emoji: string) {
		const start = input?.selectionStart ?? composer.text.length;
		const caret = composer.insert(emoji, start, input?.selectionEnd ?? start);
		await tick();
		input?.setSelectionRange(caret, caret);
	}

	async function submit() {
		if (!canSend) return;
		composer.sending = true;
		try {
			await composer.submit(draft, members);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('composer.sendFailed'));
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
	{#if popups.menu}
		<ComposerMenu id={addMenuId} {entries} {context} onpick={(entry) => popups.choose(entry, context)} onclose={(r) => popups.closeMenu(r)} />
	{/if}
	{#if popups.panel?.panel}
		{@const Panel = popups.panel.panel}
		<Panel {context} onclose={() => popups.closePanel()} />
	{/if}
	{#if popups.emoji}
		<div class="absolute bottom-full left-3 z-30 mb-2 origin-bottom-left md:left-4">
			<EmojiPanel onpick={insertEmoji} onclose={(r) => popups.closeEmoji(r)} anchor={smiley} />
		</div>
	{/if}
	{#if open}
		<MentionMenu id={menuId} {options} selected={composer.selected} onpick={pick} onhover={(i) => (composer.selected = i)} />
	{/if}
	<div class="flex flex-col gap-2 rounded-[31px] bg-fc-field p-2">
		<DraftBanner {room} {draft} />
		<PendingFiles files={composer.files} onremove={(f) => (composer.files = composer.files.filter((x) => x !== f))} />
		{#if composer.sealed}
			<div class="flex flex-wrap items-center gap-2 px-3 pt-1 text-fc-xs text-fc-fg-muted">
				<Icon icon={glyphs.capsule} size={16} class="text-fc-accent" />
				<span class="font-medium text-fc-fg">{t('capsule.on')}</span>
				<span class="ml-auto">{t('capsule.expiry')}</span>
				<Select bind:value={composer.expiry} options={expiries} aria-label={t('capsule.expiry')} class="h-8 text-fc-xs" />
			</div>
		{/if}
		<div class="flex items-end gap-2">
			<Tooltip text={t('composerMenu.open')} describe={false} class="contents">
				<button
					bind:this={plus}
					type="button"
					aria-label={t('composerMenu.open')}
					aria-haspopup="menu"
					aria-expanded={popups.menu}
					aria-controls={addMenuId}
					disabled={!!draft.editing || entries.length === 0}
					class="{ROUND} {popups.menu || popups.panel ? 'bg-fc-surface text-fc-fg' : 'text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg'}"
					onclick={() => popups.toggleMenu()}
				>
					<Icon icon={icons.plus} size={20} class="transition-transform duration-150 motion-reduce:transition-none {popups.menu ? 'rotate-45' : ''}" />
				</button>
			</Tooltip>
			<Tooltip text={t('emoji.open')} describe={false} class="contents">
				<button
					bind:this={smiley}
					type="button"
					aria-label={t('emoji.open')}
					aria-haspopup="dialog"
					aria-expanded={popups.emoji}
					class="{ROUND} -ml-1.5 {popups.emoji ? 'bg-fc-surface text-fc-fg' : 'text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg'}"
					onclick={() => popups.toggleEmoji()}
				>
					<Icon icon={glyphs.smile} size={20} />
				</button>
			</Tooltip>
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
			<Tooltip text={draft.editing ? t('composer.saveEdit') : t('composer.send')} describe={false} class="contents">
				<button
					type="button"
					aria-label={draft.editing ? t('composer.saveEdit') : t('composer.send')}
					disabled={!canSend}
					class="{ROUND} bg-fc-accent text-fc-accent-fg hover:opacity-90"
					onclick={submit}
				>
					<Icon icon={glyphs.send} size={18} />
				</button>
			</Tooltip>
		</div>
	</div>
	<input bind:this={picker} type="file" multiple hidden onchange={(e) => {
		composer.files = [...composer.files, ...(e.currentTarget.files ?? [])];
		e.currentTarget.value = '';
	}} />
</div>
