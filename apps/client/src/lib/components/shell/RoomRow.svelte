<script lang="ts">
	import type { Snippet } from 'svelte';
	import { page } from '$app/state';
	import { Icon } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import type { Badge } from '$lib/chat/unread';
	import { ui } from '$lib/state/ui.svelte';
	import CountBadge from '$lib/components/common/CountBadge.svelte';

	type Props = {
		href: string;
		roomId: string;
		name: string;
		badge: Badge;
		muted?: boolean;
		icon?: string;
		lead?: Snippet;
		label?: Snippet;
	};
	let { href, roomId, name, badge, muted = false, icon, lead, label }: Props = $props();

	const active = $derived(page.params.id ? decodeURIComponent(page.params.id) === roomId : false);
</script>

<a
	{href}
	onclick={() => (ui.drawer = false)}
	aria-current={active ? 'page' : undefined}
	class="flex min-h-11 items-center gap-2 rounded-fc-md px-2 text-fc-sm transition-colors focus-visible:outline-2 focus-visible:outline-fc-ring {active
		? 'bg-fc-accent text-fc-accent-fg'
		: badge.bold
			? 'font-semibold text-fc-fg hover:bg-fc-surface'
			: 'text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg'}"
>
	{#if lead}
		{@render lead()}
	{:else if icon}
		<Icon icon={icon} size={16} class="shrink-0" />
	{/if}
	<span class="flex min-w-0 flex-1 items-center">
		{#if label}{@render label()}{:else}<span class="truncate">{name}</span>{/if}
	</span>
	{#if muted}
		<Icon icon={glyphs.bellOff} size={14} class="shrink-0 opacity-70" aria-label="Muted" />
	{/if}
	<CountBadge {badge} />
</a>
