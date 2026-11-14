<script lang="ts">
	import {
		DownloadSection,
		ForDevelopers,
		LANDING_HEADING,
		LandingCta,
		PricingPlans,
		WorksWith,
		type DeveloperRow,
		type WorksWithApp
	} from '@facile/muse';
	import { LOGIN_PATH } from '$lib/api/prisme';
	import { DOWNLOADS } from '$lib/downloads';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import BulleTile from './BulleTile.svelte';
	import CallTile from './CallTile.svelte';
	import HeroDemo from './HeroDemo.svelte';
	import SiteShell from './SiteShell.svelte';
	import SpacesTile from './SpacesTile.svelte';

	const WIDTH = 'mx-auto w-full max-w-5xl px-4 sm:px-6';

	const apps: WorksWithApp[] = $derived([
		{ name: 'Echo', icon: glyphs.echo, line: t('landing.with.echo'), href: 'https://echo.facile.studio' },
		{ name: 'Bulle', icon: glyphs.bulle, line: t('landing.with.bulle'), href: 'https://bulle.facile.studio' }
	]);
	const developer: DeveloperRow[] = $derived([
		{ icon: glyphs.chat, title: 'Matrix', text: t('landing.dev.matrix'), href: 'https://matrix.org' },
		{ kind: 'source', text: t('landing.dev.source'), href: 'https://github.com/FacileStudio/Prisme' }
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
		<LandingCta class="mt-8 w-full sm:w-auto" signInHref={LOGIN_PATH} align="center" data-sveltekit-reload />
		<p class="mt-4 text-fc-sm text-fc-fg-muted">{t('landing.hero.note')}</p>
		<div class="mt-12 w-full max-w-[960px] sm:mt-14">
			<HeroDemo />
		</div>
	</section>

	<section id="features" aria-labelledby="features-title" class="{WIDTH} scroll-mt-20 py-16 sm:py-20">
		<h2 id="features-title" class="text-center {LANDING_HEADING}">{t('landing.features.title')}</h2>
		<div class="mt-10 grid gap-4 md:grid-cols-3">
			<SpacesTile />
			<BulleTile />
			<CallTile />
		</div>
	</section>

	<WorksWith class="{WIDTH} py-16 sm:py-20" {apps} description={t('landing.with.line')} />

	<section id="download" aria-labelledby="download-title" class="mx-auto max-w-fc-lg scroll-mt-20 px-4 py-16 sm:px-6">
		<h2 id="download-title" class="mb-10 text-center {LANDING_HEADING}">{t('downloads.heading')}</h2>
		<DownloadSection downloads={DOWNLOADS} mode="device" allHref="/downloads" headingLevel={3} />
	</section>

	<PricingPlans
		id="pricing"
		class="{WIDTH} scroll-mt-20 py-16 sm:py-20"
		title={t('landing.pricing')}
		headingClass={LANDING_HEADING}
		highlight={['ai.access']}
		highlightOnly
		notes={{ '*': t('landing.plan.included'), pro: t('landing.plan.bulle'), studio: t('landing.plan.bulle') }}
		recommended="pro"
		suiteHref="https://suite.facile.studio/#tarifs"
	/>

	<div class="{WIDTH} py-16 sm:py-20">
		<ForDevelopers rows={developer} description={t('landing.dev.line')} selfHostHref="https://suite.facile.studio/auto-hebergement" />
	</div>

	<section class="{WIDTH} flex flex-col items-center py-16 text-center sm:py-24">
		<h2 class={LANDING_HEADING}>{t('landing.cta.title')}</h2>
		<p class="mt-3 text-fc-md text-fc-fg-muted">{t('landing.cta.body')}</p>
		<LandingCta class="mt-8 w-full sm:w-auto" signInHref={LOGIN_PATH} align="center" data-sveltekit-reload />
	</section>
</SiteShell>
