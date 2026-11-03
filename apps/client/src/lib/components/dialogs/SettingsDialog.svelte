<script lang="ts">
	import { Button, Modal, Switch, icons } from '@facile/muse';
	import { app } from '$lib/state/app.svelte';
	import { notify } from '$lib/state/notify.svelte';
	import { theme, type ThemePreference } from '$lib/state/theme.svelte';
	import MxcAvatar from '$lib/components/common/MxcAvatar.svelte';

	let { onclose }: { onclose: () => void } = $props();
	const user = $derived(app.session?.user);
	const themes: { value: ThemePreference; label: string }[] = [
		{ value: 'light', label: 'Light' },
		{ value: 'dark', label: 'Dark' },
		{ value: 'system', label: 'System' }
	];
</script>

<Modal open title="Settings" showClose onClose={onclose}>
	{#if user}
		<div class="flex items-center gap-3">
			<MxcAvatar src={user.avatar_url} name={user.name || user.email} size="lg" />
			<div class="flex min-w-0 flex-col">
				<span class="truncate font-semibold">{user.name || user.email.split('@')[0]}</span>
				<span class="truncate text-fc-sm text-fc-fg-muted">{user.email}</span>
				<span class="truncate font-fc-mono text-fc-xs text-fc-fg-muted">{user.matrix_user_id}</span>
			</div>
		</div>
		<p class="text-fc-xs text-fc-fg-muted">Your name and picture come from your Facile account.</p>
	{/if}
	<div class="flex flex-col gap-2">
		<span class="text-fc-sm font-medium">Theme</span>
		<div class="flex gap-1 rounded-fc-md bg-fc-surface p-1" role="radiogroup" aria-label="Theme">
			{#each themes as option (option.value)}
				<button
					type="button"
					role="radio"
					aria-checked={theme.preference === option.value}
					class="min-h-11 flex-1 rounded-fc-sm text-fc-sm {theme.preference === option.value ? 'bg-fc-accent text-fc-accent-fg' : 'hover:bg-fc-component'}"
					onclick={() => theme.set(option.value)}
				>
					{option.label}
				</button>
			{/each}
		</div>
	</div>
	<div class="flex items-start justify-between gap-3">
		<div class="flex flex-col">
			<span class="text-fc-sm font-medium">Browser notifications</span>
			<span class="text-fc-xs text-fc-fg-muted">
				{notify.supported ? (notify.permission === 'denied' ? 'Blocked in your browser settings.' : 'For mentions and direct messages while Prisme is in the background.') : 'Not supported by this browser.'}
			</span>
		</div>
		<Switch aria-label="Browser notifications" checked={notify.enabled} disabled={!notify.supported || notify.permission === 'denied'} onchange={(e) => notify.setEnabled(e.currentTarget.checked)} />
	</div>
	{#snippet footer()}
		<Button variant="ghost-danger" icon={icons.logout} onclick={() => app.logout()}>Sign out</Button>
	{/snippet}
</Modal>
