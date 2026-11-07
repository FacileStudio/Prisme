<script lang="ts">
	import type { IContent } from 'matrix-js-sdk';
	import { Icon, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { formatBytes } from '$lib/i18n/format';
	import { t } from '$lib/i18n/index.svelte';
	import { loadMedia } from '$lib/matrix/media';

	type FileContent = IContent & { url?: string; body?: string; filename?: string; info?: { size?: number; mimetype?: string } };
	let { content }: { content: FileContent } = $props();
	let busy = $state(false);

	const name = $derived(content.filename ?? content.body ?? t('file.fallback'));
	const size = $derived(formatBytes(content.info?.size ?? 0));

	async function download() {
		if (!content.url) return;
		busy = true;
		try {
			const link = document.createElement('a');
			link.href = await loadMedia(content.url);
			link.download = name;
			link.click();
		} finally {
			busy = false;
		}
	}
</script>

<button
	type="button"
	onclick={download}
	disabled={busy}
	class="mt-1 flex min-h-11 w-full max-w-sm items-center gap-3 rounded-fc-md bg-fc-component px-3 py-2 text-left hover:bg-fc-surface focus-visible:outline-2 focus-visible:outline-fc-ring"
>
	<Icon icon={glyphs.file} size={22} class="shrink-0 text-fc-fg-muted" />
	<span class="flex min-w-0 flex-1 flex-col">
		<span class="truncate text-fc-sm font-medium">{name}</span>
		<span class="text-fc-xs text-fc-fg-muted">{size || content.info?.mimetype || t('file.kind')}</span>
	</span>
	<Icon icon={icons.download} size={18} class="shrink-0" />
</button>
