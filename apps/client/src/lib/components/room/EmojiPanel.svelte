<script lang="ts">
	import { tick } from 'svelte';
	import { fade } from 'svelte/transition';
	import { Icon, Spinner, icons, prefersReducedMotion } from '@facile/muse';
	import { pop } from '$lib/motion';
	import { glyphs } from '$lib/glyphs';
	import { locale, t } from '$lib/i18n/index.svelte';
	import { EMOJI_GROUPS, loadEmoji, searchEmoji, withTone, type Emoji } from '$lib/chat/emoji';
	import { recentEmoji, rememberEmoji, savedTone, saveTone } from '$lib/chat/emoji-prefs';

	type Props = { onpick: (emoji: string) => void; onclose: (refocus: boolean) => void; anchor?: HTMLElement; class?: string };
	let { onpick, onclose, anchor, class: className = '' }: Props = $props();

	const COLUMNS = 8;
	const TONES = ['✋', '✋🏻', '✋🏼', '✋🏽', '✋🏾', '✋🏿'];
	const GROUP_ICONS: Record<number, string> = {
		0: glyphs.smile,
		1: glyphs.emojiPeople,
		3: glyphs.emojiAnimals,
		4: glyphs.emojiFood,
		5: glyphs.emojiTravel,
		6: glyphs.emojiActivities,
		7: glyphs.emojiObjects,
		8: glyphs.emojiSymbols,
		9: glyphs.emojiFlags
	};

	let all = $state.raw<Emoji[]>([]);
	let failed = $state(false);
	let query = $state('');
	let tone = $state(savedTone());
	let toning = $state(false);
	let recent = $state(recentEmoji());
	let focused = $state(0);
	let panel: HTMLDivElement | undefined = $state();
	let scroller: HTMLDivElement | undefined = $state();
	let search: HTMLInputElement | undefined = $state();

	$effect(() => {
		const current = locale();
		loadEmoji(current)
			.then((list) => (all = list))
			.catch(() => (failed = true));
	});

	$effect(() => {
		void tick().then(() => search?.focus());
		const outside = (event: PointerEvent) => {
			const target = event.target as Node;
			if (panel && !panel.contains(target) && !anchor?.contains(target)) onclose(false);
		};
		document.addEventListener('pointerdown', outside, true);
		return () => document.removeEventListener('pointerdown', outside, true);
	});

	const byUnicode = $derived(new Map(all.map((e) => [e.unicode, e])));
	const results = $derived(searchEmoji(all, query));
	const sections = $derived.by(() => {
		const recentRow = recent.map((u) => byUnicode.get(u)).filter((e): e is Emoji => !!e);
		const groups = EMOJI_GROUPS.map((group) => ({ id: String(group), icon: GROUP_ICONS[group], emojis: all.filter((e) => e.group === group) }));
		return recentRow.length ? [{ id: 'recent', icon: icons.clock, emojis: recentRow }, ...groups] : groups;
	});
	const flat = $derived(query.trim() ? results : sections.flatMap((s) => s.emojis));

	function pick(emoji: Emoji) {
		const glyph = withTone(emoji, tone);
		recent = rememberEmoji(emoji.unicode);
		onpick(glyph);
	}

	function chooseTone(next: number) {
		tone = next;
		toning = false;
		saveTone(next);
	}

	function jump(id: string) {
		scroller?.querySelector(`[data-section="${id}"]`)?.scrollIntoView({ block: 'start', behavior: prefersReducedMotion() ? 'auto' : 'smooth' });
	}

	async function focusAt(index: number) {
		focused = Math.max(0, Math.min(index, flat.length - 1));
		await tick();
		scroller?.querySelector<HTMLButtonElement>(`[data-index="${focused}"]`)?.focus();
	}

	const MOVES: Record<string, number> = { ArrowRight: 1, ArrowLeft: -1, ArrowDown: COLUMNS, ArrowUp: -COLUMNS };

	function gridKey(event: KeyboardEvent) {
		const step = MOVES[event.key];
		if (step) {
			event.preventDefault();
			if (focused + step < 0) search?.focus();
			else void focusAt(focused + step);
		} else if (event.key === 'Escape') {
			event.preventDefault();
			onclose(true);
		}
	}

	function searchKey(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			event.preventDefault();
			if (query) query = '';
			else onclose(true);
		} else if (event.key === 'ArrowDown') {
			event.preventDefault();
			void focusAt(0);
		} else if (event.key === 'Enter' && flat[0]) {
			event.preventDefault();
			pick(flat[0]);
		}
	}

	function indexOf(emoji: Emoji, sectionIndex: number, offset: number): number {
		if (query.trim()) return offset;
		let base = 0;
		for (let i = 0; i < sectionIndex; i++) base += sections[i].emojis.length;
		return base + offset;
	}
</script>

