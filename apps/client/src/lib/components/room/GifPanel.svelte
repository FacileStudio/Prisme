<script lang="ts">
	import { tick } from 'svelte';
	import { Button, Icon, Spinner, icons, toast } from '@facile/muse';
	import { pop, rise, stagger } from '$lib/motion';
	import { prisme, type Gif } from '$lib/api/prisme';
	import { locale, t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { sendFile } from '$lib/matrix/upload';
	import type { ComposerPanelProps } from '$lib/state/composer-menu';

	let { context, onclose }: ComposerPanelProps = $props();

	let query = $state('');
	let gifs = $state.raw<Gif[]>([]);
	let next = $state(0);
	let loading = $state(true);
	let failed = $state(false);
	let sending = $state('');
	let panel: HTMLDivElement | undefined = $state();
	let search: HTMLInputElement | undefined = $state();
	let generation = 0;

	$effect(() => {
		void tick().then(() => search?.focus());
		const outside = (event: PointerEvent) => {
			if (panel && !panel.contains(event.target as Node) && !sending) onclose();
		};
		document.addEventListener('pointerdown', outside, true);
		return () => document.removeEventListener('pointerdown', outside, true);
	});

	$effect(() => {
		const q = query.trim();
		const lang = locale();
		const run = ++generation;
		loading = true;
		const timer = setTimeout(() => void load(q, lang, 0, run), q ? 300 : 0);
		return () => clearTimeout(timer);
	});

	async function load(q: string, lang: string, offset: number, run: number) {
		try {
			const page = await prisme.gifs(q, lang, offset);
			if (run !== generation) return;
			gifs = offset === 0 ? page.gifs : [...gifs, ...page.gifs];
			next = page.gifs.length > 0 ? page.next : 0;
			failed = false;
		} catch {
			if (run === generation) failed = true;
		} finally {
			if (run === generation) loading = false;
		}
	}

	function more() {
		loading = true;
		void load(query.trim(), locale(), next, generation);
	}

	async function send(gif: Gif) {
		if (sending || !app.client) return;
		sending = gif.id;
		try {
			const blob = await prisme.gifMedia(gif.full.url);
			const type = blob.type || 'image/gif';
			const name = `${(gif.title || t('gif.untitled')).slice(0, 80)}.${type.split('/')[1] ?? 'gif'}`;
			await sendFile(app.client, context.room.roomId, new File([blob], name, { type }), context.threadId);
			onclose();
		} catch {
			toast.danger(t('gif.sendFailed'));
		} finally {
			sending = '';
		}
	}

	function onkeydown(event: KeyboardEvent) {
		if (event.key !== 'Escape') return;
		event.preventDefault();
		if (query) query = '';
		else onclose();
	}

	function preview(node: HTMLImageElement, path: string) {
		let url = '';
		const controller = new AbortController();
		const observer = new IntersectionObserver((entries) => {
			if (!entries.some((entry) => entry.isIntersecting)) return;
			observer.disconnect();
			prisme
				.gifMedia(path, controller.signal)
				.then((blob) => {
					url = URL.createObjectURL(blob);
					node.onload = () => node.classList.remove('opacity-0');
					node.src = url;
				})
				.catch(() => undefined);
		}, { rootMargin: '200px' });
		observer.observe(node);
		return {
			destroy() {
				observer.disconnect();
				controller.abort();
				if (url) URL.revokeObjectURL(url);
			}
		};
	}
</script>

<div
	bind:this={panel}
	role="dialog"
	aria-label={t('gif.title')}
	tabindex="-1"
	{onkeydown}
	transition:pop
	class="absolute bottom-full left-3 z-30 mb-2 flex h-[28rem] max-h-[70dvh] w-[min(24rem,calc(100%-1.5rem))] origin-bottom-left flex-col overflow-hidden rounded-fc-lg bg-fc-component shadow-lg md:left-4"
>
	<div class="flex items-center gap-1.5 p-2">
		<label class="flex h-10 min-w-0 flex-1 items-center gap-2 rounded-fc-pill bg-fc-field px-3 focus-within:outline-2 focus-within:outline-fc-ring">
			<Icon icon={icons.search} size={16} class="shrink-0 text-fc-fg-muted" />
			<input
				bind:this={search}
				bind:value={query}
				type="search"
				placeholder={t('gif.search')}
				aria-label={t('gif.search')}
				class="min-w-0 flex-1 bg-transparent text-fc-sm text-fc-fg outline-none placeholder:text-fc-fg-muted"
			/>
		</label>
		<button
			type="button"
			class="flex size-10 shrink-0 items-center justify-center rounded-fc-pill text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-fc-ring"
			aria-label={t('panel.close')}
			onclick={() => onclose()}
		>
			<Icon icon={icons.close} size={18} />
		</button>
	</div>
	<p class="px-3 pb-1 text-fc-xs font-medium text-fc-fg-muted">{query.trim() ? t('gif.title') : t('gif.trending')}</p>
	<div class="min-h-0 flex-1 overflow-y-auto overscroll-contain px-2 pb-2">
		{#if failed && gifs.length === 0}
			<p class="px-2 py-8 text-center text-fc-sm text-fc-fg-muted">{t('gif.failed')}</p>
		{:else if !loading && gifs.length === 0}
			<p class="px-2 py-8 text-center text-fc-sm text-fc-fg-muted">{t('gif.none')}</p>
		{:else}
			<div class="columns-2 gap-1.5 transition-opacity duration-200 motion-reduce:transition-none {loading && gifs.length > 0 ? 'opacity-60' : ''}">
				{#each gifs as gif, index (gif.id)}
					{@const title = gif.title || t('gif.untitled')}
					<button
						type="button"
						class="relative mb-1.5 block w-full break-inside-avoid overflow-hidden rounded-fc-md bg-fc-surface focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring disabled:opacity-60"
						style:aspect-ratio="{gif.preview.width || 1} / {gif.preview.height || 1}"
						aria-label={t('gif.send', { title })}
						title={title}
						disabled={!!sending}
						in:rise={{ delay: stagger(index % 24) }}
						onclick={() => send(gif)}
					>
						<img use:preview={gif.preview.url} alt="" class="size-full object-cover opacity-0 transition-opacity duration-300 ease-out motion-reduce:transition-none" />
						{#if sending === gif.id}
							<span class="absolute inset-0 flex items-center justify-center bg-fc-scrim"><Spinner /></span>
						{/if}
					</button>
				{/each}
			</div>
			{#if loading}
				<div class="flex justify-center py-3"><Spinner /></div>
			{:else if next > 0}
				<div class="flex justify-center pt-1">
					<Button size="sm" variant="ghost" onclick={more}>{t('gif.more')}</Button>
				</div>
			{/if}
		{/if}
	</div>
	<p class="shrink-0 px-3 py-1.5 text-right text-[0.6875rem] font-medium tracking-wide text-fc-fg-muted">{t('gif.powered')}</p>
</div>
