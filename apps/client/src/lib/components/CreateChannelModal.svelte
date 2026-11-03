<script lang="ts">
	import { Button, Icon, IconButton, Input, icons } from '@facile/muse';
	import { backend } from '$lib/backend';
	import type { ChannelRecord } from '$lib/types';

	let {
		spaceId,
		onclose,
		oncreated
	}: {
		spaceId: number;
		onclose: () => void;
		oncreated: (channel: ChannelRecord) => void;
	} = $props();

	let name = $state('');
	let topic = $state('');
	let saving = $state(false);
	let error = $state('');

	async function handleSubmit(e: Event) {
		e.preventDefault();
		const trimmed = name.trim();
		if (!trimmed || saving) return;
		saving = true;
		error = '';
		try {
			const channel = await backend.createChannel(spaceId, trimmed, topic.trim());
			oncreated(channel);
			onclose();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create channel';
		} finally {
			saving = false;
		}
	}
</script>

<svelte:window onkeydown={(e) => { if (e.key === 'Escape') onclose(); }} />

<div
	role="presentation"
	class="fixed inset-0 z-50 flex items-center justify-center bg-fc-scrim p-4"
	onclick={(e) => { if (e.target === e.currentTarget) onclose(); }}
>
	<div class="w-full max-w-md rounded-fc-xl bg-fc-component p-6 shadow-xl text-fc-fg">
		<div class="flex items-center justify-between pb-4">
			<h2 class="text-fc-lg font-bold text-fc-fg">Create a Channel</h2>
			<IconButton variant="ghost" class="size-8" aria-label="Close" onclick={onclose}>
				<Icon icon={icons.close} size={18} />
			</IconButton>
		</div>

		<form onsubmit={handleSubmit} class="space-y-4">
			<div class="space-y-1.5">
				<label for="channel-name" class="block text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
					Channel Name
				</label>
				<div class="relative flex items-center">
					<span class="absolute left-3.5 z-10 text-fc-sm font-bold text-fc-fg-muted">#</span>
					<Input
						id="channel-name"
						bind:value={name}
						placeholder="design"
						maxlength={80}
						required
						class="pl-8"
					/>
				</div>
			</div>

			<div class="space-y-1.5">
				<label for="channel-topic" class="block text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
					Topic (optional)
				</label>
				<Input
					id="channel-topic"
					bind:value={topic}
					placeholder="What is this channel about?"
					maxlength={160}
				/>
			</div>

			{#if error}
				<div class="rounded-fc-md bg-fc-danger/10 p-3 text-fc-xs text-fc-danger">
					{error}
				</div>
			{/if}

			<div class="flex justify-end gap-2 pt-2">
				<Button variant="ghost" onclick={onclose}>
					Cancel
				</Button>
				<Button
					type="submit"
					disabled={!name.trim() || saving}
				>
					{saving ? 'Creating…' : 'Create Channel'}
				</Button>
			</div>
		</form>
	</div>
</div>
