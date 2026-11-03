<script lang="ts">
	import { goto } from '$app/navigation';
	import { Field, Input, icons, toast } from '@facile/muse';
	import type { Space } from '$lib/api/prisme';
	import { createChannel, joinChannel, roomHref } from '$lib/state/actions';
	import FormDialog from './FormDialog.svelte';

	let { space, onclose }: { space: Space; onclose: () => void } = $props();
	let name = $state('');
	let topic = $state('');
	let busy = $state(false);

	async function submit() {
		busy = true;
		try {
			const roomId = await createChannel(space.id, name.trim().replace(/^#/, ''), topic.trim());
			await joinChannel(roomId).catch(() => undefined);
			onclose();
			await goto(roomHref(roomId));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not create the channel');
		} finally {
			busy = false;
		}
	}
</script>

<FormDialog title="Create a channel in {space.name}" submitLabel="Create" icon={icons.plus} {busy} disabled={!name.trim()} onsubmit={submit} {onclose}>
	<Field label="Name">
		{#snippet children({ id })}
			<Input {id} bind:value={name} placeholder="design" required />
		{/snippet}
	</Field>
	<Field label="Topic" helper="Optional. Shown under the channel name.">
		{#snippet children({ id, describedBy })}
			<Input {id} aria-describedby={describedBy} bind:value={topic} placeholder="What this channel is for" />
		{/snippet}
	</Field>
</FormDialog>
