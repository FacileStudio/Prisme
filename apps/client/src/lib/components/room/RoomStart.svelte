<script lang="ts">
	import type { Room } from 'matrix-js-sdk';
	import { Button, Icon, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { app } from '$lib/state/app.svelte';
	import { ui } from '$lib/state/ui.svelte';
	import { call } from '$lib/call/call.svelte';
	import { sendText } from '$lib/matrix/messages';
	import { isDirect } from '$lib/matrix/rooms';
	import { t } from '$lib/i18n/index.svelte';

	let { room, empty }: { room: Room; empty: boolean } = $props();
	let busy = $state(false);

	const info = $derived.by(() => {
		void app.version;
		const direct = !!app.client && isDirect(app.client, room);
		const topic = room.currentState.getStateEvents('m.room.topic', '')?.getContent().topic as string | undefined;
		return { direct, topic: topic?.trim() ?? '' };
	});
	const title = $derived(info.direct ? t('channel.startDirect', { name: room.name }) : t('channel.startTitle', { name: room.name }));
	const body = $derived(info.topic || (info.direct ? t('channel.startDirectBody') : t('channel.startBody')));

	async function run(work: () => Promise<unknown>) {
		busy = true;
		try {
			await work();
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : String(err));
		} finally {
			busy = false;
		}
	}

	const sayHello = () => run(() => sendText(app.client!, { roomId: room.roomId }, { text: t('channel.hello'), mentions: [] }));
	const startCall = () => run(() => call.start(room.roomId));
</script>

<section
	class="mx-4 flex flex-col gap-4 rounded-fc-lg bg-fc-component p-5 sm:p-6 {empty ? 'my-auto items-center py-10 text-center sm:py-14' : 'mt-4 mb-2'}"
	aria-label={title}
>
	<span class="flex size-14 shrink-0 items-center justify-center rounded-fc-lg bg-fc-accent/10 text-fc-accent">
		<Icon icon={info.direct ? glyphs.chat : glyphs.hashtag} size={28} />
	</span>
	<div class="flex max-w-md flex-col gap-1">
		<h2 class="text-fc-lg font-semibold text-balance text-fc-fg">{title}</h2>
		<p class="text-fc-sm text-pretty text-fc-fg-muted">{body}</p>
	</div>
	{#if empty}
		<div class="flex flex-wrap justify-center gap-2">
			<Button icon={glyphs.smile} disabled={busy || !app.client} onclick={sayHello}>{t('channel.sayHi')}</Button>
			{#if !info.direct}
				<Button variant="secondary" icon={glyphs.userPlus} onclick={() => (ui.dialog = 'invite')}>{t('channel.inviteHere')}</Button>
			{/if}
			{#if app.session?.calls_enabled && !call.active}
				<Button variant="ghost" icon={glyphs.phone} disabled={busy} onclick={startCall}>{t('call.start')}</Button>
			{/if}
		</div>
	{/if}
</section>
