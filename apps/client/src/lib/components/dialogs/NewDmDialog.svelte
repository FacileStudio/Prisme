<script lang="ts">
	import { goto } from '$app/navigation';
	import { Modal, toast } from '@facile/muse';
	import type { DirectoryUser } from '$lib/api/prisme';
	import { app } from '$lib/state/app.svelte';
	import { roomHref } from '$lib/state/actions';
	import { openDm } from '$lib/matrix/directs';
	import UserPicker from './UserPicker.svelte';

	let { onclose }: { onclose: () => void } = $props();
	let busyId = $state('');

	async function start(user: DirectoryUser) {
		if (!app.client) return;
		busyId = user.user_id;
		try {
			const roomId = await openDm(app.client, user.user_id);
			onclose();
			await goto(roomHref(roomId));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not open the conversation');
		} finally {
			busyId = '';
		}
	}
</script>

<Modal open title="New message" showClose onClose={onclose}>
	<UserPicker onpick={start} {busyId} />
</Modal>
