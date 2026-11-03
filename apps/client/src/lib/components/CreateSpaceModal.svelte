<script lang="ts">
	import { Button, Icon, IconButton, Input, icons } from '@facile/muse';
	import { backend } from '$lib/backend';
	import type { SpaceRecord } from '$lib/types';

	let {
		onclose,
		oncreated
	}: {
		onclose: () => void;
		oncreated: (space: SpaceRecord) => void;
	} = $props();

	let name = $state('');
	let saving = $state(false);
	let error = $state('');

	async function handleSubmit(e: Event) {
		e.preventDefault();
		const trimmed = name.trim();
		if (!trimmed || saving) return;
		saving = true;
		error = '';
		try {
			const space = await backend.createSpace(trimmed);
			oncreated(space);
			onclose();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create space';
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
			<h2 class="text-fc-lg font-bold text-fc-fg">Create a Space</h2>
			<IconButton variant="ghost" class="size-8" aria-label="Close" onclick={onclose}>
				<Icon icon={icons.close} size={18} />
			</IconButton>
		</div>

		<form onsubmit={handleSubmit} class="space-y-4">
			<div class="space-y-1.5">
				<label for="space-name" class="block text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
					Space Name
				</label>
				<Input
					id="space-name"
					bind:value={name}
					placeholder="Facile Studio"
					maxlength={80}
					required
				/>
				<p class="text-fc-xs text-fc-fg-muted">
					Spaces group channels and agents together under one roof.
				</p>
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
					{saving ? 'Creating…' : 'Create Space'}
				</Button>
			</div>
		</form>
	</div>
</div>
