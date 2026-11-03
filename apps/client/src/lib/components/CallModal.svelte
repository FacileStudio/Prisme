<script lang="ts">
	let {
		joinUrl,
		onclose
	}: {
		joinUrl: string;
		onclose: () => void;
	} = $props();
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-4">
	<div class="flex h-[80vh] w-full max-w-4xl flex-col overflow-hidden rounded-fc-lg border border-fc-border bg-black shadow-2xl">
		<div class="flex h-12 items-center justify-between border-b border-white/10 bg-black px-4 text-white">
			<div class="flex items-center gap-2">
				<span class="inline-block size-2 rounded-full bg-fc-success animate-pulse"></span>
				<span class="text-fc-sm font-bold">Echo Call</span>
			</div>
			<div class="flex items-center gap-2">
				<a
					href={joinUrl}
					target="_blank"
					rel="noreferrer"
					class="rounded-fc-md bg-white/10 px-3 py-1 text-fc-xs font-semibold text-white hover:bg-white/20"
				>
					Open in Tab
				</a>
				<button
					type="button"
					class="rounded-fc-md bg-fc-danger px-3 py-1 text-fc-xs font-semibold text-white hover:opacity-90"
					onclick={onclose}
				>
					Leave
				</button>
			</div>
		</div>

		<div class="flex-1 bg-black">
			<iframe
				src={joinUrl}
				title="Echo Call"
				allow="camera; microphone; display-capture"
				class="size-full border-0"
			></iframe>
		</div>
	</div>
</div>
