<script lang="ts">
	import { Stage, Tile, Window, type Act } from '@facile/mise-en-scene';
	import { Blob } from '@facile/blob';
	import { Avatar, Icon, icons } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { t } from '$lib/i18n/index.svelte';

	const MENTION = '@pixel';
	const FILE = 'visuels-lancement.zip';

	const ask = $derived(t('landing.demo.ask'));
	const rows = $derived([
		{ icon: glyphs.tasks, name: t('landing.demo.taskPricing'), from: `Opus · ${t('landing.demo.friday')}` },
		{ icon: glyphs.tasks, name: t('landing.demo.taskNotes'), from: `Opus · ${t('landing.demo.friday')}` },
		{ icon: icons.folder, name: FILE, from: 'Nuage' }
	]);

	let shown = $state(0);

	function reset() {
		shown = 0;
	}

	function poster() {
		shown = rows.length + 1;
	}

	async function script(act: Act) {
		await act.wait(1000);
		for (let count = 1; count <= rows.length + 1; count++) {
			shown = count;
			await act.wait(550);
		}
		await act.wait(2600);
		reset();
		await act.wait(900);
	}
</script>

<Tile title={t('landing.card.bulle.title')} body={t('landing.card.bulle.body')}>
	{#snippet badge()}<span class="bl-badge">Pro</span>{/snippet}
	<Stage width={320} height={240} {script} {reset} {poster}>
		{#snippet children()}
			<Window style="top:20px;left:20px;width:280px;height:204px">
				<div class="bl-body">
					<div class="bl-ask">
						<Avatar name="Sam" size="sm" class="size-5" />
						<p><span class="bl-mention">{MENTION}</span>{ask.slice(MENTION.length)}</p>
					</div>
					<div class="bl-reply" class:bl-off={shown < 1}>
						<div class="bl-who">
							<span class="bl-blob"><Blob name="Pixel" size={48} background="circle" alt="" class="size-full" /></span>
							<b>Pixel</b>
							<span class="bl-tag">Bulle</span>
						</div>
						{#each rows as row, i (row.name)}
							<div class="bl-row" class:bl-off={shown < i + 2}>
								<Icon icon={row.icon} size={14} />
								<span class="bl-name">{row.name}</span>
								<span class="bl-muted">{row.from}</span>
							</div>
						{/each}
					</div>
				</div>
			</Window>
		{/snippet}
	</Stage>
</Tile>

<style>
	.bl-badge {
		padding: 2px 8px;
		border-radius: 999px;
		background: var(--color-fc-accent);
		color: var(--color-fc-accent-fg);
		font-size: var(--text-fc-xs);
		font-weight: 600;
	}
	.bl-body {
		display: flex;
		flex-direction: column;
		gap: 6px;
		padding: 8px 12px;
		color: var(--color-fc-fg);
		font-family: var(--font-fc-body);
		font-size: 12px;
	}
	.bl-ask {
		display: flex;
		align-items: center;
		gap: 8px;
		height: 26px;
	}
	.bl-ask p {
		margin: 0;
		white-space: nowrap;
	}
	.bl-mention {
		border-radius: var(--radius-fc-sm);
		background: color-mix(in oklch, var(--color-fc-accent) 12%, transparent);
	}
	.bl-reply {
		display: flex;
		flex-direction: column;
		gap: 4px;
		padding: 6px;
		border-radius: 10px;
		background: var(--color-fc-surface);
	}
	.bl-who {
		display: flex;
		align-items: center;
		gap: 6px;
		height: 22px;
	}
	.bl-blob {
		width: 22px;
		height: 22px;
	}
	.bl-tag {
		padding: 0 6px;
		border-radius: 999px;
		background: var(--color-fc-component);
		color: var(--color-fc-fg-muted);
		font-size: 10px;
	}
	.bl-row {
		display: flex;
		align-items: center;
		gap: 6px;
		height: 26px;
		padding: 0 8px;
		border-radius: 8px;
		background: var(--color-fc-component);
	}
	.bl-reply,
	.bl-row {
		transition:
			opacity 0.4s ease,
			transform 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
	}
	.bl-off {
		opacity: 0;
		transform: translateY(8px);
	}
	.bl-name {
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
	}
	.bl-muted {
		margin-left: auto;
		color: var(--color-fc-fg-muted);
		font-size: 11px;
		white-space: nowrap;
	}
</style>