{#snippet cell(emoji: Emoji, index: number)}
	<button
		type="button"
		data-index={index}
		tabindex={index === focused ? 0 : -1}
		title={emoji.label}
		aria-label={emoji.label}
		class="flex aspect-square items-center justify-center rounded-fc-md text-[1.375rem] leading-none hover:bg-fc-surface focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring"
		onfocus={() => (focused = index)}
		onclick={() => pick(emoji)}
	>
		{withTone(emoji, tone)}
	</button>
{/snippet}

<div
	bind:this={panel}
	role="dialog"
	aria-label={t('emoji.open')}
	tabindex="-1"
	transition:pop
	class="flex h-[26rem] max-h-[70dvh] w-[min(22rem,calc(100vw-1.5rem))] flex-col overflow-hidden rounded-fc-lg bg-fc-component shadow-lg {className}"
>
	<div class="flex items-center gap-1.5 p-2">
		<label class="flex h-10 min-w-0 flex-1 items-center gap-2 rounded-fc-pill bg-fc-field px-3 focus-within:outline-2 focus-within:outline-fc-ring">
			<Icon icon={icons.search} size={16} class="shrink-0 text-fc-fg-muted" />
			<input
				bind:this={search}
				bind:value={query}
				type="search"
				placeholder={t('emoji.search')}
				aria-label={t('emoji.search')}
				class="min-w-0 flex-1 bg-transparent text-fc-sm text-fc-fg outline-none placeholder:text-fc-fg-muted"
				onkeydown={searchKey}
				oninput={() => (focused = 0)}
			/>
		</label>
		<div class="relative">
			<button
				type="button"
				class="flex size-10 items-center justify-center rounded-fc-pill text-xl hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-fc-ring"
				aria-label={t('emoji.tone')}
				aria-expanded={toning}
				onclick={() => (toning = !toning)}
			>
				{TONES[tone]}
			</button>
			{#if toning}
				<div transition:pop={{ start: 0.85, duration: 150 }} class="absolute right-0 top-11 z-10 flex origin-top-right gap-0.5 rounded-fc-pill bg-fc-component p-1 shadow-lg" role="radiogroup" aria-label={t('emoji.tone')}>
					{#each TONES as glyph, index (glyph)}
						<button
							type="button"
							role="radio"
							aria-checked={index === tone}
							aria-label={`${t('emoji.tone')} ${index}`}
							class="flex size-9 items-center justify-center rounded-fc-pill text-lg hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-fc-ring {index === tone ? 'bg-fc-surface' : ''}"
							onclick={() => chooseTone(index)}
						>
							{glyph}
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>

	{#if !query.trim() && all.length > 0}
		<div class="flex shrink-0 gap-0.5 overflow-x-auto px-2 pb-1" role="toolbar" aria-label={t('emoji.open')}>
			{#each sections as section (section.id)}
				<button
					type="button"
					class="flex size-9 shrink-0 items-center justify-center rounded-fc-pill text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-fc-ring"
					aria-label={t(`emoji.group.${section.id}` as 'emoji.group.recent')}
					title={t(`emoji.group.${section.id}` as 'emoji.group.recent')}
					onclick={() => jump(section.id)}
				>
					<Icon icon={section.icon} size={18} />
				</button>
			{/each}
		</div>
	{/if}

	<div bind:this={scroller} class="min-h-0 flex-1 overflow-y-auto overscroll-contain px-2 pb-2" role="presentation" onkeydown={gridKey}>
		{#if failed}
			<p class="px-2 py-6 text-center text-fc-sm text-fc-fg-muted">{t('emoji.none')}</p>
		{:else if all.length === 0}
			<div class="flex h-full items-center justify-center" aria-label={t('emoji.loading')}><Spinner /></div>
		{:else if query.trim()}
			{#if results.length === 0}
				<p class="px-2 py-6 text-center text-fc-sm text-fc-fg-muted">{t('emoji.none')}</p>
			{:else}
				<div in:fade={{ duration: prefersReducedMotion() ? 0 : 120 }} class="grid grid-cols-8 gap-0.5 pt-1">
					{#each results as emoji, i (emoji.unicode)}
						{@render cell(emoji, indexOf(emoji, 0, i))}
					{/each}
				</div>
			{/if}
		{:else}
			{#each sections as section, s (section.id)}
				<section data-section={section.id} class="[content-visibility:auto] [contain-intrinsic-size:auto_12rem]">
					<h3 class="sticky top-0 z-[1] bg-fc-component px-1 pt-2 pb-1 text-fc-xs font-medium text-fc-fg-muted">
						{t(`emoji.group.${section.id}` as 'emoji.group.recent')}
					</h3>
					<div class="grid grid-cols-8 gap-0.5">
						{#each section.emojis as emoji, i (emoji.unicode)}
							{@render cell(emoji, indexOf(emoji, s, i))}
						{/each}
					</div>
				</section>
			{/each}
		{/if}
	</div>
</div>
