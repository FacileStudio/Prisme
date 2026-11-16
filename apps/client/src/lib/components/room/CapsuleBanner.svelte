<script lang="ts">
	import { Icon, Tooltip, icons } from '@facile/muse';
	import { EXPIRIES, type Expiry } from '$lib/chat/capsule';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';

	type Props = { expiry: Expiry; onclose: () => void };
	let { expiry = $bindable(), onclose }: Props = $props();

	const STEPS: Record<string, number> = { ArrowRight: 1, ArrowDown: 1, ArrowLeft: -1, ArrowUp: -1 };
	let group: HTMLDivElement | undefined = $state();

	function onKeydown(event: KeyboardEvent) {
		const step = STEPS[event.key];
		const edge = event.key === 'Home' ? 0 : event.key === 'End' ? EXPIRIES.length - 1 : -1;
		if (!step && edge < 0) return;
		event.preventDefault();
		const index = edge >= 0 ? edge : (EXPIRIES.indexOf(expiry) + step + EXPIRIES.length) % EXPIRIES.length;
		expiry = EXPIRIES[index];
		group?.querySelectorAll<HTMLButtonElement>('[role="radio"]')[index]?.focus();
	}
</script>

<div class="flex flex-wrap items-center gap-x-3 gap-y-2 rounded-[23px] bg-fc-accent/10 py-1.5 pr-1.5 pl-3">
	<span class="flex min-w-0 flex-1 items-center gap-2">
		<Icon icon={glyphs.capsule} size={16} class="shrink-0 text-fc-accent" />
		<span class="truncate text-fc-xs font-medium text-fc-fg">{t('capsule.on')}</span>
	</span>
	<div
		bind:this={group}
		role="radiogroup"
		tabindex="-1"
		aria-label={t('capsule.expiry')}
		class="flex items-center gap-1"
		onkeydown={onKeydown}
	>
		{#each EXPIRIES as value (value)}
			{@const checked = value === expiry}
			<button
				type="button"
				role="radio"
				aria-checked={checked}
				aria-label={t(`capsule.${value}`)}
				tabindex={checked ? 0 : -1}
				class="h-8 rounded-fc-pill px-3 text-fc-xs font-medium whitespace-nowrap transition-colors duration-150 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring motion-reduce:transition-none {checked
					? 'bg-fc-accent text-fc-accent-fg'
					: 'text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg'}"
				onclick={() => (expiry = value)}
			>
				{t(`capsule.short.${value}`)}
			</button>
		{/each}
	</div>
	<Tooltip text={t('capsule.cancel')} describe={false} class="contents">
		<button
			type="button"
			aria-label={t('capsule.cancel')}
			class="inline-flex size-8 shrink-0 items-center justify-center rounded-fc-pill text-fc-fg-muted transition-colors duration-150 hover:bg-fc-surface hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring motion-reduce:transition-none"
			onclick={onclose}
		>
			<Icon icon={icons.close} size={16} />
		</button>
	</Tooltip>
</div>
