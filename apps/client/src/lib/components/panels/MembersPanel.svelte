<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { isBot } from '$lib/chat/users';
	import Avatar from '$lib/components/common/Avatar.svelte';
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
			.sort((a, b) => Number(app.isAgent(a.userId)) - Number(app.isAgent(b.userId)) || a.name.localeCompare(b.name));
	});
	const invited = $derived.by(() => {
		void app.version;
		return room.getMembersWithMembership('invite');
	});
</script>

<ul class="flex-1 overflow-y-auto overscroll-contain p-2">
	{#each members as member (member.userId)}
		<li class="flex min-h-11 items-center gap-3 rounded-fc-md px-2">
			<Avatar userId={member.userId} name={member.name} mxc={member.getMxcAvatarUrl()} size={28} presence={!app.isAgent(member.userId)} />
			<UserName userId={member.userId} name={member.name} class="flex-1 text-fc-sm" />
		</li>
	{/each}
	{#each invited as member (member.userId)}
		<li class="flex min-h-11 items-center gap-3 rounded-fc-md px-2 text-fc-fg-muted">
			<Avatar userId={member.userId} name={member.name} mxc={member.getMxcAvatarUrl()} size={28} />
			<UserName userId={member.userId} name={t('members.invited', { name: member.name })} class="flex-1 text-fc-sm" />
		</li>
	{/each}
</ul>
