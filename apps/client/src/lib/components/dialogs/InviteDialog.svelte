<script lang="ts">
	import { Modal, toast } from '@facile/muse';
	import type { DirectoryUser, Space } from '$lib/api/prisme';
	import { app } from '$lib/state/app.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import UserPicker from './UserPicker.svelte';

	let { space, onclose }: { space: Space; onclose: () => void } = $props();
	let busyId = $state('');
	let invited = $state<string[]>([]);

	const members = $derived.by(() => {
		void app.version;
		const room = app.client?.getRoom(space.room_id);
		const ids = room ? [...room.getMembersWithMembership('join'), ...room.getMembersWithMembership('invite')].map((m) => m.userId) : [];
		return [...ids, ...invited];
	});

	async function invite(user: DirectoryUser) {
		if (!app.client) return;
		busyId = user.user_id;
		try {
			await app.client.invite(space.room_id, user.user_id);
			invited = [...invited, user.user_id];
			toast.success(t('invite.done', { name: user.name, space: space.name }));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('invite.failed'));
		} finally {
			busyId = '';
		}
	}
</script>

<Modal open title={t('invite.title', { name: space.name })} showClose onClose={onclose}>
	<p class="text-fc-sm text-fc-fg-muted">{t('invite.body')}</p>
	<UserPicker onpick={invite} {busyId} exclude={members} />
</Modal>
