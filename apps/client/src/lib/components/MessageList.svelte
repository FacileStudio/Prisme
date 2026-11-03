<script lang="ts">
	import { tick } from 'svelte';
	import type { ChatMessage } from '$lib/types';

	let {
		messages,
		currentUserId,
		onreact,
		onreplyinthread,
		onjoincall
	}: {
		messages: ChatMessage[];
		currentUserId: string;
		onreact: (messageId: string, emoji: string) => void;
		onreplyinthread: (message: ChatMessage) => void;
		onjoincall: (callId: string, joinUrl: string) => void;
	} = $props();

	let container = $state<HTMLElement | null>(null);

	$effect(() => {
		if (messages.length > 0) {
			void tick().then(() => {
				if (container) {
					container.scrollTop = container.scrollHeight;
				}
			});
		}
	});

	function formatTime(ts: number): string {
		const date = new Date(ts);
		return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
	}

	function formatBody(raw: string): string {
		return raw;
	}
</script>

<div bind:this={container} class="flex flex-1 flex-col overflow-y-auto p-4 space-y-4">
	{#if messages.length === 0}
		<div class="m-auto flex flex-col items-center justify-center text-center text-fc-fg-muted">
			<div class="mb-2 flex size-12 items-center justify-center rounded-full bg-fc-surface text-fc-fg-muted">
				💬
			</div>
			<p class="text-fc-sm font-medium">No messages yet</p>
			<p class="text-fc-xs text-fc-fg-muted">Start the conversation below.</p>
		</div>
	{/if}

	{#each messages as msg (msg.id)}
		<div class="group relative flex gap-3 rounded-fc-md p-2 transition-colors hover:bg-fc-surface/50">
			<div class="flex size-9 shrink-0 items-center justify-center overflow-hidden rounded-full bg-fc-surface text-fc-fg font-semibold border border-fc-border">
				{#if msg.senderAvatar}
					<img src={msg.senderAvatar} alt={msg.senderName} class="size-full object-cover" />
				{:else if msg.isAgent}
					🤖
				{:else}
					{msg.senderName.slice(0, 2).toUpperCase()}
				{/if}
			</div>

			<div class="flex-1 min-w-0">
				<div class="flex items-center gap-2">
					<span class="text-fc-sm font-bold text-fc-fg truncate">{msg.senderName}</span>
					{#if msg.isAgent}
						<span class="rounded bg-fc-accent/15 px-1.5 py-0.5 text-[10px] font-bold text-fc-accent tracking-wider uppercase">
							Agent
						</span>
					{/if}
					<span class="text-[11px] text-fc-fg-muted">{formatTime(msg.timestamp)}</span>
				</div>

				{#if msg.call}
					<div class="mt-2 flex max-w-md items-center justify-between gap-4 rounded-fc-md border border-fc-border bg-fc-surface p-3">
						<div class="flex items-center gap-3">
							<div class="flex size-8 items-center justify-center rounded-full bg-fc-success/20 text-fc-success">
								📞
							</div>
							<div>
								<p class="text-fc-sm font-semibold text-fc-fg">Call started</p>
								<p class="text-fc-xs text-fc-fg-muted">By {msg.senderName}</p>
							</div>
						</div>
						<button
							type="button"
							class="rounded-fc-md bg-fc-success px-3 py-1.5 text-fc-xs font-semibold text-white shadow-sm hover:opacity-90"
							onclick={() => onjoincall(msg.call!.call_id, msg.call!.join_url)}
						>
							Join
						</button>
					</div>
				{:else}
					<div class="mt-1 text-fc-sm text-fc-fg break-words whitespace-pre-wrap leading-relaxed font-sans">
						{formatBody(msg.body)}
					</div>
				{/if}

				{#if Object.keys(msg.reactions).length > 0}
					<div class="mt-2 flex flex-wrap gap-1">
						{#each Object.entries(msg.reactions) as [emoji, users]}
							<button
								type="button"
								class={`flex items-center gap-1 rounded-full border px-2 py-0.5 text-fc-xs transition-colors ${
									users.includes(currentUserId)
										? 'border-fc-accent bg-fc-accent/10 text-fc-accent'
										: 'border-fc-border bg-fc-surface text-fc-fg hover:bg-fc-surface-hover'
								}`}
								onclick={() => onreact(msg.id, emoji)}
								title={users.join(', ')}
							>
								<span>{emoji}</span>
								<span class="text-[10px] font-semibold">{users.length}</span>
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<div class="absolute right-2 top-2 hidden items-center rounded-fc-md border border-fc-border bg-fc-page p-0.5 shadow-sm group-hover:flex">
				<button
					type="button"
					class="p-1 text-fc-xs hover:bg-fc-surface-hover rounded"
					title="Thumbs up"
					onclick={() => onreact(msg.id, '👍')}
				>
					👍
				</button>
				<button
					type="button"
					class="p-1 text-fc-xs hover:bg-fc-surface-hover rounded"
					title="Heart"
					onclick={() => onreact(msg.id, '❤️')}
				>
					❤️
				</button>
				<button
					type="button"
					class="p-1 text-fc-xs hover:bg-fc-surface-hover rounded"
					title="Rocket"
					onclick={() => onreact(msg.id, '🚀')}
				>
					🚀
				</button>
				<div class="mx-0.5 h-3 w-px bg-fc-border"></div>
				<button
					type="button"
					class="p-1 text-fc-fg-muted hover:bg-fc-surface-hover hover:text-fc-fg rounded text-fc-xs"
					title="Reply in thread"
					onclick={() => onreplyinthread(msg)}
				>
					<svg width="14" height="14" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
						<path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
					</svg>
				</button>
			</div>
		</div>
	{/each}
</div>
