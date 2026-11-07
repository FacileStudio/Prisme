<script lang="ts">
	import { EmptyState } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import { app } from '$lib/state/app.svelte';
	import { spaceOf } from '$lib/state/nav';
	import { rememberRoom } from '$lib/state/last-room';
	import { DraftState } from '$lib/state/draft.svelte';
	import { ui } from '$lib/state/ui.svelte';
	import { isDirect } from '$lib/matrix/rooms';
	import RoomHeader from './RoomHeader.svelte';
	import Timeline from './Timeline.svelte';
	import Composer from './Composer.svelte';
	import Typing from './Typing.svelte';
	import InvitePrompt from './InvitePrompt.svelte';
	import ReadMarker from './ReadMarker.svelte';

	let { roomId }: { roomId: string } = $props();
	const draft = new DraftState();

	const room = $derived.by(() => {
		void app.version;
		return app.client?.getRoom(roomId) ?? null;
	});
	const membership = $derived.by(() => {
		void app.version;
		return room?.getMyMembership() ?? 'leave';
	});
	const direct = $derived(!!room && !!app.client && isDirect(app.client, room));

	$effect(() => {
		const space = spaceOf(app.spaces, roomId);
		if (space) ui.spaceId = space.room_id;
		if (membership === 'join') rememberRoom(roomId);
	});
</script>

{#if room && membership === 'join'}
	<RoomHeader {room} {direct} />
	<Timeline {room} {draft} />
	<ReadMarker {room} />
	<Typing {room} />
	<Composer {room} {draft} threadId={null} />
{:else if room && membership === 'invite'}
	<InvitePrompt {room} />
{:else}
	<div class="flex flex-1 items-center justify-center p-6">
		<EmptyState icon={glyphs.chat} title={t('room.unavailable')} description={t('room.unavailableBody')} />
	</div>
{/if}
