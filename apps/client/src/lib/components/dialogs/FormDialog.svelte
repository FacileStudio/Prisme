<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Button, Modal } from '@facile/muse';
	import { t } from '$lib/i18n/index.svelte';

	type Props = {
		title: string;
		submitLabel: string;
		busy?: boolean;
		disabled?: boolean;
		icon?: string;
		onsubmit: () => void | Promise<void>;
		onclose: () => void;
		children: Snippet;
	};
	let { title, submitLabel, busy = false, disabled = false, icon, onsubmit, onclose, children }: Props = $props();
	const formId = `form-${Math.random().toString(36).slice(2)}`;
</script>

<Modal open {title} onClose={onclose}>
	<form id={formId} class="flex flex-col gap-4" onsubmit={(e) => (e.preventDefault(), onsubmit())}>
		{@render children()}
	</form>
	{#snippet footer()}
		<Button variant="ghost" onclick={onclose}>{t('common.cancel')}</Button>
		<Button type="submit" form={formId} {icon} disabled={busy || disabled}>{busy ? t('common.working') : submitLabel}</Button>
	{/snippet}
</Modal>
