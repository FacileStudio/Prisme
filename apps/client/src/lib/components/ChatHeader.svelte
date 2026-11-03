<script lang="ts">
	import { Button, IconButton } from '@facile/muse';

	let {
		roomName,
		callsEnabled,
		hasActiveCall,
		onstartcall,
		ontogglemembers,
		onopensearch
	}: {
		roomName: string;
		callsEnabled: boolean;
		hasActiveCall: boolean;
		onstartcall: () => void;
		ontogglemembers: () => void;
		onopensearch?: () => void;
	} = $props();
</script>

<header class="flex h-14 items-center justify-between bg-fc-component px-4 select-none shrink-0">
	<div class="flex items-center gap-2 font-bold text-fc-fg">
		<span class="text-fc-base truncate">{roomName}</span>
	</div>

	<div class="flex items-center gap-2">
		{#if callsEnabled}
			<Button
				size="sm"
				onclick={onstartcall}
				title={hasActiveCall ? 'Join Call (Echo)' : 'Start Call (Echo)'}
				class={hasActiveCall ? 'bg-fc-success text-white hover:opacity-90' : 'bg-fc-surface text-fc-fg hover:bg-fc-surface/80'}
			>
				<svg width="15" height="15" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2 mr-1">
					<path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z" />
				</svg>
				<span>{hasActiveCall ? 'Join Call' : 'Call'}</span>
			</Button>
		{/if}

		{#if onopensearch}
			<IconButton
				variant="ghost"
				title="Search"
				aria-label="Search"
				class="size-8"
				onclick={onopensearch}
			>
				<svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
					<circle cx="11" cy="11" r="8" />
					<path d="m21 21-4.3-4.3" />
				</svg>
			</IconButton>
		{/if}

		<IconButton
			variant="ghost"
			title="Members"
			aria-label="Members"
			class="size-8"
			onclick={ontogglemembers}
		>
			<svg width="18" height="18" viewBox="0 0 24 24" fill="none" class="stroke-current stroke-2">
				<path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
				<circle cx="9" cy="7" r="4" />
				<path d="M23 21v-2a4 4 0 0 0-3-3.87" />
				<path d="M16 3.13a4 4 0 0 1 0 7.75" />
			</svg>
		</IconButton>
	</div>
</header>
