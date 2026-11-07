<script lang="ts">
	import { t } from '$lib/i18n/index.svelte';
	import { Avatar as Picture, twMerge } from '@facile/muse';
	import { Blob } from '@facile/blob';
	import { app } from '$lib/state/app.svelte';
	import { cachedMedia, loadMedia } from '$lib/matrix/media';

	type Size = 20 | 24 | 28 | 32 | 36 | 40 | 56;
	type Props = {
		userId?: string;
		name: string;
		mxc?: string | null;
		src?: string;
		size?: Size;
		presence?: boolean;
		class?: string;
	};
	let { userId = '', name, mxc = null, src = '', size = 32, presence = false, class: className = '' }: Props = $props();

	const BOX: Record<Size, [string, string]> = {
		20: ['size-5', 'text-[9px]'],
		24: ['size-6', 'text-[10px]'],
		28: ['size-7', 'text-fc-xs'],
		32: ['size-8', 'text-fc-xs'],
		36: ['size-9', 'text-fc-sm'],
		40: ['size-10', 'text-fc-sm'],
		56: ['size-14', 'text-fc-md']
	};

	const user = $derived.by(() => {
		void app.version;
		return userId ? (app.client?.getUser(userId) ?? null) : null;
	});
	const agent = $derived(userId && app.isAgent(userId) ? (app.agents.get(userId)?.name ?? name) : null);
	const source = $derived(agent ? null : mxc || user?.avatarUrl || null);
	const thumb = $derived(size > 40 ? 192 : 96);
	let loaded = $state<{ mxc: string; url: string } | null>(null);
	let failed = $state('');
	const url = $derived(source ? (loaded?.mxc === source ? loaded.url : cachedMedia(source, thumb)) : undefined);
	const waiting = $derived(!!source && !url && failed !== source);
	const online = $derived(user?.presence === 'online');
	const box = $derived(twMerge('relative inline-flex shrink-0', BOX[size][0], className));

	$effect(() => {
		const wanted = source;
		if (!wanted || cachedMedia(wanted, thumb)) return;
		let live = true;
		loadMedia(wanted, thumb).then(
			(blobUrl) => live && (loaded = { mxc: wanted, url: blobUrl }),
			() => live && (failed = wanted)
		);
		return () => (live = false);
	});
</script>

<span class={box}>
	{#if agent}
		<Blob name={agent} size={size * 2} background="circle" alt={name} class="size-full" />
	{:else if waiting}
		<span class="size-full rounded-fc-pill bg-fc-surface" role="img" aria-label={name}></span>
	{:else}
		<Picture src={url || src || undefined} {name} class="size-full {BOX[size][1]}" />
	{/if}
	{#if presence}
		<span
			class="absolute -bottom-0.5 -right-0.5 size-2.5 rounded-fc-pill {online ? 'bg-fc-success' : 'bg-fc-fg-muted'}"
			role="img"
			aria-label={online ? t('avatar.online') : t('avatar.offline')}
		></span>
	{/if}
</span>
