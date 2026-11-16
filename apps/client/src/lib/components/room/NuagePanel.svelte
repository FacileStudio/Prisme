<script lang="ts">
	import { tick } from 'svelte';
	import { cubicOut } from 'svelte/easing';
	import { Button, Icon, Skeleton, icons, prefersReducedMotion, toast } from '@facile/muse';
	import { pop, rise, stagger } from '$lib/motion';
	import { ApiError } from '$lib/api/http';
	import { nuage, type NuageFile, type NuageFolder, type NuagePlace, type NuageSpace } from '$lib/api/nuage';
	import { fileGlyph, fileSize, nuageBlock } from '$lib/chat/nuage';
	import { sendNuage } from '$lib/matrix/messages';
	import { app } from '$lib/state/app.svelte';
	import type { ComposerPanelProps } from '$lib/state/composer-menu';
	import { t } from '$lib/i18n/index.svelte';

	type Crumb = { label: string; place: NuagePlace };
	type View = { folders: NuageFolder[]; files: NuageFile[]; spaces: NuageSpace[] };

	let { context, onclose }: ComposerPanelProps = $props();

	let crumbs = $state<Crumb[]>([{ label: t('nuage.mine'), place: {} }]);
	let direction = $state(1);
	let query = $state('');
	let view = $state.raw<View>({ folders: [], files: [], spaces: [] });
	let loading = $state(true);
	let problem = $state('');
	let sending = $state(0);
	let panel: HTMLDivElement | undefined = $state();
	let search: HTMLInputElement | undefined = $state();
	let generation = 0;

	const here = $derived(crumbs[crumbs.length - 1]);
	const searching = $derived(query.trim().length > 0);
	const level = $derived(searching ? `search:${query.trim()}` : `${crumbs.length}:${JSON.stringify(here.place)}`);
	const empty = $derived(!loading && !problem && view.folders.length + view.files.length + view.spaces.length === 0);

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
		const place = here.place;
		const root = crumbs.length === 1;
		const run = ++generation;
		loading = true;
		const timer = setTimeout(() => void load(q, place, root, run), q ? 250 : 0);
		return () => clearTimeout(timer);
	});

	async function fetchView(q: string, place: NuagePlace, root: boolean): Promise<View> {
		if (q) return { folders: [], spaces: [], files: (await nuage.search(q)).files };
		const [listing, spaces] = await Promise.all([nuage.list(place), root ? nuage.spaces() : Promise.resolve({ spaces: [] })]);
		return { ...listing, spaces: spaces.spaces };
	}

	async function load(q: string, place: NuagePlace, root: boolean, run: number) {
		try {
			const next = await fetchView(q, place, root);
			if (run !== generation) return;
			view = next;
			problem = '';
		} catch (err) {
			if (run !== generation) return;
			problem = err instanceof ApiError && err.code === 'failed_precondition' ? t('nuage.reconnect') : t('nuage.loadFailed');
		} finally {
			if (run === generation) loading = false;
		}
	}

	function open(label: string, place: NuagePlace) {
		direction = 1;
		crumbs = [...crumbs, { label, place }];
	}

	function back(to = crumbs.length - 2) {
		if (to < 0) return;
		direction = -1;
		crumbs = crumbs.slice(0, to + 1);
	}

	async function share(file: NuageFile) {
		if (sending || !app.client) return;
		sending = file.id;
		try {
			const info = nuageBlock(await nuage.share(file.id));
			if (!info) throw new Error('not a Nuage share link');
			const last = context.threadId ? context.room.getThread(context.threadId)?.lastReply()?.getId() : undefined;
			await sendNuage(app.client, { roomId: context.room.roomId, threadId: context.threadId }, info, last);
			onclose();
		} catch {
			toast.danger(t('nuage.shareFailed'));
		} finally {
			sending = 0;
		}
	}

	function onkeydown(event: KeyboardEvent) {
		if (event.key !== 'Escape') return;
		event.preventDefault();
		if (query) query = '';
		else if (crumbs.length > 1) back();
		else onclose();
	}

	function glide(_node: Element, { dir = 1 } = {}) {
		const still = prefersReducedMotion();
		return {
			duration: still ? 0 : 200,
			easing: cubicOut,
			css: (time: number, rest: number) => `opacity: ${time}; transform: translateX(${rest * dir * 16}px);`
		};
	}

	const ROW =
		'flex min-h-11 w-full items-center gap-3 rounded-fc-md px-2 text-left text-fc-sm transition-colors duration-150 ease-fc hover:bg-fc-surface focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-fc-ring disabled:opacity-60 motion-reduce:transition-none';
</script>

