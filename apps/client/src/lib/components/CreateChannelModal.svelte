<script lang="ts">
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

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
	<div class="w-full max-w-md rounded-fc-lg border border-fc-border bg-fc-page p-6 shadow-xl">
		<div class="flex items-center justify-between pb-4">
			<h2 class="text-fc-lg font-bold text-fc-fg">Create a Channel</h2>
			<button
				type="button"
				aria-label="Close"
				class="rounded-fc-md p-1 text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg"
				onclick={onclose}
			>
				<svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
					<line x1="18" y1="6" x2="6" y2="18" />
					<line x1="6" y1="6" x2="18" y2="18" />
				</svg>
			</button>
		</div>

		<form onsubmit={handleSubmit} class="space-y-4">
			<div>
				<label for="channel-name" class="block text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
					Channel Name
				</label>
				<div class="relative mt-1.5 flex items-center">
					<span class="absolute left-3 text-fc-sm text-fc-fg-muted font-bold">#</span>
					<input
						id="channel-name"
						type="text"
						bind:value={name}
						placeholder="design"
						maxlength={80}
						required
						class="w-full rounded-fc-md border border-fc-border bg-fc-surface py-2 pr-3 pl-7 text-fc-sm text-fc-fg placeholder:text-fc-fg-muted focus:border-fc-accent focus:outline-none focus:ring-1 focus:ring-fc-accent"
					/>
				</div>
			</div>

			<div>
				<label for="channel-topic" class="block text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
					Topic (optional)
				</label>
				<input
					id="channel-topic"
					type="text"
					bind:value={topic}
					placeholder="What is this channel about?"
					maxlength={160}
					class="mt-1.5 w-full rounded-fc-md border border-fc-border bg-fc-surface px-3 py-2 text-fc-sm text-fc-fg placeholder:text-fc-fg-muted focus:border-fc-accent focus:outline-none focus:ring-1 focus:ring-fc-accent"
				/>
			</div>

			{#if error}
				<div class="rounded-fc-md bg-fc-danger/10 p-3 text-fc-xs text-fc-danger">
					{error}
				</div>
			{/if}

			<div class="flex justify-end gap-2 pt-2">
				<button
					type="button"
					class="rounded-fc-md border border-fc-border px-4 py-2 text-fc-sm font-medium text-fc-fg hover:bg-fc-surface"
					onclick={onclose}
				>
					Cancel
				</button>
				<button
					type="submit"
					disabled={!name.trim() || saving}
					class="rounded-fc-md bg-fc-accent px-4 py-2 text-fc-sm font-semibold text-fc-accent-fg shadow-sm hover:opacity-90 disabled:opacity-40"
				>
					{saving ? 'Creating…' : 'Create Channel'}
				</button>
			</div>
		</form>
	</div>
</div>
