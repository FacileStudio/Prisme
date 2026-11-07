<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Badge, Button, Card, Icon, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import { RELEASES_URL, detectPlatform, ordered, type Platform } from '$lib/downloads';
	import { isApp } from '$lib/native/platform';

	const ICONS: Record<Platform, string> = {
		android: icons.android,
		macos: icons.apple,
		windows: icons.windows,
		linux: icons.linux
	};

	const COPY = $derived({
		android: { name: t('downloads.android'), detail: t('downloads.androidDetail') },
		macos: { name: t('downloads.macos'), detail: t('downloads.macosDetail') },
		windows: { name: t('downloads.windows'), detail: t('downloads.windowsDetail') },
		linux: { name: t('downloads.linux'), detail: t('downloads.linuxDetail') }
	});

	let platform = $state<Platform | null>(null);
	const builds = $derived(ordered(platform));

	onMount(() => {
		if (isApp()) void goto('/', { replaceState: true });
		platform = detectPlatform(navigator.userAgent);
	});
</script>

<svelte:head>
	<title>{t('downloads.title')}</title>
	<meta name="description" content={t('downloads.description')} />
</svelte:head>

<div class="h-dvh overflow-y-auto bg-fc-page text-fc-fg">
	<header class="mx-auto flex max-w-fc-lg items-center justify-between px-4 py-3 sm:px-6">
		<a href="/" class="flex min-h-11 items-center gap-2">
			<Icon icon={glyphs.logo} size={24} />
			<span class="font-fc-title text-fc-lg font-semibold tracking-tight">Prisme</span>
		</a>
		<Button href="/" variant="ghost" class="min-h-11">{t('downloads.back')}</Button>
	</header>

	<main class="mx-auto flex max-w-fc-md flex-col gap-8 px-4 py-10 sm:px-6 sm:py-14">
		<section class="flex flex-col items-center gap-3 text-center">
			<h1 class="font-fc-title text-fc-3xl font-semibold tracking-tight sm:text-5xl">{t('downloads.heading')}</h1>
			<p class="max-w-md text-fc-md text-fc-fg-muted">{t('downloads.intro')}</p>
		</section>

		<ul class="flex flex-col gap-3">
			{#each builds as build (build.platform)}
				<li>
					<Card class="flex flex-col gap-4 sm:flex-row sm:items-center">
						<div class="flex min-w-0 flex-1 items-start gap-3">
							<Icon icon={ICONS[build.platform]} size={28} />
							<div class="flex min-w-0 flex-col gap-1">
								<div class="flex flex-wrap items-center gap-2">
									<h2 class="font-fc-title text-fc-lg font-semibold">{COPY[build.platform].name}</h2>
									{#if build.platform === platform}
										<Badge tone="info">{t('downloads.recommended')}</Badge>
									{/if}
								</div>
								<p class="text-fc-sm text-fc-fg-muted">{COPY[build.platform].detail}</p>
								{#if build.extra}
									<a href={build.extra.url} class="text-fc-sm underline underline-offset-2">{t('downloads.deb')}</a>
								{/if}
							</div>
						</div>
						<Button href={build.url} icon={icons.download} variant={build.platform === platform ? 'primary' : 'secondary'} class="min-h-11 shrink-0">
							{t('downloads.download')}
						</Button>
					</Card>
				</li>
			{/each}
		</ul>

		<p class="text-center text-fc-sm">
			<a href={RELEASES_URL} class="text-fc-fg-muted underline underline-offset-2">{t('downloads.allReleases')}</a>
		</p>
	</main>
</div>
