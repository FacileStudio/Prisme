<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { goto } from '$app/navigation';
	import { Button, EmptyState, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { serverName } from '$lib/state/actions';
	import { acceptInvite } from '$lib/matrix/directs';

	let { room }: { room: Room } = $props();
	let busy = $state(false);

	const inviter = $derived(room.getDMInviter() ?? '');
	const inviterName = $derived(room.getMember(inviter)?.name ?? inviter);

	async function answer(accept: boolean) {
		if (!app.client) return;
		busy = true;
		try {
			if (accept) await acceptInvite(app.client, room, serverName());
			else {
				await app.client.leave(room.roomId);
				await goto('/');
			}
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not answer the invite');
		} finally {
			busy = false;
		}
	}
</script>

<div class="flex flex-1 items-center justify-center p-6">
	<EmptyState icon={glyphs.chat} title={room.name} description={inviter ? `${inviterName} wants to talk with you.` : 'You are invited to this conversation.'}>
		<div class="flex gap-2">
			<Button variant="primary" icon={icons.check} disabled={busy} onclick={() => answer(true)}>Accept</Button>
			<Button variant="ghost" icon={icons.close} disabled={busy} onclick={() => answer(false)}>Decline</Button>
		</div>
	</EmptyState>
</div>
