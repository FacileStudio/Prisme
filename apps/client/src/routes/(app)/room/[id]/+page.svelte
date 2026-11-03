<script lang="ts">
	import { page } from '$app/state';
	import RoomView from '$lib/components/room/RoomView.svelte';

	const roomId = $derived(decodeURIComponent(page.params.id ?? ''));
</script>

{#key roomId}
	<RoomView {roomId} />
{/key}
