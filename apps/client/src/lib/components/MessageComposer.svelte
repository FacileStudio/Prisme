<script lang="ts">
	import { Button, IconButton } from '@facile/muse';

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

<div class="p-4 select-none shrink-0">
	{#if typingUsers.length > 0}
		<div class="mb-1.5 flex items-center gap-1.5 text-fc-xs text-fc-fg-muted italic">
			<span class="inline-block size-1.5 animate-pulse rounded-full bg-fc-accent"></span>
			<span>{typingUsers.map((u) => u.split(':')[0].slice(1)).join(', ')} {typingUsers.length === 1 ? 'is' : 'are'} typing…</span>
		</div>
	{/if}

	<form
		onsubmit={handleSubmit}
		class="flex flex-col bg-fc-surface rounded-fc-lg p-3 focus-within:ring-2 focus-within:ring-fc-ring transition-shadow"
	>
		<textarea
			bind:value={text}
			onkeydown={handleKeydown}
			placeholder={`Message ${roomName}…`}
			rows={2}
			class="w-full resize-none bg-transparent text-fc-sm text-fc-fg placeholder:text-fc-fg-muted focus:outline-none"
		></textarea>

		<div class="flex items-center justify-between pt-2">
			<div class="flex items-center gap-1">
				<IconButton
					variant="ghost"
					title="Code block"
					aria-label="Code block"
					class="size-7 rounded-fc-md [&_svg]:size-3.5 text-fc-xs font-mono"
					onclick={() => (text += '```\n\n```')}
				>
					<span>&lt;/&gt;</span>
				</IconButton>
				<IconButton
					variant="ghost"
					title="Bold"
					aria-label="Bold"
					class="size-7 rounded-fc-md font-bold text-fc-xs"
					onclick={() => (text += '**bold**')}
				>
					<span>B</span>
				</IconButton>
			</div>

			<Button
				type="submit"
				size="sm"
				disabled={!text.trim() || sending}
			>
				{sending ? 'Sending…' : 'Send'}
			</Button>
		</div>
	</form>
</div>
