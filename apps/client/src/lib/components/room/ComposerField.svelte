<script lang="ts">
	import type { HTMLTextareaAttributes } from 'svelte/elements';
	import type { Segment } from '$lib/chat/mentions';

	type Props = Omit<HTMLTextareaAttributes, 'value'> & {
		value: string;
		parts: Segment[];
		input?: HTMLTextAreaElement;
	};
	let { value = $bindable(), parts, input = $bindable(), onscroll, ...rest }: Props = $props();
	let mirror: HTMLDivElement | undefined = $state();

	const FIELD =
		'col-start-1 row-start-1 max-h-48 min-h-11 w-full px-1 py-2.5 text-fc-md leading-normal [scrollbar-gutter:stable] md:text-fc-sm';

	function sync() {
		if (mirror && input) mirror.scrollTop = input.scrollTop;
	}

	$effect(() => {
		void parts;
		sync();
	});
</script>

<div class="grid min-w-0 flex-1">
	<div bind:this={mirror} aria-hidden="true" class="{FIELD} pointer-events-none overflow-hidden whitespace-pre-wrap break-words text-fc-fg">
		{#each parts as part, index (index)}{#if part.mention}<span class="mention-token">{part.text}</span>{:else}{part.text}{/if}{/each}&#8203;
	</div>
	<textarea
		bind:this={input}
		bind:value
		onscroll={(event) => {
			sync();
			onscroll?.(event);
		}}
		rows="1"
		class="{FIELD} field-sizing-content resize-none bg-transparent text-transparent caret-fc-fg outline-none placeholder:text-fc-fg-muted selection:bg-fc-accent/20"
		{...rest}
	></textarea>
</div>

<style>
	.mention-token {
		border-radius: var(--radius-fc-sm);
		background: color-mix(in oklch, var(--color-fc-accent) 12%, transparent);
		box-shadow: 0 0 0 2px color-mix(in oklch, var(--color-fc-accent) 12%, transparent);
		box-decoration-break: clone;
	}
</style>
