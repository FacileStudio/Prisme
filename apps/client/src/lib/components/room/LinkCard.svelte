<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import { Button, Icon, IconButton, Skeleton, icons, prefersReducedMotion } from '@facile/muse';
	import type { LinkPreview } from '$lib/api/prisme';
	import { glyphs } from '$lib/glyphs';
	import { callLength } from '$lib/call/media';
	import { linkPreview, type SuiteLink } from '$lib/chat/links';
	import { t } from '$lib/i18n/index.svelte';

	let { link }: { link: SuiteLink } = $props();

	let preview = $state<LinkPreview | null>(null);
	let loading = $state(true);
	let playing = $state(false);
	let thumbReady = $state(false);

	$effect(() => {
		let live = true;
		loading = true;
		thumbReady = false;
		void linkPreview(link.url).then((found) => {
			if (!live) return;
			preview = found;
			loading = false;
		});
		return () => (live = false);
	});

	const still = prefersReducedMotion();
	const reveal = { duration: still ? 0 : 220 };
	const grow = { duration: still ? 0 : 260, start: 0.96, opacity: 0 };

	const fallback = $derived(t(`link.${link.app}`));
	const title = $derived(preview?.title || fallback);
	const detail = $derived.by(() => {
		if (preview?.duration_ms) return `${fallback} · ${callLength(preview.duration_ms)}`;
		return preview?.description || (preview?.title ? fallback : new URL(link.url).host);
	});
	const embed = $derived(preview?.embed_url ?? '');
	const call = $derived(link.app === 'echo');
</script>

<div class="mt-1 flex max-w-md flex-col gap-2 rounded-fc-lg bg-fc-component p-2">
	<div class="flex h-14 items-center gap-3">
		<span class="relative flex size-14 shrink-0 items-center justify-center overflow-hidden rounded-fc-md bg-fc-accent/10 text-fc-accent">
			<Icon icon={glyphs[link.app]} size={22} />
			{#if preview?.image}
				<img
					src={preview.image}
					alt=""
					loading="lazy"
					referrerpolicy="no-referrer"
					onload={() => (thumbReady = true)}
					class="absolute inset-0 size-full object-cover transition-opacity duration-300 ease-fc motion-reduce:transition-none {thumbReady
						? 'opacity-100'
						: 'opacity-0'}"
				/>
			{/if}
			{#if embed && !playing}
				<button
					type="button"
					class="absolute inset-0 flex items-center justify-center focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring"
					aria-label={t('link.play')}
					onclick={() => (playing = true)}
				>
					<span class="flex size-8 items-center justify-center rounded-fc-pill bg-fc-bg/85 text-fc-fg shadow-lg">
						<Icon icon={glyphs.play} size={16} />
					</span>
				</button>
			{/if}
		</span>
		<div class="grid h-10 min-w-0 flex-1">
			{#if loading}
				<div class="col-start-1 row-start-1 flex flex-col justify-center gap-1.5" out:fade={reveal}>
					<Skeleton class="h-3.5 w-40 max-w-full" />
					<Skeleton class="h-3 w-24" />
				</div>
			{:else}
				<div class="col-start-1 row-start-1 flex min-w-0 flex-col justify-center" in:fade={reveal}>
					<span class="truncate text-fc-sm font-semibold">{title}</span>
					<span class="truncate text-fc-xs text-fc-fg-muted">{detail}</span>
				</div>
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
	{#if playing && embed}
		<div class="relative origin-top overflow-hidden rounded-fc bg-fc-surface" transition:scale={grow}>
			<iframe
				src={embed}
				{title}
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
	{/if}
</div>
