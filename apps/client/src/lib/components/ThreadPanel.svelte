<script lang="ts">
	import type { ChatMessage } from '$lib/types';

	let {
		rootMessage,
		replies,
		currentUserId,
		onclose,
		onsendreply
	}: {
		rootMessage: ChatMessage;
		replies: ChatMessage[];
		currentUserId: string;
		onclose: () => void;
		onsendreply: (text: string) => Promise<void>;
	} = $props();

	let replyText = $state('');
	let sending = $state(false);

	async function handleSend(e?: Event) {
		if (e) e.preventDefault();
		const trimmed = replyText.trim();
		if (!trimmed || sending) return;
		sending = true;
		try {
			await onsendreply(trimmed);
			replyText = '';
		} finally {
			sending = false;
		}
	}

	function formatTime(ts: number): string {
		return new Date(ts).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
	}
</script>

<aside class="flex h-full w-80 flex-col overflow-hidden rounded-fc-lg bg-fc-component select-none shrink-0">
	<div class="flex h-14 items-center justify-between px-4 font-bold text-fc-fg">
		<span class="text-fc-sm">Thread</span>
		<button
			type="button"
			aria-label="Close"
			class="flex size-7 items-center justify-center rounded-fc-md text-fc-fg-muted hover:bg-fc-surface hover:text-fc-fg"
			onclick={onclose}
		>
			<svg width="16" height="16" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
				<line x1="18" y1="6" x2="6" y2="18" />
				<line x1="6" y1="6" x2="18" y2="18" />
			</svg>
		</button>
	</div>

	<div class="flex-1 overflow-y-auto p-4 space-y-4">
		<div class="rounded-fc-md bg-fc-surface p-3">
			<div class="flex items-center gap-2">
				<span class="text-fc-xs font-bold text-fc-fg">{rootMessage.senderName}</span>
				{#if rootMessage.isAgent}
					<span class="rounded bg-fc-accent/15 px-1.5 py-0.5 text-[9px] font-bold text-fc-accent uppercase">
						Agent
					</span>
				{/if}
				<span class="text-[10px] text-fc-fg-muted">{formatTime(rootMessage.timestamp)}</span>
			</div>
			<div class="mt-1 text-fc-sm text-fc-fg break-words whitespace-pre-wrap">
				{rootMessage.body}
			</div>
		</div>

		<div class="space-y-3">
			{#each replies as reply (reply.id)}
				<div class="flex gap-2">
					<div class="flex size-7 shrink-0 items-center justify-center rounded-full bg-fc-surface text-[11px] font-bold">
						{reply.senderName.slice(0, 2).toUpperCase()}
					</div>
					<div class="flex-1 min-w-0">
						<div class="flex items-center gap-1.5">
							<span class="text-fc-xs font-bold text-fc-fg">{reply.senderName}</span>
							{#if reply.isAgent}
								<span class="rounded bg-fc-accent/15 px-1 py-0.2 text-[9px] font-bold text-fc-accent uppercase">
									Agent
								</span>
							{/if}
							<span class="text-[10px] text-fc-fg-muted">{formatTime(reply.timestamp)}</span>
						</div>
						<div class="mt-0.5 text-fc-sm text-fc-fg break-words whitespace-pre-wrap">
							{reply.body}
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>

	<form onsubmit={handleSend} class="p-3">
		<div class="flex gap-2">
			<input
				type="text"
				bind:value={replyText}
				placeholder="Reply in thread…"
				class="flex-1 rounded-fc-md bg-fc-surface px-3 py-1.5 text-fc-sm text-fc-fg placeholder:text-fc-fg-muted focus:outline-none focus:ring-1 focus:ring-fc-ring"
			/>
			<button
				type="submit"
				disabled={!replyText.trim() || sending}
				class="rounded-fc-md bg-fc-accent px-3 py-1.5 text-fc-xs font-semibold text-fc-accent-fg disabled:opacity-40"
			>
				Reply
			</button>
		</div>
	</form>
</aside>
