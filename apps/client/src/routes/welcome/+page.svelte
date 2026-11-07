<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Alert, Button, Icon, Spinner, icons } from '@facile/muse';
	import { prisme } from '$lib/api/prisme';
	import { t } from '$lib/i18n/index.svelte';
	import { markWelcomed, welcomed } from '$lib/native/intro';
	import { isApp } from '$lib/native/platform';
	import { notify } from '$lib/state/notify.svelte';

	let ready = $state(false);
	let asking = $state(false);

	async function open() {
		markWelcomed();
		await goto('/', { replaceState: true });
	}

	async function enable() {
		asking = true;
		await notify.setEnabled(true);
		asking = false;
		if (notify.enabled) await open();
	}

	onMount(async () => {
		const signedIn = await prisme.session().then(() => true, () => false);
		if (!signedIn) {
			await goto('/login', { replaceState: true });
			return;
		}
		if (!isApp() || welcomed() || notify.enabled) {
			await open();
			return;
		}
		ready = true;
	});
</script>

<svelte:head>
	<title>{t('welcome.title')} · Prisme</title>
</svelte:head>

{#if !ready}
	<main class="flex min-h-dvh items-center justify-center bg-fc-page"><Spinner /></main>
{:else}
	<main class="fixed inset-0 flex flex-col bg-fc-page text-fc-fg">
		<section class="flex min-h-0 flex-1 flex-col items-center justify-center gap-8 px-8 text-center">
			<div class="flex size-36 items-center justify-center rounded-[2rem] bg-fc-component shadow-sm">
				<Icon icon={notify.enabled ? icons.check : icons.notification} size={64} />
			</div>
			<div class="flex max-w-xs flex-col gap-3">
				<h1 class="font-fc-title text-fc-3xl leading-tight font-semibold tracking-tight">{t('welcome.title')}</h1>
				<p class="text-fc-md leading-relaxed text-fc-fg-muted">{t('welcome.body')}</p>
			</div>
			{#if notify.permission === 'denied'}
				<div class="w-full max-w-sm" role="alert"><Alert tone="warning">{t('welcome.blocked')}</Alert></div>
			{/if}
		</section>
		<footer class="flex flex-col items-center gap-2 px-6 pt-4 pb-[max(1.5rem,env(safe-area-inset-bottom))]">
			{#if notify.permission === 'denied'}
				<Button size="lg" class="h-12 w-full max-w-sm rounded-fc-md" onclick={open}>{t('welcome.open')}</Button>
			{:else}
				<Button size="lg" class="h-12 w-full max-w-sm rounded-fc-md" icon={asking ? undefined : icons.notification} disabled={asking} onclick={enable}>
					{#if asking}<Spinner size="sm" />{/if}
					{t('welcome.enable')}
				</Button>
				<Button variant="ghost" class="min-h-11 w-full max-w-sm" disabled={asking} onclick={open}>{t('welcome.later')}</Button>
			{/if}
		</footer>
	</main>
{/if}
