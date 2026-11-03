<script lang="ts">
	import { onMount } from 'svelte';
	import { Button, OptionCards, ProfileCard, SettingsRow, SettingsSection, Skeleton, Switch, icons } from '@facile/muse';
	import { prisme, type Agent } from '$lib/api/prisme';
	import { app } from '$lib/state/app.svelte';
	import { handle } from '$lib/chat/users';
	import { notify } from '$lib/state/notify.svelte';
	import { theme, type ThemePreference } from '$lib/state/theme.svelte';
	import Avatar from '$lib/components/common/Avatar.svelte';

	const modes = [
		{ value: 'system', label: 'System', icon: icons.monitor },
		{ value: 'light', label: 'Light', icon: icons.sun },
		{ value: 'dark', label: 'Dark', icon: icons.moon }
	];

	const user = $derived(app.session?.user);
	const name = $derived(user ? user.name || user.email.split('@')[0] : '');

	let agents = $state<Agent[] | null>(null);
	let agentsError = $state('');
	let signingOut = $state(false);

	const notifyHint = $derived.by(() => {
		if (!notify.supported) return 'Not supported by this browser.';
		if (notify.permission === 'denied') return 'Blocked in your browser settings.';
		return 'For mentions and direct messages while Prisme is in the background.';
	});

	onMount(() => {
		prisme.agents().then(
			(list) => (agents = list),
			(err) => {
				agents = [];
				agentsError = err instanceof Error ? err.message : 'Could not load your agents.';
			}
		);
	});

	async function signOut() {
		signingOut = true;
		await app.logout();
	}
</script>

<svelte:head>
	<title>Settings · Prisme</title>
</svelte:head>

<div class="mx-auto flex w-full max-w-3xl flex-col gap-10 px-4 py-6 md:px-6 md:py-10">
	<div class="flex flex-col gap-2">
		<h1 class="text-fc-2xl font-semibold text-fc-fg">Settings</h1>
		<p class="text-fc-sm text-fc-fg-muted">You, this browser, and the agents that chat on your behalf.</p>
	</div>

	{#if user}
		<SettingsSection title="Profile" description="Your name and picture come from your Facile account." bare>
			<ProfileCard
				{name}
				email={user.email}
				avatar={user.avatar_url || undefined}
				meta={[{ label: 'Matrix ID', value: user.matrix_user_id }]}
			/>
		</SettingsSection>
	{/if}

	<SettingsSection title="Appearance" description="Stored in this browser. Every device keeps its own choice.">
		<SettingsRow label="Colour scheme" description="System follows your operating system while Prisme is open." stacked>
			<OptionCards
				options={modes}
				value={theme.preference}
				name="theme-mode"
				label="Colour scheme"
				onSelect={(value) => theme.set(value as ThemePreference)}
			/>
		</SettingsRow>
	</SettingsSection>

	<SettingsSection title="Notifications">
		<SettingsRow label="Browser notifications" description={notifyHint}>
			<Switch
				aria-label="Browser notifications"
				checked={notify.enabled}
				disabled={!notify.supported || notify.permission === 'denied'}
				onchange={(e) => notify.setEnabled(e.currentTarget.checked)}
			/>
		</SettingsRow>
	</SettingsSection>

	<SettingsSection title="Agents" description="Bots you own. Bulle creates them and adds them to your spaces.">
		{#if agents === null}
			<Skeleton class="h-11 w-full" />
		{:else if agentsError}
			<p class="text-fc-sm text-fc-danger">{agentsError}</p>
		{:else if agents.length === 0}
			<p class="text-fc-sm text-fc-fg-muted">You have no agents yet.</p>
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
							{agent.room_ids.length === 1 ? 'In 1 room' : `In ${agent.room_ids.length} rooms`}
						</span>
					</li>
				{/each}
			</ul>
		{/if}
	</SettingsSection>

	<SettingsSection title="Session" bare>
		<div>
			<Button size="lg" variant="ghost-danger" icon={icons.logout} disabled={signingOut} onclick={signOut}>Sign out</Button>
		</div>
	</SettingsSection>
</div>
