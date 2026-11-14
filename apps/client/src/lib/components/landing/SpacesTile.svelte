<script lang="ts">
	import { Stage, Tile, Window, type Act } from '@facile/mise-en-scene';
	import { Avatar, Icon, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';

	const MEMBERS = ['Camille', 'Léa', 'Hugo', 'Sam'];

	let joined = $state(1);

	function reset() {
		joined = 1;
	}

	function poster() {
		joined = MEMBERS.length;
	}

	async function script(act: Act) {
		await act.wait(900);
		for (let count = 2; count <= MEMBERS.length; count++) {
			joined = count;
			await act.wait(650);
		}
		await act.wait(2600);
		reset();
		await act.wait(900);
	}
</script>

<Tile title={t('landing.card.spaces.title')} body={t('landing.card.spaces.body')}>
	<Stage width={320} height={240} {script} {reset} {poster}>
		{#snippet children()}
			<Window style="top:20px;left:20px;width:280px;height:204px">
				<div class="sp-body">
					<div class="sp-head">
						<Icon icon={glyphs.users} size={16} />
						<b>{t('landing.demo.space')}</b>
						<span class="sp-chip"><Icon icon={icons.check} size={12} />{t('landing.demo.account')}</span>
					</div>
					{#each MEMBERS as name, i (name)}
						<div class="sp-row" class:sp-off={i >= joined}>
							<Avatar {name} size="sm" class="size-5" />
							<b>{name}</b>
							<span class="sp-muted">{t('landing.demo.member')}</span>
						</div>
					{/each}
				</div>
			</Window>
		{/snippet}
	</Stage>
</Tile>

<style>
	.sp-body {
		display: flex;
		flex-direction: column;
		gap: 5px;
		padding: 8px 12px;
		color: var(--color-fc-fg);
		font-family: var(--font-fc-body);
		font-size: 12px;
	}
	.sp-head {
		display: flex;
		align-items: center;
		gap: 6px;
		height: 28px;
	}
	.sp-chip {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		margin-left: auto;
		padding: 2px 8px;
		border-radius: 999px;
		background: var(--color-fc-accent);
		color: var(--color-fc-accent-fg);
		font-size: 10px;
		white-space: nowrap;
	}
	.sp-row {
		display: flex;
		align-items: center;
		gap: 8px;
		height: 27px;
		padding: 0 8px;
		border-radius: 10px;
		background: var(--color-fc-surface);
		transition:
			opacity 0.4s ease,
			transform 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
	}
	.sp-off {
		opacity: 0;
		transform: translateY(8px);
	}
	.sp-muted {
		margin-left: auto;
		color: var(--color-fc-fg-muted);
		font-size: 11px;
	}
</style>
