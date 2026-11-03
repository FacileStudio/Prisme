<script lang="ts">
	import type { IContent } from 'matrix-js-sdk';
	import { Skeleton } from '@facile/muse';
	import { loadMedia } from '$lib/matrix/media';

	type ImageContent = IContent & { url?: string; body?: string; info?: { w?: number; h?: number; mimetype?: string } };
	let { content }: { content: ImageContent } = $props();

	let thumb = $state('');
	let failed = $state(false);

	const ratio = $derived(content.info?.w && content.info?.h ? content.info.w / content.info.h : 4 / 3);

	$effect(() => {
		const url = content.url;
		if (!url) return;
		let live = true;
		loadMedia(url, 800).then(
			(blob) => live && (thumb = blob),
			() => live && (failed = true)
		);
		return () => (live = false);
	});

	async function open() {
		if (!content.url) return;
		const full = await loadMedia(content.url).catch(() => thumb);
		window.open(full, '_blank', 'noopener');
	}
</script>

<button
	type="button"
	onclick={open}
	class="mt-1 block max-w-full overflow-hidden rounded-fc-md border border-fc-border focus-visible:outline-2 focus-visible:outline-fc-ring"
	style="width: min(100%, {Math.min(content.info?.w ?? 360, 360)}px); aspect-ratio: {ratio};"
	aria-label="Open {content.body ?? 'image'}"
>
	{#if thumb}
		<img src={thumb} alt={content.body ?? ''} class="size-full object-cover" />
	{:else if failed}
		<span class="flex size-full items-center justify-center bg-fc-surface text-fc-xs text-fc-fg-muted">Image unavailable</span>
	{:else}
		<Skeleton class="size-full" />
	{/if}
</button>
