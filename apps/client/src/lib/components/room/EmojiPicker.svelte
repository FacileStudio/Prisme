<script lang="ts">
	import { Icon, icons } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	const EMOJIS = ['👍', '❤️', '😂', '🎉', '👀', '🙏', '🔥', '✅', '😮', '😢', '🚀', '👏'];

	let { onpick, onmore }: { onpick: (emoji: string) => void; onmore?: () => void } = $props();
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
	{#if onmore}
		<button
			type="button"
			role="menuitem"
			class="col-span-6 flex h-9 items-center justify-center gap-1.5 rounded-fc-md text-fc-xs font-medium text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-fc-ring"
			onclick={onmore}
		>
			<Icon icon={icons.plus} size={14} />{t('emoji.more')}
		</button>
	{/if}
</div>
