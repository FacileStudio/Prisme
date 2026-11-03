<script lang="ts">
	import { Icon, IconButton, icons } from '@facile/muse';

	type Props = { title: string; action?: string; onaction?: () => void };
	let { title, action, onaction }: Props = $props();
</script>

<div class="flex min-h-11 items-center pl-4 pr-1 pt-2">
	<h3 class="flex-1 text-fc-xs font-semibold uppercase tracking-wide text-fc-fg-muted">{title}</h3>
	{#if action && onaction}
		<IconButton variant="ghost" aria-label={action} title={action} onclick={onaction}>
			<Icon icon={icons.plus} size={18} />
		</IconButton>
	{/if}
</div>
