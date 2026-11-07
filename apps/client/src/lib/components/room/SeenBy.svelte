<script lang="ts">
	import type { MatrixEvent, Room } from 'matrix-js-sdk';
	import { app } from '$lib/state/app.svelte';
	import { t } from '$lib/i18n/index.svelte';
	import { formatList } from '$lib/i18n/format';
	import { seenBy } from '$lib/matrix/timeline';
	import Avatar from '$lib/components/common/Avatar.svelte';

	let { room, event }: { room: Room; event: MatrixEvent } = $props();

	const readers = $derived.by(() => {
		void app.version;
		return seenBy(room, event, app.client?.getUserId() ?? '').map((id) => {
			const member = room.getMember(id);
			return { id, name: member?.name ?? id, avatar: member?.getMxcAvatarUrl() ?? null };
		});
	});
	const label = $derived(t('seen.byNames', { names: formatList(readers.map((r) => r.name)) }));
</script>

{#if readers.length > 0}
	<div class="flex items-center justify-end gap-1 px-4 pt-1" title={label} aria-label={label} role="img">
		<span class="text-fc-xs text-fc-fg-muted">{t('seen.by')}</span>
		{#each readers.slice(0, 5) as reader (reader.id)}
			<Avatar userId={reader.id} mxc={reader.avatar} name={reader.name} size={20} />
		{/each}
		{#if readers.length > 5}<span class="text-fc-xs text-fc-fg-muted">+{readers.length - 5}</span>{/if}
	</div>
{/if}
