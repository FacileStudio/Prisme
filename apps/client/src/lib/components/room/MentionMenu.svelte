<script lang="ts">
	import type { Mentionable } from '$lib/chat/mentions';
	import { handle } from '$lib/chat/users';
	import UserName from '$lib/components/common/UserName.svelte';
	import Avatar from '$lib/components/common/Avatar.svelte';

	type Props = {
		id: string;
		options: Mentionable[];
		selected: number;
		onpick: (index: number) => void;
		onhover: (index: number) => void;
	};
	let { id, options, selected, onpick, onhover }: Props = $props();
	let list: HTMLUListElement | undefined = $state();

	$effect(() => {
		list?.querySelector(`#${id}-${selected}`)?.scrollIntoView({ block: 'nearest' });
	});
</script>

<ul
	bind:this={list}
	{id}
	class="absolute bottom-full left-3 right-3 z-20 mb-1 flex max-h-72 flex-col overflow-y-auto overscroll-contain rounded-fc-lg border border-fc-border bg-fc-component p-1 shadow-lg md:left-4 md:right-auto md:w-80"
	role="listbox"
	aria-label="Mention someone"
>
	{#each options as option, index (option.userId)}
		<li
			id="{id}-{index}"
			role="option"
			aria-selected={index === selected}
			onpointerdown={(e) => e.preventDefault()}
			onpointermove={() => index !== selected && onhover(index)}
			onclick={() => onpick(index)}
			onkeydown={(e) => e.key === 'Enter' && onpick(index)}
			class="flex min-h-11 cursor-pointer items-center gap-2 rounded-fc-md px-2 text-fc-sm {index === selected ? 'bg-fc-surface' : ''}"
		>
			<Avatar userId={option.userId} name={option.name} size={28} />
			<UserName userId={option.userId} name={option.name} class="shrink-0 font-medium" />
			<span class="ml-auto min-w-0 truncate text-fc-xs text-fc-fg-muted">{handle(option.userId)}</span>
		</li>
	{/each}
</ul>
