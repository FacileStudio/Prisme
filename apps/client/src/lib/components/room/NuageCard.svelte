<script lang="ts">
	import { Button, Icon, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { fileGlyph, fileSize, type NuageInfo } from '$lib/chat/nuage';
	import { t } from '$lib/i18n/index.svelte';

	let { file }: { file: NuageInfo } = $props();
</script>

<div class="mt-1 flex max-w-md items-center gap-3 rounded-fc-lg bg-fc-component p-2">
	<span class="relative flex size-12 shrink-0 items-center justify-center rounded-fc-md bg-fc-accent/10 text-fc-accent">
		<Icon icon={fileGlyph(file.mime)} size={22} />
		<span class="absolute -right-1 -bottom-1 flex size-5 items-center justify-center rounded-fc-pill bg-fc-component text-fc-fg-muted">
			<Icon icon={glyphs.nuage} size={12} />
		</span>
	</span>
	<div class="flex min-w-0 flex-1 flex-col">
		<span class="truncate text-fc-sm font-semibold text-fc-fg">{file.name}</span>
		<span class="truncate text-fc-xs text-fc-fg-muted">{t('nuage.card')} · {fileSize(file.size)}</span>
	</div>
	<Button size="sm" variant="secondary" icon={icons.external} href={file.url} target="_blank" rel="noopener noreferrer">{t('nuage.open')}</Button>
</div>
