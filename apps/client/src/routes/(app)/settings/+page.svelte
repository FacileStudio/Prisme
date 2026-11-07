<script lang="ts">
	import { onMount } from 'svelte';
	import { Button, OptionCards, ProfileCard, SettingsRow, SettingsSection, Skeleton, Switch, icons } from '@facile/muse';
	import { prisme, type Agent } from '$lib/api/prisme';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { handle } from '$lib/chat/users';
	import { notify } from '$lib/state/notify.svelte';
	import { theme, type ThemePreference } from '$lib/state/theme.svelte';
	import Avatar from '$lib/components/common/Avatar.svelte';

	const modes = $derived([
		{ value: 'system', label: t('settings.modeSystem'), icon: icons.monitor },
		{ value: 'light', label: t('settings.modeLight'), icon: icons.sun },
		{ value: 'dark', label: t('settings.modeDark'), icon: icons.moon }
	]);

	const user = $derived(app.session?.user);
	const name = $derived(user ? user.name || user.email.split('@')[0] : '');

	let agents = $state<Agent[] | null>(null);
	let agentsError = $state('');
	let signingOut = $state(false);

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

	async function signOut() {
		signingOut = true;
		await app.logout();
	}
</script>

<svelte:head>
	<title>{t('settings.title')}</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-3xl flex-col gap-10 px-4 py-6 md:px-6 md:py-10">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">{t('settings.heading')}</h1>
		<p class="text-fc-sm text-fc-fg-muted">{t('settings.subtitle')}</p>
	</div>

	{#if user}
		<SettingsSection title={t('settings.profile')} description={t('settings.profileBody')} bare>
			<ProfileCard
				{name}
				email={user.email}
				avatar={user.avatar_url || undefined}
				meta={[{ label: t('settings.matrixId'), value: user.matrix_user_id }]}
			/>
		</SettingsSection>
	{/if}

	<SettingsSection title={t('settings.appearance')} description={t('settings.appearanceBody')}>
		<SettingsRow label={t('settings.scheme')} description={t('settings.schemeBody')} stacked>
			<OptionCards
				options={modes}
				value={theme.preference}
				name="theme-mode"
				label={t('settings.scheme')}
				onSelect={(value) => theme.set(value as ThemePreference)}
			/>
		</SettingsRow>
	</SettingsSection>

	<SettingsSection title={t('settings.notifications')}>
		<SettingsRow label={t('settings.browserNotifications')} description={notifyHint}>
			<Switch
				aria-label={t('settings.browserNotifications')}
				checked={notify.enabled}
				disabled={!notify.supported || notify.permission === 'denied'}
				onchange={(e) => notify.setEnabled(e.currentTarget.checked)}
			/>
		</SettingsRow>
	</SettingsSection>

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
						<span class="shrink-0 text-fc-xs text-fc-fg-muted">
							{t('settings.agentRooms', { count: agent.room_ids.length })}
						</span>
					</li>
				{/each}
			</ul>
		{/if}
	</SettingsSection>

	<SettingsSection title={t('settings.session')} bare>
		<div>
			<Button size="lg" variant="ghost-danger" icon={icons.logout} disabled={signingOut} onclick={signOut}>{t('settings.signOut')}</Button>
		</div>
	</SettingsSection>
</div>
