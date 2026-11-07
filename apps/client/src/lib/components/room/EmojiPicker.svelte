<script lang="ts">
	import { t } from '$lib/i18n/index.svelte';

	const EMOJIS = ['👍', '❤️', '😂', '🎉', '👀', '🙏', '🔥', '✅', '😮', '😢', '🚀', '👏'];

	let { onpick }: { onpick: (emoji: string) => void } = $props();
</script>

<div class="grid w-max grid-cols-6 gap-0.5 rounded-fc-lg bg-fc-component p-1 shadow-lg" role="menu" aria-label={t('emoji.pick')}>
	{#each EMOJIS as emoji (emoji)}
		<button
			type="button"
			role="menuitem"
			class="flex size-10 items-center justify-center rounded-fc-md text-fc-lg hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-fc-ring"
			aria-label={t('emoji.react', { emoji })}
			onclick={() => onpick(emoji)}
		>
			{emoji}
		</button>
	{/each}
</div>
