<script lang="ts">
	import type { Badge } from '$lib/chat/unread';
	import CountBadge from '$lib/components/common/CountBadge.svelte';

	type Props = { name: string; active?: boolean; invite?: boolean; badge: Badge; onclick: () => void };
	let { name, active = false, invite = false, badge, onclick }: Props = $props();

	const initials = $derived(
		name
			.split(/\s+/)
			.map((w) => w[0] ?? '')
			.join('')
			.slice(0, 2)
			.toUpperCase() || '?'
	);
</script>

<div class="relative">
	<button
		type="button"
		{onclick}
		title={invite ? `Join ${name}` : name}
		aria-label={invite ? `Accept the invite to ${name}` : name}
		aria-current={active ? 'true' : undefined}
		class="flex size-11 items-center justify-center rounded-fc-lg text-fc-sm font-semibold transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring {active
			? 'bg-fc-accent text-fc-accent-fg'
			: 'bg-fc-surface text-fc-fg hover:bg-fc-border'} {invite ? 'border-2 border-dashed border-fc-fg-muted' : ''}"
	>
		{initials}
	</button>
	{#if badge.count > 0}
		<span class="absolute -right-1.5 -top-1.5 flex">
			<CountBadge {badge} />
		</span>
	{:else if badge.bold && !active}
		<span class="absolute -left-2 top-1/2 h-2 w-1 -translate-y-1/2 rounded-fc-pill bg-fc-fg" aria-hidden="true"></span>
	{/if}
</div>
