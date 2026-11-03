<script lang="ts">
	import { Modal, toast } from '@facile/muse';
	import type { DirectoryUser, Space } from '$lib/api/prisme';
	import { app } from '$lib/state/app.svelte';
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
			toast.success(`${user.name} is invited to ${space.name}`);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not send the invite');
		} finally {
			busyId = '';
		}
	}
</script>

<Modal open title="Invite people to {space.name}" showClose onClose={onclose}>
	<p class="text-fc-sm text-fc-fg-muted">Invited people join the space and its #general channel when they accept.</p>
	<UserPicker onpick={invite} {busyId} exclude={members} />
</Modal>
