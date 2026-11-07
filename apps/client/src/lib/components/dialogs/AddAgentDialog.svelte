<script lang="ts">
	import { Button, EmptyState, Modal, Spinner, icons, toast } from '@facile/muse';
	import { glyphs } from '$lib/glyphs';
	import { prisme, type Agent } from '$lib/api/prisme';
	import { t } from '$lib/i18n/index.svelte';
	import UserName from '$lib/components/common/UserName.svelte';

	let { roomId, onclose }: { roomId: string; onclose: () => void } = $props();
	let agents = $state<Agent[] | null>(null);
	let error = $state('');
	let busyId = $state('');

	$effect(() => {
		prisme.agents().then(
			(list) => (agents = list),
			(err) => (error = err instanceof Error ? err.message : t('addAgent.listFailed'))
		);
	});

	async function add(agent: Agent) {
		busyId = agent.user_id;
		try {
			const result = await prisme.setAgentRooms(agent.user_id, [...agent.room_ids, roomId]);
			agents = (agents ?? []).map((a) => (a.user_id === agent.user_id ? { ...a, room_ids: result.room_ids } : a));
			toast.success(t('addAgent.joined', { name: agent.name }));
		} catch (err) {
			toast.danger(err instanceof Error ? err.message : t('addAgent.failed'));
		} finally {
			busyId = '';
		}
	}
</script>

<Modal open title={t('addAgent.title')} showClose onClose={onclose}>
	{#if error}
		<p class="text-fc-sm text-fc-danger" role="alert">{error}</p>
	{:else if !agents}
		<div class="flex justify-center py-4"><Spinner label={t('addAgent.loading')} /></div>
	{:else if agents.length === 0}
		<EmptyState icon={glyphs.mention} title={t('addAgent.none')} description={t('addAgent.noneBody')} bare />
	{:else}
		<ul class="flex flex-col gap-1">
			{#each agents as agent (agent.user_id)}
				{@const here = agent.room_ids.includes(roomId)}
				<li class="flex min-h-11 items-center gap-3 rounded-fc-md px-2">
					<UserName userId={agent.user_id} name={agent.name} class="flex-1 text-fc-sm font-medium" />
					<Button size="sm" variant={here ? 'ghost' : 'primary'} icon={here ? icons.check : icons.plus} disabled={here || busyId === agent.user_id} onclick={() => add(agent)}>
						{here ? t('addAgent.here') : t('addAgent.add')}
					</Button>
				</li>
			{/each}
		</ul>
	{/if}
</Modal>
