<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { DownloadSection } from '@facile/muse';
	import SiteShell from '$lib/components/landing/SiteShell.svelte';
	import { CHECKSUMS_URL, DOWNLOADS, RELEASES_URL } from '$lib/downloads';
	import { t } from '$lib/i18n/index.svelte';
	import { isApp } from '$lib/native/platform';

	onMount(() => {
		if (isApp()) void goto('/', { replaceState: true });
	});
</script>

<svelte:head>
	<title>{t('downloads.title')}</title>
	<meta name="description" content={t('downloads.description')} />
</svelte:head>

<SiteShell>
	<div class="mx-auto flex max-w-fc-lg flex-col gap-10 px-4 py-14 sm:px-6 sm:py-20">
		<section class="flex flex-col items-center gap-3 text-center">
			<h1 class="font-fc-title text-fc-3xl font-bold tracking-tight sm:text-5xl">{t('downloads.heading')}</h1>
			<p class="max-w-md text-fc-md text-fc-fg-muted">{t('downloads.intro')}</p>
		</section>
		<DownloadSection downloads={DOWNLOADS} releasesHref={RELEASES_URL} checksumsHref={CHECKSUMS_URL} />
	</div>
</SiteShell>
