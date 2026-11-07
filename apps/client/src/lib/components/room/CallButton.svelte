<script lang="ts">
	import { Icon, IconButton, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { call } from '$lib/call/call.svelte';
	import { t } from '$lib/i18n/index.svelte';

	let { roomId, join = false }: { roomId: string; join?: boolean } = $props();
	let busy = $state(false);

	const enabled = $derived(app.session?.calls_enabled ?? false);
	const label = $derived(enabled ? t('call.start') : t('call.unavailable'));


	async function startCall() {
		busy = true;
		try {
			await call.start(roomId, join);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('call.startFailed'));
		} finally {
			busy = false;
		}
	}
</script>

<IconButton variant="ghost" aria-label={label} title={label} disabled={!enabled || busy} onclick={startCall}>
	<Icon icon={glyphs.phone} size={18} />
</IconButton>
