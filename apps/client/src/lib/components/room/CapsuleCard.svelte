<script lang="ts">
	import { Button, Icon } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { prisme } from '$lib/api/prisme';
	import { openCapsule, type CapsuleInfo } from '$lib/chat/capsule';
	import { t } from '$lib/i18n/index.svelte';

	let { capsule, mine }: { capsule: CapsuleInfo; mine: boolean } = $props();

	type Phase = 'checking' | 'sealed' | 'opening' | 'opened' | 'gone' | 'failed';
	let phase = $state<Phase>('checking');
	let text = $state('');

	$effect(() => {
		const id = capsule.id;
		let live = true;
		prisme
			.capsuleMeta(id)
			.then((meta) => live && phase === 'checking' && (phase = meta.exists ? 'sealed' : 'gone'))
			.catch(() => live && phase === 'checking' && (phase = 'sealed'));
		return () => (live = false);
	});

	async function open() {
		phase = 'opening';
		try {
			text = await openCapsule(capsule);
			phase = 'opened';
		} catch (err) {
			phase = (err as { status?: number }).status === 404 ? 'gone' : 'failed';
		}
	}

	const caption = $derived.by(() => {
		if (phase === 'gone') return t('capsule.gone');
		if (phase === 'failed') return t('capsule.failed');
		if (phase === 'opened') return t('capsule.opened');
		return mine ? t('capsule.mine') : t('capsule.warn');
	});
</script>

<div class="mt-1 flex max-w-md flex-col gap-2 rounded-fc-lg bg-fc-component p-3">
	<div class="flex items-center gap-3">
		<span
			class="flex size-10 shrink-0 items-center justify-center rounded-fc-pill {phase === 'gone'
				? 'bg-fc-surface text-fc-fg-muted'
				: 'bg-fc-accent/10 text-fc-accent'}"
		>
			<Icon icon={glyphs.capsule} size={20} />
		</span>
		<div class="flex min-w-0 flex-1 flex-col">
			<span class="text-fc-sm font-semibold">{t('capsule.sealed')}</span>
			<span class="text-fc-xs text-fc-fg-muted">{caption}</span>
		</div>
		{#if !mine && (phase === 'sealed' || phase === 'opening' || phase === 'failed')}
			<Button size="sm" disabled={phase === 'opening'} onclick={open}>{t('capsule.open')}</Button>
		{/if}
	</div>
	{#if phase === 'opened'}
		<p class="whitespace-pre-wrap break-words rounded-fc bg-fc-surface p-3 text-fc-sm">{text}</p>
	{/if}
</div>
