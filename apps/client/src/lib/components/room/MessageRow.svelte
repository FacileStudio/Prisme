<script lang="ts">
	import type { MatrixEvent, Room } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';
	import type { DraftState } from '$lib/state/draft.svelte';
	import { replyTarget } from '$lib/matrix/timeline';
	import { mentionsUser } from '$lib/chat/mentions';
	import Avatar from '$lib/components/common/Avatar.svelte';
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
		const member = room.getMember(sender) ?? event.sender;
		const me = app.client?.getUserId() ?? '';
		return {
			sender,
			name: member?.name ?? sender,
			avatar: member?.getMxcAvatarUrl() ?? event.sender?.getMxcAvatarUrl() ?? null,
			forMe: sender !== me && !event.isRedacted() && mentionsMe(me),
			time: new Date(event.getTs()).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' }),
			edited: !!event.replacingEventId(),
			pending: event.isSending(),
			reply: replyTarget(event),
			deleted: event.isRedacted()
		};
	});

	function mentionsMe(me: string): boolean {
		if (mentionsUser(event.getContent(), me)) return true;
		return !!app.client?.getPushActionsForEvent(event)?.tweaks?.highlight;
	}
</script>

<div
	class="group relative flex gap-3 px-4 {view.forMe
		? 'bg-fc-warning/10 shadow-[inset_2px_0_0_var(--color-fc-warning)]'
		: 'hover:bg-fc-surface/60 focus-within:bg-fc-surface/60'} {grouped ? 'py-0.5' : 'mt-2 pt-1.5 pb-0.5'}"
	data-event-id={event.getId()}
>
	<div class="w-9 shrink-0">
		{#if !grouped}
			<Avatar userId={view.sender} name={view.name} mxc={view.avatar} size={36} />
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
