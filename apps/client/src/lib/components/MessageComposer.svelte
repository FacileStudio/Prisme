<script lang="ts">
	let {
		roomName,
		typingUsers,
		onsend,
		ontyping
	}: {
		roomName: string;
		typingUsers: string[];
		onsend: (text: string) => Promise<void>;
		ontyping: (typing: boolean) => void;
	} = $props();

	let text = $state('');
	let sending = $state(false);
	let typingTimer: any = null;

	async function handleSubmit(e?: Event) {
		if (e) e.preventDefault();
		const trimmed = text.trim();
		if (!trimmed || sending) return;
		sending = true;
		try {
			await onsend(trimmed);
			text = '';
			ontyping(false);
		} finally {
			sending = false;
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			void handleSubmit();
			return;
		}
		ontyping(true);
		clearTimeout(typingTimer);
		typingTimer = setTimeout(() => {
			ontyping(false);
		}, 3000);
	}
</script>

<div class="border-t border-fc-border bg-fc-page p-4 select-none shrink-0">
	{#if typingUsers.length > 0}
		<div class="mb-1.5 flex items-center gap-1.5 text-fc-xs text-fc-fg-muted italic">
			<span class="inline-block size-1.5 animate-pulse rounded-full bg-fc-accent"></span>
			<span>{typingUsers.map((u) => u.split(':')[0].slice(1)).join(', ')} {typingUsers.length === 1 ? 'is' : 'are'} typing…</span>
		</div>
	{/if}

	<form onsubmit={handleSubmit} class="flex flex-col rounded-fc-md border border-fc-border bg-fc-surface focus-within:border-fc-accent focus-within:ring-1 focus-within:ring-fc-accent">
		<textarea
			bind:value={text}
			onkeydown={handleKeydown}
			placeholder={`Message ${roomName}…`}
			rows={2}
			class="w-full resize-none bg-transparent p-3 text-fc-sm text-fc-fg placeholder:text-fc-fg-muted focus:outline-none"
		></textarea>

		<div class="flex items-center justify-between border-t border-fc-border/40 px-3 py-2">
			<div class="flex items-center gap-1 text-fc-fg-muted text-fc-xs">
				<button
					type="button"
					class="rounded p-1 hover:bg-fc-surface-hover hover:text-fc-fg"
					title="Code block"
					onclick={() => (text += '```\n\n```')}
				>
					&lt;/&gt;
				</button>
				<button
					type="button"
					class="rounded p-1 hover:bg-fc-surface-hover hover:text-fc-fg font-bold"
					title="Bold"
					onclick={() => (text += '**bold**')}
				>
					B
				</button>
			</div>

			<button
				type="submit"
				disabled={!text.trim() || sending}
				class="inline-flex items-center justify-center rounded-fc-md bg-fc-accent px-3 py-1.5 text-fc-xs font-semibold text-fc-accent-fg shadow-sm transition-opacity hover:opacity-90 disabled:opacity-40 disabled:cursor-not-allowed"
			>
				{sending ? 'Sending…' : 'Send'}
			</button>
		</div>
	</form>
</div>
