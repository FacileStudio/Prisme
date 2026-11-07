<script lang="ts">
	import { untrack } from 'svelte';
	import { Field, Input, icons, toast } from '@facile/muse';
	import { app } from '$lib/state/app.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import FormDialog from './FormDialog.svelte';

	type Props = { roomId: string; kind: 'channel' | 'space'; onclose: () => void };
	let { roomId, kind, onclose }: Props = $props();
	let name = $state(untrack(() => app.client?.getRoom(roomId)?.name ?? ''));
	let busy = $state(false);

	async function submit() {
		if (!app.client) return;
		busy = true;
		try {
			await app.client.setRoomName(roomId, name.trim());
			toast.success(t(kind === 'space' ? 'rename.doneSpace' : 'rename.doneChannel'));
			onclose();
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t(kind === 'space' ? 'rename.failedSpace' : 'rename.failedChannel'));
		} finally {
			busy = false;
		}
	}
</script>

<FormDialog title={t(kind === 'space' ? 'rename.titleSpace' : 'rename.titleChannel')} submitLabel={t('rename.submit')} icon={icons.edit} {busy} disabled={!name.trim()} onsubmit={submit} {onclose}>
	<Field label={t('common.name')}>
		{#snippet children({ id })}
			<Input {id} bind:value={name} maxlength={80} required />
		{/snippet}
	</Field>
</FormDialog>
