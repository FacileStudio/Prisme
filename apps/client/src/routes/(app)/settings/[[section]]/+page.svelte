<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import {
		NotificationsSection,
		PreferencesSection,
		ProfileSection,
		SettingsLayout,
		SettingsRow,
		SettingsSection,
		Skeleton,
		Switch,
		icons,
		orderSections,
		sectionHref,
		settingsRoute
	} from '@facile/muse';
	import { prisme, type Agent } from '$lib/api/prisme';
	import { t, rememberLocale } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { handle } from '$lib/chat/users';
	import { notify } from '$lib/state/notify.svelte';
	import { theme, type ThemePreference } from '$lib/state/theme.svelte';
	import Avatar from '$lib/components/common/Avatar.svelte';

	const BASE = '/settings';

	const user = $derived(app.session?.user);
	const name = $derived(user ? user.name || user.email.split('@')[0] : '');

	const sections = $derived([
		{ id: 'profile', label: t('settings.profile'), icon: icons.userCircle, href: sectionHref(BASE, 'profile') },
		{ id: 'preferences', label: t('settings.preferences'), icon: icons.palette, href: sectionHref(BASE, 'preferences') },
		{ id: 'notifications', label: t('settings.notifications'), icon: icons.notification, href: sectionHref(BASE, 'notifications') },
		{ id: 'agents', label: t('settings.agents'), icon: icons.usersGroup, href: sectionHref(BASE, 'agents') }
	]);

	const requested = $derived(page.params.section);
	const route = $derived(settingsRoute(requested, orderSections(sections)));

	$effect(() => {
		if (requested === 'appearance') void goto(sectionHref(BASE, 'preferences'), { replaceState: true });
		else if (!route.known) void goto(BASE, { replaceState: true });
	});

	let agents = $state<Agent[] | null>(null);
	let agentsError = $state('');

	const notifyHint = $derived.by(() => {
		if (!notify.supported) return t('settings.notifyUnsupported');
		if (notify.permission === 'denied') return t('settings.notifyBlocked');
		return t('settings.notifyHint');
	});

	onMount(() => {
		prisme.agents().then(
			(list) => (agents = list),
			(err) => {
				agents = [];
				agentsError = err instanceof Error ? err.message : t('settings.agentsFailed');
			}
		);
	});
</script>

<svelte:head>
	<title>{t('settings.title')}</title>
</svelte:head>

<div class="mx-auto w-full max-w-5xl px-4 py-6 md:px-6 md:py-10">
	<SettingsLayout
		{sections}
		active={requested}
		backHref={BASE}
		title={t('settings.heading')}
		description={t('settings.subtitle')}
	>
		{#snippet children(id)}
			{#if id === 'profile' && user}
				<ProfileSection {name} email={user.email} avatar={user.avatar_url || undefined}>
					<p class="text-fc-sm text-fc-fg-muted">
						{t('settings.matrixId')}: <span class="font-fc-mono">{user.matrix_user_id}</span>
					</p>
				</ProfileSection>
			{:else if id === 'preferences'}
				<PreferencesSection
					theme={theme.preference}
					onThemeChange={(next) => theme.set(next as ThemePreference)}
					onLocaleChange={rememberLocale}
				/>
			{:else if id === 'notifications'}
				<NotificationsSection>
					<SettingsRow label={t('settings.browserNotifications')} description={notifyHint}>
						<Switch
							aria-label={t('settings.browserNotifications')}
							checked={notify.enabled}
							disabled={!notify.supported || notify.permission === 'denied'}
							onchange={(e) => notify.setEnabled(e.currentTarget.checked)}
						/>
					</SettingsRow>
				</NotificationsSection>
			{:else if id === 'agents'}
				<SettingsSection title={t('settings.agents')} description={t('settings.agentsBody')}>
					{#if agents === null}
						<Skeleton class="h-11 w-full" />
					{:else if agentsError}
						<p class="text-fc-sm text-fc-danger">{agentsError}</p>
					{:else if agents.length === 0}
						<p class="text-fc-sm text-fc-fg-muted">{t('settings.agentsNone')}</p>
					{:else}
						<ul class="flex flex-col gap-4">
							{#each agents as agent (agent.user_id)}
								<li class="flex min-h-11 items-center gap-3">
									<Avatar userId={agent.user_id} name={agent.name} size={40} />
									<div class="flex min-w-0 flex-1 flex-col">
										<span class="truncate text-fc-sm font-medium text-fc-fg">{agent.name}</span>
										<span class="truncate font-fc-mono text-fc-xs text-fc-fg-muted">{handle(agent.user_id)}</span>
									</div>
									<span class="shrink-0 whitespace-nowrap text-fc-xs text-fc-fg-muted">
										{t('settings.agentRooms', { count: agent.room_ids.length })}
									</span>
								</li>
							{/each}
						</ul>
					{/if}
				</SettingsSection>
			{/if}
		{/snippet}
	</SettingsLayout>
</div>
