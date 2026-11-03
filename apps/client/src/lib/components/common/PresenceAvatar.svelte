<script lang="ts">
	import Avatar from './Avatar.svelte';

	type Props = { userId: string; name: string; src?: string; size?: 'sm' | 'md' | 'lg'; class?: string };
	let { userId, name, src = '', size = 'sm', class: className = '' }: Props = $props();

	const PX = { sm: 24, md: 40, lg: 56 } as const;
</script>

<Avatar {userId} {name} {src} size={PX[size]} presence class={className} />
