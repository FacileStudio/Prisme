<script lang="ts">
	import { Badge, Button, Icon, IconButton, Input, icons } from '@facile/muse';
	import { backend } from '$lib/backend';

	let {
		currentRoomId,
		onclose,
		oncreated
	}: {
		currentRoomId: string | null;
		onclose: () => void;
		oncreated: (agentUserId: string) => void;
	} = $props();

	let name = $state('');
	let checking = $state(false);
	let available = $state<boolean | null>(null);
	let reason = $state('');
	let saving = $state(false);
	let submitError = $state('');

	let checkTimer: ReturnType<typeof setTimeout> | null = null;

	function handleInput() {
		available = null;
		reason = '';
		submitError = '';
		if (checkTimer) clearTimeout(checkTimer);
		const trimmed = name.trim();
		if (!trimmed) {
			checking = false;
			return;
		}
		checking = true;
		checkTimer = setTimeout(async () => {
			try {
				const res = await backend.checkAgentName(trimmed);
				available = res.available;
				reason = res.reason || '';
			} catch {
				available = false;
				reason = 'Could not verify name availability';
			} finally {
				checking = false;
			}
		}, 300);
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		const trimmed = name.trim();
		if (!trimmed || saving || available === false) return;
		saving = true;
		submitError = '';
		try {
			const login = await backend.createAgent(trimmed);
			if (currentRoomId) {
				await backend.setAgentRooms(login.user_id, [currentRoomId]);
			}
			oncreated(login.user_id);
			onclose();
		} catch (err) {
			submitError = err instanceof Error ? err.message : 'Failed to create agent';
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
			<h2 class="text-fc-lg font-bold text-fc-fg">Add an Agent</h2>
			<IconButton variant="ghost" class="size-8" aria-label="Close" onclick={onclose}>
				<Icon icon={icons.close} size={18} />
			</IconButton>
		</div>

		<form onsubmit={handleSubmit} class="space-y-4">
			<div class="space-y-1.5">
				<label for="agent-name" class="block text-fc-xs font-semibold uppercase tracking-wider text-fc-fg-muted">
					Agent Name
				</label>
				<div class="relative flex items-center">
					<span class="absolute left-3.5 z-10 text-fc-sm font-bold text-fc-fg-muted">@</span>
					<Input
						id="agent-name"
						bind:value={name}
						oninput={handleInput}
						placeholder="pixel"
						maxlength={80}
						autocomplete="off"
						required
						class="pl-8"
					/>
				</div>

				<div class="pt-1">
					{#if checking}
						<span class="text-fc-xs text-fc-fg-muted">Checking availability…</span>
					{:else if available === true}
						<Badge tone="success">✓ @{name.trim().toLowerCase()} is available</Badge>
					{:else if available === false}
						<Badge tone="danger">✗ {reason || 'The name is already taken'}</Badge>
					{:else}
						<span class="text-fc-xs text-fc-fg-muted">Names are unique across the workspace.</span>
					{/if}
				</div>
			</div>

			{#if submitError}
				<div class="rounded-fc-md bg-fc-danger/10 p-3 text-fc-xs text-fc-danger">
					{submitError}
				</div>
			{/if}

			<div class="flex justify-end gap-2 pt-2">
				<Button variant="ghost" onclick={onclose}>
					Cancel
				</Button>
				<Button
					type="submit"
					disabled={!name.trim() || checking || available === false || saving}
				>
					{saving ? 'Creating…' : 'Create Agent'}
				</Button>
			</div>
		</form>
	</div>
</div>
