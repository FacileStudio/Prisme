<script lang="ts">
	import { Icon, IconButton, icons } from '@facile/muse';
	import { app } from '$lib/state/app.svelte';
	import { ui } from '$lib/state/ui.svelte';
	import PresenceAvatar from '$lib/components/common/PresenceAvatar.svelte';

	const user = $derived(app.session?.user);
	const name = $derived(user?.name || user?.email.split('@')[0] || 'You');
</script>

{#if user}
	<div class="flex items-center gap-2 border-t border-fc-border py-1 pl-3 pr-1">
		<PresenceAvatar userId={user.matrix_user_id} {name} src={user.avatar_url} size="md" />
		<div class="flex min-w-0 flex-1 flex-col">
			<span class="truncate text-fc-sm font-medium">{name}</span>
			<span class="truncate text-fc-xs text-fc-fg-muted">{user.email}</span>
		</div>
		<IconButton variant="ghost" aria-label="Settings" title="Settings" onclick={() => (ui.dialog = 'settings')}>
			<Icon icon={icons.settings} size={18} />
		</IconButton>
	</div>
{/if}
