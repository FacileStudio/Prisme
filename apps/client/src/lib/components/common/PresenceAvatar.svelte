<script lang="ts">
	import { app } from '$lib/state/app.svelte';
	import MxcAvatar from './MxcAvatar.svelte';

	type Props = { userId: string; name: string; src?: string; size?: 'sm' | 'md' | 'lg'; class?: string };
	let { userId, name, src = '', size = 'sm', class: className = '' }: Props = $props();

	const user = $derived.by(() => {
		void app.version;
		return userId ? (app.client?.getUser(userId) ?? null) : null;
	});
	const online = $derived(user?.presence === 'online');
</script>

<span class="relative inline-flex shrink-0 {className}">
	<MxcAvatar mxc={user?.avatarUrl ?? null} {src} {name} {size} class={size === 'sm' ? 'size-6 text-[10px]' : ''} />
	<span
		class="absolute -bottom-0.5 -right-0.5 size-2.5 rounded-fc-pill border-2 border-fc-component {online ? 'bg-fc-success' : 'bg-fc-border'}"
		role="img"
		aria-label={online ? 'Online' : 'Offline'}
	></span>
</span>
