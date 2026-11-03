<script lang="ts">
	import { Badge, Button, Icon, IconButton, icons } from '@facile/muse';

	let {
		joinUrl,
		onclose
	}: {
		joinUrl: string;
		onclose: () => void;
	} = $props();
</script>

<svelte:window onkeydown={(e) => { if (e.key === 'Escape') onclose(); }} />

<div
	role="presentation"
	class="fixed inset-0 z-50 flex items-center justify-center bg-fc-scrim p-4"
	onclick={(e) => { if (e.target === e.currentTarget) onclose(); }}
>
	<div class="flex h-[80vh] w-full max-w-4xl flex-col rounded-fc-xl bg-fc-component p-6 shadow-xl text-fc-fg">
		<div class="flex items-center justify-between pb-4">
			<div class="flex items-center gap-2.5">
				<Badge tone="success">Live</Badge>
				<h2 class="text-fc-md font-bold text-fc-fg">Echo Call</h2>
			</div>
			<div class="flex items-center gap-2">
				<Button
					href={joinUrl}
					target="_blank"
					rel="noreferrer"
					variant="ghost"
					size="sm"
				>
					Open in Tab
				</Button>
				<Button
					variant="danger"
					size="sm"
					onclick={onclose}
				>
					Leave
				</Button>
				<IconButton
					variant="ghost"
					class="size-8 ml-1"
					aria-label="Close"
					onclick={onclose}
				>
					<Icon icon={icons.close} size={18} />
				</IconButton>
			</div>
		</div>

		<div class="flex-1 min-h-0 overflow-hidden rounded-fc-lg bg-black">
			<iframe
				src={joinUrl}
				title="Echo Call"
				allow="camera; microphone; display-capture"
				class="size-full border-0"
			></iframe>
		</div>
	</div>
</div>
