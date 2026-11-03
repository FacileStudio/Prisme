<script lang="ts">
	import { Button, Icon, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { call as callState } from '$lib/call/call.svelte';
	import { callLength } from '$lib/call/media';
	import type { CallInfo } from '$lib/matrix/timeline';

	let { call, roomId }: { call: CallInfo; roomId: string } = $props();
	let busy = $state(false);

	const starter = $derived(app.client?.getRoom(roomId)?.getMember(call.started_by)?.rawDisplayName ?? call.started_by);
	const ended = $derived(!!call.ended_at);
	const here = $derived(!ended && callState.active && callState.roomId === roomId);

	async function join() {
		busy = true;
		try {
			await callState.start(roomId, true);
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : 'Could not join the call');
		} finally {
			busy = false;
		}
	}
</script>

<div class="mt-1 flex max-w-sm items-center gap-3 rounded-fc-lg border border-fc-border bg-fc-component p-3">
	<span
		class="flex size-10 shrink-0 items-center justify-center rounded-fc-pill {ended
			? 'bg-fc-surface text-fc-fg-muted'
			: 'bg-fc-success/10 text-fc-success'}"
	>
		<Icon icon={ended ? glyphs.endCall : glyphs.phone} size={20} />
	</span>
	<div class="flex min-w-0 flex-1 flex-col">
		<span class="text-fc-sm font-semibold">{ended ? 'Call ended' : 'Call started'}</span>
		<span class="truncate text-fc-xs text-fc-fg-muted">
			{ended ? `${callLength(call.duration_ms ?? 0)} · started by ${starter}` : `by ${starter}`}
		</span>
	</div>
	{#if !ended}
		<Button size="sm" variant={here ? 'outline' : 'primary'} disabled={busy || here || !app.session?.calls_enabled} onclick={join}>
			{here ? 'Joined' : 'Join'}
		</Button>
	{/if}
</div>
