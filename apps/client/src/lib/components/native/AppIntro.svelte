<script lang="ts">
	import { onMount } from 'svelte';
	import { Alert, Button, Icon, Spinner, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import { startAppLogin } from '$lib/native/signin';
	import { introSeen, markIntroSeen } from '$lib/native/intro';

	let { error = $bindable('') }: { error?: string } = $props();

	const slides = $derived([
		{ id: 'spaces', icon: glyphs.logo, title: t('intro.spaces.title'), body: t('intro.spaces.body') },
		{ id: 'people', icon: icons.usersGroup, title: t('intro.people.title'), body: t('intro.people.body') },
		{ id: 'calm', icon: icons.notification, title: t('intro.calls.title'), body: t('intro.calls.body') },
		{ id: 'sign-in', icon: glyphs.logo, title: t('intro.signIn.title'), body: t('intro.signIn.body') }
	]);

	let track = $state<HTMLElement | null>(null);
	let index = $state(0);
	let opening = $state(false);

	const last = $derived(index === slides.length - 1);

	function show(target: number, smooth = true) {
		const motion = smooth && !window.matchMedia('(prefers-reduced-motion: reduce)').matches;
		track?.scrollTo({ left: target * track.clientWidth, behavior: motion ? 'smooth' : 'auto' });
	}

	function onscroll() {
		if (!track) return;
		index = Math.round(track.scrollLeft / track.clientWidth);
		if (index === slides.length - 1) markIntroSeen();
	}

	async function signIn() {
		if (opening) return;
		error = '';
		opening = true;
		try {
			await startAppLogin();
		} catch {
			error = t('app.signInOpenFailed');
		} finally {
			opening = false;
		}
	}

	onMount(() => {
		if (introSeen() || error) requestAnimationFrame(() => show(slides.length - 1, false));
	});
</script>

<div class="fixed inset-0 flex flex-col overflow-hidden bg-fc-page text-fc-fg">
	<div class="ambient" aria-hidden="true">
		<span class="orb orb-a"></span>
		<span class="orb orb-b"></span>
	</div>

	<header class="relative z-10 flex h-14 items-center justify-between px-5">
		<div class="flex items-center gap-2">
			<Icon icon={glyphs.logo} size={22} />
			<span class="font-fc-title text-fc-md font-semibold tracking-tight">Prisme</span>
		</div>
		{#if !last}
			<Button variant="ghost" size="sm" class="min-h-11 px-3" onclick={() => show(slides.length - 1)}>{t('app.skip')}</Button>
		{/if}
	</header>

	<div
		bind:this={track}
		{onscroll}
		class="relative z-10 flex min-h-0 flex-1 snap-x snap-mandatory overflow-x-auto overscroll-x-contain [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
		aria-roledescription="carousel"
		aria-label={t('app.slides')}
	>
		{#each slides as slide, i (slide.id)}
			<section
				class="flex w-full shrink-0 snap-center flex-col items-center justify-center gap-8 px-8 text-center"
				aria-roledescription="slide"
				aria-label="{i + 1} / {slides.length}"
				aria-hidden={i !== index}
			>
				<div class="flex size-36 items-center justify-center rounded-[2rem] bg-fc-component shadow-sm">
					<Icon icon={slide.icon} size={64} />
				</div>
				<div class="flex max-w-xs flex-col gap-3">
					<h1 class="font-fc-title text-fc-3xl leading-tight font-semibold tracking-tight">{slide.title}</h1>
					<p class="text-fc-md leading-relaxed text-fc-fg-muted">{slide.body}</p>
				</div>
			</section>
		{/each}
	</div>

	<footer class="relative z-10 flex flex-col items-center gap-4 px-6 pt-4 pb-[max(1.5rem,env(safe-area-inset-bottom))]">
		<div class="flex items-center gap-1" role="tablist" aria-label={t('app.slides')}>
			{#each slides as slide, i (slide.id)}
				<button
					type="button"
					role="tab"
					aria-selected={i === index}
					aria-label={t('app.goToSlide', { n: i + 1 })}
					class="flex size-11 items-center justify-center"
					onclick={() => show(i)}
				>
					<span
						class="h-1.5 rounded-fc-pill transition-all duration-300 motion-reduce:transition-none {i === index
							? 'w-6 bg-fc-fg'
							: 'w-1.5 bg-fc-fg/25'}"
					></span>
				</button>
			{/each}
		</div>
		{#if last}
			{#if error}
				<div class="w-full max-w-sm" role="alert"><Alert tone="danger">{error}</Alert></div>
			{/if}
			<Button size="lg" class="h-12 w-full max-w-sm rounded-fc-md" icon={opening ? undefined : icons.key} disabled={opening} onclick={signIn}>
				{#if opening}<Spinner size="sm" />{/if}
				{t('app.continue')}
			</Button>
			<p class="max-w-xs text-center text-fc-xs text-fc-fg-muted">{t('app.browserHint')}</p>
		{:else}
			<Button size="lg" class="h-12 w-full max-w-sm rounded-fc-md" iconRight={icons.arrow} onclick={() => show(index + 1)}>
				{t('app.next')}
			</Button>
		{/if}
	</footer>
</div>

<style>
	.ambient {
		position: absolute;
		inset: 0;
		overflow: hidden;
		pointer-events: none;
	}

	.orb {
		position: absolute;
		aspect-ratio: 1;
		border-radius: 9999px;
		filter: blur(90px);
		opacity: 0.12;
		background-color: var(--color-fc-accent);
		animation: drift 26s var(--ease-fc) infinite alternate;
	}

	.orb-a {
		width: min(70vw, 520px);
		top: -14%;
		left: -16%;
	}

	.orb-b {
		width: min(60vw, 440px);
		bottom: -16%;
		right: -14%;
		opacity: 0.08;
		animation-duration: 32s;
		animation-delay: -8s;
	}

	@keyframes drift {
		to {
			transform: translate3d(6vw, 5vh, 0) scale(1.12);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.orb {
			animation: none;
		}
	}
</style>