<div
	bind:this={panel}
	role="dialog"
	aria-label={t('nuage.title')}
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
				placeholder={t('nuage.search')}
				aria-label={t('nuage.search')}
				class="min-w-0 flex-1 bg-transparent text-fc-sm text-fc-fg outline-none placeholder:text-fc-fg-muted"
			/>
		</label>
		<button
			type="button"
			class="flex size-10 shrink-0 items-center justify-center rounded-fc-pill text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-fc-ring"
			aria-label={t('nuage.close')}
			onclick={() => onclose()}
		>
			<Icon icon={icons.close} size={18} />
		</button>
	</div>

	{#if !searching}
		<nav aria-label={t('nuage.title')} class="flex min-h-8 items-center gap-1 overflow-x-auto px-2 pb-1 text-fc-xs">
			{#if crumbs.length > 1}
				<button type="button" class="flex size-7 shrink-0 items-center justify-center rounded-fc-pill text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-fc-ring" aria-label={t('nuage.back')} onclick={() => back()}>
					<Icon icon={icons.chevronLeft} size={14} />
				</button>
			{/if}
			{#each crumbs as crumb, i (i)}
				{#if i > 0}<span class="text-fc-fg-muted" aria-hidden="true">/</span>{/if}
				{#if i === crumbs.length - 1}
					<span class="truncate px-1 font-medium text-fc-fg" aria-current="location">{crumb.label}</span>
				{:else}
					<button type="button" class="truncate rounded-fc-sm px-1 text-fc-fg-muted hover:text-fc-fg focus-visible:outline-2 focus-visible:outline-fc-ring" onclick={() => back(i)}>{crumb.label}</button>
				{/if}
			{/each}
		</nav>
	{/if}

	<div class="relative min-h-0 flex-1 overflow-hidden">
		{#key level}
			<div class="absolute inset-0 overflow-y-auto overscroll-contain px-2 pb-2" in:glide={{ dir: direction }} out:glide={{ dir: -direction }}>
				{#if loading}
					<div class="flex flex-col gap-1 pt-1" aria-busy="true">
						{#each { length: 6 } as _, i (i)}
							<div class="flex min-h-11 items-center gap-3 px-2">
								<Skeleton class="size-8 rounded-fc-md" />
								<Skeleton class="h-3.5 flex-1" />
							</div>
						{/each}
					</div>
				{:else if problem}
					<div class="flex h-full flex-col items-center justify-center gap-3 px-6 text-center">
						<p class="text-fc-sm text-fc-fg-muted">{problem}</p>
						<Button size="sm" variant="secondary" icon={icons.refresh} onclick={() => (crumbs = [...crumbs])}>{t('nuage.retry')}</Button>
					</div>
				{:else if empty}
					<div class="flex h-full flex-col items-center justify-center gap-2 px-6 text-center">
						<span class="flex size-12 items-center justify-center rounded-fc-pill bg-fc-surface text-fc-fg-muted">
							<Icon icon={searching ? icons.search : icons.folder} size={22} />
						</span>
						<p class="text-fc-sm text-fc-fg-muted">{searching ? t('nuage.noMatch') : t('nuage.empty')}</p>
					</div>
				{:else}
					<ul class="flex flex-col">
						{#each view.folders as folder, i (`f${folder.id}`)}
							<li in:rise={{ delay: stagger(i) }}>
								<button type="button" class={ROW} aria-label={t('nuage.openFolder', { name: folder.name })} onclick={() => open(folder.name, { folderId: folder.id, spaceId: here.place.spaceId })}>
									<span class="flex size-8 shrink-0 items-center justify-center rounded-fc-md bg-fc-accent/10 text-fc-accent"><Icon icon={icons.folder} size={18} /></span>
									<span class="min-w-0 flex-1 truncate font-medium text-fc-fg">{folder.name}</span>
								</button>
							</li>
						{/each}
						{#each view.files as file, i (`d${file.id}`)}
							<li in:rise={{ delay: stagger(view.folders.length + i) }}>
								<button type="button" class={ROW} disabled={!!sending} aria-label={t('nuage.share', { name: file.name })} onclick={() => share(file)}>
									<span class="flex size-8 shrink-0 items-center justify-center rounded-fc-md bg-fc-surface text-fc-fg-muted"><Icon icon={fileGlyph(file.mime_type)} size={18} /></span>
									<span class="flex min-w-0 flex-1 flex-col">
										<span class="truncate text-fc-fg">{file.name}</span>
										<span class="text-fc-xs text-fc-fg-muted">{fileSize(file.size)}</span>
									</span>
									{#if sending === file.id}<Icon icon={icons.refresh} size={16} class="shrink-0 animate-spin text-fc-fg-muted motion-reduce:animate-none" />{/if}
								</button>
							</li>
						{/each}
					</ul>
					{#if view.spaces.length > 0}
						<p class="px-2 pt-3 pb-1 text-fc-xs font-medium text-fc-fg-muted">{t('nuage.spaces')}</p>
						<ul class="flex flex-col">
							{#each view.spaces as space, i (`s${space.id}`)}
								<li in:rise={{ delay: stagger(view.folders.length + view.files.length + i) }}>
									<button type="button" class={ROW} aria-label={t('nuage.openFolder', { name: space.name })} onclick={() => open(space.name, { spaceId: space.id })}>
										<span class="flex size-8 shrink-0 items-center justify-center rounded-fc-md bg-fc-accent/10 text-fc-accent"><Icon icon={icons.usersGroup} size={18} /></span>
										<span class="min-w-0 flex-1 truncate font-medium text-fc-fg">{space.name}</span>
									</button>
								</li>
							{/each}
						</ul>
					{/if}
				{/if}
			</div>
		{/key}
	</div>
</div>
