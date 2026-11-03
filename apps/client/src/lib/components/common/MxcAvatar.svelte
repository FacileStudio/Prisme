<script lang="ts">
	import Avatar from './Avatar.svelte';

	type Props = { mxc?: string | null; src?: string; name: string; size?: 'sm' | 'md' | 'lg'; class?: string };
	let { mxc = null, src = '', name, size = 'sm', class: className = '' }: Props = $props();

	const PX = { sm: 32, md: 40, lg: 56 } as const;
</script>

<Avatar {mxc} {src} {name} size={PX[size]} class={className} />
