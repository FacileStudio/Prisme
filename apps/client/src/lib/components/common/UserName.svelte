<script lang="ts">
	import { Badge } from '@facile/muse';
	import { isAgent } from '$lib/chat/users';

	type Props = { userId: string; name: string; class?: string };
	let { userId, name, class: className = '' }: Props = $props();
</script>

<span class="inline-flex min-w-0 items-center gap-1.5 {className}">
	<span class="truncate">{name}</span>
	{#if isAgent(userId)}
		<Badge tone="admin" class="shrink-0">Agent</Badge>
	{/if}
</span>
