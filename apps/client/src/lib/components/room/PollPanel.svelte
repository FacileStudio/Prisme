<script lang="ts">
	import type { ComposerPanelProps } from '$lib/state/composer-menu';
	import PollDialog from './PollDialog.svelte';

	let { context, onclose }: ComposerPanelProps = $props();
</script>

<PollDialog room={context.room} threadId={context.threadId} {onclose} />
