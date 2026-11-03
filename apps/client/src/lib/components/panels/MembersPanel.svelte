<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';
	import { isAgent, isBot } from '$lib/chat/users';
	import PresenceAvatar from '$lib/components/common/PresenceAvatar.svelte';
	import UserName from '$lib/components/common/UserName.svelte';

	let { room }: { room: Room } = $props();

	$effect(() => {
		void room.loadMembersIfNeeded();
	});

	const members = $derived.by(() => {
		void app.version;
		return room
			.getJoinedMembers()
			.filter((m) => !isBot(m.userId))
			.sort((a, b) => Number(isAgent(a.userId)) - Number(isAgent(b.userId)) || a.name.localeCompare(b.name));
	});
	const invited = $derived.by(() => {
		void app.version;
		return room.getMembersWithMembership('invite');
	});
</script>

<ul class="flex-1 overflow-y-auto overscroll-contain p-2">
	{#each members as member (member.userId)}
		<li class="flex min-h-11 items-center gap-3 rounded-fc-md px-2">
			<PresenceAvatar userId={member.userId} name={member.name} />
			<UserName userId={member.userId} name={member.name} class="flex-1 text-fc-sm" />
		</li>
	{/each}
	{#each invited as member (member.userId)}
		<li class="flex min-h-11 items-center gap-3 rounded-fc-md px-2 text-fc-fg-muted">
			<PresenceAvatar userId={member.userId} name={member.name} />
			<UserName userId={member.userId} name="{member.name} (invited)" class="flex-1 text-fc-sm" />
		</li>
	{/each}
</ul>
