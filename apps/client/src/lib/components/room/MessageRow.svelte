<script lang="ts">
	import type { MatrixEvent, Room } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';
	import type { DraftState } from '$lib/state/draft.svelte';
	import { replyTarget } from '$lib/matrix/timeline';
	import MxcAvatar from '$lib/components/common/MxcAvatar.svelte';
	import UserName from '$lib/components/common/UserName.svelte';
	import MessageBody from './MessageBody.svelte';
	import MessageActions from './MessageActions.svelte';
	import Reactions from './Reactions.svelte';
	import ReplyQuote from './ReplyQuote.svelte';
	import ThreadSummary from './ThreadSummary.svelte';

	type Props = { room: Room; event: MatrixEvent; grouped: boolean; draft: DraftState; inThread?: boolean };
	let { room, event, grouped, draft, inThread = false }: Props = $props();

	const view = $derived.by(() => {
		void app.version;
		const sender = event.getSender() ?? '';
		const member = room.getMember(sender);
		return {
			sender,
			name: member?.name ?? sender,
			avatar: member?.getMxcAvatarUrl() ?? null,
			time: new Date(event.getTs()).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' }),
			edited: !!event.replacingEventId(),
			pending: event.isSending(),
			reply: replyTarget(event),
			deleted: event.isRedacted()
		};
	});
</script>

<div
	class="group relative flex gap-3 px-4 hover:bg-fc-surface/60 focus-within:bg-fc-surface/60 {grouped ? 'py-0.5' : 'mt-2 pt-1.5 pb-0.5'}"
	data-event-id={event.getId()}
>
	<div class="w-9 shrink-0">
		{#if !grouped}
			<MxcAvatar mxc={view.avatar} name={view.name} class="size-9" />
		{/if}
	</div>
	<div class="min-w-0 flex-1" class:opacity-60={view.pending}>
		{#if !grouped}
			<div class="flex items-baseline gap-2">
				<UserName userId={view.sender} name={view.name} class="text-fc-sm font-semibold" />
				<time class="shrink-0 text-fc-xs text-fc-fg-muted" datetime={new Date(event.getTs()).toISOString()}>{view.time}</time>
			</div>
		{/if}
		{#if view.reply}
			<ReplyQuote {room} eventId={view.reply} />
		{/if}
		<MessageBody {event} />
		{#if view.edited && !view.deleted}
			<span class="text-fc-xs text-fc-fg-muted">(edited)</span>
		{/if}
		{#if !view.deleted}
			<Reactions {room} {event} />
			{#if !inThread}
				<ThreadSummary {room} {event} />
			{/if}
		{/if}
	</div>
	{#if !view.deleted && !view.pending}
		<MessageActions {room} {event} {draft} {inThread} />
	{/if}
</div>
