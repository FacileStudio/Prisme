<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { ConfirmModal, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { BROWSE_PATH } from '$lib/state/nav';
	import { t } from '$lib/i18n/index.svelte';

	let { roomId, onclose }: { roomId: string; onclose: () => void } = $props();
	let open = $state(true);

	const name = $derived(app.client?.getRoom(roomId)?.name ?? '');

	async function leave() {
		try {
			await app.client?.leave(roomId);
			toast.success(t('channel.left', { name }));
			const here = page.params.id ? decodeURIComponent(page.params.id) === roomId : false;
			if (here) await goto(BROWSE_PATH);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('channel.leaveFailed'));
		} finally {
			onclose();
		}
	}

	$effect(() => {
		if (!open) onclose();
	});
</script>

<ConfirmModal
	bind:open
	tone="danger"
	icon={glyphs.exit}
	title={t('channel.leaveTitle', { name })}
	description={t('channel.leaveBody')}
	confirmLabel={t('channel.leave')}
	cancelLabel={t('common.cancel')}
	onConfirm={leave}
/>
