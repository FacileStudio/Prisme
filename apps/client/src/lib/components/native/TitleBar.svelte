<script lang="ts">
	import { getCurrentWindow } from '@tauri-apps/api/window';
	import { Tooltip } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';
	import { hasNativeTrafficLights } from '$lib/native/platform';

	const drawn = !hasNativeTrafficLights();
	const win = getCurrentWindow();

	const controls = $derived([
		{ label: t('app.windowClose'), color: 'bg-[#ff5f57]', glyph: 'M3.5 3.5l5 5m0-5l-5 5', run: () => win.close() },
		{ label: t('app.windowMinimize'), color: 'bg-[#febc2e]', glyph: 'M3 6h6', run: () => win.minimize() },
		{ label: t('app.windowMaximize'), color: 'bg-[#28c840]', glyph: 'M3.5 6h5M6 3.5v5', run: () => win.toggleMaximize() }
	]);
</script>

<div data-tauri-drag-region class="fixed inset-x-0 top-0 z-[60] flex h-[var(--titlebar-h)] items-center bg-fc-page px-3.5 select-none">
	{#if drawn}
		<div class="group flex items-center gap-2" role="group" aria-label={t('app.window')}>
			{#each controls as control (control.glyph)}
				<Tooltip text={control.label} describe={false} class="contents">
					<button
						type="button"
						aria-label={control.label}
						onclick={() => void control.run()}
						class="flex size-3 items-center justify-center rounded-full {control.color} ring-1 ring-black/10 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fc-ring"
					>
						<svg viewBox="0 0 12 12" class="size-2 opacity-0 transition-opacity group-hover:opacity-100" aria-hidden="true">
							<path d={control.glyph} stroke="rgb(0 0 0 / 0.6)" stroke-width="1.3" stroke-linecap="round" fill="none" />
						</svg>
					</button>
				</Tooltip>
			{/each}
		</div>
	{/if}
</div>
