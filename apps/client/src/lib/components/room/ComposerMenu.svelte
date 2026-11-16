<script lang="ts">
	import { tick } from 'svelte';
	import { Icon } from '@facile/muse';
	import { pop, rise, stagger } from '$lib/motion';
	import { t } from '$lib/i18n/index.svelte';
	import type { ComposerContext, ComposerEntry } from '$lib/state/composer-menu';

	type Props = {
		id: string;
		entries: ComposerEntry[];
		context: ComposerContext;
		onpick: (entry: ComposerEntry) => void;
		onclose: (refocus: boolean) => void;
	};
	let { id, entries, context, onpick, onclose }: Props = $props();
	let menu: HTMLDivElement | undefined = $state();

	const KEYS: Record<string, (index: number, count: number) => number> = {
		ArrowRight: (i, n) => (i + 1) % n,
		ArrowDown: (i, n) => (i + 1) % n,
		ArrowLeft: (i, n) => (i - 1 + n) % n,
		ArrowUp: (i, n) => (i - 1 + n) % n,
		Home: () => 0,
		End: (_, n) => n - 1
	};

	const items = () => Array.from(menu?.querySelectorAll<HTMLButtonElement>('[role="menuitem"]') ?? []);

	$effect(() => {
		void tick().then(() => items()[0]?.focus());
		const outside = (event: PointerEvent) => {
			const target = event.target as Node;
			if (menu && !menu.contains(target) && !(target instanceof Element && target.closest(`[aria-controls="${id}"]`))) onclose(false);
		};
		document.addEventListener('pointerdown', outside, true);
		return () => document.removeEventListener('pointerdown', outside, true);
	});

	function onkeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' || event.key === 'Tab') {
			if (event.key === 'Escape') event.preventDefault();
			onclose(event.key === 'Escape');
			return;
		}
		const move = KEYS[event.key];
		if (!move) return;
		event.preventDefault();
		const list = items();
		const current = list.indexOf(document.activeElement as HTMLButtonElement);
		list[move(Math.max(current, 0), list.length)]?.focus();
	}
</script>

<div
	bind:this={menu}
	{id}
	role="menu"
	tabindex="-1"
	aria-label={t('composerMenu.label')}
	aria-orientation="horizontal"
	{onkeydown}
	transition:pop
	class="absolute bottom-full left-3 z-30 mb-2 flex max-w-[calc(100%-1.5rem)] origin-bottom-left items-center gap-1 overflow-x-auto rounded-fc-pill bg-fc-component p-1.5 shadow-lg md:left-4"
>
	{#each entries as entry, index (entry.id)}
		{@const pressed = entry.active?.(context) ?? false}
		<button
			type="button"
			role="menuitem"
			tabindex="-1"
			in:rise|global={{ delay: 40 + stagger(index, 28), distance: 4 }}
			class="inline-flex h-10 shrink-0 items-center gap-2 rounded-fc-pill px-3.5 text-fc-sm font-medium focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring {pressed
				? 'bg-fc-accent/15 text-fc-accent'
				: 'text-fc-fg hover:bg-fc-surface'}"
			onclick={() => onpick(entry)}
		>
			<Icon icon={entry.icon} size={18} class={pressed ? '' : 'text-fc-fg-muted'} />
			<span class="whitespace-nowrap">{entry.label(context)}</span>
		</button>
	{/each}
</div>
