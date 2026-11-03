<script lang="ts">
	import { Icon, IconButton, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { call } from '$lib/call/call.svelte';

	let { roomId, join = false }: { roomId: string; join?: boolean } = $props();
	let busy = $state(false);

	const enabled = $derived(app.session?.calls_enabled ?? false);
	const label = $derived(enabled ? 'Start a call' : 'Calls are not available on this server');


	async function startCall() {
		busy = true;
		try {
			await call.start(roomId, join);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not start the call');
		} finally {
			busy = false;
		}
	}
</script>

<IconButton variant="ghost" aria-label={label} title={label} disabled={!enabled || busy} onclick={startCall}>
	<Icon icon={glyphs.phone} size={18} />
</IconButton>
