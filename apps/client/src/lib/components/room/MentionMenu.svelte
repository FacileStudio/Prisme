<script lang="ts">
	import type { Mentionable } from '$lib/chat/mentions';
	import { app } from '$lib/state/app.svelte';
	import UserName from '$lib/components/common/UserName.svelte';
	import MxcAvatar from '$lib/components/common/MxcAvatar.svelte';

	type Props = { options: Mentionable[]; selected: number; onpick: (index: number) => void };
	let { options, selected, onpick }: Props = $props();
</script>

<ul
	class="absolute bottom-full left-3 right-3 z-20 mb-1 flex max-h-72 flex-col overflow-y-auto rounded-fc-lg border border-fc-border bg-fc-component p-1 shadow-lg md:left-4 md:right-auto md:w-80"
	role="listbox"
	aria-label="Mention someone"
>
	{#each options as option, index (option.userId)}
		<li role="option" aria-selected={index === selected}>
			<button
				type="button"
				onmousedown={(e) => e.preventDefault()}
				onclick={() => onpick(index)}
				class="flex min-h-11 w-full items-center gap-2 rounded-fc-md px-2 text-left text-fc-sm {index === selected ? 'bg-fc-surface' : 'hover:bg-fc-surface'}"
			>
				<MxcAvatar mxc={app.client?.getUser(option.userId)?.avatarUrl ?? null} name={option.name} class="size-6 text-[10px]" />
				<UserName userId={option.userId} name={option.name} class="shrink-0 font-medium" />
				<span class="ml-auto min-w-0 truncate text-fc-xs text-fc-fg-muted">{option.userId}</span>
			</button>
		</li>
	{/each}
</ul>
