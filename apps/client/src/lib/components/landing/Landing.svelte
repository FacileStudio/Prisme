<script lang="ts">
	import { Button, DownloadSection, Icon, SignInWithFacile, icons } from '@facile/muse';
	import { LOGIN_PATH } from '$lib/api/prisme';
	import { DOWNLOADS } from '$lib/downloads';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import HeroDemo from './HeroDemo.svelte';
	import SiteShell from './SiteShell.svelte';

	const features = $derived([
		{ icon: glyphs.users, title: t('landing.features.spaces.title'), body: t('landing.features.spaces.body') },
		{ icon: glyphs.mention, title: t('landing.features.agents.title'), body: t('landing.features.agents.body') },
		{ icon: glyphs.lock, title: t('landing.features.calls.title'), body: t('landing.features.calls.body') },
		{ icon: glyphs.chat, title: t('landing.features.matrix.title'), body: t('landing.features.matrix.body') }
	]);
</script>

<svelte:head>
	<title>{t('landing.title')}</title>
	<meta name="description" content={t('landing.description')} />
</svelte:head>

<SiteShell>
	<section class="mx-auto flex max-w-fc-xl flex-col items-center px-4 pt-16 pb-12 text-center sm:px-6 sm:pt-24">
		<h1 class="max-w-3xl font-fc-title text-fc-3xl font-bold tracking-tight text-balance sm:text-5xl">{t('landing.hero.title')}</h1>
		<p class="mt-5 max-w-xl text-fc-md text-balance text-fc-fg-muted">{t('landing.hero.body')}</p>
		<div class="mt-8 flex flex-wrap justify-center gap-3">
			<SignInWithFacile href={LOGIN_PATH} />
			<Button href="#download" variant="secondary" size="lg" icon={icons.download}>{t('downloads.menu')}</Button>
		</div>
		<div class="mt-14 w-full max-w-[960px]">
			<HeroDemo />
		</div>
	</section>

	<section id="features" aria-labelledby="features-title" class="mx-auto max-w-fc-xl scroll-mt-20 px-4 py-16 sm:px-6">
		<h2 id="features-title" class="text-center font-fc-title text-fc-2xl font-bold tracking-tight">{t('landing.features.title')}</h2>
		<ul class="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
			{#each features as feature (feature.icon)}
				<li class="flex flex-col gap-2 rounded-fc-lg bg-fc-surface p-6">
					<span class="mb-2 flex size-10 items-center justify-center rounded-fc-md bg-fc-component">
						<Icon icon={feature.icon} size={20} />
					</span>
					<h3 class="text-fc-lg font-bold">{feature.title}</h3>
					<p class="text-fc-sm text-fc-fg-muted">{feature.body}</p>
				</li>
			{/each}
		</ul>
	</section>

	<section id="download" aria-labelledby="download-title" class="mx-auto max-w-fc-lg scroll-mt-20 px-4 py-16 sm:px-6">
		<h2 id="download-title" class="mb-10 text-center font-fc-title text-fc-2xl font-bold tracking-tight">{t('downloads.heading')}</h2>
		<DownloadSection downloads={DOWNLOADS} mode="device" allHref="/downloads" headingLevel={3} />
	</section>

	<section class="mx-auto flex max-w-fc-xl flex-col items-center px-4 py-20 text-center sm:px-6">
		<h2 class="font-fc-title text-fc-2xl font-bold tracking-tight">{t('landing.cta.title')}</h2>
		<p class="mt-3 text-fc-md text-fc-fg-muted">{t('landing.cta.body')}</p>
		<div class="mt-8"><SignInWithFacile href={LOGIN_PATH} size="lg" /></div>
	</section>
</SiteShell>
