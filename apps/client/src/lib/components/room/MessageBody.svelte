<script lang="ts">
	import { MsgType, type MatrixEvent } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';
	import { sanitizeMessage } from '$lib/chat/sanitize';
	import { escapeHtml, stripReplyFallback } from '$lib/chat/markdown';
	import { callInfo, capsuleInfo } from '$lib/matrix/timeline';
	import { suiteLinks } from '$lib/chat/links';
	import { pollInfo } from '$lib/chat/poll';
	import { voiceInfo } from '$lib/matrix/voice';
	import { t } from '$lib/i18n/index.svelte';
	import CallCard from './CallCard.svelte';
	import CapsuleCard from './CapsuleCard.svelte';
	import ImageMessage from './ImageMessage.svelte';
	import FileMessage from './FileMessage.svelte';
	import LinkCard from './LinkCard.svelte';
	import PollCard from './PollCard.svelte';
	import VoiceMessage from './VoiceMessage.svelte';

	let { event }: { event: MatrixEvent } = $props();

	const content = $derived.by(() => {
		void app.version;
		return event.getContent();
	});
	const call = $derived(callInfo(event));
	const capsule = $derived(capsuleInfo(event));
	const cards = $derived(suiteLinks(stripReplyFallback(String(content.body ?? ''))));
	const poll = $derived.by(() => {
		void app.version;
		return pollInfo(event);
	});
	const voice = $derived(voiceInfo(event));
	const mine = $derived(event.getSender() === app.client?.getUserId());
	const html = $derived.by(() => {
		const me = app.client?.getUserId() ?? '';
		if (content.format === 'org.matrix.custom.html' && typeof content.formatted_body === 'string') {
			return sanitizeMessage(content.formatted_body, me);
		}
		const body = stripReplyFallback(String(content.body ?? ''));
		return escapeHtml(body).replace(/\n/g, '<br>');
	});
</script>

{#if event.isRedacted()}
	<p class="text-fc-sm italic text-fc-fg-muted">{t('message.deleted')}</p>
{:else if capsule}
	<CapsuleCard {capsule} {mine} />
{:else if poll}
	<PollCard {event} {poll} />
{:else if voice}
	<VoiceMessage {voice} {mine} />
{:else if call}
	<CallCard {call} roomId={event.getRoomId() ?? ''} />
{:else if content.msgtype === MsgType.Image}
	<ImageMessage {content} />
{:else if content.msgtype === MsgType.File || content.msgtype === MsgType.Video || content.msgtype === MsgType.Audio}
	<FileMessage {content} />
{:else}
	<div class="prose-msg text-fc-sm leading-relaxed {content.msgtype === MsgType.Emote ? 'italic' : ''} {content.msgtype === MsgType.Notice ? 'text-fc-fg-muted' : ''}">
		{@html html}
	</div>
	{#each cards as link (link.url)}
		<LinkCard {link} />
	{/each}
{/if}

<style>
	.prose-msg :global(a.mention-pill) {
		text-decoration: none;
		border-radius: var(--radius-fc-pill);
		padding: 0 0.35rem;
		font-weight: 500;
		color: var(--color-fc-fg);
		background: color-mix(in oklch, var(--color-fc-accent) 12%, transparent);
	}

	.prose-msg :global(a.mention-pill.mention-me) {
		color: var(--color-fc-accent-fg);
		background: var(--color-fc-accent);
	}
</style>
