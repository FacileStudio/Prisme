<script lang="ts">
	const EMOJIS = ['👍', '❤️', '😂', '🎉', '👀', '🙏', '🔥', '✅', '😮', '😢', '🚀', '👏'];

	let { onpick }: { onpick: (emoji: string) => void } = $props();
</script>

<div class="grid w-max grid-cols-6 gap-0.5 rounded-fc-lg border border-fc-border bg-fc-component p-1 shadow-lg" role="menu" aria-label="Pick a reaction">
	{#each EMOJIS as emoji (emoji)}
		<button
			type="button"
			role="menuitem"
			class="flex size-10 items-center justify-center rounded-fc-md text-fc-lg hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-fc-ring"
			aria-label="React with {emoji}"
			onclick={() => onpick(emoji)}
		>
			{emoji}
		</button>
	{/each}
</div>
