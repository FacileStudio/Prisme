<script lang="ts">
	import { Button, Icon, IconButton, Skeleton, icons } from '@facile/muse';
	import type { LinkPreview } from '$lib/api/prisme';
	import { glyphs } from '$lib/glyphs';
	import { callLength } from '$lib/call/media';
	import { linkPreview, type SuiteLink } from '$lib/chat/links';
	import { t } from '$lib/i18n/index.svelte';

	let { link }: { link: SuiteLink } = $props();

	let preview = $state<LinkPreview | null>(null);
	let loading = $state(true);
	let playing = $state(false);

	$effect(() => {
		let live = true;
		loading = true;
		void linkPreview(link.url).then((found) => {
			if (!live) return;
			preview = found;
			loading = false;
		});
		return () => (live = false);
	});

	const fallback = $derived(t(`link.${link.app}`));
	const title = $derived(preview?.title || fallback);
	const detail = $derived.by(() => {
		if (preview?.duration_ms) return `${fallback} · ${callLength(preview.duration_ms)}`;
		return preview?.description || (preview?.title ? fallback : new URL(link.url).host);
	});
	const embed = $derived(preview?.embed_url ?? '');
	const call = $derived(link.app === 'echo');
</script>

<div class="mt-1 flex max-w-md flex-col gap-2 rounded-fc-lg bg-fc-component p-3">
	{#if playing && embed}
		<div class="relative overflow-hidden rounded-fc bg-fc-surface">
			<iframe
				src={embed}
				title={title}
				class="aspect-video w-full"
				allow="autoplay; fullscreen; picture-in-picture"
				allowfullscreen
				referrerpolicy="no-referrer"
				sandbox="allow-scripts allow-same-origin allow-presentation allow-popups"
			></iframe>
			<IconButton class="absolute top-2 right-2" variant="secondary" aria-label={t('link.close')} tooltip={t('link.close')} onclick={() => (playing = false)}>
				<Icon icon={icons.close} size={16} />
			</IconButton>
		</div>
	{:else if preview?.image}
		<div class="relative overflow-hidden rounded-fc bg-fc-surface">
			<img src={preview.image} alt="" loading="lazy" referrerpolicy="no-referrer" class="aspect-video w-full object-cover" />
			{#if embed}
				<button
					type="button"
					class="absolute inset-0 flex items-center justify-center focus-visible:outline-2 focus-visible:-outline-offset-4 focus-visible:outline-fc-ring"
					aria-label={t('link.play')}
					onclick={() => (playing = true)}
				>
					<span class="flex size-14 items-center justify-center rounded-fc-pill bg-fc-bg/85 text-fc-fg shadow-lg">
						<Icon icon={glyphs.play} size={24} />
					</span>
				</button>
			{/if}
		</div>
	{/if}
	<div class="flex items-center gap-3">
		<span class="flex size-10 shrink-0 items-center justify-center rounded-fc-md bg-fc-accent/10 text-fc-accent">
			<Icon icon={glyphs[link.app]} size={20} />
		</span>
		<div class="flex min-w-0 flex-1 flex-col">
			{#if loading}
				<Skeleton class="h-4 w-40" />
				<Skeleton class="mt-1 h-3 w-24" />
			{:else}
				<span class="truncate text-fc-sm font-semibold">{title}</span>
				<span class="truncate text-fc-xs text-fc-fg-muted">{detail}</span>
			{/if}
		</div>
		<Button
			size="sm"
			variant={call ? 'primary' : 'secondary'}
			icon={call ? glyphs.camera : icons.external}
			href={link.url}
			target="_blank"
		>
			{call ? t('link.join') : t('link.open')}
		</Button>
	</div>
</div>
