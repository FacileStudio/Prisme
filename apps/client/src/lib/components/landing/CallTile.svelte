<script lang="ts">
	import { Stage, Tile, Window, type Act } from '@facile/mise-en-scene';
	import { Avatar, Icon } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';

	const PEOPLE = ['Sam', 'Léa', 'Tom', 'Hugo'];

	let pressed = $state(false);
	let live = $state(false);

	function reset() {
		pressed = false;
		live = false;
	}

	function poster() {
		pressed = true;
		live = true;
	}

	async function script(act: Act) {
		await act.wait(1200);
		pressed = true;
		await act.wait(450);
		live = true;
		await act.wait(3200);
		reset();
		await act.wait(900);
	}
</script>

<Tile title={t('landing.card.call.title')} body={t('landing.card.call.body')}>
	<Stage width={320} height={240} {script} {reset} {poster}>
		{#snippet children()}
			<Window style="top:20px;left:20px;width:280px;height:204px">
				<div class="cl-body">
					<div class="cl-head">
						<Icon icon={glyphs.hashtag} size={14} />
						<b>launch</b>
						<span class="cl-call" class:cl-pressed={pressed}><Icon icon={glyphs.phone} size={14} /></span>
					</div>
					<div class="cl-msg">
						<Avatar name="Sam" size="sm" class="size-6" />
						<span><b>Sam</b><br />{t('landing.demo.callAsk')}</span>
					</div>
					<div class="cl-room" class:cl-off={!live}>
						<div class="cl-grid">
							{#each PEOPLE as name (name)}
								<span class="cl-cell"><Avatar {name} size="sm" class="size-7" /></span>
							{/each}
						</div>
						<div class="cl-bar">
							<span class="cl-lock"><Icon icon={glyphs.lock} size={12} />{t('landing.demo.encrypted')} · Echo</span>
							<span class="cl-btn"><Icon icon={glyphs.mic} size={12} /></span>
							<span class="cl-btn"><Icon icon={glyphs.camera} size={12} /></span>
							<span class="cl-btn cl-end"><Icon icon={glyphs.endCall} size={12} /></span>
						</div>
					</div>
				</div>
			</Window>
		{/snippet}
	</Stage>
</Tile>

<style>
	.cl-body {
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 8px;
		height: 100%;
		padding: 10px 12px;
		color: var(--color-fc-fg);
		font-family: var(--font-fc-body);
		font-size: 12px;
	}
	.cl-head {
		display: flex;
		align-items: center;
		gap: 6px;
		height: 28px;
	}
	.cl-call {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		margin-left: auto;
		border-radius: 999px;
		background: var(--color-fc-surface);
		transition:
			background-color 0.2s ease,
			color 0.2s ease,
			transform 0.2s ease;
	}
	.cl-pressed {
		background: var(--color-fc-accent);
		color: var(--color-fc-accent-fg);
		transform: scale(0.92);
	}
	.cl-msg {
		display: flex;
		align-items: flex-start;
		gap: 8px;
		line-height: 1.4;
	}
	.cl-room {
		position: absolute;
		top: 46px;
		right: 12px;
		bottom: 10px;
		left: 12px;
		display: flex;
		flex-direction: column;
		gap: 6px;
		padding: 8px;
		border-radius: 10px;
		background: var(--color-fc-surface);
		transition:
			opacity 0.4s ease,
			transform 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
	}
	.cl-off {
		opacity: 0;
		transform: translateY(10px) scale(0.97);
	}
	.cl-grid {
		display: grid;
		flex: 1;
		grid-template-columns: 1fr 1fr;
		gap: 6px;
	}
	.cl-cell {
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 8px;
		background: var(--color-fc-component);
	}
	.cl-bar {
		display: flex;
		align-items: center;
		gap: 6px;
		height: 26px;
	}
	.cl-lock {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		margin-right: auto;
		color: var(--color-fc-fg-muted);
		font-size: 11px;
	}
	.cl-btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 26px;
		height: 26px;
		border-radius: 999px;
		background: var(--color-fc-component);
	}
	.cl-end {
		background: var(--color-fc-danger);
		color: var(--color-fc-danger-fg);
	}
</style>
