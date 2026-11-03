<script lang="ts">
	import type { Room, RoomMember } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';
	import { typingSentence } from '$lib/chat/users';
	import { watchTyping } from '$lib/matrix/typing';
	import Avatar from '$lib/components/common/Avatar.svelte';

	let { room }: { room: Room } = $props();
	let typers = $state.raw<RoomMember[]>([]);

	$effect(() => {
		const client = app.client;
		if (!client) return;
		return watchTyping(client, room, (members) => (typers = members));
	});

	const text = $derived(typingSentence(typers.map((m) => m.name)));
</script>

<div class="flex h-6 items-center gap-1.5 px-4 text-fc-xs text-fc-fg-muted" aria-live="polite">
	{#if typers.length > 0}
		<span class="flex -space-x-1" aria-hidden="true">
			{#each typers.slice(0, 3) as member (member.userId)}
				<Avatar userId={member.userId} name={member.name} mxc={member.getMxcAvatarUrl()} size={20} class="ring-2 ring-fc-page" />
			{/each}
		</span>
		<span class="truncate">{text}</span>
	{/if}
</div>
