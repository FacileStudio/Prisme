<script lang="ts">
	import { Icon, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';

	type Props = { files: File[]; onremove: (file: File) => void };
	let { files, onremove }: Props = $props();

	const previews = $derived(files.map((file) => ({ file, url: file.type.startsWith('image/') ? URL.createObjectURL(file) : '' })));

	$effect(() => {
		const urls = previews.map((p) => p.url).filter(Boolean);
		return () => urls.forEach((url) => URL.revokeObjectURL(url));
	});
</script>

{#if previews.length > 0}
	<ul class="flex gap-2 overflow-x-auto border-b border-fc-border p-2" aria-label="Attachments to send">
		{#each previews as preview (preview.file)}
			<li class="relative flex h-20 w-20 shrink-0 items-center justify-center overflow-hidden rounded-fc-md border border-fc-border bg-fc-surface">
				{#if preview.url}
					<img src={preview.url} alt={preview.file.name} class="size-full object-cover" />
				{:else}
					<span class="flex flex-col items-center gap-1 p-1 text-center text-fc-xs text-fc-fg-muted">
						<Icon icon={glyphs.file} size={20} />
						<span class="line-clamp-2 break-all">{preview.file.name}</span>
					</span>
				{/if}
				<button
					type="button"
					class="absolute right-0.5 top-0.5 flex size-6 items-center justify-center rounded-fc-pill bg-fc-scrim text-white"
					aria-label="Remove {preview.file.name}"
					onclick={() => onremove(preview.file)}
				>
					<Icon icon={icons.close} size={14} />
				</button>
			</li>
		{/each}
	</ul>
{/if}
