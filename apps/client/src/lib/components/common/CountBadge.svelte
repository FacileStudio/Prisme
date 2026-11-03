<script lang="ts">
	import { formatCount, type Badge } from '$lib/chat/unread';

	type Props = { badge: Badge; label?: string };
	let { badge, label = 'unread mentions' }: Props = $props();
</script>

{#if badge.count > 0}
	<span
		class="ml-auto inline-flex h-5 min-w-5 shrink-0 items-center justify-center rounded-fc-pill bg-fc-danger px-1.5 text-fc-xs font-semibold text-fc-danger-fg"
		aria-label="{badge.count} {label}"
	>
		{formatCount(badge.count)}
	</span>
{/if}
