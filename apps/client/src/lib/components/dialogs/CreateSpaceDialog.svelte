<script lang="ts">
	import { goto } from '$app/navigation';
	import { Field, Input, icons, toast } from '@facile/muse';
	import { createSpace, joinChannel, roomHref } from '$lib/state/actions';
	import { ui } from '$lib/state/ui.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import FormDialog from './FormDialog.svelte';

	let { onclose }: { onclose: () => void } = $props();
	let name = $state('');
	let busy = $state(false);

	async function submit() {
		busy = true;
		try {
			const space = await createSpace(name.trim());
			ui.spaceId = space.room_id;
			onclose();
			const general = space.channels[0];
			if (general) {
				await joinChannel(general.room_id).catch(() => undefined);
				await goto(roomHref(general.room_id));
			}
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('createSpace.failed'));
		} finally {
			busy = false;
		}
	}
</script>

<FormDialog title={t('createSpace.title')} submitLabel={t('common.create')} icon={icons.plus} {busy} disabled={!name.trim()} onsubmit={submit} {onclose}>
	<p class="text-fc-sm text-fc-fg-muted">{t('createSpace.body')}</p>
	<Field label={t('common.name')}>
		{#snippet children({ id })}
			<Input {id} bind:value={name} maxlength={80} placeholder="Facile Studio" required />
		{/snippet}
	</Field>
</FormDialog>
