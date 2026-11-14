<script lang="ts">
	import { Bar, Cursor, Stage, Stream, Window, cursorAt, type Act } from '@facile/mise-en-scene';
	import { Blob } from '@facile/blob';
	import { Icon, Avatar as Picture, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';
	import './hero-demo.css';

	const MENTION = '@pixel';

	const prompt = $derived(t('landing.demo.prompt'));
	const answer = $derived(t('landing.demo.answer'));
	const words = $derived(answer.split(' ').length);
	const rest = $derived(prompt.slice(MENTION.length));

	let draft = $state('');
	let sent = $state(false);
	let thinking = $state(false);
	let shown = $state(0);
	let cursor = $state(cursorAt());
	let send = $state<HTMLElement>();
	let askHeight = $state(0);
	let replyHeight = $state(0);

	const lift = $derived(shown > 0 || thinking ? 0 : sent ? replyHeight : askHeight + replyHeight);

	function reset() {
		draft = '';
		sent = false;
		thinking = false;
		shown = 0;
		cursor.visible = false;
	}

	function poster() {
		sent = true;
		shown = words;
	}

	async function script(act: Act) {
		await act.wait(700);
		await act.type(prompt, (v) => (draft = v), 45);
		const target = act.point(send, 0.55, 0.6);
		await act.move(cursor, { x: target.x - 90, y: target.y + 60 }, 0);
		cursor.visible = true;
		await act.move(cursor, target);
		await act.click(cursor, () => {
			sent = true;
			draft = '';
		});
		await act.wait(300);
		cursor.visible = false;
		thinking = true;
		await act.wait(1100);
		thinking = false;
		await act.count(words, (v) => (shown = v), 70);
		await act.wait(2400);
		reset();
		await act.wait(600);
	}
</script>

{#snippet sidebar()}
	<div class="hd-side">
		<p class="hd-space">{t('landing.demo.space')}</p>
		<div class="hd-row"><Icon icon={glyphs.hashtag} size={14} />general</div>
		<div class="hd-row hd-active"><Icon icon={glyphs.hashtag} size={14} />launch</div>
		<div class="hd-row"><Icon icon={glyphs.hashtag} size={14} />design</div>
		<div class="hd-gap"></div>
		<div class="hd-row"><span class="hd-dot"></span><Bar w="56%" /></div>
		<div class="hd-row"><span class="hd-dot"></span><Bar w="44%" /></div>
	</div>
{/snippet}

{#snippet message(name: string, lines: number[])}
	<div class="hd-msg">
		<Picture {name} class="size-9 shrink-0 text-fc-sm" />
		<div class="hd-body">
			<div class="hd-meta"><b>{name}</b><Bar w={36} h={8} /></div>
			{#each lines as w, i (i)}
				<Bar w="{w}%" />
			{/each}
		</div>
	</div>
{/snippet}

<Stage width={960} height={584} compact={{ below: 640, width: 358, height: 556 }} label={t('landing.demo.label')} {script} {reset} {poster}>
	{#snippet children({ compact })}
		<Window
			sidebar={compact ? undefined : sidebar}
			style={compact ? 'top:12px;left:12px;width:334px;height:516px' : 'top:20px;left:40px;width:880px;height:520px'}
		>
			<div class="hd-main">
				<div class="hd-head">
					<Icon icon={glyphs.hashtag} size={16} class="text-fc-fg-muted" />
					<b>launch</b>
					<span class="hd-fill"></span>
					<Icon icon={glyphs.phone} size={18} class="text-fc-fg-muted" />
					<Icon icon={glyphs.users} size={18} class="text-fc-fg-muted" />
				</div>
				<div class="hd-thread">
					<div class="hd-stack" style="transform:translateY({lift}px)">
						{@render message('Tom', compact ? [76, 40] : [54, 30])}
						{@render message('Léa', compact ? [64] : [44])}
						{@render message('Tom', compact ? [70] : [48])}
						<div class="hd-pad">
							<div class="hd-msg">
								<Picture name="Léa" class="size-9 shrink-0 text-fc-sm" />
								<div class="hd-body">
									<div class="hd-meta"><b>Léa</b><Bar w={36} h={8} /></div>
									<p>{t('landing.demo.first')}</p>
								</div>
							</div>
						</div>
						<div class="hd-pad hd-slot" class:hd-on={sent} bind:offsetHeight={askHeight}>
							<div class="hd-msg">
								<Picture name="Sam" class="size-9 shrink-0 text-fc-sm" />
								<div class="hd-body">
									<div class="hd-meta"><b>Sam</b><span>{t('landing.demo.now')}</span></div>
									<p><span class="hd-mention">{MENTION}</span>{rest}</p>
								</div>
							</div>
						</div>
						<div class="hd-pad hd-slot" class:hd-on={thinking || shown > 0} bind:offsetHeight={replyHeight}>
							<div class="hd-msg">
								<span class="hd-blob"><Blob name="Pixel" size={72} background="circle" alt="" class="size-full" /></span>
								<div class="hd-body">
									<div class="hd-meta"><b>Pixel</b><span class="hd-tag">Bulle</span><span>{t('landing.demo.now')}</span></div>
									<div class="hd-reply">
										<Stream text={answer} count={shown} class="hd-answer" />
										<span class="hd-typing" class:hd-on={thinking}><i></i><i></i><i></i></span>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				<div class="hd-composer">
					<span class="hd-round hd-muted"><Icon icon={icons.plus} size={20} /></span>
					<span class="hd-field" class:hd-placeholder={!draft}>
						{#if draft}
							<span class="hd-mention">{draft.slice(0, MENTION.length)}</span>{draft.slice(MENTION.length)}
						{:else}
							{t('composer.placeholder', { name: '#launch' })}
						{/if}
					</span>
					<span bind:this={send} class="hd-round hd-send" class:hd-idle={!draft}><Icon icon={glyphs.send} size={18} /></span>
				</div>
			</div>
		</Window>
		<Cursor {cursor} />
	{/snippet}
</Stage>
