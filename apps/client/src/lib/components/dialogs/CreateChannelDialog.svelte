<script lang="ts">
	import { goto } from '$app/navigation';
	import { Field, Input, icons, toast } from '@facile/muse';
	import type { Space } from '$lib/api/prisme';
	import { createChannel, joinChannel, roomHref } from '$lib/state/actions';
	import { t } from '$lib/i18n/index.svelte';
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
			toast.danger(err instanceof Error ? err.message : t('createChannel.failed'));
		} finally {
			busy = false;
		}
	}
</script>

<FormDialog title={t('createChannel.title', { name: space.name })} submitLabel={t('common.create')} icon={icons.plus} {busy} disabled={!name.trim()} onsubmit={submit} {onclose}>
	<Field label={t('common.name')}>
		{#snippet children({ id })}
			<Input {id} bind:value={name} placeholder="design" required />
		{/snippet}
	</Field>
	<Field label={t('createChannel.topic')} helper={t('createChannel.topicHelp')}>
		{#snippet children({ id, describedBy })}
			<Input {id} aria-describedby={describedBy} bind:value={topic} placeholder={t('createChannel.topicPlaceholder')} />
		{/snippet}
	</Field>
</FormDialog>
